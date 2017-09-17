package hls

import (
	"context"
	"fmt"
	"net/http"

	"strings"
)

var (
	pToBitRate = map[string]string{
		"1440": "2500",
		"1080": "2000",
		"720":  "950",
		"480":  "700",
		"360":  "550",
		"240":  "350",
		"144":  "190",
	}

	pToSize = map[string]string{
		"1440": "2560x1440",
		"1080": "1920x1080",
		"720":  "1280x720",
		"480":  "854x480",
		"360":  "640x360",
		"240":  "426x240",
		"144":  "256x144",
	}
)

func convertPToSize(p string) string {
	return pToSize[p]
}

func convertPToBitRate(p string) string {
	return pToBitRate[p]
}

func (c *Controller) serveMasterM3U8(ctx context.Context, req *requestData, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "#EXTM3U")
	for _, p := range strings.Split(resolutions.String(), ",") {
		fmt.Fprintf(w, "#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=%s, RESOLUTION=%s\n", convertPToBitRate(p), convertPToSize(p))
		fmt.Fprintf(w, "%s/index.m3u8\n", p)
	}
}
