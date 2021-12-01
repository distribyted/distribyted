package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func TestTemplateConfig(t *testing.T) {
	t.Parallel()

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

func TestDefaults(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	r := &Root{}
	dr := AddDefaults(r)
	require.NotNil(dr)
	require.NotNil(dr.Fuse)
	require.NotNil(dr.HTTPGlobal)
	require.NotNil(dr.Log)
	require.NotNil(dr.Torrent)
}
