package config

// Root is the main yaml config object
type Root struct {
	MaxCacheSize   int64  `yaml:"max-cache-size,omitempty"`
	MetadataFolder string `yaml:"metadata-folder-name,omitempty"`

	MountPoints []*MountPoint `yaml:"mountPoints"`
}

type MountPoint struct {
	AllowOther bool   `yaml:"fuse-allow-other,omitempty"`
	Path       string `yaml:"path"`
	Torrents   []struct {
		MagnetURI   string `yaml:"magnetUri,omitempty"`
		TorrentPath string `yaml:"torrentPath,omitempty"`
		FolderName  string `yaml:"folderName,omitempty"`
	} `yaml:"torrents"`
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
