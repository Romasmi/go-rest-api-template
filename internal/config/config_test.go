package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	const ConfigPath = "/tmp"
	const ConfigFile = "config.yaml"
	const OverrideConfigFile = "override.yaml"

	configFile := `
database:
  url: postgres://user:pass@localhost:5432/app
  maxConnections: 10
  minConnection: 2
  maxConnectionLifetime: 60 # seconds
  maxConnectionIdleTime: 1 # minutes
server:
  port: 8080
`

	err := createConfigFile(ConfigPath+"/"+ConfigFile, []byte(configFile))
	if err != nil {
		t.Fatalf("Error while %v creation: %v", ConfigFile, err)
	}

	overrideFile := `
server:
  port: 80
`

	err = createConfigFile(ConfigPath+"/"+OverrideConfigFile, []byte(overrideFile))
	if err != nil {
		t.Fatalf("Error while %v creation: %v", OverrideConfigFile, err)
	}

	config, err := LoadConfig("/tmp")
	if err != nil {
		t.Fatalf("Error while loading config: %v", err)
	}

	expected := uint(10)
	if config.Database.MaxConnections != expected {
		t.Errorf("Wrong config value,  expected: %v, actual: %v", expected, config.Database.MaxConnections)
	}

	expected = uint(80)
	if config.Server.Port != expected {
		t.Errorf("Wrong overide value,  expected: %v, actual: %v", expected, config.Server.Port)
	}
}

func createConfigFile(filePath string, content []byte) error {
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}
