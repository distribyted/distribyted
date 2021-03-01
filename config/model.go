package config

// Root is the main yaml config object
type Root struct {
	MaxCacheSize   int64  `yaml:"max-cache-size,omitempty"`
	MetadataFolder string `yaml:"metadata-folder-name,omitempty"`
	AllowOther     bool   `yaml:"fuse-allow-other,omitempty"`

	MountPoints []*MountPoint `yaml:"mountPoints"`
	WebDAV      *WebDAV       `yaml:"webDav"`
}

type WebDAV struct {
	Torrents []*Torrent `yaml:"torrents"`
}

type MountPoint struct {
	Path     string     `yaml:"path"`
	Torrents []*Torrent `yaml:"torrents"`
}

type Torrent struct {
	MagnetURI   string `yaml:"magnetUri,omitempty"`
	TorrentPath string `yaml:"torrentPath,omitempty"`
	FolderName  string `yaml:"folderName,omitempty"`
}

func AddDefaults(r *Root) *Root {
	if r.MaxCacheSize == 0 {
		r.MaxCacheSize = 1024 // 1GB
	}
	if r.MetadataFolder == "" {
		r.MetadataFolder = "./distribyted-data/metadata"
	}

	return r
}
