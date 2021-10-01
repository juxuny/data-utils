package ffmpeg

import (
	"fmt"
	"strconv"
	"strings"
)

type StringNumber int64

func (t *StringNumber) UnmarshalJSON(data []byte) error {
	in := strings.Trim(string(data), "\"")
	v, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return err
	}
	*t = StringNumber(v)
	return nil
}

func (t StringNumber) MarshalJSON() ([]byte, error) {
	return []byte("\"" + fmt.Sprintf("%v", t) + "\""), nil
}

type StringBool bool

func (t *StringBool) UnmarshalJSON(data []byte) error {
	in := strings.Trim(string(data), "\"")
	*t = in == "true" || in == "1"
	return nil
}

func (t StringBool) MarshalJSON() ([]byte, error) {
	return []byte("\"" + fmt.Sprintf("%v", t) + "\""), nil
}

type IntBool int

const (
	IntBoolTrue  = IntBool(1)
	IntBoolFalse = IntBool(0)
)

func (t *IntBool) UnmarshalJSON(data []byte) error {
	in := strings.Trim(string(data), "\"")
	v, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return err
	}
	*t = IntBool(v)
	return nil
}

func (t IntBool) MarshalJSON() ([]byte, error) {
	return []byte("\"" + fmt.Sprintf("%v", t) + "\""), nil
}

type StringFloat float64

func (t *StringFloat) UnmarshalJSON(data []byte) error {
	in := strings.Trim(string(data), "\"")
	v, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	*t = StringFloat(v)
	return nil
}

func (t StringFloat) MarshalJSON() ([]byte, error) {
	return []byte("\"" + fmt.Sprintf("%v", t) + "\""), nil
}
