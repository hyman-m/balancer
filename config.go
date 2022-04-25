package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var ascii = `
___ _ _  _ _   _ ___  ____ _    ____ _  _ ____ ____ ____ 
 |  | |\ |  \_/  |__] |__| |    |__| |\ | |    |___ |__/ 
 |  | | \|   |   |__] |  | |___ |  | | \| |___ |___ |  \                                        
`

// Config .
type Config struct {
	SSLCertificateKey string      `yaml:"ssl_certificate_key"`
	Location          []*Location `yaml:"location"`
	Schema            string      `yaml:"schema"`
	Port              int         `yaml:"port"`
	SSLCertificate    string      `yaml:"ssl_certificate"`
}

// Location .
type Location struct {
	Pattern     string   `yaml:"pattern"`
	ProxyPass   []string `yaml:"proxy_pass"`
	BalanceMode string   `yaml:"balance_mode"`
}

// ReadConfig .
func ReadConfig(fileName string) (*Config, error) {
	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Print .
func (c *Config) Print() {
	fmt.Printf("%s\n", ascii)
	fmt.Printf("Schema: %s\nPort: %d\nLocation:\n", c.Schema, c.Port)
	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\tProxyPass: %s\n\tMode: %s\n",
			l.Pattern, l.ProxyPass, l.BalanceMode)
	}
}
