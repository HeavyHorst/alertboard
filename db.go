package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type store interface {
	putAlert(a alertData) error
	getAlert(id string) []byte
	deleteAlert(id string) error
	getAlertsByPrefix(prefix string) ([]byte, error)
}

type boltStore struct {
	db *bolt.DB
}

func newBoltStore() (*boltStore, error) {
	db, err := bolt.Open("alertboard.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	b := &boltStore{db}

	//create the bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("alerts"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return b, err
}

func (b *boltStore) putAlert(a alertData) error {
	if a.Time.IsZero() {
		a.Time = time.Now()
	}

	if a.Status == "" {
		a.Status = "Open"
	}

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("alerts"))
		err := b.Put([]byte(a.ID), data)
		return err
	})
	return err
}

func (b *boltStore) getAlert(id string) []byte {
	var alert []byte
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("alerts"))
		alert = b.Get([]byte(id))
		return nil
	})
	return alert
}

func (b *boltStore) deleteAlert(id string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("alerts"))
		err := b.Delete([]byte(id))
		return err
	})
	return err
}

func (b *boltStore) getAlertsByPrefix(prefix string) ([]byte, error) {
	alerts := make([]alertData, 0, 0)
	var alert alertData

	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("alerts"))
		c := b.Cursor()
		p := []byte(prefix)
		for k, v := c.Seek(p); bytes.HasPrefix(k, p) && k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &alert)
			if err != nil {
				return err
			}
			alerts = append(alerts, alert)
		}
		return nil
	})

	data, _ := json.Marshal(alerts)
	return data, err
}

func (b *boltStore) Close() {
	b.db.Close()
}
