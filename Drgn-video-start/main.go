package main


import (
"gocv.io/x/gocv"

)

func main() {
	videofile, _ := gocv.VideoCaptureFile("drgnvideo.mp4") // package VideoCaptureFile takes in a .mp4 .avi file I have not played with any other formats. These can also be changed into .jpegs for each frame
	window := gocv.NewWindow("Drgn-Auth-Video") // just init the new window with any name that the file will play in
	img := gocv.NewMat() // this creates a new image mat(matrix) a 2 dimension 3 channel array. (x,y) position with (BGR)color. 2 dim 1 channel, 3 dim 4 channels are also possible.

	for {
		videofile.Read(&img)
		window.IMShow(img)
		window.WaitKey(25) // default on this value is 1; (0) means forever; 30 is the frames per-second that the video was captured; by the program I used
	}
}