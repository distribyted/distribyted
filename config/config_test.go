package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestTemplateConfig(t *testing.T) {
	require := require.New(t)
	f, err := os.Open("../templates/config_template.yaml")
	require.NoError(err)

	b, err := ioutil.ReadAll(f)
	require.NoError(err)

	conf := &Root{}
	err = yaml.Unmarshal(b, conf)
	require.NoError(err)

	require.Equal(DefaultConfig(), conf)
}
