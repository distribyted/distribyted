package loader

import (
	"path"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/dgraph-io/badger/v3"
	dlog "github.com/distribyted/distribyted/log"
	"github.com/rs/zerolog/log"
)

var _ LoaderAdder = &DB{}

const routeRootKey = "/route/"

type DB struct {
	db *badger.DB
}

func NewDB(path string) (*DB, error) {
	l := log.Logger.With().Str("component", "torrent-store").Logger()
	db, err := badger.Open(badger.DefaultOptions(path).WithLogger(&dlog.Badger{L: l}))
	if err != nil {
		return nil, err
	}

	err = db.RunValueLogGC(0.5)
	if err != nil && err != badger.ErrNoRewrite {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func (l *DB) AddMagnet(r, m string) error {
	err := l.db.Update(func(txn *badger.Txn) error {
		spec, err := metainfo.ParseMagnetUri(m)
		if err != nil {
			return err
		}

		ih := spec.InfoHash.HexString()

		rp := path.Join(routeRootKey, ih, r)
		return txn.Set([]byte(rp), []byte(m))
	})

	if err != nil {
		return err
	}

	return l.db.Sync()
}

func (l *DB) RemoveFromHash(r, h string) (bool, error) {
	tx := l.db.NewTransaction(true)
	defer tx.Discard()

	var mh metainfo.Hash
	if err := mh.FromHexString(h); err != nil {
		return false, err
	}

	rp := path.Join(routeRootKey, h, r)
	if _, err := tx.Get([]byte(rp)); err != nil {
		return false, nil
	}

	if err := tx.Delete([]byte(rp)); err != nil {
		return false, err
	}

	return true, tx.Commit()
}

func (l *DB) ListMagnets() (map[string][]string, error) {
	tx := l.db.NewTransaction(false)
	defer tx.Discard()

	it := tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(routeRootKey)
	out := make(map[string][]string)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		_, r := path.Split(string(it.Item().Key()))
		i := it.Item()
		if err := i.Value(func(v []byte) error {
			out[r] = append(out[r], string(v))
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (l *DB) ListTorrentPaths() (map[string][]string, error) {
	return nil, nil
}

func (l *DB) Close() error {
	return l.db.Close()
}
