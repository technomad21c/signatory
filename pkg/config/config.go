package config

import (
	"io/ioutil"
	"net"

	"gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v3"
)

// ServerConfig contains the information necessary to the tezos signing server
type ServerConfig struct {
	Address        string `yaml:"address" validate:"required,hostport"`
	UtilityAddress string `yaml:"utility_address" validate:"required,hostport"`
}

// TezosConfig contains the configuration related to tezos network
type TezosConfig map[string]*TezosPolicy

// TezosPolicy contains policy definition for a specific address
type TezosPolicy struct {
	AllowedOperations []string `yaml:"allowed_operations" validate:"dive,oneof=generic block endorsement"`
	AllowedKinds      []string `yaml:"allowed_kinds" validate:"dive,oneof=endorsement seed_nonce_revelation activate_account ballot reveal transaction origination delegation"`
	LogPayloads       bool     `yaml:"log_payloads"`
    Backend           []string `yaml:"backend"`
}

// VaultConfig represents single vault instance
type VaultConfig struct {
	Driver string    `yaml:"driver" validate:"required"`
	Config yaml.Node `yaml:"config"`
}

// Config contains all the configuration necessary to run the signatory
type Config struct {
	Vaults map[string]*VaultConfig `yaml:"vaults" validate:"gt=0,dive,required"`
	Tezos  TezosConfig             `yaml:"tezos" validate:"dive,keys,startswith=tz1|startswith=tz2|startswith=tz3,len=36,endkeys,dive"`
	Server ServerConfig            `yaml:"server"`
}

// Read read the config from a file
func (c *Config) Read(file string) error {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return err
	}

	return nil
}

func Validator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("hostport", func(fl validator.FieldLevel) bool {
		_, _, err := net.SplitHostPort(fl.Field().String())
		return err == nil
	})
	return validate
}
