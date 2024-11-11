# Image Diff Tool

This is a simple command-line tool written in Go to compare two images and generate an animated GIF highlighting the differences.

### Features
- Compares two images pixel-by-pixel to identify differences.
- Saves a 3-frame animated GIF showing the original, target, and diff images (with a 20x20 pixel buffer around any differing pixels).

### Usage
To run the tool, specify the paths to the source and target images. Optionally specify the output path.
```bash
./imgdiff --source path/to/source.png --target path/to/target.png [--output output.gif]
```

#### Arguments
- `--source`: Path to the source image to compare.
- `--target`: Path to the target image to compare against the source.
- `--output` (optional): Output path for the generated GIF (defaults to `output.gif`).

### Notes
- Both images must have the same dimensions and be no larger than 4000x4000 pixels.
- Currently, only `.png`, `.jpeg`, and `.gif` formats are supported for input images.

