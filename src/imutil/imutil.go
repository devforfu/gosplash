package imutil

import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "log"
    "os"
    "path"
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

// ThumbnailMaker is a convenience wrapper to convert full-size image into smaller thumbnail.
type ThumbnailMaker struct {
    // In stores a path to the original image.
    In string
    // Out stores a path to the generated thumbnail.
    Out string
    // Width defines created thumbnail width.
    Width int
    // Height defines created thumbnail height.
    Height int
}

func NewThumbnailMaker(infile string) *ThumbnailMaker {
    ext := filepath.Ext(infile)
    outfile := strings.TrimSuffix(infile, ext) + ".thumb" + ext
    return &ThumbnailMaker{In:infile, Out:outfile, Width:128, Height:128}
}

// Create converts an image into smaller thumbnail and saves to the output file.
func (t *ThumbnailMaker) Create(format FileFormat) (err error) {
    in, err := os.Open(t.In)
    if err != nil { return }
    defer in.Close()

    out, err := os.Create(t.Out)
    if err != nil { return }

    if err := t.Convert(out, in, format); err != nil {
        out.Close()
        return fmt.Errorf("scaling %s to %s: %s", t.In, t.Out, err)
    }

    return out.Close()
}

func (t *ThumbnailMaker) Convert(w io.Writer, r io.Reader, format FileFormat) error {
    src, _, err := image.Decode(r)
    if err != nil { return err }
    dst := CreateThumbnailImage(src, t.Width, t.Height)
    switch format {
    case JPEG:
        return jpeg.Encode(w, dst, nil)
    case PNG:
        return png.Encode(w, dst)
    default:
        return fmt.Errorf("unknown image format")
    }
}

func CreateThumbnailImage(src image.Image, width, height int) image.Image {
    xs := src.Bounds().Size().X
    ys := src.Bounds().Size().Y

    //width, height := 128, 128
    if aspect := float64(xs)/float64(ys); aspect < 1.0 {
        width = int(float64(width) * aspect)
    } else {
        height = int(float64(height) / aspect)
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

// ThumbnailsFromFolder traverses a directory and converts all images with formats from pipe-separated
// patterns string into thumbnails.
func ThumbnailsFromFolder(dirname string, patterns string, format FileFormat) (thumbs []string, err error) {
    for _, pattern := range strings.Split(patterns, "|") {
        extensions:= []string{strings.ToLower(pattern), strings.ToUpper(pattern)}
        for _, ext := range extensions {
            glob := path.Join(dirname, fmt.Sprintf("*.%s", ext))
            files, err := filepath.Glob(glob)
            if err != nil {
                log.Printf("cannot discover *.%s images: %s", ext, err)
                continue
            }
            for _, file := range files {
                maker := NewThumbnailMaker(file)
                err := maker.Create(format)
                if err != nil {
                    log.Printf("cannot create thumbnail image from file: %s", file)
                    return nil, err
                }
                thumbs = append(thumbs, maker.Out)
            }
        }
    }
    return thumbs, nil
}

type ImageComposer struct {
    Images []image.Image
    Result string
}

func NewImageComposer(files []string) (composer *ImageComposer, err error) {
    return nil, fmt.Errorf("not implemented")
}