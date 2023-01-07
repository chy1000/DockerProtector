package rules

import (
	path2 "DockerProtector/lib/path"
	"gopkg.in/yaml.v3"
	"os"
)

type Rule struct{
	AllowIp string `yaml:"ip"`
	AllowPort int `yaml:"port"`
}

type Rules map[string][]Rule

func ReadYamlConfig() (*Rules, error) {
	path := path2.GetCurrentAbPath() + "rules.yaml"
	var rules Rules
	file, err := os.Open(path)
	if  err != nil {
		return nil, err
	}
	defer file.Close()
	yaml.NewDecoder(file).Decode(&rules)
	return &rules, nil
}