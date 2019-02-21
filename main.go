package main

import (
    "./src/unsplash"
    "encoding/json"
    "flag"
    "github.com/mitchellh/go-homedir"
    "io/ioutil"
    "log"
    "path"
)

func init() {
    log.SetFlags(0)
}

func main()  {
    conf := parseArgs()
    client := unsplash.Client{AccessKey: conf["accessKey"], SecretKey: conf["secretKey"]}
    err := client.DownloadRandomPhotos(conf["output"], 5)
    if err != nil { log.Fatal(err) }
    log.Printf("Images downloaded into folder: %s", conf["output"])
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
