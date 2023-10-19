package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Command          string
	WorkingDirectory string
	Args             []string
}

func CreateDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	return nil
}

func ExecuteWithNoOutput(command Command) (string, error) {
	cmd := exec.Command(command.Command, command.Args...)
	if command.WorkingDirectory != "" {
		cmd.Dir = command.WorkingDirectory
	}

	var stdOut, stdIn, stderr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stderr
	cmd.Stdin = &stdIn

	if err := cmd.Run(); err != nil {
		if stderr.String() != "" {
			return stdOut.String(), fmt.Errorf("%v, err: %v", stderr.String(), err.Error())
		} else {
			return stdOut.String(), fmt.Errorf("empty output, err: %v", err.Error())
		}
	}

	return stdOut.String(), nil
}

func RemoveFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to remove folder: %v", err)
	}

	return nil
}

func MoveFolder(source string, destination string) error {
	// if !helper.DirectoryExists(destination) {
	// 	return fmt.Errorf("destination folder does not exist")
	// }

	// if !helper.DirectoryExists(source) {
	// 	return fmt.Errorf("source folder does not exist")
	// }

	err := os.Rename(source, destination)
	if err != nil {
		return err
	}
	return nil
}

func ToHCL(m map[string]interface{}, indent int) string {
	var lines []string
	for k, v := range m {
		switch v := v.(type) {
		case string:
			lines = append(lines, fmt.Sprintf("%s = \"%s\"", k, v))
		case bool:
			lines = append(lines, fmt.Sprintf("%s = %t", k, v))
		case []string:
			lines = append(lines, fmt.Sprintf("%s = [", k))
			for idx, item := range v {
				line := strings.Repeat(" ", 2*(indent+1))
				line = fmt.Sprintf("%v\"%s\"", line, item)
				if idx >= 0 && idx < len(v)-1 {
					line = fmt.Sprintf("%s,", line)
				}
				lines = append(lines, line)
			}
			if indent > 0 {
				lines = append(lines, fmt.Sprintf("%s%s", strings.Repeat(" ", 2*(indent)), "]"))
			} else {
				lines = append(lines, "]")
			}

		case []interface{}:
			lines = append(lines, fmt.Sprintf("%s = [", k))
			for idx, item := range v {
				line := ""
				switch item := item.(type) {
				case string:
					line = fmt.Sprintf("\"%s\"", item)
				case bool:
					line = fmt.Sprintf("%t", item)
				case map[string]interface{}:
					line = ToHCL(item, indent+1)
				default:
					line = fmt.Sprintf("%v", item)
				}
				if idx >= 0 && idx < len(v)-1 {
					line = fmt.Sprintf("%s,", line)
				}
				lines = append(lines, fmt.Sprintf("%s%s", strings.Repeat(" ", 2*(indent+1)), line))
			}
			if indent > 0 {
				lines = append(lines, fmt.Sprintf("%s%s", strings.Repeat(" ", 2*(indent)), "],"))
			} else {
				lines = append(lines, "]")
			}
		case map[string]interface{}:
			lines = append(lines, fmt.Sprintf("%s = {", k))
			count := 0
			for k2, v2 := range v {
				line := strings.Repeat(" ", 2*(indent+1))
				line = fmt.Sprintf("%v%v", line, ToHCL(map[string]interface{}{k2: v2}, indent+1))
				if count >= 0 && count < len(v)-1 {
					line = fmt.Sprintf("%s,", line)
				}
				lines = append(lines, line)
				count = count + 1
			}
			if indent > 0 {
				lines = append(lines, fmt.Sprintf("%s%s", strings.Repeat(" ", 2*(indent)), "}"))
			} else {
				lines = append(lines, "}")
			}
		default:
			lines = append(lines, fmt.Sprintf("%s = %v", k, v))
		}
	}
	result := strings.Join(lines, "\n")
	if strings.HasSuffix(result, ",") {
		result = result[:len(result)-1]
	}
	return result
}
