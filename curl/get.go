package curl

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func Get(url string, headers ...http.Header) (data []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf(resp.Status)
	}
	data, err = ioutil.ReadAll(resp.Body)
	return
}
