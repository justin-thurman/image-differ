package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"os"

	"github.com/ericpauley/go-quantize/quantize"

	_ "image/jpeg"
	_ "image/png"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parsing CLI args
	var sourcePath string
	var targetPath string
	var outputPath string

	flag.StringVar(&sourcePath, "source", "", "Path to the source image to diff against.")
	flag.StringVar(&targetPath, "target", "", "Path to the target image to diff against the source image.")
	flag.StringVar(&outputPath, "output", "output.gif", "Path to save the output gif (defaults to ./output.gif)")

	flag.Parse()

	if sourcePath == "" || targetPath == "" {
		log.Fatal("source and target are required")
	}

	// Loading image files
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer sourceFile.Close()

	targetFile, err := os.Open(targetPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer targetFile.Close()

	// Decoding images
	sourceConfig, _, err := image.DecodeConfig(sourceFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = sourceFile.Seek(0, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	targetConfig, _, err := image.DecodeConfig(targetFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = targetFile.Seek(0, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	if sourceConfig.Width > 4000 || sourceConfig.Height > 4000 {
		log.Fatal("source image too large; ensure image height and width <= 4000")
	}
	if targetConfig.Width > 4000 || targetConfig.Height > 4000 {
		log.Fatal("target image too large; ensure image height and width <= 4000")
	}
	if sourceConfig.Width != targetConfig.Width || sourceConfig.Height != targetConfig.Height {
		// TODO: A way to handle diffs of this kind would be nice. Perhaps a method of normalizing the 2D
		// coordinate plane? But such normalization would be lossy in cases where the number of pixels do
		// not evenly divide into dimensions of the normalized plane. Would need a way to handle that.
		log.Fatal("source and target images must have identical dimensions")
	}

	sourceImage, _, err := image.Decode(sourceFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	targetImage, _, err := image.Decode(targetFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Compare images
	bounds := sourceImage.Bounds()
	differingPixels := make(map[image.Point]struct{})

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sourcePixel := sourceImage.At(x, y)
			targetPixel := targetImage.At(x, y)
			if sourcePixel != targetPixel {
				// Draw a 20x20 square around each differing pixel
				for dy := -10; dy <= 10; dy++ {
					for dx := -10; dx <= 10; dx++ {
						differingPixels[image.Point{dx + x, dy + y}] = struct{}{}
					}
				}
			}
		}
	}
	if len(differingPixels) == 0 {
		log.Println("Images are identical")
		os.Exit(0)
	}

	// Create the gif showing the diff
	quantizer := quantize.MedianCutQuantizer{}
	sourcePalette := quantizer.Quantize(make([]color.Color, 0, 256), sourceImage)
	targetPalette := quantizer.Quantize(make([]color.Color, 0, 256), targetImage)

	// Use targetPalette for diff since we render pixels from the target image only in the gif
	diff := image.NewPaletted(bounds, targetPalette)

	for point := range differingPixels {
		diffColor := targetImage.At(point.X, point.Y)
		diff.Set(point.X, point.Y, diffColor)
	}

	sourcePaletted := image.NewPaletted(bounds, sourcePalette)
	draw.Draw(sourcePaletted, sourcePaletted.Rect, sourceImage, bounds.Min, draw.Over)

	targetPaletted := image.NewPaletted(bounds, targetPalette)
	draw.Draw(targetPaletted, targetPaletted.Rect, targetImage, bounds.Min, draw.Over)

	diffGif := gif.GIF{
		Image: []*image.Paletted{sourcePaletted, targetPaletted, diff},
		Delay: []int{100, 100, 100},
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	gif.EncodeAll(outputFile, &diffGif)
}
