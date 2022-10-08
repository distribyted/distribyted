package loader

import (
	"io/fs"
	"path"
	"path/filepath"

	"github.com/distribyted/distribyted/config"
)

var _ Loader = &Folder{}

type Folder struct {
	c []*config.Route
}

func NewFolder(r []*config.Route) *Folder {
	return &Folder{
		c: r,
	}
}

func (f *Folder) ListMagnets() (map[string][]string, error) {
	return nil, nil
}

func (f *Folder) ListTorrentPaths() (map[string][]string, error) {
	out := make(map[string][]string)
	for _, r := range f.c {
		if r.TorrentFolder == "" {
			continue
		}

		err := filepath.WalkDir(r.TorrentFolder, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}
			if path.Ext(p) == ".torrent" {
				out[r.Name] = append(out[r.Name], p)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
