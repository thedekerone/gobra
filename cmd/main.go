package main

import (
	"fmt"
	"math/rand"

	"github.com/thedekerone/gobra/internal/video"
)

func main() {

	video2 := video.ImageToVideo("assets/story1.jpg", 4)

	video1 := video.ImageToVideo("assets/story2.jpg", 4)

	video3 := video.ImageToVideo("assets/story3.jpg", 4)

	video4 := video.ImageToVideo("assets/story4.jpg", 4)

	videos := []*video.Video{&video1, &video2, &video3, &video4}

	for i, _ := range videos {
		videos[i] = videos[i].Crop(900, 900)
	}

	custom := video.MergeVideos(&video1, &video2, &video3, &video4)

	custom.SaveVideo(fmt.Sprintf("test%d.mp4", rand.Intn(100)))
}
