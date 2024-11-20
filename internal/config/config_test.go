package config_test

import (
	"strings"
	"testing"

	"github.com/gruyaume/lesvieux/internal/config"
)

func TestGoodConfigSuccess(t *testing.T) {
	conf, err := config.Validate("testdata/valid.yaml")
	if err != nil {
		t.Fatalf("Error occurred: %s", err)
	}

	if conf.DBPath == "" {
		t.Fatalf("No database path was configured for server")
	}

	if conf.Port != 8000 {
		t.Fatalf("Port was not configured correctly")
	}
}

func TestBadConfigFail(t *testing.T) {
	cases := []struct {
		Name               string
		ConfigYAMLFilePath string
		ExpectedError      string
	}{
		{"no db path", "testdata/invalid_no_db.yaml", "`db_path` is empty"},
		{"invalid yaml", "testdata/invalid_yaml.yaml", "unmarshal errors"},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := config.Validate(tc.ConfigYAMLFilePath)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			if !strings.Contains(err.Error(), tc.ExpectedError) {
				t.Errorf("Expected error: %s, got: %s", tc.ExpectedError, err)
			}
		})
	}
}
