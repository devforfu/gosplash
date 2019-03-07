package main

//import (
//    "./src/cli"
//    "./src/imutil"
//    "./src/unsplash"
//    "log"
//    "os"
//    "strconv"
//)
//
//func init() {
//    log.SetFlags(0)
//}
//
//func main() {
//    if len(os.Args) == 1 {
//       log.Fatal("usage: gosplash <command> [<args>]")
//    }
//    params := cli.Parse(os.Args)
//    switch cmd := os.Args[1]; cmd {
//    case "download":
//        DownloadImages(params)
//    case "thumb":
//        MakeThumbnails(params)
//    case "canvas":
//        ConcatIntoSingle(params)
//    }
//}

//func DownloadImages(params cli.Params) {
//    log.Println("Downloading random images...")
//    jsonConfig := MustLoadJSON(params["conf"].String())
//    client := unsplash.Client{
//        AccessKey:jsonConfig["accessKey"],
//        SecretKey:jsonConfig["secretKey"]}
//    err := client.DownloadRandomPhotos(params["out"].String(), params["n"].Integer())
//    if err != nil { log.Fatal(err) }
//    log.Printf("Images downlaoded into folder: %s\n", params["out"].String())
//}

//func MakeThumbnails(params cli.Params) {
//    log.Println("Converting images into thumbnails...")
//    maker := imutil.NewThumbnailMaker(imutil.PNG, 128, 128)
//    thumbs, err := maker.ConvertFolder(params["dir"].String(), params["p"].String())
//    if err != nil { log.Fatal(err) }
//    log.Println("Created thumbnails:")
//    n := len(strconv.Itoa(len(thumbs)))
//    for i, t := range thumbs { log.Printf("\t%0*d: %s", n, i, t) }
//}

//func ConcatIntoSingle(params cli.Params) {
//    log.Fatal("Not implemented!")
//}

//func MustLoadJSON(filename string) map[string]string {
//    data, err := ioutil.ReadFile(filename)
//    if err != nil { log.Fatal(err) }
//    config := make(map[string]string)
//    err = json.Unmarshal(data, &config)
//    if err != nil { log.Fatal(err) }
//    return config
//}
