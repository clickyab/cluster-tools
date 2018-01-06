package hls

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"io"

	"encoding/binary"

	"time"

	"github.com/clickyab/services/kv"
	"github.com/sirupsen/logrus"
)

var videoSuffixes = []string{".mp4", ".avi", ".mkv", ".flv", ".wmv", ".mov", ".mpg", ".m4v", ".jbh"}

// VideoInfo is the video information
type VideoInfo struct {
	tag      string  // this field is not required in the outside, and is just for caching
	Duration float64 `json:"duration"`
}

// Encode is the cache encoder
func (vi VideoInfo) Encode(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, vi.Duration)
}

// Decode is the cache decoder
func (vi *VideoInfo) Decode(r io.Reader) error {
	var d float64
	if err := binary.Read(r, binary.BigEndian, &d); err != nil {
		return err
	}
	vi.Duration = d
	return nil
}

// String is the cache tag
func (vi VideoInfo) String() string {
	return vi.tag
}

func validExt(name string) bool {
	for _, suffix := range videoSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func getRawFFMPEGInfo(path string) ([]byte, error) {
	logrus.WithField("path", path).Debug("Executing ffprobe")
	cmd := exec.Command(ffprobe.String(), "-v", "quiet", "-print_format", "json", "-show_format", ""+path+"")
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing ffprobe for file '%v': %s", path, err)
	}
	return data, nil
}

func getFFMPEGJson(path string) (map[string]interface{}, error) {
	data, err := getRawFFMPEGInfo(path)
	if err != nil {
		return nil, err
	}
	var info map[string]interface{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON from ffprobe output for file '%v': %s", path, err)
	}
	return info, nil
}

// GetVideoInformation return the video information from executing the ffprobe
func GetVideoInformation(path string) (*VideoInfo, error) {
	var res VideoInfo
	if err := kv.Hit(path, &res); err == nil {
		return &res, nil
	}
	if !validExt(path) {
		return nil, fmt.Errorf("invalid ext")
	}

	info, err := getFFMPEGJson(path)
	if err != nil {
		return nil, err
	}
	logrus.WithField("path", path).WithField("result", info).Debug("ffprobe returned")
	if _, ok := info["format"]; !ok {
		return nil, fmt.Errorf("ffprobe data for '%v' does not contain format info", path)
	}
	format := info["format"].(map[string]interface{})
	if _, ok := format["duration"]; !ok {
		return nil, fmt.Errorf("ffprobe format data for '%v' does not contain duration", path)
	}
	duration, err := strconv.ParseFloat(format["duration"].(string), 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse duration (%v) of '%v': %s", format["duration"].(string), path, err)
	}
	res = VideoInfo{tag: path, Duration: duration}
	_ = kv.Do(res.String(), &res, time.Hour*60, nil)

	return &res, nil
}
