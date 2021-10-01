package ffmpeg

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

func FFProbe(fileName string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", fileName)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "start ffprobe failed")
	}
	return buf.Bytes(), nil
}

func parseData(data []byte) (info *VideoInfo, err error) {
	info = &VideoInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, errors.Wrap(err, "parse json failed")
	}
	return
}
