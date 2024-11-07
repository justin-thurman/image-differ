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
	var sourcePath string
	var targetPath string

	flag.StringVar(&sourcePath, "source", "", "Path to the source image to diff against.")
	flag.StringVar(&targetPath, "target", "", "Path to the target image to diff against the source image.")

	flag.Parse()

	if sourcePath == "" || targetPath == "" {
		log.Fatal("source and target are required")
	}

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
