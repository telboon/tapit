package main

import (
    "github.com/jinzhu/gorm"
    "net/http"
    "net/url"
    "strings"
    "log"
    "encoding/json"
    "io/ioutil"
    "time"
)

type TwilioProvider struct {
    gorm.Model
    AccountSID string
    AuthToken string
}

type TwilioProviderJson struct {
    AccountSID string               `json:"accountSID"`
    AuthToken string                `json:"authToken"`
}

func (tapit *Tapit) handleTwilioProvider(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.getTwilioProvider(w, r)
    } else if strings.ToUpper(r.Method) == "POST" {
        tapit.updateTwilioProvider(w, r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) getTwilioProvider(w http.ResponseWriter, r *http.Request) {
    var twilioProvider TwilioProvider
    tapit.db.Last(&twilioProvider)
    jsonResults, err := json.Marshal(twilioProviderToJson(twilioProvider))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}

func (tapit *Tapit) updateTwilioProvider(w http.ResponseWriter, r *http.Request) {
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    var newTwilioProviderJson TwilioProviderJson
    err = json.Unmarshal(requestBody, &newTwilioProviderJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // first check if already exist
    var twilioProvider TwilioProvider
    tapit.db.Last(&twilioProvider)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // update twilioProvider
    twilioProvider.AccountSID = newTwilioProviderJson.AccountSID
    twilioProvider.AuthToken = newTwilioProviderJson.AuthToken

    // does not exist
    if twilioProvider.ID == 0 {
        tapit.db.NewRecord(&twilioProvider)
        tapit.db.Create(&twilioProvider)

        if twilioProvider.ID == 0 {
            notifyPopup(w, r, "failure", "Failed to create Twilio Provider", nil)
            return
        } else {
            notifyPopup(w, r, "success", "Twilio provider updated", newTwilioProviderJson)
            return
        }
    } else {
    // exists
        tapit.db.Save(&twilioProvider)
        notifyPopup(w, r, "success", "Twilio provider updated", newTwilioProviderJson)
        return
    }
}

func twilioProviderToJson(tProvider TwilioProvider) TwilioProviderJson {
    var results TwilioProviderJson
    results.AccountSID = tProvider.AccountSID
    results.AuthToken = tProvider.AuthToken
    return results
}

func (tapit *Tapit) twilioSend(accSid string, accToken string, bodyText string, fromNum string, toNum string) []byte {
    // if burp proxy is necessary
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    method1 := "POST"
    url1 := "https://api.twilio.com/2010-04-01/Accounts/"+accSid+"/Messages.json"
    // making body
    params := url.Values{}
    params.Add("Body", bodyText)
    params.Add("From", fromNum)
    params.Add("To", toNum)
    body1 := strings.NewReader(params.Encode())
    log.Println(params.Encode())
    // making request
    newRequest1, err := http.NewRequest(method1, url1, body1)
    if err != nil {
        log.Fatal("Error in creating request")
    }

    //basic auth with token
    newRequest1.SetBasicAuth(accSid, accToken)

    //set headers
    newRequest1.Header.Add("Content-Type","application/x-www-form-urlencoded; charset=UTF-8")

    // sending request
    res, err := client.Do(newRequest1)
    retriesLeft := tapit.globalSettings.MaxRequestRetries
    for err != nil && retriesLeft > 0 {
        log.Println("Error in sending request")
        res, err = client.Do(newRequest1)
        retriesLeft -= 1
        time.Sleep(time.Duration(tapit.globalSettings.WaitBeforeRetry) * time.Millisecond)
    }

    // exit gracefully if can't
    if err!= nil {
        var emptyBytes []byte
        return emptyBytes
    }
    outputStr, _ := ioutil.ReadAll(res.Body)
    log.Println(string(outputStr))
    return outputStr
}

func (tapit *Tapit) twilioCheck(accSid string, accToken string, messageSid string) []byte {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    method1 := "GET"
    url1 := "https://api.twilio.com/2010-04-01/Accounts/"+accSid+"/Messages/"+messageSid+".json"
    body1 := strings.NewReader("")
    newRequest1, err := http.NewRequest(method1, url1, body1)

    // authenticate
    newRequest1.SetBasicAuth(accSid, accToken)

    // sending request
    res, err := client.Do(newRequest1)
    retriesLeft := tapit.globalSettings.MaxRequestRetries
    for err != nil && retriesLeft > 0 {
        log.Println("Error in sending request")
        res, err = client.Do(newRequest1)
        retriesLeft -= 1
        time.Sleep(time.Duration(tapit.globalSettings.WaitBeforeRetry) * time.Millisecond)
    }

    // exit gracefully if can't
    if err!= nil {
        var emptyBytes []byte
        return emptyBytes
    }
    outputStr, _ := ioutil.ReadAll(res.Body)
    log.Println(string(outputStr))
    return outputStr
}

func (tapit *Tapit) workerTwilioChecker() {
    // infinite loop to keep checking for queued jobs to check delivery status
    for true {
        // sleep 5 second per cycle
        time.Sleep(5000 * time.Millisecond)
        var pendJobs []Job

        tapit.db.Where("provider_tag = ? AND (current_status = ? OR current_status = ?)", "twilio", "Queued", "Sent").Find(&pendJobs)

        for _, job := range pendJobs {
            // sleep 100ms per job
            time.Sleep(100 * time.Millisecond)
            resultJson := tapit.twilioCheck(job.AccSID, job.AuthToken, job.MessageSid)
            job.ResultStr = string(resultJson)

            var twilioResult TwilioMessageJson
            err := json.Unmarshal(resultJson, &twilioResult)
            if err != nil {
                log.Println(err)
                job.CurrentStatus = "Failed"
            } else if twilioResult.Status == "queued" {
                job.MessageSid = twilioResult.Sid
                job.CurrentStatus = "Queued"
            } else if twilioResult.Status == "sent" {
                job.MessageSid = twilioResult.Sid
                job.CurrentStatus = "Sent"
            } else if twilioResult.Status == "delivered" {
                job.MessageSid = twilioResult.Sid
                job.CurrentStatus = "Delivered"
                job.TimeSent = time.Now()
            }
            tapit.db.Save(&job)
        }
    }
}

