package dt

import "github.com/pkg/errors"

var (
	ErrConnectFailed  = errors.Errorf("connect failed")
	ErrSendDataFailed = errors.Errorf("send request failed")
	ErrReadDataFailed = errors.Errorf("read response failed")
	ErrEmptyResponse  = errors.Errorf("empty response")
	ErrEOF            = errors.Errorf("EOF")
	ErrNotFound       = errors.Errorf("NotFound")
)

func IsEof(err error) bool {
	if err != nil {
		return err.Error() == ErrEOF.Error()
	}
	return false
}
