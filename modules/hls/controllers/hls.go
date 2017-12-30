package hls

import (
	"io"
	"io/ioutil"
	"os/exec"
	"syscall"
	"time"

	"regexp"

	"fmt"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework/controller"
	"github.com/clickyab/services/framework/router"
	"github.com/sirupsen/logrus"
)

var (
	rootPath   = config.RegisterString("mantis.hls.root_folder", "/home/develop/go/src/clickyab.com/cluster-tools/videos/", "the cdn root")
	cachePath  = config.RegisterString("mantis.hls.chache_folder", "/home/develop/go/cache", "the cache root")
	ffmpeg     = config.RegisterString("mantis.hls.ffmpeg", "/usr/bin/ffmpeg", "the ffmpeg binary")
	ffprobe    = config.RegisterString("mantis.hls.ffmpeg", "/usr/bin/ffprobe", "the ffprobe binary")
	segmentLen = config.RegisterDuration("mantis.hls.segment_len", 6*time.Second, "segment len")

	// comma separated strings
	resolutions = config.RegisterString("mantis.hls.resolutions", "360,240,144", "video resolution")
)

func serveCommand(cmd *exec.Cmd, w io.Writer, streamResultToWriter bool) ([]byte, error) {
	var result []byte

	stdout, err := cmd.StdoutPipe()

	defer stdout.Close()
	// defer cmd.Process.Kill()

	if err != nil {
		logrus.Errorf("Error opening stdout of command: %v", err)
		return result, err
	}

	err = cmd.Start()
	if err != nil {
		logrus.Errorf("Error starting command: %v", err)
		return result, err
	}

	

	if streamResultToWriter {
		_, err = io.Copy(w, stdout)

		if err != nil {
			logrus.Errorf("Error copying data to client: %v", err)
			// Ask the process to exit
			cmd.Process.Signal(syscall.SIGKILL)
			cmd.Process.Wait()

			return result, err
		}

	}
	
	result, err = ioutil.ReadAll(stdout)

	return result, err
}

// Controller is the controller for hls
type Controller struct {
	controller.Base
}

// Routes is the controller registration function
func (c *Controller) Routes(r router.Mux) {
	seg := fmt.Sprint(int(segmentLen.Duration().Seconds()))
	index = regexp.MustCompile("/(.*)/" + seg + "/manifest[.]m3u8$")
	subIndex = regexp.MustCompile("/(.*)/" + seg + "/([0-9]+)/index[.]m3u8$")
	subTs = regexp.MustCompile("/(.*)/" + seg + "/([0-9]+)/([0-9]+)[.]ts$")
	
	// TODO : need the quality in route. there is no quality flag here, also a master manifest contain all qualities
	r.GET("/static/*path", c.serveStatic)
	r.GET("/fly/*path", c.entry)
}

func init() {
	router.Register(&Controller{})
}
