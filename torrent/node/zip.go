package node

import (
	"archive/zip"
	"context"
	"io"
	"log"
	"syscall"

	"github.com/ajnavarro/distribyted/iio"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/panjf2000/ants/v2"
)

var _ fs.NodeOnAdder = &ZipFolder{}
var _ fs.NodeGetattrer = &ZipFolder{}

type ZipFolder struct {
	fs.Inode
	readerFunc ReaderFunc
	size       int64
	zr         *zip.Reader
	pool       *ants.Pool
	name       string
}

func NewZipFolder(pool *ants.Pool, readerFunc ReaderFunc, size int64, name string) *ZipFolder {
	return &ZipFolder{
		readerFunc: readerFunc,
		size:       size,
		pool:       pool,
		name:       name,
	}
}

func (folder *ZipFolder) OnAdd(ctx context.Context) {
	err := folder.pool.Submit(func() {
		reader, err := folder.readerFunc()
		if err != nil {
			log.Println("error opening reader for zip file", err, "NAME", folder.name)
		}

		zr, err := zip.NewReader(reader, folder.size)
		if err != nil {
			log.Println("error opening zip file:", err, "NAME", folder.name)
			return
		}
		folder.zr = zr

		for _, file := range folder.zr.File {
			file := file
			rf := func() (io.ReaderAt, error) {
				zfr, err := file.Open()
				if err != nil {
					return nil, err
				}

				return iio.NewUnbufferedReaderAt(zfr), nil
			}

			LoadNodeByPath(
				ctx,
				folder.pool,
				file.Name,
				rf,
				&folder.Inode,
				file.FileInfo().Size(),
				0,
				0,
			)
		}
	})
	if err != nil {
		log.Println("error on pool task", err)
	}
}

func (folder *ZipFolder) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = syscall.S_IFDIR & 07777
	return fs.OK
}
