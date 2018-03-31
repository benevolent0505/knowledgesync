package command

import (
	"os"
	"io"
	"io/ioutil"
	"path/filepath"
	"fmt"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host string
	PrivateToken string
	LocalRoot string
}

type knowledgeConfig struct {
	Knowledge struct {
		Host string `yaml:"host"`
		PrivateToken string `yaml:"private_token"`
	}
	LocalRoot string `yaml:"local_root"`
}

func loadConfig(r io.Reader) (*Config, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	knowledges := knowledgeConfig{}
	err = yaml.Unmarshal(bytes, &knowledges)
	if err != nil {
		return nil, err
	}
	c := &Config{
		Host: knowledges.Knowledge.Host,
		PrivateToken: knowledges.Knowledge.PrivateToken,
		LocalRoot: knowledges.LocalRoot,
	}

	return c, nil
}

func loadSingleConfigFile(fname string) (*Config, error) {
	if _, err := os.Stat(fname); err != nil {
		return nil, err
	}
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return loadConfig(f)
}

func LoadConfigFile() (*Config, error) {
	var conf *Config
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	conf, err = loadSingleConfigFile(filepath.Join(pwd, "config.yml"))
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, fmt.Errorf("no config files found")
	}
	return conf, nil
}
