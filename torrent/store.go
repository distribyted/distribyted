package torrent

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/anacrolix/dht/v2/bep44"
	"github.com/dgraph-io/badger/v3"
	dlog "github.com/distribyted/distribyted/log"
	"github.com/rs/zerolog/log"
)

var _ bep44.Store = &FileItemStore{}

type FileItemStore struct {
	ttl time.Duration
	db  *badger.DB
}

func NewFileItemStore(path string, itemsTTL time.Duration) (*FileItemStore, error) {
	l := log.Logger.With().Str("component", "item-store").Logger()

	opts := badger.DefaultOptions(path).
		WithLogger(&dlog.Badger{L: l}).
		WithValueLogFileSize(1<<26 - 1)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	err = db.RunValueLogGC(0.5)
	if err != nil && err != badger.ErrNoRewrite {
		return nil, err
	}

	return &FileItemStore{
		db:  db,
		ttl: itemsTTL,
	}, nil
}

func (fis *FileItemStore) Put(i *bep44.Item) error {
	tx := fis.db.NewTransaction(true)
	defer tx.Discard()

	key := i.Target()
	var value bytes.Buffer

	enc := gob.NewEncoder(&value)
	if err := enc.Encode(i); err != nil {
		return err
	}

	e := badger.NewEntry(key[:], value.Bytes()).WithTTL(fis.ttl)
	if err := tx.SetEntry(e); err != nil {
		return err
	}

	return tx.Commit()
}

func (fis *FileItemStore) Get(t bep44.Target) (*bep44.Item, error) {
	tx := fis.db.NewTransaction(false)
	defer tx.Discard()

	dbi, err := tx.Get(t[:])
	if err == badger.ErrKeyNotFound {
		return nil, bep44.ErrItemNotFound
	}
	if err != nil {
		return nil, err
	}
	valb, err := dbi.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(valb)
	dec := gob.NewDecoder(buf)
	var i *bep44.Item
	if err := dec.Decode(&i); err != nil {
		return nil, err
	}

	return i, nil
}

func (fis *FileItemStore) Del(t bep44.Target) error {
	// ignore this
	return nil
}

func (fis *FileItemStore) Close() error {
	return fis.db.Close()
}
