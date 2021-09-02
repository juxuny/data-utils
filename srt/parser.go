package srt

import (
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Block struct {
	Id        int64
	StartTime time.Time
	EndTime   time.Time
	Subtitle  *Node
}

func parseBlock(blockData string) (Block, error) {
	var ret Block
	return ret, nil
}

func Parse(data []byte) ([]Block, error) {
	blocks := strings.Split(string(data), "\n\n")
	ret := make([]Block, 0)
	for index, item := range blocks {
		if b, err := parseBlock(item); err != nil {
			ret = append(ret, b)
		} else {
			return nil, errors.Wrapf(err, "parse block failed, index=%d", index)
		}
	}
	return ret, nil
}
