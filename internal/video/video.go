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

func ImageToVideo(path string, duration int, width int, height int) Video {
	c := Video{}
	if duration < 0 {
		panic("duration can't be less than 0")
	}

	c.stream = ffmpeg.
		Input(path, ffmpeg.KwArgs{"loop": 1, "t": duration, "framerate": 1}).
		Crop(0, 0, width, height).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", width, height)})

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
		videos = append(videos, v.stream)
	}
	merged.stream = ffmpeg.Concat(videos)

	return merged
}

func (s *Video) VFlip() {
	s.stream = s.stream.VFlip()
}
