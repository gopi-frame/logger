package file

import (
	"encoding/json"
	"os"

	"github.com/mitchellh/mapstructure"
)

type Config struct {
	File string      `json:"file" yaml:"file" toml:"file"`
	Perm os.FileMode `json:"perm" yaml:"perm" toml:"perm"`
}

func (c *Config) UnmarshalYAML(unmarshal func(any) error) error {
	var o struct {
		File string `json:"file" yaml:"file" toml:"file"`
	}
	if err := unmarshal(&o); err == nil {
		c.File = o.File
		return nil
	}
	var file string
	if err := unmarshal(&file); err != nil {
		return err
	}
	c.File = file
	return nil
}

func (c *Config) UnmarshalJSON(data []byte) error {
	return c.UnmarshalYAML(func(a any) error {
		return json.Unmarshal(data, a)
	})
}

func UnmarshalOptions(options map[string]any) (*Config, error) {
	cfg := new(Config)
	err := mapstructure.WeakDecode(options, cfg)
	return cfg, err
}
