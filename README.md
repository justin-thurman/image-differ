# Description

A work-in-progress command-line image diffing utility. Current functionality:
1. If the two input images differ, logs to STDOUT the first pixel (from top left) that differs. Exit status: 1.
2. If the two input images do not differ, logs this to STDOUT. Exit status: 0.

# Usage
`image-differ --source path/to/base/image.png --target path/to/target/image.png`

# Limitations
- Images must be no more than 4000x4000 pixels
- Images must have identical dimensions
- Images must be in png, jpg, or jpeg format
