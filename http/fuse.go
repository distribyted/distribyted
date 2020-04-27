package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var _ fs.NodeOnAdder = &HttpRoot{}
var _ fs.NodeGetattrer = &HttpRoot{}

type HttpRoot struct {
	fs.Inode

	URLs []string

	m      sync.Mutex
	loaded bool
}

func (r *HttpRoot) OnAdd(ctx context.Context) {
	r.m.Lock()
	defer r.m.Unlock()
	if !r.loaded {
		for _, u := range r.URLs {
			fu, err := url.Parse(u)
			if err != nil {
				log.Println("ERROR FORMATTING URL", u)
				panic("BUH")
			}

			ok := r.AddChild(path.Base(fu.Path), r.NewPersistentInode(ctx, NewHttpFile(u), fs.StableAttr{
				Mode: syscall.S_IFREG & 07777,
			}), true)
			if !ok {
				log.Println("Problem adding node child with name", u)
			}

		}
		log.Println("ALL LOADED")
		r.loaded = true
	}
}

func (r *HttpRoot) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	log.Println("GET ATTR FOLDER")
	out.Mode = syscall.S_IFDIR & 07777

	return fs.OK
}

var _ fs.NodeGetattrer = &HttpFile{}
var _ fs.NodeOpener = &HttpFile{}
var _ fs.NodeReader = &HttpFile{}

type HttpFile struct {
	fs.Inode

	len uint64
	url string

	c *http.Client
}

func NewHttpFile(url string) *HttpFile {
	roundTripper := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		MaxConnsPerHost: 1,
	}
	httpClient := &http.Client{
		Transport: roundTripper,
		Timeout:   60 * time.Second,
	}

	return &HttpFile{c: httpClient, url: url}
}

func (f *HttpFile) getLen() (uint64, error) {
	if f.len != 0 {
		return f.len, nil
	}

	res, err := f.c.Head(f.url)
	if err != nil {
		return 0, err
	}

	lStr := res.Header.Get("content-length")
	len, err := strconv.Atoi(lStr)
	if err != nil {
		return 0, err
	}

	f.len = uint64(len)

	return f.len, nil
}

func (f *HttpFile) Getattr(ctx context.Context, fi fs.FileHandle, out *fuse.AttrOut) syscall.Errno {

	len, err := f.getLen()
	if err != nil {
		log.Println("error getting len", err)
		return syscall.EIO
	}

	out.Mode = syscall.S_IFREG & 07777
	out.Nlink = 1
	out.Size = len
	// out.Blksize
	// out.Blocks

	return fs.OK
}

func (f *HttpFile) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	log.Println("OPEN FILE", f.url)

	return nil, fuse.FOPEN_KEEP_CACHE, fs.OK
	//return nil, fuse.FOPEN_DIRECT_IO, fs.OK
}

func (f *HttpFile) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	log.Println("READDDD FROM", off, "TO", int64(len(dest))+off, "TOTAL", len(dest))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil)
	if err != nil {
		log.Println("error generating request from url", err, f.url)
		return nil, syscall.EIO
	}

	l, err := f.getLen()
	if err != nil {
		log.Println("error getting length", err, f.url)
		return nil, syscall.EIO
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", off, int64(len(dest))+off))

	res, err := f.c.Do(req)
	if err != nil {
		log.Println("error sending request", err, f.url)
		return nil, syscall.EIO
	}
	if res.StatusCode != 200 && res.StatusCode != 206 {
		log.Println("ERROR GETTING RESPONSE FROM SERVER", res.StatusCode)
		return nil, syscall.EIO
	}

	defer res.Body.Close()

	buf := dest[:int(math.Min(float64(len(dest)), float64(int64(l)-off)))]
	n, err := io.ReadFull(res.Body, buf)
	if err != nil && err != io.EOF {
		log.Println("error readd fully data", err)

		return nil, syscall.EIO
	}
	buf = buf[:n]

	return fuse.ReadResultData(buf), fs.OK
}
