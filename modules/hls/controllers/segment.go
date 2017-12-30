package hls

import (
	"context"
	"net/http"
	"fmt"
	"os/exec"
	"strconv"
	"bytes"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (c *Controller) serveSegment(ctx context.Context, req *requestData, w http.ResponseWriter, r *http.Request) {
	p := strconv.Itoa(req.BitRate)
	bitrate := convertPToBitRate(p)
	if bitrate == "" {
		w.Write([]byte("wrong resolution"))
	}

	segmentData, err := c.getSegmentDataFromPlaylist(ctx, req, req.TSIndex)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error to find segment file: %v", err), 500)
		
		logrus.Errorf("Error to find segment file: %v", err)

		return
	}

	/*param := []string{
		"-i", req.FullPath,
		"-ss", fmt.Sprint(segmentData["from"]),
		"-to", fmt.Sprint(segmentData["to"]),
		"-vcodec", "libx264",
		"-profile:v", "baseline",
		"-level", "3.0",
		"-vf", fmt.Sprintf("scale=-2:%d", req.BitRate),
		"-b:v", bitrate + "k",
		// "-force_key_frames", "5.040000,10.080000,15.120000,20.120000,25.160000,30.200000,35.200000,40.240000,45.280000,50.320000",
		// "-segment_times", "5.040000,10.080000,15.120000,20.120000,25.160000,30.200000,35.200000,40.240000,45.280000,50.320000",
		// "-hls_list_size", "10",
		// "-hls_list_size", "0",
		// "-hls_init_time" , fmt.Sprint(segmentData["from"]),
		// "-hls_time", fmt.Sprint(segmentData["to"] - segmentData["from"]),
		// "-hls_time", "5",
		// "-hls_flags", "split_by_time", 
		// "-hls_segment_filename", fmt.Sprintf("%d.ts",req.TSIndex), 
		"-hls_flags", "omit_endlist",
		"-hls_flags", "single_file",
		// "-hls_flags", "append_list",
		"-f", "hls",
		"-"}*/

	newParams := []string{
		"-ss", fmt.Sprintf("%v", segmentData["from"] + 0.1),
		"-t", "5",
		"-i", req.FullPath,
		"-vcodec", "libx264",
		"-strict", "experimental",
		"-acodec", "aac",
		"-pix_fmt", "yuv420p",
		"-r", "25",
		"-profile:v", "baseline",
		"-b:v", "2000k",
		"-maxrate", "2500k",
		"-f", "mpegts",
		"-"}

	logrus.Infof("%s %s", ffmpeg.String(), strings.Join(newParams, " "))
	cmd := exec.Command(ffmpeg.String(), newParams...)

	data, err := serveCommand(cmd, w, false)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening stdout of command: %v", err), 500)
		
		logrus.Errorf("Error opening stdout of command: %v", err)

		return
	}
	
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%d.ts",req.TSIndex))
	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeContent(w, r, fmt.Sprintf("attachment; filename=%d.ts",req.TSIndex), time.Now(), bytes.NewReader(data))
}
