package imutil

import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "os"
    "path/filepath"
    "strings"
)

type FileFormat uint8

const (
    Unknown FileFormat = iota
    JPEG
    PNG
)

func Format(ext string) FileFormat {
    switch ext {
    case ".jpg", ".jpeg":
        return JPEG
    case ".png":
        return PNG
    default:
        return Unknown
    }
}

type ThumbnailMaker struct {
    In, Out string
    Format FileFormat
}

func NewThumbnailMaker(infile string, format FileFormat) *ThumbnailMaker {
    ext := filepath.Ext(infile)
    outfile := strings.TrimSuffix(infile, ext) + ".thumb" + ext
    return &ThumbnailMaker{In:infile, Out:outfile, Format:format}
}

func (t *ThumbnailMaker) Create(format FileFormat) (err error) {
    in, err := os.Open(t.In)
    if err != nil { return }
    defer in.Close()

    out, err := os.Create(t.Out)
    if err != nil { return }

    if err := t.Convert(out, in); err != nil {
        out.Close()
        return fmt.Errorf("scaling %s to %s: %s", t.In, t.Out, err)
    }

    return out.Close()
}

func (t *ThumbnailMaker) Convert(w io.Writer, r io.Reader) error {
    src, _, err := image.Decode(r)
    if err != nil { return err }
    dst := CreateThumbnailImage(src)
    switch t.Format{
    case JPEG:
        return jpeg.Encode(w, dst, nil)
    case PNG:
        return png.Encode(w, dst)
    default:
        return fmt.Errorf("unknown image format")
    }
}

func CreateThumbnailImage(src image.Image) image.Image {
    xs := src.Bounds().Size().X
    ys := src.Bounds().Size().Y

    width, height := 128, 128
    if aspect := float64(xs)/float64(ys); aspect < 1.0 {
        width = int(128 * aspect)
    } else {
        height = int(128 / aspect)
    }

    xScale := float64(xs)/float64(width)
    yScale := float64(ys)/float64(height)
    dst := image.NewRGBA(image.Rect(0, 0, width, height))

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            xSrc := int(float64(x)*xScale)
            ySrc := int(float64(y)*yScale)
            dst.Set(x, y, src.At(xSrc, ySrc))
        }
    }

    return dst
}