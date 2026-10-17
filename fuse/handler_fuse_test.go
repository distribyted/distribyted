//go:build !nofuse && cgo

package fuse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuseSupported(t *testing.T) {
	require.True(t, Supported)
}
