package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/distribyted/distribyted"
	"gopkg.in/yaml.v3"
)

type EventFunc func(event string)
type ReloadFunc func(*Root, EventFunc) error

type Handler struct {
	p string
}

func NewHandler(path string) *Handler {
	return &Handler{p: path}
}

func (c *Handler) createFromTemplateFile() ([]byte, error) {
	t, err := distribyted.Templates.Open("templates/config_template.yaml")
	if err != nil {
		return nil, err
	}
	defer t.Close()

	tb, err := ioutil.ReadAll(t)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(c.p), 0744); err != nil {
		return nil, fmt.Errorf("error creating path for configuration file: %s, %w", c.p, err)
	}
	return tb, ioutil.WriteFile(c.p, tb, 0644)
}

func (c *Handler) GetRaw() ([]byte, error) {
	f, err := ioutil.ReadFile(c.p)
	if os.IsNotExist(err) {
		fmt.Println("configuration file does not exist, creating from template file:", c.p)
		return c.createFromTemplateFile()
	}

	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %w", err)
	}

	return f, nil
}

func (c *Handler) Get() (*Root, error) {
	b, err := c.GetRaw()
	if err != nil {
		return nil, err
	}

	conf := &Root{}
	if err := yaml.Unmarshal(b, conf); err != nil {
		return nil, fmt.Errorf("error parsing configuration file: %w", err)
	}

	conf = AddDefaults(conf)

	return conf, nil
}
