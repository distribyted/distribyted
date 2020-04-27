package node

import (
	"context"
	"io"
	"log"
	"math"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var _ fs.NodeGetattrer = &File{}
var _ fs.NodeOpener = &File{}
var _ fs.NodeReader = &File{}

// File is a fuse node for files inside a torrent
type File struct {
	fs.Inode

	f         ReaderFunc
	r         io.ReaderAt
	len       int64
	pieceLen  int32
	numPieces int64
}

func NewFile(readerFunc ReaderFunc, len int64) *File {
	return &File{
		f:   readerFunc,
		len: len,
	}
}

func NewFileWithBlocks(readerFunc ReaderFunc, len int64, pieceLen int32, numPieces int64) *File {
	return &File{
		f:         readerFunc,
		len:       len,
		pieceLen:  pieceLen,
		numPieces: numPieces,
	}
}

func (tr *File) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = syscall.S_IFREG & 07777
	out.Nlink = 1
	out.Size = uint64(tr.len)
	if tr.pieceLen != 0 {
		out.Blksize = uint32(tr.pieceLen)
		out.Blocks = uint64(tr.numPieces)
	}

	return fs.OK
}

func (tr *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	// TODO use filehandle
	if tr.r == nil {
		r, err := tr.f()
		if err != nil {
			log.Println("error opening reader for file", err)
			return nil, fuse.FOPEN_KEEP_CACHE, syscall.EIO
		}
		log.Println("OPEN", r)
		tr.r = r
	}

	return nil, fuse.FOPEN_KEEP_CACHE, fs.OK
}

func (tr *File) Read(ctx context.Context, f fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	log.Println("READ", len(dest), off)

	end := int(math.Min(float64(len(dest)), float64(int64(tr.len)-off)))

	buf := dest[:end]

	n, err := tr.r.ReadAt(buf, off)

	if err != nil && err != io.EOF {
		log.Println("error read data", err)
		return nil, syscall.EIO
	}
	log.Println("READ AT ERR", err, n)
	buf = buf[:n]
	return fuse.ReadResultData(buf), fs.OK
}
