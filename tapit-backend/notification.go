package main

import (
    "encoding/json"
    "net/http"
)

type NotificationJson struct {
    ResultType string           `json:"resultType"`           // success/failure/info
    Text string                 `json:"text"`
    Payload interface{}         `json:"payload"`
}

type Payload interface{}

func notifyPopup(w http.ResponseWriter, r *http.Request, resultType string, text string, payload Payload) {
    messageOutput := NotificationJson{
        ResultType:   resultType,
        Text: text,
        Payload: payload,
        }
    jsonResults, err := json.Marshal(messageOutput)
    if err!=nil {
        http.Error(w, "Internal server error", 500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    http.Error(w, string(jsonResults), 200)
    return
}
