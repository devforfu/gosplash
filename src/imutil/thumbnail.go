package imutil

import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "os"
    "strings"
)

// ThumbnailMaker converts an image or a set of images into thumbnails.
type ThumbnailMaker struct {
    // Format defines generated thumbnail format.
    Format ImageFormat
    // Width and Height define thumbnail maximal sizes.
    Width, Height int
}

func NewThumbnailMaker(format ImageFormat, width, height int) *ThumbnailMaker {
    return &ThumbnailMaker{format, width, height}
}

// Create takes infile image, converts it into thumbnail and saves into outfile.
func (t *ThumbnailMaker) Create(infile, outfile string) (err error) {
    in, err := os.Open(infile)
    if err != nil { return }
    defer in.Close()

    out, err := os.Create(outfile)
    if err != nil { return }

    if err := t.Convert(out, in); err != nil {
        out.Close()
        return fmt.Errorf("scaling %s to %s: %s", infile, outfile, err)
    }

    return out.Close()
}

// Convert reads image from r and writes thumbnail into w.
func (t *ThumbnailMaker) Convert(w io.Writer, r io.Reader) (err error) {
    src, _, err := image.Decode(r)
    if err != nil { return }
    dst := CreateThumbnailImage(src, t.Width, t.Height)
    switch t.Format {
    case JPEG:
        return jpeg.Encode(w, dst, nil)
    case PNG:
        return png.Encode(w, dst)
    default:
        return fmt.Errorf("unknown image format")
    }
}

// ConvertFolder traverses the dirname folder and converts all images discovered
// images into thumbnails. The patterns string define which image extensions to
// look for and should be a pipe-separated string, like "jpeg|png".
func (t *ThumbnailMaker) ConvertFolder(dirname string, patterns string) {

}

// CreateThumbnailImage rescales the src image into thumbnail with maximal
// dimensions defined by width and height.
func CreateThumbnailImage(src image.Image, width, height int) image.Image {
   xs := src.Bounds().Size().X
   ys := src.Bounds().Size().Y

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

// SplitPattern converts pipe-separated pattern into array of strings.
// For example, the string "jpeg|png" is converted into array ["jpeg", "png"].
func SplitPattern(pattern string) []string {
    result := make([]string, 0)
    for _, p := range strings.Split(pattern, "|") {
        for _, ext := range []string{strings.ToUpper(p), strings.ToLower(p)} {
            result = append(result, ext)
        }
    }
    return result
}

