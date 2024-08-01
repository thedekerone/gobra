package video

import (
	"bytes"
	"fmt"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type CustomStream struct {
	Stream *ffmpeg.Stream
}

func ReadVideo(path string) *ffmpeg.Stream {
	return ffmpeg.Input(path)
}

func ImageToVideo(path string, duration int) CustomStream {
	c := CustomStream{}
	if duration < 0 {
		panic("duration can't be less than 0")
	}

	c.Stream = ffmpeg.Input(path, ffmpeg.KwArgs{"loop": 1, "t": duration, "framerate": 1}).Crop(0, 0, 500, 500).Filter("scale", ffmpeg.Args{"500:500"})

	return c

}

func (s *CustomStream) GetVideoOutput(name string) *CustomStream {
	s.Stream = s.Stream.Output(fmt.Sprintf("%s", name), ffmpeg.KwArgs{"c:v": "libx264", "r": 25, "framerate": 1})

	return s
}

func (s *CustomStream) SaveVideo(name string) {
	buf := bytes.NewBuffer(nil)

	s.GetVideoOutput(name).Stream.WithOutput(buf, os.Stdout).Run()

}

func MergeVideos(streams ...*CustomStream) CustomStream {
	merged := CustomStream{}
	videos := []*ffmpeg.Stream{}

	for _, v := range streams {
		videos = append(videos, v.Stream)
	}
	merged.Stream = ffmpeg.Concat(videos)

	return merged
}
