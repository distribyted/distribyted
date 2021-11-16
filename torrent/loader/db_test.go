package loader

import (
	"os"
	"testing"

	"github.com/anacrolix/torrent/storage"
	"github.com/stretchr/testify/require"
)

const m1 = "magnet:?xt=urn:btih:c9e15763f722f23e98a29decdfae341b98d53056"

func TestDB(t *testing.T) {
	require := require.New(t)

	tmpService, err := os.MkdirTemp("", "service")
	require.NoError(err)
	tmpStorage, err := os.MkdirTemp("", "storage")
	require.NoError(err)

	cs := storage.NewFile(tmpStorage)
	defer cs.Close()

	s, err := NewDB(tmpService)
	require.NoError(err)
	defer s.Close()

	err = s.AddMagnet("route1", "WRONG MAGNET")
	require.Error(err)

	err = s.AddMagnet("route1", m1)
	require.NoError(err)

	err = s.AddMagnet("route2", m1)
	require.NoError(err)

	l, err := s.ListMagnets()
	require.NoError(err)
	require.Len(l, 2)
	require.Len(l["route1"], 1)
	require.Equal(l["route1"][0], m1)
	require.Len(l["route2"], 1)
	require.Equal(l["route2"][0], m1)

	removed, err := s.RemoveFromHash("other", "c9e15763f722f23e98a29decdfae341b98d53056")
	require.NoError(err)
	require.False(removed)

	removed, err = s.RemoveFromHash("route1", "c9e15763f722f23e98a29decdfae341b98d53056")
	require.NoError(err)
	require.True(removed)

	l, err = s.ListMagnets()
	require.NoError(err)
	require.Len(l, 1)
	require.Len(l["route2"], 1)
	require.Equal(l["route2"][0], m1)

	require.NoError(s.Close())
	require.NoError(cs.Close())

}
