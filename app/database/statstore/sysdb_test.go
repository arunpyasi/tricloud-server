package statstore

import (
	"testing"
	"time"
)

func Test_dbstore(t *testing.T) {

	s := "Hellllllllo"
	tm := time.Now()

	//initDB("testing.db")

	StoreStat("abc", tm.UnixNano(), []byte(s))

	output := GetStats("abc", 100, tm.UnixNano())

	t.Log(output)

	for _, v := range output {
		t.Log(v)
	}

}
func Test_panictest(t *testing.T) {

	s := "Hellllllllo"
	tm := time.Now()

	StoreStat("abc", tm.UnixNano(), []byte(s))

	GetStats("gggg", 100, tm.UnixNano())

}
