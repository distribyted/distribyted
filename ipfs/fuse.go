package distribyted

import (
	"context"
	"io"
	"log"
	"math"
	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	format "github.com/ipfs/go-ipld-format"
	uio "github.com/ipfs/go-unixfs/io"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

var _ fs.NodeGetattrer = &IPFSRoot{}
var _ fs.NodeOpener = &IPFSRoot{}
var _ fs.NodeReader = &IPFSRoot{}

// var _ fs.NodeFlusher = &IPFSRoot{}
var _ fs.NodeOnAdder = &IPFSRoot{}

type IPFSRoot struct {
	api  coreiface.CoreAPI
	Node format.Node

	addChildrens sync.Once
	ds           fs.DirStream

	fs.Inode
}

func NewIPFSRoot(api coreiface.CoreAPI, node format.Node) *IPFSRoot {
	return &IPFSRoot{
		api:  api,
		Node: node,
	}
}

func (ir *IPFSRoot) OnAdd(ctx context.Context) {
	if fuseNodeType(ir.Node) != syscall.S_IFDIR {
		return
	}
	ir.addChildrens.Do(func() {
		log.Println("ADDING LINKS TO NODE", ir.Node.Cid())
		for _, ll := range ir.Node.Links() {
			l := ll
			go func() {
				node, err := ir.api.Dag().Get(ctx, l.Cid)
				if err != nil {
					log.Println("ERROR GETTING CHILD NODE", err)
				}
				log.Println("ADDING NEW NODE", node.Cid())
				ok := ir.AddChild(l.Name, ir.NewPersistentInode(ctx, NewIPFSRoot(ir.api, node), fs.StableAttr{
					Mode: fuseNodeType(node),
				}), true)
				if !ok {
					log.Println("Problem adding node child with name", l.Name)
				}
			}()
		}
	})
}

func fuseNodeType(n format.Node) uint32 {
	for _, l := range n.Links() {
		if l.Name != "" {
			// it is a folder, links with names
			return syscall.S_IFDIR
		}
	}

	return syscall.S_IFREG
}

// func (ir *IPFSRoot) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
// 	// log.Println("lookup", name)
// 	node := ir.getChildByName(name)

// 	if node == nil {
// 		return nil, syscall.ENOENT
// 	}
// 	log.Println("LOOKUP NODE ID", node.Cid())
// 	n := NewIPFSRoot(ir.api, node)

// 	s, err := node.Stat()
// 	if err != nil {
// 		log.Println("err getting stats", err)
// 		return nil, syscall.ENOENT
// 	}

// 	out.Mode = uint32(fuseNodeType(node)) & 07777
// 	//out.Nlink = 1
// 	//out.Mtime = uint64(zf.file.ModTime().Unix())
// 	//out.Atime = out.Mtime
// 	//out.Ctime = out.Mtime
// 	out.Size = uint64(s.DataSize)
// 	out.Blksize = uint32(s.BlockSize)
// 	out.Blocks = uint64((s.DataSize + s.BlockSize - 1) / s.BlockSize)

// 	return ir.NewPersistentInode(ctx, n, fs.StableAttr{Mode: uint32(fuseNodeType(node))}), fs.OK
// }

func (ir *IPFSRoot) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	s, err := ir.Node.Stat()
	if err != nil {
		log.Println("err getting stats", err)
		return syscall.ENODATA
	}

	out.Mode = uint32(fuseNodeType(ir.Node)) & 07777
	out.Nlink = 1
	//out.Mtime = uint64(zf.file.ModTime().Unix())
	//out.Atime = out.Mtime
	//out.Ctime = out.Mtime
	out.Size = uint64(s.CumulativeSize)
	out.Blksize = uint32(s.BlockSize)
	out.Blocks = uint64((s.CumulativeSize + s.BlockSize - 1) / s.BlockSize)

	return fs.OK
}

func (ir *IPFSRoot) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Println("OPEN")

	if flags&(syscall.O_RDWR) != 0 || flags&syscall.O_WRONLY != 0 {
		return nil, 0, syscall.EPERM
	}

	return nil, fuse.FOPEN_KEEP_CACHE, fs.OK
}

// func (f *IPFSRoot) Flush(ctx context.Context, fh fs.FileHandle) syscall.Errno {
// 	log.Println("FLUSH")

// 	return 0
// }

func (ir *IPFSRoot) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	log.Println("READDDD", off)

	dr, err := uio.NewDagReader(ctx, ir.Node, ir.api.Dag())
	if err != nil {
		log.Println("error reading data", err)
		return nil, syscall.EIO
	}

	_, err = dr.Seek(off, io.SeekStart)
	if err != nil {
		log.Println("error seeking data", err)
		return nil, syscall.EIO
	}

	buf := dest[:int(math.Min(float64(len(dest)), float64(int64(dr.Size())-off)))]
	n, err := io.ReadFull(dr, buf)
	if err != nil && err != io.EOF {
		log.Println("error readd fully data", err)

		return nil, syscall.EIO
	}
	buf = buf[:n]

	return fuse.ReadResultData(buf), fs.OK
}
