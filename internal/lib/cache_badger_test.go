package lib_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
)

func TestBadger(t *testing.T) {
	// 创建临时目录
	dir := "./badger-test"
	defer os.RemoveAll(dir)

	opts := badger.DefaultOptions(dir).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	key := []byte("hello")
	value := []byte("world")

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := db.View(func(txn *badger.Txn) error {
				item, err := txn.Get(key)
				if err == badger.ErrKeyNotFound {
					t.Log("Key not found.")
					return nil
				}
				if err != nil {
					return err
				}

				val, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				t.Log("Read value:", string(val))
				return nil
			})

			if err != nil {
				log.Println("Read error:", err)
			}
		}
	}()
	t.Log("准备写入")
	if err = db.Update(func(txn *badger.Txn) error {
		txn.Set(key, value)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	t.Log("准备删除")
	time.Sleep(1 * time.Second)
	db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		t.Log("准备删除:", string(val))
		txn.Delete(key)
		time.Sleep(10 * time.Second)
		return nil
	})
}
