package hls

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Controller) servePlaylist(ctx context.Context, req *requestData, w http.ResponseWriter, r *http.Request) {
	duration := req.VInfo.Duration

	//w.Header()["Content-Type"] = []string{"application/vnd.apple.mpegurl"}

	fmt.Fprint(w, "#EXTM3U\n")
	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprint(w, "#EXT-X-ALLOW-CACHE:YES\n")
	fmt.Fprint(w, "#EXT-X-TARGETDURATION:5\n")
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := duration
	segmentIndex := 0

	for leftover > 0 {
		if leftover > segmentLen.Duration().Seconds() {
			fmt.Fprintf(w, "#EXTINF: %f,\n", segmentLen.Duration().Seconds())
		} else {
			fmt.Fprintf(w, "#EXTINF: %f,\n", leftover)
		}
		fmt.Fprint(w, fmt.Sprintf("%d.ts", segmentIndex)+"\n")
		segmentIndex++
		leftover = leftover - segmentLen.Duration().Seconds()
	}
	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
}
