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

	flag.StringVar(&sourcePath, "source", "", "Path to the source image to diff against.")
	flag.StringVar(&targetPath, "target", "", "Path to the target image to diff against the source image.")

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
	quantizer := quantize.MedianCutQuantizer{}

	sourcePalette := quantizer.Quantize(make([]color.Color, 0, 256), sourceImage)
	targetPalette := quantizer.Quantize(make([]color.Color, 0, 256), targetImage)

	// Use targetPalette for diff since we render pixels from the target image only in the gif
	diff := image.NewPaletted(bounds, targetPalette)
	imagesDiffer := false
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			sourcePixel := sourceImage.At(x, y)
			targetPixel := targetImage.At(x, y)
			if sourcePixel != targetPixel {
				imagesDiffer = true
				for dy := -10; dy <= 10; dy++ {
					for dx := -10; dx <= 10; dx++ {
						diffColor := targetImage.At(dx+x, dy+y)
						diff.Set(dx+x, dy+y, diffColor)
					}
				}
				diff.Set(x, y, targetPixel)
			}
		}
	}
	if !imagesDiffer {
		log.Println("Images are identical")
		os.Exit(0)
	}

	// TODO: 1. Initialize a gif
	// 2. Convert source and target to Paletted
	// 3. Add source and target to the gif
	// 4. Create a blank Paletted and add differing pixels to it
	// 5. Add the diff to the gif
	// 6. Write the gif out

	sourcePaletted := image.NewPaletted(bounds, sourcePalette)
	draw.Draw(sourcePaletted, sourcePaletted.Rect, sourceImage, bounds.Min, draw.Over)

	targetPaletted := image.NewPaletted(bounds, targetPalette)
	draw.Draw(targetPaletted, targetPaletted.Rect, targetImage, bounds.Min, draw.Over)

	diffGif := gif.GIF{
		Image: []*image.Paletted{sourcePaletted, targetPaletted, diff},
		Delay: []int{100, 100, 100},
	}
	outputFile, err := os.Create("output.gif")
	if err != nil {
		log.Fatal(err.Error())
	}
	gif.EncodeAll(outputFile, &diffGif)
}
