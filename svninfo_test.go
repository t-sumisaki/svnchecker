package svnchecker

import (
	"fmt"
	"testing"
)

func TestGetInfo(t *testing.T) {

	path := "H:\\mgc\\trunk\\mgc\\mgc.uproject"

	info, err := GetInfo(path)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(info)
}
