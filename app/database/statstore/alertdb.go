package statstore

import (
	"encoding/binary"
	"github.com/indrenicloud/tricloud-server/app/logg"
		bolt "go.etcd.io/bbolt"
)

var alertsBucket = []byte("alerts")

func StoreAlert(alertbyte []byte, agentname []byte, timestamp int64)error{
	logg.Debug("storealert")

	db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists(alertsBucket)
		if err != nil {
			logg.Info("systemstatus bkt error")
			return err
		}
		b2, err := bkt.CreateBucketIfNotExists(agentname)
		if err != nil {
			logg.Info("systemstatus bkt error")
			return err
		}
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(timestamp))
		
		err = b2.Put(b, alertbyte)
		if err != nil {
			logg.Info("agent bucket err")
			return err
		}

		return nil
	})

	return nil
}

func GetAlert(agentname []byte)([][]byte, error){
	logg.Debug("getalert")
	var outbytes [][]byte
	err := db.View(func(tx *bolt.Tx) error {

		bkt := tx.Bucket([]byte(alertsBucket))

		if bkt == nil {
			logg.Debug("bucket not created")
			return ErrNotExists
		}
		
		b2 := bkt.Bucket(agentname)

		c := b2.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outbytes = append(outbytes, v)
		}
		return nil
	})
	
	
	return outbytes, err
}
