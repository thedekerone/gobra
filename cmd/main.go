package main

import (
	"fmt"

	"github.com/thedekerone/gobra/internal/video"
)

func main() {
	testVid()
}

func testVid() {
	config := video.Config{
		Width:       1080,
		Height:      1920,
		Fps:         60,
		AspectRatio: 1.25,
	}
	video2 := video.CreateZoomPanVideoFromImage("assets/story1.jpg", 4, 1.25, config)

	video1 := video.CreateZoomPanVideoFromImage("assets/story2.jpg", 4, 1.25, config)

	video3 := video.CreateZoomPanVideoFromImage("assets/story3.jpg", 4, 1.25, config)

	video4 := video.CreateZoomPanVideoFromImage("assets/story4.jpg", 4, 1.25, config)

	videos := []*video.Video{&video1, &video2, &video3, &video4}

	for i, _ := range videos {
		videos[i] = videos[i].Crop(900, 1200)
		videos[i] = videos[i].AddFadeIn(0.5)
		videos[i] = videos[i].AddFadeOut(0.5)
	}

	custom := video.MergeVideos(&video1, &video2, &video3, &video4)

	custom.AddSubtitles("assets/test_sub.srt")

	custom.SaveVideo(fmt.Sprintf("test.mp4"))

}
