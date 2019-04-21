package statstore

import (
	"encoding/binary"
	"errors"
	"log"

	"github.com/indrenicloud/tricloud-server/app/logg"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

var (
	ErrNotExists = errors.New("Aagent db doesnot exist")
)

func init() { InitDB("sysstat.db") }

func InitDB(path string) {
	dbcon, err := bolt.Open(path, 0666, nil)
	if err != nil {
		logg.Error(err)
	}
	db = dbcon
}

func StoreStat(agentname string, t int64, value []byte) {
	db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(agentname))
		if err != nil {
			log.Println("systemstatus bkt error")
			return err
		}

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(t))

		err = bkt.Put(b, value)
		if err != nil {
			log.Println("agent bucket err")
			return err
		}

		return nil
	})
}

func GetStats(agentname string, noOfEntries, offset int64) ([][]byte, [][]byte, error) {

	var outkeybytes [][]byte
	var outvalbytes [][]byte

	err := db.View(func(tx *bolt.Tx) error {

		bkt := tx.Bucket([]byte(agentname))

		if bkt == nil {
			logg.Warn("bucket not created")
			return ErrNotExists
		}
		c := bkt.Cursor()

		offbyte := make([]byte, 8)
		binary.LittleEndian.PutUint64(offbyte, uint64(offset))

		var count int64 = 0

		for k, v := c.Seek(offbyte); k != nil; k, v = c.Next() {
			outkeybytes = append(outkeybytes, k)
			outvalbytes = append(outvalbytes, v)

			count++
			if count >= noOfEntries {
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return outkeybytes, outvalbytes, nil
}
