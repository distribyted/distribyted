package node

import (
	"archive/zip"
	"context"
	"io"
	"syscall"

	"github.com/ajnavarro/distribyted/iio"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
)

var _ fs.NodeGetattrer = &Zip{}
var _ fs.NodeOpendirer = &Zip{}

type Zip struct {
	fs.Inode

	reader ReaderFunc
	size   int64
	files  []*zip.File
}

func NewZip(reader ReaderFunc, size int64) *Zip {
	return &Zip{
		reader: reader,
		size:   size,
	}
}

func (z *Zip) Opendir(ctx context.Context) syscall.Errno {
	if z.files == nil {
		r, err := z.reader()
		if err != nil {
			log.WithError(err).Error("error opening reader for zip")
			return syscall.EIO
		}
		zr, err := zip.NewReader(r, z.size)
		if err != nil {
			log.WithError(err).Error("error getting zip reader from reader")
			return syscall.EIO
		}

		for _, f := range zr.File {
			f := f
			if f.FileInfo().IsDir() {
				continue
			}
			LoadNodeByPath(
				ctx,
				f.Name,
				func() (io.ReaderAt, error) {
					zfr, err := f.Open()
					if err != nil {
						log.WithError(err).Error("error opening zip file")
						return nil, err
					}

					return iio.NewDiskTeeReader(zfr)
				},
				&z.Inode,
				int64(f.UncompressedSize64),
				0,
				0,
			)
		}

		z.files = zr.File
	}

	return fs.OK
}

func (z *Zip) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = syscall.S_IFDIR & 0555
	out.Size = uint64(z.size)

	return fs.OK
}
