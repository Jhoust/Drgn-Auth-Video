// What it does:
//
// This example uses the CascadeClassifier class to detect faces,
// and draw a rectangle around each of them, before displaying them within a Window.
//
// How to run:
//
// facedetect [camera ID] [classifier XML file]
//
// USE THIS COMMAND--->>>>   go run main.go 0 haarcascade_frontalface_default.xml tensorflow_inception_graph.pb imagenet_comp_graph_label_strings.txt opencv cpu

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"

	"gocv.io/x/gocv"
)

 go func main() {
	if len(os.Args) < 0 {
		fmt.Println("How to run:\n\tfacedetect [videofile] [classifier XML file] [modelfile] [descriptionsfile]")
		return
	}

	// parse args
	videofile := os.Args[1]
	xmlFile := os.Args[2]
	model := os.Args[3]
	descr := os.Args[4]
	//--------------------------------------------------------------------------------------------
	descriptions, err := readDescriptions(descr)
	if err != nil {
		fmt.Printf("Error reading descriptions file: %v\n", descr)
		return
	}
	//--------------------------------------------------------------------------------------------

	// NetBackendType is the type for the various different kinds of DNN backends.
	backend := gocv.NetBackendDefault
	if len(os.Args) > 4 {
		backend = gocv.ParseNetBackend(os.Args[4])
	}
	//-----------------------------------------------------------------------------------------

	//NetTargetCPU is the default CPU device target.
	target := gocv.NetTargetCPU
	if len(os.Args) > 6 {
		target = gocv.ParseNetTarget(os.Args[6])
	}

	//-------------------------------------------------------------------------------------------
	// open video
	video, err := gocv.VideoCaptureFile("drgnvideo.mp4")
	if err != nil {
		fmt.Printf("error opening video capture device: %v\n", videofile)
		return
	}
	defer video.Close()
	//---------------------------------------------------------------------------------------
	// open display window
	window := gocv.NewWindow("TF Joe Face Detect")
	defer window.Close()
	//------------------------------------------------------------------------------------
	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()
	//-----------------------------------------------------------------------------------
	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier :=  gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}
	//---------------------------------------------------------------------------------
//***********************************************************************
	 //open DNN classifier
	 net := gocv.ReadNet(model, "")
	 if net.Empty() {
		 fmt.Printf("Error reading network model : %v\n", model)
		 return
	 }
	 defer net.Close()

	 net.SetPreferableBackend(gocv.NetBackendType(backend))
	 net.SetPreferableTarget(gocv.NetTargetType(target))

	 //********************************************

	//----------------------------------------------------------------------------------
	status := "Ready"


	// read videofile
	fmt.Printf("Start reading video: %v\n", videofile)

	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("video closed: %v\n", videofile)
			return
		}
		if img.Empty() {
			continue
		}

		//-------------------------------------------------------------------------------------------
		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text "status". status is the tf desrc output
		for _, r := range rects {
			//if v != 0{
			//	panic(rects)
			//}
			gocv.Rectangle(&img, r, blue, 3)

		// these value determine the position the text appears above the rectangle relative to the size (r.Min.X/9) I changed the value 9 from 3
			size := gocv.GetTextSize(status, gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/9)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, status, pt, gocv.FontHersheyPlain, 1.2, blue, 2)


		}
		//--------------------------------------------------------------------------------------------------
		for _, r := range rects {
			imgFace := img.Region(r)
			//if v != 0{
				// imgFace.Close()
				//panic(imgFace)
			//}
			defer imgFace.Close()



			//---------------------------------------------------------------------------------------------
			// convert imgFace Mat to 224x224 blob that the classifier can analyze hopefully
			blob := gocv.BlobFromImage(imgFace, 1.0, image.Pt(224,224), gocv.NewScalar(0, 0, 0, 0), true, false)

			// feed the blob into the classifier
			net.SetInput(blob, "input")

			// run a forward pass thru the network
			prob := net.Forward("softmax2")

			// reshape the results into a 1x1000 matrix
			probMat := prob.Reshape(1, 1)

			// determine the most probable classification
			_, maxVal, _, maxLoc := gocv.MinMaxLoc(probMat)

			// display classification
			desc := "Unknown"
			if maxLoc.X < 1000 {
				desc = descriptions[maxLoc.X]
			}

			status = fmt.Sprintf("Description: %v, maxVal: %v\n", desc, maxVal)

			blob.Close()
			prob.Close()
			probMat.Close()

			// show the image in the window, and wait 1 millisecond
			window.IMShow(img)
			if window.WaitKey(1) >= 0 {
				break
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------
go func readDescriptions(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

//------------------------------------------------------------------------------------------------
