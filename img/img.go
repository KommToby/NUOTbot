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

// tournament images
var tournamentToImageMap = map[int]string{
	1:  "NUOT1.jpg",
	2:  "NUOT2.jpg",
	3:  "NUOT3B1.jpg",
	4:  "NUOT3B2.jpg",
	5:  "NUOT3B3.jpg",
	6:  "NUOT3B4.jpg",
	7:  "NUOT4.jpg",
	8:  "NUOT5.jpg",
	9:  "NUOT6.jpg",
	10: "NUOT7.jpg",
	12: "NUOT9.jpg",
	13: "NUOT91.jpg",
	14: "NUOT92.jpg",
	16: "NUOT93.jpg",
	17: "NUOT94B1.jpg",
	18: "NUOT94B2.jpg",
	19: "NUOT94B3.jpg",
	20: "NUOT95.jpg",
}

func LoadImagesForTournaments(tournamentIDs []int) ([]image.Image, error) {
	var images []image.Image

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error determining current directory")
		return nil, fmt.Errorf("error determining current directory")
	}
	currentDir := filepath.Dir(currentFile)
	fmt.Println("Current Directory:", currentDir)

	for _, tID := range tournamentIDs {
		filename, exists := tournamentToImageMap[tID]
		if !exists {
			fmt.Println("Tournament ID does not have associated image:", tID)
			continue
		}

		file := filepath.Join(currentDir, filename)

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
	numImages := len(images)
	rows := (numImages + 4) / 5
	totalWidth := 0
	if numImages < 5 {
		totalWidth = gap*numImages + sumWidths(images)
	} else {
		totalWidth = 5*gap + sumWidths(images[0:5])
	}

	// If there are less than 5 images, set the total width to 960 pixels
	if numImages < 5 {
		totalWidth = 960
	}

	maxHeightSingleRow := maxImageHeight(images)
	totalHeight := rows*maxHeightSingleRow + (rows-1)*gap // Total banner height

	banner := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Set banner to be transparent
	draw.Draw(banner, banner.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	x, y := 0, 0
	counter := 0 // To keep track of number of images in the current row

	for _, img := range images {
		if counter >= 5 { // If we have drawn 5 images on the current row, reset x and increase y
			x = 0
			y += maxHeightSingleRow + gap
			counter = 0
		}

		r := image.Rect(x, y, x+img.Bounds().Dx(), y+img.Bounds().Dy())

		// Resize the image if its height is not equal to maxHeightSingleRow
		if img.Bounds().Dy() != maxHeightSingleRow {
			img = resize.Resize(0, uint(maxHeightSingleRow), img, resize.Lanczos3)
		}

		draw.Draw(banner, r, img, image.Point{}, draw.Over)
		x += img.Bounds().Dx() + gap
		counter++
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
