package main

import (
	"flag"
	"image"
	"log"
	"os"

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

	sourceImage, _, err := image.Decode(sourceFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	targetImage, _, err := image.Decode(targetFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(sourceImage, targetImage)
}
