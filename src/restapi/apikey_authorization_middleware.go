package restapi

import (
	"context"
	"encoding/base64"

	"github.com/Parallels/pd-api-service/basecontext"
	"github.com/Parallels/pd-api-service/constants"
	"github.com/Parallels/pd-api-service/errors"
	"github.com/Parallels/pd-api-service/helpers"
	"github.com/Parallels/pd-api-service/models"
	"github.com/Parallels/pd-api-service/serviceprovider"

	"net/http"
	"strings"
)

type ApiKeyHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ApiKeyAuthorizationMiddlewareAdapter(roles []string, claims []string) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			baseCtx := basecontext.NewBaseContextFromRequest(r)
			var authorizationContext *basecontext.AuthorizationContext
			authCtxFromRequest := baseCtx.GetAuthorizationContext()
			if authCtxFromRequest == nil {
				authorizationContext = basecontext.InitAuthorizationContext()
			} else {
				authorizationContext = authCtxFromRequest
			}

			// If the authorization context is already authorized we will skip this middleware
			if authorizationContext.IsAuthorized || HasAuthorizationHeader(r) {
				baseCtx.LogDebug("No Api Key was found in the request, skipping")
				next.ServeHTTP(w, r)
				return
			}

			baseCtx.LogInfo("ApiKey Authorization layer started")
			authError := models.OAuthErrorResponse{
				Error:            models.OAuthUnauthorizedClient,
				ErrorDescription: "The Api Key is not valid",
			}

			apiKey, err := extractApiKey(r.Header)
			if err != nil {
				authError.ErrorDescription = err.Error()
				authorizationContext.AuthorizationError = &authError
				baseCtx.LogInfo("No Api Key was found in the request, skipping")
				next.ServeHTTP(w, r)
				return
			}
			isValid := true
			db := serviceprovider.Get().JsonDatabase
			_ = db.Connect(baseCtx)
			dbApiKey, err := db.GetApiKey(baseCtx, apiKey.Key)

			if err != nil || dbApiKey == nil {
				isValid = false
			}

			if isValid {
				if dbApiKey.Revoked {
					isValid = false
					authError.ErrorDescription = "Api Key has been revoked"
				}
			}
			if isValid {
				hashedSecret := helpers.Sha256Hash(apiKey.Value)
				if dbApiKey.Secret != hashedSecret {
					isValid = false
					authError.ErrorDescription = "Api Key is not Valid"
				}
			}

			if !isValid {
				baseCtx.LogInfo("The Api Key is not valid")
				authorizationContext.IsAuthorized = false
				authorizationContext.AuthorizationError = &authError

				ctx := context.WithValue(r.Context(), constants.AUTHORIZATION_CONTEXT_KEY, authorizationContext)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			authorizationContext.IsAuthorized = true
			authorizationContext.IsMicroService = true
			authorizationContext.AuthorizedBy = "ApiKeyAuthorization"
			authorizationContext.AuthorizationError = nil
			ctx := context.WithValue(r.Context(), constants.AUTHORIZATION_CONTEXT_KEY, authorizationContext)
			baseCtx.LogInfo("ApiKey Authorization layer finished")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractApiKey(headers http.Header) (*ApiKeyHeader, error) {
	authHeader := headers.Get("X-Api-Key")
	if authHeader == "" {
		err := errors.New("No Api Key was found in the request")
		return nil, err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(authHeader)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(decodedKey), ":")
	if len(parts) != 2 {
		err := errors.New("The Api Key is not in the correct format")
		return nil, err
	}

	return &ApiKeyHeader{
		Key:   parts[0],
		Value: parts[1],
	}, nil
}
