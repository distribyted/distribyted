package node

import (
	"context"
	"io"
	"syscall"

	"github.com/ajnavarro/distribyted/iio"
	"github.com/anacrolix/torrent"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/panjf2000/ants/v2"
)

var _ fs.NodeOnAdder = &Folder{}
var _ fs.NodeGetattrer = &Folder{}

type Folder struct {
	fs.Inode
	t *torrent.Torrent

	pool *ants.Pool
}

func (folder *Folder) OnAdd(ctx context.Context) {
	for _, file := range folder.t.Files() {
		file := file
		LoadNodeByPath(
			ctx,
			folder.pool,
			file.Path(),
			func() (io.ReaderAt, error) { return iio.NewReadAtWrapper(file.NewReader()), nil },
			&folder.Inode,
			file.Length(),
			int32(file.Torrent().Info().PieceLength),
			int64(file.Torrent().Info().NumPieces()),
		)
	}
}

func (folder *Folder) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = syscall.S_IFDIR & 07777

	return fs.OK
}
