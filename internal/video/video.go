package video

import (
	"bytes"
	"fmt"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Video struct {
	stream   *ffmpeg.Stream
	duration int
	config   Config
}

type Config struct {
	Width       int
	Height      int
	Fps         int
	AspectRatio float64
}

func ReadVideo(path string, config Config) Video {
	c := Video{}
	c.stream = ffmpeg.Input(path)

	c.config = config

	return c
}

func (s *Video) TrimVideo(start float32, end float32) *Video {
	if start < 0 || end < 0 {
		panic("start and end can't be less than 0")
	}
	s.stream = s.stream.Trim(ffmpeg.KwArgs{"start": start, "end": end})
	s.duration = int(end - start)
	return s
}

func (s *Video) ScaleVideo(width int, height int) *Video {
	if width < 0 || height < 0 {
		panic("width and height can't be less than 0")
	}
	s.stream = s.stream.Filter("scale", ffmpeg.Args{fmt.Sprintf("%d", width), fmt.Sprintf("%d", height)})
	s.config.Width = width
	s.config.Height = height
	return s
}

func ImageToVideo(path string, duration int, fps int) Video {
	c := Video{}
	if duration < 0 {
		panic("duration can't be less than 0")
	}

	c.duration = duration
	c.config.Fps = fps
	c.stream = ffmpeg.Input(path, ffmpeg.KwArgs{"loop": 1, "t": duration, "framerate": 30})

	return c

}

func (s *Video) GetVideoOutput(name string) *Video {
	s.stream = s.stream.Output(fmt.Sprintf("%s", name), ffmpeg.KwArgs{"c:a": "copy", "r": s.config.Fps, "framerate": 1})

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
	merged.config.Fps = streams[0].config.Fps
	merged.duration = newDuration
	return merged
}

func (s *Video) VFlip() {
	s.stream = s.stream.VFlip()
}

func (s *Video) Crop(width int, height int) *Video {
	ratio := 1.0
	w := fmt.Sprintf("%d", width)
	h := fmt.Sprintf("%d", height)

	fortmattedWidth := fmt.Sprintf("if(gt(iw,ih), min(iw*%f, iw), iw)", ratio)
	fortmattedHeight := fmt.Sprintf("if(gt(ih,iw), min(ih*%f, ih), ih)", ratio)

	s.stream = s.stream.
		Filter("crop", ffmpeg.Args{fortmattedWidth, fortmattedHeight}).
		Filter("scale", ffmpeg.Args{w, h}).
		Filter("setsar", ffmpeg.Args{fmt.Sprintf("%f", ratio)})

	s.config.Width = width
	s.config.Height = height
	s.config.AspectRatio = ratio
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

func (s *Video) AddFadeIn(duration float32) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}
	s.stream = s.stream.
		Filter("fade", ffmpeg.Args{"t=in", fmt.Sprintf("d=%f", duration)})

	return s
}

func (s *Video) AddFadeOut(duration float32) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}
	s.stream = s.stream.
		Filter("fade", ffmpeg.Args{"t=out", fmt.Sprintf("d=%f", duration), fmt.Sprintf("st=%f", float32(s.duration)-duration)})

	return s
}

func (s *Video) AddZoomIn(zoom float64) *Video {
	if zoom < 1 {
		panic("duration can't be less than 1")
	}
	s.stream = s.stream.
		Filter("scale", ffmpeg.Args{fmt.Sprintf("w=iw*%d:h=ih*%d", 4, 4)}).
		Filter("zoompan",
			ffmpeg.Args{
				fmt.Sprintf("z=min(max(pzoom,zoom) + 0.001,%f)", zoom),
				fmt.Sprintf("fps=%d", s.config.Fps),
				fmt.Sprintf("d=%d*%d", s.duration, s.config.Fps),
				"x=iw/2-(iw/zoom/2)",
			})

	return s
}

func CreateZoomPanVideoFromImage(path string, duration int, zoom float32, config Config) Video {
	if zoom < 1 || duration < 0 {
		panic("duration can't be less than 0 and zoom can't be less than 1")
	}

	v := Video{}
	v.config = config
	i := ffmpeg.Input(path).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d", v.config.Width*4), fmt.Sprintf("%d", v.config.Height*4)}).
		Filter("zoompan",
			ffmpeg.Args{
				fmt.Sprintf("z=min(max(pzoom,zoom) + 0.001,%f)", zoom),
				fmt.Sprintf("fps=%d", v.config.Fps),
				fmt.Sprintf("d=%d*%d", duration, v.config.Fps),
				"x=iw/2-(iw/zoom/2)",
			})
	v.stream = i
	v.Crop(v.config.Width, v.config.Height)
	v.duration = duration

	return v
}

func (s *Video) AddSubtitles(path string) *Video {
	s.stream = s.stream.
		Filter("subtitles", ffmpeg.Args{path})

	return s
}

func (s *Video) AddOverlayVideo(v *Video, x string, y string) *Video {
	s.stream = s.stream.
		Overlay(v.stream,
			"pass",
			ffmpeg.KwArgs{
				"x": x,
				"y": y,
			})

	return s
}

func MergeVideoWithAudio(a *Audio, v *Video) {
	buf := bytes.NewBuffer(nil)

	ffmpeg.Input("test.mp4", ffmpeg.KwArgs{"i": "assets/subway.mp3"}).
		Output("test222.mov", ffmpeg.KwArgs{"shortest": "", "c:v": "libx264", "c:a": "aac"}).
		WithOutput(buf, os.Stdout).
		OverWriteOutput().
		Run()

}

type Audio struct {
	stream *ffmpeg.Stream
}

func ReadAudio(path string) Audio {
	c := Audio{}
	c.stream = ffmpeg.Input(path)

	return c
}
