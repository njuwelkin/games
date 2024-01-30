package gobang

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	BucketName = "menual"
)

type GobangDB struct {
	db *bolt.DB
}

func OpenDB(name string) (*GobangDB, error) {
	ret := GobangDB{}
	db, err := bolt.Open(name, 0600, nil)
	if err != nil {
		return nil, err
	}
	ret.db = db
	return &ret, nil
}

func (gbdb *GobangDB) Init() error {
	err := gbdb.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("failed to create bucket %v", err)
		}
		return nil
	})
	return err
}

func (gbdb *GobangDB) Put(key string, value MenualItem) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = gbdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			panic("bucket not exists")
		}
		if err := b.Put([]byte(key), data); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (gbdb *GobangDB) Get(key string) (*MenualItem, error) {
	ret := MenualItem{}
	err := gbdb.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket %s not found!", BucketName)
		}

		val := bucket.Get([]byte(key))
		err := json.Unmarshal(val, &ret)
		if err != nil {
			return fmt.Errorf("cannot unmarshal %s", string(val))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (gbdb *GobangDB) Close() {
	gbdb.db.Close()
}
