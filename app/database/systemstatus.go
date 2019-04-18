package database

import (
	"encoding/binary"
	"errors"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

var SystemStatusBucketName = []byte("systemstatus")

var (
	ErrNoBucket = errors.New("bucket doesnot exist")
)

func UpdateSystemStatus(agentname string, sysinfo map[string]string) {

	_ = DB.conn.Update(func(tx *bolt.Tx) error {

		rootbkt, err := tx.CreateBucketIfNotExists(SystemStatusBucketName)
		if err != nil {
			log.Println("systemstatus bkt error")
			return err
		}
		abkt, err := rootbkt.CreateBucketIfNotExists([]byte(agentname))
		if err != nil {
			log.Println("agent bucket err")
			return err
		}

		t := time.Now()

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(t.UnixNano()))

		valByte, err := Encode(sysinfo)

		err = abkt.Put(b, valByte)
		if err != nil {
			log.Println("agent bucket err")
			return err
		}

		return nil
	})
}

func GetSystemStatus(agentname string, noOfEntries, offset int) {
	// TODO later conever map[string]string to more efficient type

	outkeybytes := make([][]byte, noOfEntries)
	outvalbytes := make([][]byte, noOfEntries)

	err := DB.conn.View(func(tx *bolt.Tx) error {

		rootbkt := tx.Bucket(SystemStatusBucketName)
		if rootbkt == nil {
			return ErrNoBucket
		}

		abkt := rootbkt.Bucket([]byte(agentname))
		if rootbkt == nil {
			return ErrNoBucket
		}
		c := abkt.Cursor()

		offbyte := make([]byte, 8)
		binary.LittleEndian.PutUint64(offbyte, uint64(offset))

		count := 0

		for k, v := c.Seek(offbyte); k != nil || count >= noOfEntries; k, v = c.Next() {
			count++
			outkeybytes = append(outkeybytes, k)
			outvalbytes = append(outvalbytes, v)
		}

		return nil
	})

	if err != nil {
		//err
	}

	// key value pair stored in another map with another timestamp as key
	// like { "time monday 4:45:33": {"cpu":"4%", mem:"33"}}
	systeminfomap := map[int64]map[string]string{}

	for index, val := range outvalbytes {

		i := int64(binary.LittleEndian.Uint64(outkeybytes[index]))

		m := map[string]string{}

		err = Decode(val, &m)
		if err != nil {
			//err
		}

		systeminfomap[i] = m
	}

}
