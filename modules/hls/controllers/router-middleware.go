package hls

import (
	"net/http"

	"context"

	"regexp"

	"strconv"

	"fmt"

	"path/filepath"

	"github.com/clickyab/services/assert"
	"github.com/rs/xmux"
	"github.com/sirupsen/logrus"
)

var (
	index    *regexp.Regexp
	subIndex *regexp.Regexp
	subTs    *regexp.Regexp
)

type requestData struct {
	MasterIndex bool // this is manifest request
	BitRate     int
	Index       bool // this is a bit rate specific request
	TSIndex     int  // this is the index of ts requested

	FullPath string
	File     string
	VInfo    *VideoInfo
}

func getRequestData(f string) (*requestData, error) {
	var (
		file requestData
		err  error
	)
	if parts := index.FindStringSubmatch(f); len(parts) == 2 {
		// Index request
		file.File = parts[1]
		file.MasterIndex = true
	} else if parts := subIndex.FindStringSubmatch(f); len(parts) == 3 {
		file.File = parts[1]
		file.BitRate, err = strconv.Atoi(parts[2])
		assert.Nil(err)
		file.Index = true
	} else if parts := subTs.FindStringSubmatch(f); len(parts) == 4 {
		file.File = parts[1]
		file.BitRate, err = strconv.Atoi(parts[2])
		assert.Nil(err)
		file.TSIndex, err = strconv.Atoi(parts[3])
		assert.Nil(err)
	} else {
		return nil, fmt.Errorf("invalid request")
	}
	file.FullPath = filepath.Join(rootPath.String(), file.File)
	file.VInfo, err = GetVideoInformation(file.FullPath)
	if err != nil {
		return nil, err
	}

	logrus.Infof("%+v", file)

	return &file, nil
}

// router handler try to get data from request and add that into context
func (c *Controller) entry(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	path := xmux.Param(ctx, "path")
	logrus.WithField("path", path).Debug("new request")
	req, err := getRequestData(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if req.MasterIndex {
		c.serveMasterM3U8(ctx, req, w, r)
	} else if req.Index {
		c.servePlaylist(ctx, req, w, r)
	} else {
		c.serveSegment(ctx, req, w, r)
	}
}
