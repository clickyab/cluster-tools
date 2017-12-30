package hls

import (
	"os"
	"os/exec"
	"crypto/sha1"
	"io"
	"strconv"
	"strings"	
	"context"
	"fmt"
	"bufio"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *Controller) servePlaylist(ctx context.Context, req *requestData, w http.ResponseWriter, r *http.Request) {
	playlistPath, err := checkAndGeneratePlaylist(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening stdout of command: %v", err), 404)
		
		logrus.Errorf("Error opening stdout of command: %v", err)

		return
	}

	logrus.WithField("playlist path", playlistPath).Debug("request")

	//Check if file exists and open
	Openfile, err := os.Open(playlistPath)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file

	//Get the file size
	fileStat, _ := Openfile.Stat()                     //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string
	
	filename := fmt.Sprintf("%s_%d_index.m3u8", req.VInfo.Name, req.BitRate)
	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Content-Length", fileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client

	return
}

func (c *Controller) getSegmentDataFromPlaylist(ctx context.Context, req *requestData, segmentNumber int) (map[string]float64, error) {
	manifestFilePath := fmt.Sprintf("%s/%s/%d_index.m3u8", req.CacheFolder, sha(req.FullPath), req.BitRate)

	_,segmentsData, err := readM3u8(manifestFilePath)

	if err != nil {
		return nil, err
	}

	return segmentsData[segmentNumber], nil
}

func checkAndGeneratePlaylist(ctx context.Context, req *requestData) (string, error) {
	filePath, err := checkAndGetPlaylistPath(req)

	if filePath == "" {
		result, err := generatePlaylist(ctx, req)

		return result, err
	}

	return filePath, err
}

func checkAndGetPlaylistPath(req *requestData) (string, error) {
	manifestFilePath := fmt.Sprintf("%s/%s/%d_index.m3u8", req.CacheFolder, sha(req.FullPath), req.BitRate)

	_, err := os.Stat(manifestFilePath);

	if err == nil {
		return manifestFilePath, nil
	}

	return "", err
}

func generatePlaylist(ctx context.Context, req *requestData) (string, error) {
	keyframesList,err := getVideoKeyframesList(ctx, req)

	mkdirError := os.MkdirAll(req.CacheFolder + "/" + sha(req.FullPath), 0777)
	if mkdirError != nil {
		return "", mkdirError
	}

	manifestFilePath := fmt.Sprintf("%s/%s/%d_index.m3u8", req.CacheFolder, sha(req.FullPath), req.BitRate)

	out, err := os.Create(manifestFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

    fileWriter := bufio.NewWriter(out)
	defer fileWriter.Flush()

	fmt.Fprint(fileWriter, "#EXTM3U\n")
	fmt.Fprint(fileWriter, "#EXT-X-VERSION:3\n")
	fmt.Fprint(fileWriter, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprint(fileWriter, "#EXT-X-ALLOW-CACHE:YES\n")
	fmt.Fprint(fileWriter, "#EXT-X-TARGETDURATION:6\n")
	// fmt.Fprint(fileWriter, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	segmentIndex := 0
	var lastTime float64
	for _,time := range keyframesList {
		floatTime, _ := strconv.ParseFloat(time, 64)

		if floatTime - lastTime > 5 {

			fmt.Fprintf(fileWriter, "#EXTINF: %f,\n", (floatTime - lastTime))
			fmt.Fprint(fileWriter, fmt.Sprintf("%d.ts", segmentIndex)+"\n")

			lastTime = floatTime
			segmentIndex++
		}
	}

	fmt.Fprint(fileWriter, "#EXT-X-ENDLIST\n")

	return manifestFilePath, err
}


func getVideoKeyframesList(ctx context.Context, req *requestData) ([]string, error) {
	param := []string{
		"-loglevel", "error",
		"-select_streams", "v:0",
		"-show_entries", "frame=key_frame,pkt_pts_time",
		"-of", "csv=print_section=0",
		req.FullPath}

	logrus.Infof("%s %s", ffprobe.String(), strings.Join(param, " "))
	cmd := exec.Command(ffprobe.String(), param...)

	var emptyWriter io.Writer
	result, err := serveCommand(cmd, emptyWriter,false)
	
	keyframesList := strings.Split(string(result), "\n0,")

	return keyframesList, err
}

func readM3u8(path string) (map[string]string, []map[string]float64, error) {
	var parsedSegmentsData []map[string]float64
	parsedMetaData := make(map[string]string)
	var data []string

	fileLines, err := readLines(path)
	total := 0.0

	for i := 0; i < len(fileLines); i++ {

		if strings.Contains(fileLines[i], "#EXT-X-") && fileLines[i] != "#EXT-X-ENDLIST" {
			
			data = strings.Split(fileLines[i], ":")

			parsedMetaData[data[0][7:len(data[0]) - 1]] = data[1]

		} else if strings.Contains(fileLines[i], "#EXTINF:") {
			data = strings.Split(fileLines[i], "#EXTINF:")

			segmentDuration, err := strconv.ParseFloat(strings.TrimSpace(data[1][:len(data[1]) - 2]), 64)

			if err == nil {
				// segmentDuration = strconv.FormatFloat(segmentDuration, 'f', 6, 64)
				total = total + segmentDuration

				from := 0.0
				if len(parsedSegmentsData) > 0 {
					from = parsedSegmentsData[len(parsedSegmentsData) - 1]["to"]
				}

				tmp := map[string]float64{"from": from, "to": total}

				parsedSegmentsData = append(parsedSegmentsData, tmp)
			}
		}
	}

	return parsedMetaData, parsedSegmentsData, err
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }

  return lines, scanner.Err()
}

// Sha1 is the sha1 generation func
func sha(k string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(k))

	return fmt.Sprintf("%x", h.Sum(nil))
}