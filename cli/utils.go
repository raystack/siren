package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/odpf/siren/pkg/errors"
	"gopkg.in/yaml.v3"
)

// TODO need to test this
func parseFile(filePath string, v interface{}) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	switch filepath.Ext(filePath) {
	case ".json":
		if err := json.Unmarshal(b, v); err != nil {
			return fmt.Errorf("invalid json: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(b, v); err != nil {
			return fmt.Errorf("invalid yaml: %w", err)
		}
	default:
		return errors.New("unsupported file type")
	}

	return nil
}
