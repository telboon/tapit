package main

import (
    "github.com/jinzhu/gorm"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

// GlobalSettings contains basic common settings across the app
type GlobalSettings struct {
    gorm.Model
    SecretRegistrationCode string
    ThreadsPerCampaign int
    BcryptCost int
    MaxRequestRetries int
    WaitBeforeRetry int
    WebTemplatePrefix string
    WebFrontPlaceholder string
}

type GlobalSettingsJson struct {
    SecretRegistrationCode string `json:"secretRegistrationCode"`
    ThreadsPerCampaign int `json:"threadsPerCampaign"`
    BcryptCost int `json:"bcryptCost"`
    MaxRequestRetries int `json:"maxRequestRetries"`
    WaitBeforeRetry int `json:"waitBeforeRetry"`
    WebTemplatePrefix string `json:"webTemplatePrefix"`
    WebFrontPlaceholder string `json:"webFrontPlaceholder"`
}


func (tapit *Tapit) handleGlobalSettings(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "PUT" {
        tapit.updateGlobalSettings(w, r)
    } else if strings.ToUpper(r.Method) == "GET" {
        tapit.getGlobalSettings(w,r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) updateGlobalSettings(w http.ResponseWriter, r *http.Request) {
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    var globalSettingsJson GlobalSettingsJson
    var globalSettings GlobalSettings
    err = tapit.db.Last(&globalSettings).Error
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    err = json.Unmarshal(requestBody, &globalSettingsJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if globalSettingsJson.SecretRegistrationCode != "" && globalSettingsJson.ThreadsPerCampaign != 0 && globalSettingsJson.BcryptCost != 0 && globalSettingsJson.WebTemplatePrefix != "" {
        globalSettings.SecretRegistrationCode = globalSettingsJson.SecretRegistrationCode
        globalSettings.ThreadsPerCampaign = globalSettingsJson.ThreadsPerCampaign
        globalSettings.BcryptCost = globalSettingsJson.BcryptCost
        globalSettings.MaxRequestRetries = globalSettingsJson.MaxRequestRetries
        globalSettings.WaitBeforeRetry = globalSettingsJson.WaitBeforeRetry
        globalSettings.WebTemplatePrefix = globalSettingsJson.WebTemplatePrefix
        globalSettings.WebFrontPlaceholder = globalSettingsJson.WebFrontPlaceholder
        err = tapit.db.Save(&globalSettings).Error

        if err != nil {
            notifyPopup(w, r, "failure", "Failed to update global settings", nil)
            return
        } else {
            tapit.globalSettings = globalSettings
            notifyPopup(w, r, "success", "Successfully updated global settings", globalSettingsJson)
            return
        }
    } else {
        notifyPopup(w, r, "failure", "Failed to update global settings", nil)
        return
    }
}

func (tapit *Tapit) getGlobalSettings(w http.ResponseWriter, r *http.Request) {
    var globalSettingsJson GlobalSettingsJson
    var globalSettings GlobalSettings
    err := tapit.db.Last(&globalSettings).Error

    if err == nil {
        globalSettingsJson.SecretRegistrationCode = globalSettings.SecretRegistrationCode
        globalSettingsJson.ThreadsPerCampaign = globalSettings.ThreadsPerCampaign
        globalSettingsJson.BcryptCost = globalSettings.BcryptCost
        globalSettingsJson.MaxRequestRetries = globalSettings.MaxRequestRetries
        globalSettingsJson.WaitBeforeRetry = globalSettings.WaitBeforeRetry
        globalSettingsJson.WebTemplatePrefix = globalSettings.WebTemplatePrefix
        globalSettingsJson.WebFrontPlaceholder = globalSettings.WebFrontPlaceholder

        jsonResults, err := json.Marshal(globalSettingsJson)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.Write(jsonResults)
            return
        }
    } else {
        http.Error(w, "Bad request", 400)
        return
    }
}
