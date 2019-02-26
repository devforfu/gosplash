package main

import (
    "./src/cli"
    "./src/imutil"
    "./src/unsplash"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
    "strconv"
)

func init() {
    log.SetFlags(0)
}

func main() {
    if len(os.Args) == 1 {
        log.Fatal("usage: gosplash <command> [<args>]")
    }

    params := cli.Parse(os.Args)
    switch cmd := os.Args[1]; cmd {
    case "download":
        DownloadImages(params)
    case "thumb":
        MakeThumbnails(params)
    case "canvas":
        ConcatIntoSingle(params)
    }
}

func DownloadImages(params cli.Params) {
    log.Println("Downloading random images...")
    jsonConfig := MustLoadJSON(params["conf"].String())
    client := unsplash.Client{
        AccessKey:jsonConfig["accessKey"],
        SecretKey:jsonConfig["secretKey"]}
    err := client.DownloadRandomPhotos(params["out"].String(), params["n"].Integer())
    if err != nil { log.Fatal(err) }
    log.Printf("Images downlaoded into folder: %s\n", params["out"].String())
}

func MakeThumbnails(params cli.Params) {
    log.Println("Converting images into thumbnails...")
    maker := imutil.NewThumbnailMaker(imutil.PNG, 128, 128)
    thumbs, err := maker.ConvertFolder(params["dir"].String(), params["p"].String())
    if err != nil { log.Fatal(err) }
    log.Println("Created thumbnails:")
    n := len(strconv.Itoa(len(thumbs)))
    for i, t := range thumbs { log.Printf("\t%0*d: %s", n, i, t) }
}

func ConcatIntoSingle(params cli.Params) {
    log.Fatal("Not implemented!")
}

func MustLoadJSON(filename string) map[string]string {
    data, err := ioutil.ReadFile(filename)
    if err != nil { log.Fatal(err) }
    config := make(map[string]string)
    err = json.Unmarshal(data, &config)
    if err != nil { log.Fatal(err) }
    return config
}

//    switch cmd := os.Args[1]; cmd {
//    case "download":
//        log.Println("Downloading random images...")
//        params = ParseDownload()
//        client := unsplash.Client{
//            AccessKey:params["accessKey"].(string),
//            SecretKey:params["secretKey"].(string)}
//        output := params["output"].(string)
//        err := client.DownloadRandomPhotos(output, params["count"].(int))
//        log.Printf("Images downloaded into folder: %s\n", output)
//        if err != nil { log.Fatal(err) }
//    case "thumb":
//        log.Fatal("Converting images into thumbnails...")
//        params = ParseThumbnails()
//        maker := imutil.NewThumbnailMaker(imutil.PNG, 128, 128)
//        output := params["output"].(string)
//        thumbs, err := maker.ConvertFolder(output, params["p"].(string))
//        if err != nil { log.Fatal(err) }
//        log.Println("Created thumbnails:")
//        n := len(strconv.Itoa(len(thumbs)))
//        for i, t := range thumbs {
//            log.Printf("\t%0*d: %s", n, i, t)
//        }
//    case "canvas":
//        log.Fatal("Making a single image from list of thumbnails...")
//        // params := ParseCanvas()
//    }
//}

//type Params map[string]interface{}
//
//func ParseDownload() Params {
//    confPath := flag.String(
//        "-conf",
//        "unsplash.key.json",
//        "path to the file with Unsplash API keys")
//
//    numOfImages := flag.Int(
//        "-n",
//        5,
//        "number of images to download")
//
//    outputPath := flag.String(
//        "-out",
//        path.Join(mustHomeDir(), "Unsplash"),
//        "path to the output folder")
//
//    flag.Parse()
//    data, err := ioutil.ReadFile(*confPath)
//    if err != nil { log.Fatal(err) }
//
//    params := make(Params)
//    err = json.Unmarshal(data, &params)
//    if err != nil { log.Fatal(err) }
//    params["count"] = *numOfImages
//    params["output"] = *outputPath
//    return params
//}
//
//func ParseThumbnails() Params {
//    dirname := flag.String(
//        "-dir",
//        mustWorkDir(),
//        "path to the folder with images")
//
//    pattern := flag.String(
//        "-p",
//        imutil.ImageFormats,
//        "image extensions pattern")
//
//    flag.Parse()
//    params := make(Params)
//    params["dirname"] = *dirname
//    params["pattern"] = *pattern
//    return params
//}
//
//func ParseCanvas() Params {
//    dirname := flag.String("-dir", mustWorkDir(), "path to folder with images")
//    flag.Parse()
//    params := make(Params)
//    params["dirname"] = *dirname
//    return params
//}
//
//func mustHomeDir() string {
//    home, err := homedir.Dir()
//    if err != nil { log.Fatal(err) }
//    return home
//}
//
//func mustWorkDir() string {
//    workdir, err := os.Getwd()
//    if err != nil { log.Fatal(err) }
//    return workdir
//}