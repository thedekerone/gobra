package gobra

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

func NewVideo(path string, config Config) *Video {
	return &Video{
		stream: ffmpeg.Input(path),
		config: config,
	}
}

func NewVideoWithAudio(path, audioPath string, config Config) *Video {
	return &Video{
		stream: ffmpeg.Input(path, ffmpeg.KwArgs{"i": audioPath}),
		config: config,
	}
}

func (v *Video) Trim(start, end float32) *Video {
	if start < 0 || end < 0 {
		panic("start and end can't be less than 0")
	}
	v.stream = v.stream.Trim(ffmpeg.KwArgs{"start": start, "end": end})
	v.duration = int(end - start)
	return v
}

func (v *Video) Scale(width, height int) *Video {
	if width < 0 || height < 0 {
		panic("width and height can't be less than 0")
	}
	v.stream = v.stream.Filter("scale", ffmpeg.Args{fmt.Sprintf("%d", width), fmt.Sprintf("%d", height)})
	v.config.Width = width
	v.config.Height = height
	return v
}

func NewImageVideo(path string, duration, fps int) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}

	return &Video{
		duration: duration,
		config: Config{
			Fps: fps,
		},
		stream: ffmpeg.Input(path, ffmpeg.KwArgs{"loop": 1, "t": duration, "framerate": 30}),
	}
}

func (v *Video) Output(name string) *Video {
	v.stream = v.stream.Output(fmt.Sprintf("%s", name))
	return v
}

func (v *Video) OutputWithSubtitles(name, subtitlesPath string) *Video {
	v.stream = v.stream.Output(fmt.Sprintf("%s", name), ffmpeg.KwArgs{"vf": fmt.Sprintf("subtitles=%s:force_style='Alignment=10'", subtitlesPath)})
	return v
}

func (v *Video) Save(name string) {
	buf := bytes.NewBuffer(nil)
	v.Output(name).stream.WithOutput(buf, os.Stdout).OverWriteOutput().Run()
}

func (v *Video) SaveWithSubtitles(name, subtitlesPath string) {
	buf := bytes.NewBuffer(nil)
	v.OutputWithSubtitles(name, subtitlesPath).stream.WithOutput(buf, os.Stdout).OverWriteOutput().Run()
}

func MergeVideos(videos ...*Video) *Video {
	merged := &Video{}
	videoStreams := []*ffmpeg.Stream{}
	newDuration := 0

	for _, v := range videos {
		videoStreams = append(videoStreams, v.stream)
		newDuration += v.duration
	}
	merged.stream = ffmpeg.Concat(videoStreams)
	merged.config.Fps = videos[0].config.Fps
	merged.duration = newDuration
	return merged
}

func (v *Video) VFlip() *Video {
	v.stream = v.stream.VFlip()
	return v
}

func (v *Video) Crop(width, height int) *Video {
	ratio := 1.0
	w := fmt.Sprintf("%d", width)
	h := fmt.Sprintf("%d", height)

	formattedWidth := fmt.Sprintf("if(gt(iw,ih), min(iw*%f, iw), iw)", ratio)
	formattedHeight := fmt.Sprintf("if(gt(ih,iw), min(ih*%f, ih), ih)", ratio)

	v.stream = v.stream.
		Filter("crop", ffmpeg.Args{formattedWidth, formattedHeight}).
		Filter("scale", ffmpeg.Args{w, h}).
		Filter("setsar", ffmpeg.Args{fmt.Sprintf("%f", ratio)})

	v.config.Width = width
	v.config.Height = height
	v.config.AspectRatio = ratio
	return v
}

func (v *Video) Filter(name string, args ffmpeg.Args) *Video {
	v.stream = v.stream.Filter(name, args)
	return v
}

func (v *Video) Input(path string, kwargs ffmpeg.KwArgs) *Video {
	v.stream = ffmpeg.Input(path, kwargs)
	return v
}

func (v *Video) AddFadeIn(duration float32) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}
	v.stream = v.stream.
		Filter("fade", ffmpeg.Args{"t=in", fmt.Sprintf("d=%f", duration)})

	return v
}

func (v *Video) AddFadeOut(duration float32) *Video {
	if duration < 0 {
		panic("duration can't be less than 0")
	}
	v.stream = v.stream.
		Filter("fade", ffmpeg.Args{"t=out", fmt.Sprintf("d=%f", duration), fmt.Sprintf("st=%f", float32(v.duration)-duration)})

	return v
}

func (v *Video) AddZoomIn(zoom float64) *Video {
	if zoom < 1 {
		panic("zoom can't be less than 1")
	}
	v.stream = v.stream.
		Filter("scale", ffmpeg.Args{fmt.Sprintf("w=iw*%d:h=ih*%d", 4, 4)}).
		Filter("zoompan",
			ffmpeg.Args{
				fmt.Sprintf("z=min(max(pzoom,zoom) + 0.001,%f)", zoom),
				fmt.Sprintf("fps=%d", v.config.Fps),
				fmt.Sprintf("d=%d*%d", v.duration, v.config.Fps),
				"x=iw/2-(iw/zoom/2)",
			})

	return v
}

func NewZoomPanVideoFromImage(path string, duration int, zoom float32, config Config) *Video {
	if zoom < 1 || duration < 0 {
		panic("duration can't be less than 0 and zoom can't be less than 1")
	}

	v := &Video{
		config: config,
	}
	v.stream = ffmpeg.Input(path).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d", v.config.Width*4), fmt.Sprintf("%d", v.config.Height*4)}).
		Filter("zoompan",
			ffmpeg.Args{
				fmt.Sprintf("z=min(max(pzoom,zoom) + 0.001,%f)", zoom),
				fmt.Sprintf("fps=%d", v.config.Fps),
				fmt.Sprintf("d=%d*%d", duration, v.config.Fps),
				"x=iw/2-(iw/zoom/2)",
			})
	v.Crop(v.config.Width, v.config.Height)
	v.duration = duration

	return v
}

func (v *Video) AddSubtitles(path string) *Video {
	v.stream = v.stream.Video().Filter("subtitles", ffmpeg.Args{path})
	return v
}

func (v *Video) AddOverlay(overlay *Video, x, y string) *Video {
	v.stream = v.stream.
		Overlay(overlay.stream,
			"pass",
			ffmpeg.KwArgs{
				"x": x,
				"y": y,
			})

	return v
}

type Audio struct {
	stream *ffmpeg.Stream
}

func NewAudio(path string) *Audio {
	return &Audio{
		stream: ffmpeg.Input(path),
	}
}
