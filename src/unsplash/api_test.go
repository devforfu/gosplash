package unsplash

import (
    "encoding/json"
    "net/http"
    "strconv"
    "testing"
)

func TestClient_GetRandomPhotos(t *testing.T) {
    sendRequest = func(req *http.Request, token string) (data []byte, err error) {
        var results []Result
        var urls = map[string]string{"MockURL": "https://api.mock"}

        count, _ := strconv.Atoi(req.URL.Query().Get("count"))
        for i := 0; i < count; i++ {
            result := Result{strconv.Itoa(i), urls}
            results = append(results, result)
        }
        return json.Marshal(results)
    }

    c := Client{"mock", "mock"}
    results, _ := c.GetRandomPhotos(3)
    if len(results) != 3 {
        t.Errorf("Invalid number of items retrieved")
    }
}