package main

import (
    "./src/unsplash"
    "./src/imutil"
    "encoding/json"
    "flag"
    "github.com/mitchellh/go-homedir"
    "io/ioutil"
    "log"
    "path"
    "strconv"
)

func init() {
    log.SetFlags(0)
}

func main()  {
    conf := parseArgs()

    client := unsplash.Client{AccessKey: conf["accessKey"], SecretKey: conf["secretKey"]}
    output := conf["output"]
    err := client.DownloadRandomPhotos(output, 5)
    if err != nil { log.Fatal(err) }

    log.Printf("Images downloaded into folder: %s\n", output)
    thumbs, err := imutil.ThumbnailsFromFolder(output, "jpeg|jpg|png", imutil.PNG)
    if err != nil { log.Fatal(err) }

    log.Println("Created thumbnails:")
    n := len(strconv.Itoa(len(thumbs)))
    for i, t := range thumbs {
        log.Printf("\t%0*d: %s", n, i, t)
    }
}

func parseArgs() map[string]string {
    confPath := flag.String(
        "-conf",
        "unsplash.key.json",
        "path to the file with Unsplash API keys")

    outputPath := flag.String(
        "-out",
        path.Join(mustHomeDir(), "Unsplash"),
        "path to the output folder")

    flag.Parse()
    data, err := ioutil.ReadFile(*confPath)
    if err != nil { log.Fatal(err) }

    config := make(map[string]string)
    err = json.Unmarshal(data, &config)
    if err != nil { log.Fatal(err) }

    config["output"] = *outputPath
    return config
}

func mustHomeDir() string {
    home, err := homedir.Dir()
    if err != nil { log.Fatal(err) }
    return home
}
