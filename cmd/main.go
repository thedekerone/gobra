package main

import (
	"fmt"
	"math/rand"

	"github.com/thedekerone/gobra/internal/video"
)

func main() {
	video2 := video.ImageToVideo("assets/jojos.jpg", 2, 400, 400)

	video1 := video.ImageToVideo("assets/jojos.jpg", 3, 400, 400)

	video1.VFlip()

	custom := video.MergeVideos(&video1, &video2)

	custom.SaveVideo(fmt.Sprintf("test%d.mp4", rand.Intn(100)))
}
