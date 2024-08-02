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
		videos = append(videos, v.stream)
	}
	merged.stream = ffmpeg.Concat(videos)

	return merged
}

func (s *Video) VFlip() {
	s.stream = s.stream.VFlip()
}

func (s *Video) Crop(width int, height int) *Video {
	ratio := float64(width) / float64(height)
	w := fmt.Sprintf("%d", width)
	h := fmt.Sprintf("%d", height)

	fortmattedWidth := fmt.Sprintf("if(gt(iw,ih), min(iw*%f, iw), iw)", ratio)
	fortmattedHeight := fmt.Sprintf("if(gt(ih,iw), min(ih*%f, ih), ih)", ratio)

	s.stream = s.stream.
		Filter("crop", ffmpeg.Args{fortmattedWidth, fortmattedHeight}).
		Filter("scale", ffmpeg.Args{w, h}).
		Filter("setsar", ffmpeg.Args{fmt.Sprintf("%f", ratio)})
	return s
}

// if more custom filters are needed
func (s *Video) Filter(name string, args ffmpeg.Args) *Video {
	s.stream = s.stream.Filter(name, args)
	return s
}

// if custom input is needed
func (s *Video) Input(path string, kwargs ffmpeg.KwArgs) *Video {
	s.stream = ffmpeg.Input(path, kwargs)
	return s
}

func (s *Video) Output(name string, kwargs ffmpeg.KwArgs) *Video {
	s.stream = s.stream.Output(name, kwargs)
	return s
}
