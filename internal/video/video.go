package video

import (
	"bytes"
	"fmt"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Video struct {
	stream *ffmpeg.Stream
}

func ReadVideo(path string) *ffmpeg.Stream {
	return ffmpeg.Input(path)
}

func ImageToVideo(path string, duration int) Video {
	c := Video{}
	if duration < 0 {
		panic("duration can't be less than 0")
	}

	c.stream = ffmpeg.
		Input(path, ffmpeg.KwArgs{"loop": 1, "t": duration, "framerate": 1})

	return c

}

func (s *Video) GetVideoOutput(name string) *Video {
	s.stream = s.stream.Output(fmt.Sprintf("%s", name), ffmpeg.KwArgs{"c:v": "libx264", "r": 25, "framerate": 1})

	return s
}

func (s *Video) SaveVideo(name string) {
	buf := bytes.NewBuffer(nil)

	s.GetVideoOutput(name).stream.WithOutput(buf, os.Stdout).Run()

}

func MergeVideos(streams ...*Video) Video {
	merged := Video{}
	videos := []*ffmpeg.Stream{}

	for _, v := range streams {
		videos = append(videos, v.Crop(900, 1600).stream)
	}
	merged.stream = ffmpeg.Concat(videos)

	return merged
}

func (s *Video) VFlip() {
	s.stream = s.stream.VFlip()
}

func (s *Video) Crop(width int, height int) *Video {
	w := fmt.Sprintf("%d", width)
	h := fmt.Sprintf("%d", height)
	ratio := float64(width) / float64(height)

	s.stream = s.stream.
		Filter("crop", ffmpeg.Args{fmt.Sprintf("min(iw*%f,%s)", ratio, w), fmt.Sprintf("min(ih*%f,%s)", ratio, h)}).
		Filter("scale", ffmpeg.Args{w, h}).
		Filter("setsar", ffmpeg.Args{fmt.Sprintf("%d", width/height)})
	return s
}
