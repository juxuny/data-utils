package lib

import (
	"github.com/pkg/errors"
	"os"
)

func TouchDir(path string, perm os.FileMode) error {
	stat, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return os.Mkdir(path, perm)
	}
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.Errorf("%s is not a directory", path)
	}
	return nil
}
