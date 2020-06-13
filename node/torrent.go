package node

import (
	"context"
	"io"
	"syscall"

	"github.com/ajnavarro/distribyted/iio"
	"github.com/anacrolix/torrent"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var _ fs.NodeGetattrer = &Torrent{}
var _ fs.NodeOpendirer = &Torrent{}

type Torrent struct {
	fs.Inode
	t *torrent.Torrent
}

func (folder *Torrent) Opendir(ctx context.Context) syscall.Errno {
	<-folder.t.GotInfo()

	for _, file := range folder.t.Files() {
		file := file
		LoadNodeByPath(
			ctx,
			file.Path(),
			func() (io.ReaderAt, error) { return iio.NewReadAtWrapper(file.NewReader()), nil },
			&folder.Inode,
			file.Length(),
			int32(file.Torrent().Info().PieceLength),
			int64(file.Torrent().Info().NumPieces()),
		)
	}
	return fs.OK
}

func (folder *Torrent) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = syscall.S_IFDIR & 0555

	return fs.OK
}
