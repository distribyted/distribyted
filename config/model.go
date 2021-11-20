package config

// Root is the main yaml config object
type Root struct {
	HTTPGlobal *HTTPGlobal    `yaml:"http"`
	WebDAV     *WebDAVGlobal  `yaml:"webdav"`
	Torrent    *TorrentGlobal `yaml:"torrent"`
	Fuse       *FuseGlobal    `yaml:"fuse"`
	Log        *Log           `yaml:"log"`

	Routes []*Route `yaml:"routes"`
}

type Log struct {
	Debug      bool   `yaml:"debug"`
	MaxBackups int    `yaml:"max_backups"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	Path       string `yaml:"path"`
}

type TorrentGlobal struct {
	AddTimeout      int    `yaml:"add_timeout,omitempty"`
	GlobalCacheSize int64  `yaml:"global_cache_size,omitempty"`
	MetadataFolder  string `yaml:"metadata_folder,omitempty"`
	DisableIPv6     bool   `yaml:"disable_ipv6,omitempty"`
}

type WebDAVGlobal struct {
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type HTTPGlobal struct {
	Port int `yaml:"port"`
}

type FuseGlobal struct {
	AllowOther bool   `yaml:"allow_other,omitempty"`
	Path       string `yaml:"path"`
}

type Route struct {
	Name     string     `yaml:"name"`
	Torrents []*Torrent `yaml:"torrents"`
}

type Torrent struct {
	MagnetURI   string `yaml:"magnet_uri,omitempty"`
	TorrentPath string `yaml:"torrent_path,omitempty"`
}

func AddDefaults(r *Root) *Root {
	if r.Torrent == nil {
		r.Torrent = &TorrentGlobal{}
	}
	if r.Torrent.GlobalCacheSize == 0 {
		r.Torrent.GlobalCacheSize = 1024 // 1GB
	}

	if r.Torrent.MetadataFolder == "" {
		r.Torrent.MetadataFolder = metadataFolder
	}

	if r.Fuse == nil {
		r.Fuse = &FuseGlobal{}
	}
	if r.Fuse.Path == "" {
		r.Fuse.Path = mountFolder
	}

	return r
}
