package provider

import (
	//	"fmt"
	"testing"
)

func TestInitAppFilesDir(t *testing.T) {
	if err := InitAppFilesDir("files"); err != nil {
		t.Error(err)
	}
}
