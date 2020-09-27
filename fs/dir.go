package fs

var _ File = &Dir{}

type Dir struct {
}

func (d *Dir) Size() int64 {
	return 0
}

func (d *Dir) IsDir() bool {
	return true
}

func (d *Dir) Close() error {
	return nil
}

func (d *Dir) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (d *Dir) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}
