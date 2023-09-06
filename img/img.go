package img

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nfnt/resize"
)

const gap = 10 // The gap between images

func LoadImages() ([]image.Image, error) {
	var images []image.Image

	// Determine the directory of this file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error determining current directory")
		return nil, fmt.Errorf("error determining current directory")
	}
	currentDir := filepath.Dir(currentFile)
	fmt.Println("Current Directory:", currentDir)

	files, err := filepath.Glob(filepath.Join(currentDir, "*.jpg"))
	if err != nil {
		fmt.Println("Error during filepath.Glob:", err)
		return nil, err
	}

	fmt.Println("Found Files:", files)

	for _, file := range files {
		infile, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		defer infile.Close()

		var img image.Image

		switch filepath.Ext(file) {
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(infile)
			if err != nil {
				fmt.Println("Error decoding jpeg:", err)
				return nil, err
			}
		case ".png":
			img, err = png.Decode(infile)
			if err != nil {
				fmt.Println("Error decoding png:", err)
				return nil, err
			}
		default:
			err = fmt.Errorf("unsupported file type: %s", file)
			fmt.Println(err)
			return nil, err
		}

		images = append(images, img)
	}

	return images, nil
}

func CreateBanner(images []image.Image) (image.Image, error) {
	totalWidth := (len(images)-1)*gap + sumWidths(images)
	maxHeight := maxImageHeight(images)

	banner := image.NewRGBA(image.Rect(0, 0, totalWidth, maxHeight))

	// Set banner to be transparent
	draw.Draw(banner, banner.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	x := 0
	for _, img := range images {
		r := image.Rect(x, 0, x+img.Bounds().Dx(), maxHeight)

		// Resize the image if its height is not equal to maxHeight
		if img.Bounds().Dy() != maxHeight {
			img = resize.Resize(0, uint(maxHeight), img, resize.Lanczos3)
		}

		draw.Draw(banner, r, img, image.Point{}, draw.Over)
		x += img.Bounds().Dx() + gap
	}

	return banner, nil
}

func sumWidths(images []image.Image) int {
	total := 0
	for _, img := range images {
		total += img.Bounds().Dx()
	}
	return total
}

func maxImageHeight(images []image.Image) int {
	maxHeight := 0
	for _, img := range images {
		if h := img.Bounds().Dy(); h > maxHeight {
			maxHeight = h
		}
	}
	return maxHeight
}
