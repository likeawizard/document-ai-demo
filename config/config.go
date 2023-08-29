package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	SCHEMA_DOC_INT     = "docu-intel"
	SCHEMA_DOCUMENT_AI = "document-ai"
)

const (
	CUR_CURR_API = "currencyapi"
)

type Config struct {
	App       AppCfg        `yaml:"app"`
	Store     StorageCfg    `yaml:"store"`
	Db        DbCfg         `yaml:"database"`
	Currency  CurrencyCfg   `yaml:"currency"`
	DocuAI    DocumentAICfg `yaml:"document-ai"`
	DocuIntel DocuIntelCfg  `yaml:"docu-intel"`
	Processor ProcessorCfg
}

type AppCfg struct {
	Debug           bool   `yaml:"debug"`
	Secret          string `yaml:"secret"`
	ProcessorDriver string `yaml:"processor-driver"`
}

type StorageCfg struct {
	Driver   string `yaml:"driver"`
	Location string `yaml:"location"`
}

type DbCfg struct {
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type CurrencyCfg struct {
	Service  string `yaml:"service"`
	Endpoint string `yaml:"endpoint"`
	AuthKey  string `yaml:"auth-key"`
}

type ProcessorCfg interface {
	Driver() string
}

type DocumentAICfg struct {
	ProjectId   string `yaml:"project-id"`
	ProcessorId string `yaml:"processor-id"`
	Location    string `yaml:"location"`
	CredsFile   string `yaml:"credsfile"`
}

type DocuIntelCfg struct {
	Endpoint   string `yaml:"endpoint"`
	Key        string `yaml:"key"`
	ModelId    string `yaml:"model-id"`
	ApiVersion string `yaml:"api-version"`
}

const CONFIG_PATH = "config.yml"

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

	switch cfg.App.ProcessorDriver {
	case SCHEMA_DOCUMENT_AI:
		cfg.Processor = &cfg.DocuAI
	case SCHEMA_DOC_INT:
		cfg.Processor = &cfg.DocuIntel
	}

	return cfg, nil

}

func (cfg *DocumentAICfg) Driver() string {
	return SCHEMA_DOCUMENT_AI
}

func (cfg *DocuIntelCfg) Driver() string {
	return SCHEMA_DOC_INT
}
