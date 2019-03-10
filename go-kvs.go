package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var (
	// kvp.db contains one bucket for now
	bucketName = []byte("kvp")
)

func main() {

	// open database for reading or writing
	db, err := bolt.Open("kvp.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// make sure bucket exists in db
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})

	// process args (arg 0 is always command)
	if len(os.Args) == 2 {
		// arg1 is key, so give value
		key := os.Args[1]
		err := db.View(func(tx *bolt.Tx) error {
			c := tx.Bucket(bucketName).Cursor()
			if kb, vb := c.Seek([]byte(key)); (string(kb) == key) && (kb != nil) && (vb != nil) {
				fmt.Printf("value of %s is %s\n", string(kb), string(vb))
				return nil
			} else {
				return errors.New("error: pair not found")
			}
		})
		if err != nil {
			fmt.Println(err)
		}

		// TODO -- get value from DB
	} else if len(os.Args) == 3 {
		key := os.Args[1]
		value := os.Args[2]
		// arg1 is key, arg2 is value, so insert kvp
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			err := b.Put([]byte(key), []byte(value))
			return err
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("created pair %s -> %s\n", key, value)
	} else {
		fmt.Printf("go-kvs: store and retrieve key value pairs!\n")
		fmt.Printf("usage:\n")
		fmt.Printf("go-kvs <key> <value>\n create a pair in the db.\n")
		fmt.Printf("go-kvs <key>\n request value at given key\n")
	}

}
