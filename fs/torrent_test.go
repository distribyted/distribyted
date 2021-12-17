package fs

import (
	"os"
	"testing"

	"github.com/anacrolix/torrent"

	"github.com/stretchr/testify/require"
)

const testMagnet = "magnet:?xt=urn:btih:a88fda5954e89178c372716a6a78b8180ed4dad3&dn=The+WIRED+CD+-+Rip.+Sample.+Mash.+Share&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F&xs=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fwired-cd.torrent"

var Cli *torrent.Client

func TestMain(m *testing.M) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = os.TempDir()

	// disable webseeds to avoid a panic when closing client on tests
	cfg.DisableWebseeds = true

	client, err := torrent.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	Cli = client

	exitVal := m.Run()

	client.Close()

	os.Exit(exitVal)
}

func TestTorrentFilesystem(t *testing.T) {
	require := require.New(t)

	to, err := Cli.AddMagnet(testMagnet)
	require.NoError(err)

	tfs := NewTorrent(600)
	tfs.AddTorrent(to)

	files, err := tfs.ReadDir("/")
	require.NoError(err)
	require.Len(files, 1)
	require.Contains(files, "The WIRED CD - Rip. Sample. Mash. Share")

	files, err = tfs.ReadDir("/The WIRED CD - Rip. Sample. Mash. Share")
	require.NoError(err)
	require.Len(files, 18)

	f, err := tfs.Open("/The WIRED CD - Rip. Sample. Mash. Share/not_existing_file.txt")
	require.Equal(os.ErrNotExist, err)
	require.Nil(f)

	f, err = tfs.Open("/The WIRED CD - Rip. Sample. Mash. Share/01 - Beastie Boys - Now Get Busy.mp3")
	require.NoError(err)
	require.NotNil(f)
	require.Equal(f.Size(), int64(1964275))

	b := make([]byte, 10)

	n, err := f.Read(b)
	require.NoError(err)
	require.Equal(10, n)
	require.Equal([]byte{0x49, 0x44, 0x33, 0x3, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x76}, b)

	n, err = f.ReadAt(b, 10)
	require.NoError(err)
	require.Equal(10, n)

	n, err = f.ReadAt(b, 10000)
	require.NoError(err)
	require.Equal(10, n)

	tfs.RemoveTorrent(to.InfoHash().String())
	files, err = tfs.ReadDir("/")
	require.NoError(err)
	require.Len(files, 0)

	require.NoError(f.Close())
}

func TestReadAtTorrent(t *testing.T) {
	require := require.New(t)

	to, err := Cli.AddMagnet(testMagnet)
	require.NoError(err)

	<-to.GotInfo()
	torrFile := to.Files()[0]

	tf := torrentFile{
		readerFunc: torrFile.NewReader,
		len:        torrFile.Length(),
		timeout:    500,
	}

	defer tf.Close()

	toRead := make([]byte, 5)
	n, err := tf.ReadAt(toRead, 6)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal([]byte{0x0, 0x0, 0x1f, 0x76, 0x54}, toRead)

	n, err = tf.ReadAt(toRead, 0)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal([]byte{0x49, 0x44, 0x33, 0x3, 0x0}, toRead)
}

func TestReadAtWrapper(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	to, err := Cli.AddMagnet(testMagnet)
	require.NoError(err)

	<-to.GotInfo()
	torrFile := to.Files()[0]

	r := newReadAtWrapper(torrFile.NewReader(), 10)
	defer r.Close()

	toRead := make([]byte, 5)
	n, err := r.ReadAt(toRead, 6)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal([]byte{0x0, 0x0, 0x1f, 0x76, 0x54}, toRead)

	n, err = r.ReadAt(toRead, 0)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal([]byte{0x49, 0x44, 0x33, 0x3, 0x0}, toRead)
}
