package hls

import (
	"io"
	"os/exec"
	"syscall"
	"time"

	"regexp"

	"fmt"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework/controller"
	"github.com/clickyab/services/framework/router"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"github.com/sirupsen/logrus"
)

var (
	rootPath   = config.RegisterString("mantis.hls.root_folder", "/home/f0rud/", "the cdn root")
	ffmpeg     = config.RegisterString("mantis.hls.ffmpeg", "/usr/bin/ffmpeg", "the ffmpeg binary")
	ffprobe    = config.RegisterString("mantis.hls.ffmpeg", "/usr/bin/ffprobe", "the ffprobe binary")
	segmentLen = config.RegisterDuration("mantis.hls.segment_len", 7*time.Second, "segment len")

	// comma separated strings
	resolutions = config.RegisterString("mantis.hls.resolutions", "360,240,144", "video resolution")
)

func serveCommand(cmd *exec.Cmd, w io.Writer) error {
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		logrus.Errorf("Error opening stdout of command: %v", err)
		return err
	}
	err = cmd.Start()
	if err != nil {
		logrus.Errorf("Error starting command: %v", err)
		return err
	}
	_, err = io.Copy(w, stdout)
	if err != nil {
		logrus.Errorf("Error copying data to client: %v", err)
		// Ask the process to exit
		cmd.Process.Signal(syscall.SIGKILL)
		cmd.Process.Wait()
		return err
	}
	return cmd.Wait()
}

// Controller is the controller for hls
type Controller struct {
	controller.Base
}

// Routes is the controller registration function
func (c *Controller) Routes(r *xmux.Mux, mountPoint string) {
	seg := fmt.Sprint(int(segmentLen.Duration().Seconds()))
	index = regexp.MustCompile("/(.*)/" + seg + "/manifest[.]m3u8$")
	subIndex = regexp.MustCompile("/(.*)/" + seg + "/([0-9]+)/index[.]m3u8$")
	subTs = regexp.MustCompile("/(.*)/" + seg + "/([0-9]+)/([0-9]+)[.]ts$")

	// TODO : need the quality in route. there is no quality flag here, also a master manifest contain all qualities
	logrus.Debugf(mountPoint + "/*path")
	r.GET(mountPoint+"/*path", xhandler.HandlerFuncC(c.entry))
}

func init() {
	router.Register(&Controller{})
}
