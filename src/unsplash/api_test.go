package unsplash

import (
    "encoding/json"
    "net/http"
    "strconv"
    "testing"
)

func TestClient_GetRandomPhotos(t *testing.T) {
    counts := []struct{
        param, expected int
        error bool
    }{
        {-1,0, true},
        {0,0,true},
        {1, 1, false},
        {5, 5, false},
    }
    c := Client{"mock", "mock"}
    for _, test := range counts {
        results, err := c.GetRandomPhotos(test.param)
        if test.error && (err == nil) {
            t.Errorf("Error was expected for case: %v", test)
        } else {
            got := len(results)
            if got != test.expected {
                t.Errorf("Invalid number of items retrieved: %d != %d", got, test.expected)
            }
        }
    }
}

func init() {
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
}