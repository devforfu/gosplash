package imutil

import (
    "image"
    "os"
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

func ReadImages(filenames []string) (images []image.Image, err error) {
    for _, filename := range filenames {
        file, err := os.Open(filename)
        if err != nil { return }
        img, _, err := image.Decode(file)
        if err != nil { return }
        images = append(images, img)
    }
    return images, nil
}