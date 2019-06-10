package statstore

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/indrenicloud/tricloud-server/app/config"
	"github.com/indrenicloud/tricloud-server/app/logg"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

var (
	ErrNotExists = errors.New("Agent db doesnot exist")
)

func init() { InitDB(config.GetConfig().StatDBpath) }

func InitDB(path string) {
	dbcon, err := bolt.Open(path, 0666, nil)
	if err != nil {
		logg.Error(err)
	}
	db = dbcon
}

func StoreStat(agentname string, t int64, value []byte) {

	logg.Info(string(value))

	db.Update(func(tx *bolt.Tx) error {

		bkt, err := tx.CreateBucketIfNotExists([]byte(agentname))
		if err != nil {
			logg.Info("systemstatus bkt error")
			return err
		}

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(t))

		err = bkt.Put(b, value)
		if err != nil {
			logg.Info("agent bucket err")
			return err
		}

		return nil
	})
}

func GetStats(agentname string, noOfEntries, offset int64) map[int64]map[string]interface{} {
	m := make(map[int64]map[string]interface{})

	keybytes, valbytes, err := getStats(agentname, noOfEntries, offset)
	if err != nil {
		return m
	}

	for index, key := range keybytes {

		vb := valbytes[index]
		v := make(map[string]interface{})
		err = json.Unmarshal(vb, &v)
		if err != nil {
			logg.Warn("decoding stats error!")
			return m
		}

		m[int64(binary.LittleEndian.Uint64(key))] = v
	}

	return m
}

func getStats(agentname string, noOfEntries, offset int64) ([][]byte, [][]byte, error) {

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
		if offset == 0 {
			k, _ := c.First()
			offbyte = k
		} else {
			binary.LittleEndian.PutUint64(offbyte, uint64(offset))
		}

		if noOfEntries == 0 {
			noOfEntries = 10
		}
		count := int64(0)

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
