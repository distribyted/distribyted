//go:build nofuse || !cgo

package fuse

import (
	"testing"

	"github.com/distribyted/distribyted/fs"
	"github.com/stretchr/testify/require"
)

func TestNoFuseSupported(t *testing.T) {
	require.False(t, Supported)
}

func TestNoFuseHandler(t *testing.T) {
	h := NewHandler(false, "/tmp/test")
	require.NotNil(t, h)
	require.False(t, h.fuseAllowOther)
	require.Equal(t, "/tmp/test", h.path)

	err := h.Mount(map[string]fs.Filesystem{})
	require.ErrorIs(t, err, ErrFuseNotSupported)

	mem := fs.NewMemory()
	err = h.Mount(map[string]fs.Filesystem{"/mem": mem})
	require.ErrorIs(t, err, ErrFuseNotSupported)

	require.NotPanics(t, func() {
		h.Unmount()
	})
}

func TestNoFuseHandlerAllowOther(t *testing.T) {
	h := NewHandler(true, "/mnt/other")
	require.NotNil(t, h)
	require.True(t, h.fuseAllowOther)
	require.Equal(t, "/mnt/other", h.path)

	err := h.Mount(nil)
	require.ErrorIs(t, err, ErrFuseNotSupported)

	require.NotPanics(t, func() {
		h.Unmount()
	})
}
