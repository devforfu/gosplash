package imutil

import (
    "fmt"
    "image"
    "os"
    "path"
    "path/filepath"
    "strings"
)

type ImageFormat uint8

const (
   Unknown ImageFormat = iota
   JPEG
   PNG
)

func Format(ext string) ImageFormat {
   switch ext {
   case ".jpg", ".jpeg":
       return JPEG
   case ".png":
       return PNG
   default:
       return Unknown
   }
}

const ImageFormats = "jpeg|jpg|png|bmp"

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

func ReadImages(filenames []string) ([]image.Image, error) {
    images := make([]image.Image, 0)
    for _, filename := range filenames {
        file, err := os.Open(filename)
        if err != nil { return nil, err }
        img, _, err := image.Decode(file)
        if err != nil { return nil, err }
        images = append(images, img)
    }
    return images, nil
}

func DiscoverImages(dirname, formats string) ([]string, error) {
    files := make([]string, 0)
    for _, ext := range SplitPattern(formats) {
        glob := path.Join(dirname, fmt.Sprintf("*.%s", ext))
        matched, err := filepath.Glob(glob)
        if err != nil { return nil, err }
        files = append(files, matched...)
    }
    return files, nil
}