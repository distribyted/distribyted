package fs

type ContainerFs struct {
	s *storage
}

func NewContainerFs(fss map[string]Filesystem) (*ContainerFs, error) {
	s := newStorage(SupportedFactories)
	for p, fs := range fss {
		if err := s.AddFS(fs, p); err != nil {
			return nil, err
		}
	}

	return &ContainerFs{s: s}, nil
}

func (fs *ContainerFs) Open(filename string) (File, error) {
	return fs.s.Get(filename)
}

func (fs *ContainerFs) ReadDir(path string) (map[string]File, error) {
	return fs.s.Children(path)
}
