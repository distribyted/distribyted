package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var dummyFactories = map[string]FsFactory{
	".test": func(f File) (Filesystem, error) {
		return &DummyFs{}, nil
	},
}

func TestStorage(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	s := newStorage(dummyFactories)

	err := s.Add(&Dummy{}, "/path/to/dummy/file.txt")
	require.NoError(err)

	err = s.Add(&Dummy{}, "/path/to/dummy/file2.txt")
	require.NoError(err)

	contains := s.Has("/path")
	require.True(contains)

	contains = s.Has("/path/to/dummy/")
	require.True(contains)

	file, err := s.Get("/path/to/dummy/file.txt")
	require.NoError(err)
	require.Equal(&Dummy{}, file)

	file, err = s.Get("/path/to/dummy/file3.txt")
	require.Error(err)
	require.Nil(file)

	files, err := s.Children("/path/to/dummy/")
	require.NoError(err)
	require.Len(files, 2)
	require.Contains(files, "file.txt")
	require.Contains(files, "file2.txt")

	err = s.Add(&Dummy{}, "/path/to/dummy/folder/file.txt")
	require.NoError(err)

	files, err = s.Children("/path/to/dummy/")
	require.NoError(err)
	require.Len(files, 3)
	require.Contains(files, "file.txt")
	require.Contains(files, "file2.txt")
	require.Contains(files, "folder")

	err = s.Add(&Dummy{}, "path/file4.txt")
	require.NoError(err)

	require.True(s.Has("/path/file4.txt"))

	files, err = s.Children("/")
	require.NoError(err)
	require.Len(files, 1)

	err = s.Add(&Dummy{}, "/path/special_file.test")
	require.NoError(err)

	file, err = s.Get("/path/special_file.test/dir/here/file1.txt")
	require.NoError(err)
	require.Equal(&Dummy{}, file)

	files, err = s.Children("/path/special_file.test")
	require.NoError(err)
	require.NotNil(files)

	files, err = s.Children("/path/special_file.test/dir/here")
	require.NoError(err)
	require.Len(files, 2)

	err = s.Add(&Dummy{}, "/path/to/__special__path/file3.txt")
	require.NoError(err)

	file, err = s.Get("/path/to/__special__path/file3.txt")
	require.NoError(err)
	require.Equal(&Dummy{}, file)

	s.Clear()
}

func TestStorageWindowsPath(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	s := newStorage(dummyFactories)

	err := s.Add(&Dummy{}, "\\path\\to\\dummy\\file.txt")
	require.NoError(err)

	file, err := s.Get("\\path\\to\\dummy\\file.txt")
	require.NoError(err)
	require.Equal(&Dummy{}, file)

	file, err = s.Get("/path/to/dummy/file.txt")
	require.NoError(err)
	require.Equal(&Dummy{}, file)
}

func TestStorageAddFs(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	s := newStorage(dummyFactories)

	err := s.AddFS(&DummyFs{}, "/test")
	require.NoError(err)

	f, err := s.Get("/test/dir/here/file1.txt")
	require.NoError(err)
	require.NotNil(f)

	err = s.AddFS(&DummyFs{}, "/test")
	require.Error(err)
}

func TestSupportedFactories(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	require.Contains(SupportedFactories, ".zip")
	require.Contains(SupportedFactories, ".rar")
	require.Contains(SupportedFactories, ".7z")

	fs, err := SupportedFactories[".zip"](&Dummy{})
	require.NoError(err)
	require.NotNil(fs)

	fs, err = SupportedFactories[".rar"](&Dummy{})
	require.NoError(err)
	require.NotNil(fs)

	fs, err = SupportedFactories[".7z"](&Dummy{})
	require.NoError(err)
	require.NotNil(fs)
}

var _ Filesystem = &DummyFs{}

type DummyFs struct {
}

func (d *DummyFs) Open(filename string) (File, error) {
	return &Dummy{}, nil
}

func (d *DummyFs) ReadDir(path string) (map[string]File, error) {
	if path == "/dir/here" {
		return map[string]File{
			"file1.txt": &Dummy{},
			"file2.txt": &Dummy{},
		}, nil
	}

	return nil, os.ErrNotExist
}

var _ File = &Dummy{}

type Dummy struct {
}

func (d *Dummy) Size() int64 {
	return 0
}

func (d *Dummy) IsDir() bool {
	return false
}

func (d *Dummy) Close() error {
	return nil
}

func (d *Dummy) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (d *Dummy) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}
