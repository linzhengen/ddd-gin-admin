package structure

import (
	"github.com/jinzhu/copier"
)

func Copy(s, ts interface{}) {
	err := copier.Copy(ts, s)
	if err != nil {
		panic(err)
	}
}
