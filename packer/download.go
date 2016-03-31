package packer

import (
	"io/ioutil"
	"os"
)

func TemporaryDir() (string, error) {
	tmp, err := ioutil.TempDir(os.TempDir(), "bakery")
	if err != nil {
		return "", err
	}

	return tmp, nil
}
