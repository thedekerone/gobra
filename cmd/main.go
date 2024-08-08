package main

import (
	"github.com/thedekerone/gobra/internal/video"
)

func main() {
	v := video.ReadVideo("assets/subway.mp4", video.Config{Width: 1920, Height: 1080, Fps: 30, AspectRatio: 16 / 9})

	v.TrimVideo(0, 300)

	v.SaveVideo("assets/subway_trim.mp4")
}

func testVid() {
	audio1 := video.ReadAudio("assets/subway.mp3")
	video1 := video.ReadVideo("assets/subway.mp4", video.Config{Width: 1920, Height: 1080, Fps: 30, AspectRatio: 16 / 9})
	video.MergeVideoWithAudio(&audio1, &video1)

}
