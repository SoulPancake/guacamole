package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	//"os"
	"strings"

	"gocv.io/x/gocv"
)

// Download and decode image from URL
func downloadImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	var img image.Image
	switch {
	case strings.Contains(contentType, "png"):
		img, err = png.Decode(resp.Body)
	case strings.Contains(contentType, "jpeg"), strings.Contains(contentType, "jpg"):
		img, err = jpeg.Decode(resp.Body)
	default:
		err = fmt.Errorf("unsupported content type: %s", contentType)
	}
	return img, err
}

// Convert image.Image to gocv.Mat
func imageToMat(img image.Image) (gocv.Mat, error) {
	mat, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		return gocv.Mat{}, err
	}
	return mat, nil
}

// Extract QR corners from Mat returned by QRCodeDetector
func extractQRCorners(points gocv.Mat) ([]image.Point, error) {
	if points.Empty() || points.Rows() != 4 || points.Cols() < 2 {
		return nil, fmt.Errorf("unexpected QR corner shape: %dx%d", points.Rows(), points.Cols())
	}

	corners := make([]image.Point, 4)
	for i := 0; i < 4; i++ {
		x := points.GetFloatAt(i, 0)
		y := points.GetFloatAt(i, 1)
		corners[i] = image.Pt(int(x), int(y))
	}
	return corners, nil
}

func main() {
	inputPath := "input.jpg"
	outputPath := "output.jpg"
	qrURL := "qr.png" // ← your QR code image URL

	// Read input image
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		log.Fatalf("Failed to read image: %s", inputPath)
	}
	defer img.Close()

	// Detect QR
	detector := gocv.NewQRCodeDetector()
	defer detector.Close()

	points := gocv.NewMat()
	defer points.Close()

	found := detector.Detect(img, &points)
	if !found {
		log.Fatal("QR code not found")
	}

	corners, err := extractQRCorners(points)
	if err != nil {
		log.Fatalf("Failed to extract corners: %v", err)
	}

	// Load replacement QR code
	replacementMat := gocv.IMRead(qrURL, gocv.IMReadColor)
	if replacementMat.Empty() {
		log.Fatalf("Failed to read replacement QR image: %s", qrURL)
	}
	defer replacementMat.Close()

	// Wrap corners in PointVector
	dstVector := gocv.NewPointVectorFromPoints(corners)
	defer dstVector.Close()

	// Source points from the replacement QR image corners
	srcPoints := []image.Point{
		{0, 0},
		{replacementMat.Cols(), 0},
		{replacementMat.Cols(), replacementMat.Rows()},
		{0, replacementMat.Rows()},
	}
	srcVector := gocv.NewPointVectorFromPoints(srcPoints)
	defer srcVector.Close()

	// Get perspective transform using your wrapper
	homography := gocv.GetPerspectiveTransform(srcVector, dstVector)
	defer homography.Close()

	// Warp and overlay
	warped := gocv.NewMat()
	defer warped.Close()
	// convert img.Size to image.Point
	imgSize := image.Point{X: img.Cols(), Y: img.Rows()}
	gocv.WarpPerspective(replacementMat, &warped, homography, imgSize)

	// Masked paste (optional)
	mask := gocv.NewMatWithSize(warped.Rows(), warped.Cols(), gocv.MatTypeCV8UC1)
	defer mask.Close()
	gocv.CvtColor(warped, &mask, gocv.ColorBGRToGray)

	warped.CopyToWithMask(&img, mask)

	// Save final image
	if ok := gocv.IMWrite(outputPath, img); !ok {
		log.Fatalf("Failed to write output: %s", outputPath)
	}

	log.Println("✅ Replacement done:", outputPath)
}
