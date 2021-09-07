package canvas

import (
	"github.com/pkg/errors"
	"io/ioutil"
)

var fontCache map[string][]byte

func loadFont(file string) ([]byte, error) {
	if data, b := fontCache[file]; b {
		return data, nil
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "load font failed")
	}
	return data, nil
}
