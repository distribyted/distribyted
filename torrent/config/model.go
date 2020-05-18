package config

// Root is the main yaml config object
type Root struct {
	MaxCacheSize   int64  `yaml:"max-cache-size,omitempty"`
	MetadataFolder string `yaml:"metadata-folder-name,omitempty"`

	MountPoints []*MountPoint `yaml:"mountPoints"`
}

type MountPoint struct {
	Path     string `yaml:"path"`
	Torrents []struct {
		MagnetURI   string `yaml:"magnetUri"`
		TorrentPath string `yaml:"torrentPath"`
		FolderName  string `yaml:"folderName,omitempty"`
	} `yaml:"torrents"`
}

func AddDefaults(r *Root) *Root {
	if r.MaxCacheSize == 0 {
		r.MaxCacheSize = 1024 // 1GB
	}
	if r.MetadataFolder == "" {
		r.MetadataFolder = "./metadata"
	}

	return r
}
