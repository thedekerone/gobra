package audio

import ffmpeg_go "github.com/u2takey/ffmpeg-go"

type Audio struct {
	stream *ffmpeg_go.Stream
}

func ReadAudio(path string) Audio {
	c := Audio{}
	c.stream = ffmpeg_go.Input(path, ffmpeg_go.KwArgs{"a": 1, "v": 1})

	return c
}
