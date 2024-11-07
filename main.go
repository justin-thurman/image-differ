package main

import (
	"flag"
	"log"
	"os"
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

	targetFile, err := os.Open(targetPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(sourceFile, targetFile)
}
