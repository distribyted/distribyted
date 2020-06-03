package node

import (
	"context"
	"path/filepath"
	"syscall"

	"github.com/anacrolix/torrent"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var _ fs.NodeOnAdder = &Root{}
var _ fs.NodeGetattrer = &Root{}

type Root struct {
	fs.Inode
	torrents []*torrent.Torrent
}

func NewRoot(torrents []*torrent.Torrent) *Root {
	return &Root{torrents: torrents}
}

func (root *Root) OnAdd(ctx context.Context) {
	for _, torrent := range root.torrents {
		root.AddChild(
			filepath.Clean(torrent.Name()),
			root.NewPersistentInode(ctx, &Torrent{t: torrent}, fs.StableAttr{
				Mode: syscall.S_IFDIR,
			}), true)
	}
}

func (root *Root) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = syscall.S_IFDIR & 07777

	return fs.OK
}
