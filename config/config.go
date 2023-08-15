package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App    AppCfg        `yaml:"app"`
	Store  StorageCfg    `yaml:"store"`
	Db     DbCfg         `yaml:"database"`
	DocuAI DocumentAICfg `yaml:"document-ai"`
}

type AppCfg struct {
	Debug  bool   `yaml:"debug"`
	Secret string `yaml:"secret"`
}

type StorageCfg struct {
	Driver   string `yaml:"driver"`
	Location string `yaml:"location"`
}

type DbCfg struct {
	Driver string `yaml:"driver"`
	Name   string `yaml:"name"`
}

type DocumentAICfg struct {
	ProjectId   string `yaml:"project-id"`
	ProcessorId string `yaml:"processor-id"`
	Location    string `yaml:"location"`
	CredsFile   string `yaml:"credsfile"`
}

var Store StorageCfg
var App AppCfg
var Db DbCfg
var DocumentAI DocumentAICfg

const CONFIG_PATH = "config.yml"

func Init(cfg Config) {
	Store = cfg.Store
	App = cfg.App
	Db = cfg.Db
	DocumentAI = cfg.DocuAI
}

func LoadConfig() (Config, error) {
	var cfg Config
	cfgFile, err := os.Open(CONFIG_PATH)
	if err != nil {
		return cfg, err
	}
	defer cfgFile.Close()

	d := yaml.NewDecoder(cfgFile)
	err = d.Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil

}
