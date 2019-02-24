package imutil

import (
    "fmt"
    "image"
    "math"
)

const autoSize = -1

type ImageGrid struct {
    NumRows int
    NumCols int
    Format ImageFormat
}

func NewImageGrid(format ImageFormat) *ImageGrid {
    return &ImageGrid{autoSize, autoSize, format}
}

func (g *ImageGrid) SetRows(n int) {
    g.NumRows = n
    g.NumCols = autoSize
}

func (g *ImageGrid) SetCols(n int) {
    g.NumRows = autoSize
    g.NumCols = n
}

func (g *ImageGrid) CreateFromArray(imageFiles []string) (image.Image, error) {
    var nRows, nCols int
    n := len(imageFiles)

    if g.NumRows != autoSize && g.NumCols != autoSize && n != (g.NumRows * g.NumCols) {
        return nil, fmt.Errorf("incomatible layout: nrows=%d, ncols=%d", g.NumRows, g.NumCols)
    }

    if g.NumRows == autoSize && g.NumCols == autoSize {
        nCols = int(math.Sqrt(float64(n)))
        nRows = int(math.Ceil(float64(n) / float64(nCols)))
    } else if g.NumRows == autoSize {
        nCols = int(math.Ceil(float64(n) / float64(g.NumRows)))
        nRows = g.NumRows
    } else if g.NumRows == autoSize {
        nCols = g.NumCols
        nRows = int(math.Ceil(float64(n) / float64(g.NumCols)))
    }

    images, err := ReadImages(imageFiles)
    if err != nil { return nil, err }

    _ := CreateCanvas(images, nRows, nCols)
    return nil, fmt.Errorf("not implemented")
}

func CreateCanvas(images []image.Image, nRows, nCols int) image.Image {
    var maxWidth, currWidth int
    for i := 0; i < nRows; i++ {
        for j := 0; j < nCols; j++ {
            currWidth += images[i*nRows + j].Bounds().Size().X
        }
        maxWidth = maxInt(maxWidth, currWidth)
        currWidth = 0
    }

    var maxHeight, currHeight int
    for j := 0; j < nCols; j++ {
        for i := 0; i < nRows; i++ {
            currHeight += images[i*nCols + j].Bounds().Size().Y
        }
        maxHeight = maxInt(maxHeight, currHeight)
        currHeight = 0
    }

    return image.NewRGBA(image.Rect(0,0, maxWidth, maxHeight))
}

func maxInt(a, b int) int {
    if a > b { return a }
    return b
}