package unsplash

import (
    "bytes"
    "encoding/json"
    "fmt"
    "image"
    "image/jpeg"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "strconv"
)

type EndpointName uint32

var Schema string = "https://api.unsplash.com/"

const (
    BaseURL      EndpointName = 1 << (32 - 1 - iota)
    RandomPhotos
    SearchPhotos
)

// URL converts endpoint name into its string representation.
func URL(name EndpointName) string {
    switch name {
    case BaseURL:
        return Schema
    case RandomPhotos:
        return Schema + "photos/random"
    case SearchPhotos:
        return Schema + "search/photos"
    }
    return ""
}

var sendRequest = func(req *http.Request, token string) (data []byte, err error) {
    req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", token))
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil { return }
    defer resp.Body.Close()
    data, err = ioutil.ReadAll(resp.Body)
    if err != nil { return }
    return data, nil
}

// Client wraps Unsplash REST API with convenient interface.
type Client struct {
    AccessKey, SecretKey string
}

type Result struct {
    Id string `json:"id"`
    URLs map[string]string `json:"urls"`
}

func (c *Client) DownloadRandomPhotos(dirname string, count int) (err error) {
    log.Printf("Downloading %d image(s) into folder %s\n", count, dirname)
    result, err := c.GetRandomPhotos(count)
    if err != nil { return }

    log.Println("Creating folder")
    err = os.MkdirAll(dirname, os.ModePerm)
    if err != nil { return }

    n := len(result)
    for i, item := range result {
        fmt.Printf("Downloading image %d of %d...\r", i+1, n)
        imageURL := item.URLs["regular"]
        img, err := FetchImage(imageURL)
        if err != nil { return err }

        filename := path.Join(dirname, fmt.Sprintf("%s.jpeg", item.Id))
        fmt.Printf("Saving downloaded image into file: %s\n", filename)
        file, err := os.Create(filename)
        if err != nil { return err }

        err = jpeg.Encode(file, img, nil)
        if err != nil { return err }

        _ = file.Close()
    }

    return nil
}

func (c *Client) GetRandomPhotos(count int) (result []Result, err error) {
    if count <= 0 {
        err = fmt.Errorf("number of photos should be >= 1 but %d received", count)
        return
    }
    req, _ := http.NewRequest("GET", URL(RandomPhotos), nil)
    values := req.URL.Query()
    values.Set("client_id", c.AccessKey)
    values.Set("count", strconv.Itoa(count))
    req.URL.RawQuery = values.Encode()
    data, err := sendRequest(req, c.SecretKey)
    if err != nil { return }
    return MustDecodeArray(data), nil
}

// MustDecodeArray is sure that data contains a valid Unsplash response and panics otherwise.
func MustDecodeArray(data []byte) (result []Result) {
    if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&result); err != nil {
        panic(fmt.Sprintf("Cannot decode Unsplash response: %s", err))
    }
    return result
}

func FetchImage(url string) (image image.Image, err error) {
    resp, err := http.Get(url)
    if err != nil { return }
    defer resp.Body.Close()
    image, err = jpeg.Decode(resp.Body)
    if err != nil { return }
    return image, nil
}