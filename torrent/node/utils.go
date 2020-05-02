package node

import (
	"context"
	"io"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/panjf2000/ants/v2"
)

type ReaderFunc func() (io.ReaderAt, error)

func LoadNodeByPath(ctx context.Context, pool *ants.Pool, fp string, reader ReaderFunc, parent *fs.Inode, fileLength int64, pieceLen int32, numPieces int64) {
	p := parent
	dir, base := filepath.Split(filepath.Clean(fp))
	for i, component := range strings.Split(dir, "/") {
		if i == 0 {
			continue
		}

		if len(component) == 0 {
			continue
		}

		ch := p.GetChild(component)
		if ch == nil {
			ch = p.NewPersistentInode(ctx, &fs.Inode{},
				fs.StableAttr{Mode: fuse.S_IFDIR})
			p.AddChild(component, ch, true)
		}

		p = ch
	}

	ext := path.Ext(base)
	switch ext {
	default:
		n := NewFileWithBlocks(
			reader,
			fileLength,
			pieceLen,
			numPieces,
		)
		p.AddChild(
			base,
			p.NewPersistentInode(ctx, n, fs.StableAttr{
				Mode: syscall.S_IFREG,
			}), true)
	}
}
