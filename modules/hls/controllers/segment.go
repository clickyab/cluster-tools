package hls

import (
	"context"
	"net/http"

	"fmt"
	"os/exec"

	"strconv"

	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Controller) serveSegment(ctx context.Context, req *requestData, w http.ResponseWriter, r *http.Request) {
	startTime := float64(req.TSIndex) * segmentLen.Duration().Seconds()
	logrus.WithField("start_time", startTime).WithField("path", req.FullPath).Debug("Streaming...")
	p := strconv.Itoa(req.BitRate)
	bitrate := convertPToBitRate(p)
	if bitrate == "" {
		w.Write([]byte("wrong resolution"))
	}
	param := []string{
		"-i", req.FullPath,
		"-ss", fmt.Sprint(startTime),
		"-t", fmt.Sprint(int(segmentLen.Duration().Seconds())),
		"-vcodec", "libx264",
		"-profile:v", "baseline",
		"-level", "3.0",
		"-vf", fmt.Sprintf("scale=-2:%d", req.BitRate),
		"-b:v", bitrate + "k",
		"-f", "mpegts",
		"-"}

	logrus.Infof("%s %s", ffmpeg.String(), strings.Join(param, " "))
	cmd := exec.Command(ffmpeg.String(), param...)
	serveCommand(cmd, w)
	//assert.Nil(err)

}
