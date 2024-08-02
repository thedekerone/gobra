package main

import (
	"fmt"

	"github.com/thedekerone/gobra/internal/video"
)

func main() {

	video2 := video.ImageToVideo("assets/story1.jpg", 7, 30)

	video1 := video.ImageToVideo("assets/story2.jpg", 7, 30)

	video3 := video.ImageToVideo("assets/story3.jpg", 7, 30)

	video4 := video.ImageToVideo("assets/story4.jpg", 7, 30)

	videos := []*video.Video{&video1, &video2, &video3, &video4}

	for i, _ := range videos {
		videos[i] = videos[i].Crop(900, 900)
		videos[i] = videos[i].AddFadeIn(1)
		videos[i] = videos[i].AddFadeOut(1)
	}

	custom := video.MergeVideos(&video1, &video2, &video3, &video4)

	custom.SaveVideo(fmt.Sprintf("test.mp4"))
}
