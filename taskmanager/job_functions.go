package taskmanager

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
)

func doWebCall(jn JobNotification) (output map[string]interface{}, err error) {
    var resp *http.Response
    var respBody []byte
    var rawOutput interface{}

    switch jn.Type {
    case "post-to-web":
        resp, err = http.Post(jn.Context["url"], jn.Context["content-type"], bytes.NewReader([]byte(jn.Context["data"])))
    case "get-from-web":
        resp, err = http.Get(jn.Context["url"])
    }

    if err != nil {
        return
    }

    defer resp.Body.Close()

    if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
        return
    }

    if err = json.Unmarshal([]byte(respBody), &rawOutput); err != nil {
        return
    }

    output, _ = rawOutput.(map[string]interface{})
    return
}

func logOutput(jn JobNotification) (output map[string]interface{}, err error) {
    log.Println(jn.Context["message"])
    return
}
