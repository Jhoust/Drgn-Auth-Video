
// What it does:
//
// This example uses the CascadeClassifier class to detect faces,
// and draw a rectangle around each of them, before displaying them within a Window.
//
// How to run:
//
//  facedetect [videofile] [classifier XML file]
//	use this command----->>>	go run main.go 0 haarcascade_frontalface_default.xml
//
//

package main

import (
"fmt"
"image"
"image/color"
"os"

"gocv.io/x/gocv"
)

 func main() {
	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tfacedetect [videofile] [classifier XML file]")
		return
	}

	// parse args
	videofile := os.Args[1]
	xmlFile := os.Args[2]

	// open webcam
	video, err := gocv.VideoCaptureFile("drgnvideo.mp4")
	if err != nil {
		fmt.Printf("error opening video capture device: %v\n", videofile)
		return
	}
	defer video.Close()

	// open display window
	window := gocv.NewWindow("Joe Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("Start reading device: %v\n", videofile)
	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", videofile)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifing as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		for _, r := range rects {
			imgFace := img.Region(r)
			defer imgFace.Close()

			//blur face
			gocv.GaussianBlur(imgFace, &imgFace, image.Pt(75, 75), 0, 0, gocv.BorderDefault)//this code needs to change jared
		}


		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(24) >= 0 {
			break
		}
	}

}