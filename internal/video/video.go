package video

import (
	"bytes"
	"fmt"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Video struct {
	stream   *ffmpeg.Stream
	fps      int
	duration int
}

func ReadVideo(path string) *ffmpeg.Stream {
	return ffmpeg.Input(path)
}

func ImageToVideo(path string, duration int, fps int) Video {
	c := Video{}
	if duration < 0 {
		panic("duration can't be less than 0")
	}

	c.duration = duration
	c.fps = fps
	c.stream = ffmpeg.Input(path, ffmpeg.KwArgs{"loop": 1, "t": duration, "framerate": 30})

	return c

}

func (s *Video) GetVideoOutput(name string) *Video {
	s.stream = s.stream.Output(fmt.Sprintf("%s", name), ffmpeg.KwArgs{"c:v": "libx264", "r": s.fps, "framerate": 1})

	return s
}

func (s *Video) SaveVideo(name string) {
	buf := bytes.NewBuffer(nil)

	s.GetVideoOutput(name).stream.WithOutput(buf, os.Stdout).Run()

}

func MergeVideos(streams ...*Video) Video {
	merged := Video{}
	videos := []*ffmpeg.Stream{}
	newDuration := 0

	for _, v := range streams {
		videos = append(videos, v.stream)
		newDuration += v.duration
	}
	merged.stream = ffmpeg.Concat(videos)
	merged.fps = streams[0].fps
	merged.duration = newDuration
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

func (s *Video) AddFadeIn(duration int) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}
	s.stream = s.stream.
		Filter("fade", ffmpeg.Args{"t=in", fmt.Sprintf("d=%d", duration)})

	return s
}

func (s *Video) AddFadeOut(duration int) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}
	s.stream = s.stream.
		Filter("fade", ffmpeg.Args{"t=out", fmt.Sprintf("d=%d", duration), fmt.Sprintf("st=%d", s.duration-duration)})

	return s
}

func (s *Video) AddZoomIn(start int, duration int) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}
	s.stream = s.stream.
		Filter("zoompan", ffmpeg.Args{"z=min(max(zoom,pzoom)+0.001,1.25)", "d=1", "fps=60", "s=1800x1800"})

	return s
}
