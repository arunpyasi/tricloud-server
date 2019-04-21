package statstore

import (
	"testing"
	"time"
)

func Test_dbstore(t *testing.T) {

	s := "Hellllllllo"
	tm := time.Now()

	initDB("testing.db")

	StoreStat("abc", tm.UnixNano(), []byte(s))

	_, output, err := GetStats("abc", 100, tm.UnixNano())

	if err != nil {
		t.Error(err)
	}

	t.Log(string(output[0]))

	for _, v := range output {
		t.Log(string(v))
	}

	if string(output[0]) != s {
		t.Error(" error :(")
	}

}
func Test_panictest(t *testing.T) {

	s := "Hellllllllo"
	tm := time.Now()

	initDB("testing.db")

	StoreStat("abc", tm.UnixNano(), []byte(s))

	_, _, err := GetStats("gggg", 100, tm.UnixNano())

	if err == ErrNotExists {
		t.Log("panic test pass")
	}

}
