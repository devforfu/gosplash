package unsplash

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
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
