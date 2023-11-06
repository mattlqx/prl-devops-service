package models

type CatalogManifest struct {
	Name               string                        `json:"name"`
	ID                 string                        `json:"id"`
	Type               string                        `json:"type"`
	Tags               []string                      `json:"tags,omitempty"`
	Size               string                        `json:"size,omitempty"`
	Provider           *RemoteVirtualMachineProvider `json:"provider,omitempty"`
	CreatedAt          string                        `json:"created_at,omitempty"`
	UpdatedAt          string                        `json:"updated_at,omitempty"`
	RequiredClaims     []string                      `json:"required_claims,omitempty"`
	RequiredRoles      []string                      `json:"required_roles,omitempty"`
	LastDownloadedAt   string                        `json:"last_downloaded_at,omitempty"`
	LastDownloadedUser string                        `json:"last_downloaded_user,omitempty"`
	PackContents       []CatalogManifestPackItem     `json:"pack_contents,omitempty"`
}

type RemoteVirtualMachineProvider struct {
	Type     string            `json:"type,omitempty"`
	Host     string            `json:"host,omitempty"`
	Port     string            `json:"port,omitempty"`
	Username string            `json:"user,omitempty"`
	Password string            `json:"password,omitempty"`
	ApiKey   string            `json:"api_key,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
}

type CatalogManifestPackItem struct {
	IsDir bool   `json:"is_dir,omitempty"`
	Name  string `json:"name,omitempty"`
	Path  string `json:"path,omitempty"`
}

type PullCatalogManifestResponse struct {
	ID          string           `json:"id"`
	LocalPath   string           `json:"local_path"`
	MachineName string           `json:"machine_name"`
	Manifest    *CatalogManifest `json:"manifest"`
}

type ImportCatalogManifestResponse struct {
	ID string `json:"id"`
}