package lib

import (
	"github.com/dgraph-io/badger/v4"
)

var Badger *badger.DB

func InitBadger(path string) error {
	var err error
	Badger, err = badger.Open(badger.DefaultOptions(path))
	return err
}
