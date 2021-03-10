package main

import (
    "github.com/jinzhu/gorm"
    "github.com/gorilla/mux"
    "sync"
    "time"
    "net/http"
    "strings"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "log"
)

type Campaign struct {
    gorm.Model
    Name string
    FromNumber string
    Size int
    CurrentStatus string        // enum Running, Paused, Completed, Not Started
    PhonebookId uint
    TextTemplateId uint
    WebTemplateId uint
    ProviderTag string
    Jobs []Job              `gorm:"foreignkey:CampaignId"`
}

type CampaignComms struct {
    Campaign Campaign
    Action string           // enum run, stop
}

type JobComms struct {
    Job Job
    Action string           // enum run, stop
}

type CampaignJson struct {
    Id uint                     `json:"id"`
    Name string                 `json:"name"`
    FromNumber string           `json:"fromNumber"`
    Size int                    `json:"size"`
    CurrentStatus string        `json:"currentStatus"`
    CreateDate time.Time        `json:"createDate"`
    PhonebookId uint            `json:"phoneBookId"`
    TextTemplateId uint         `json:"textTemplateId"`
    WebTemplateId uint          `json:"webTemplateId"`
    ProviderTag string          `json:"providerTag"`
    Jobs []JobJson              `json:"jobs"`
}

type Job struct {
    gorm.Model
    CampaignId uint
    CurrentStatus string        // enum Failed, Queued, Sent, Delivered, Not Started
    WebStatus string            // enum Not Visited, xx visits
    TimeSent time.Time
    ProviderTag string
    AccSID string
    AuthToken string
    BodyText string
    FromNum string
    ToNum string
    ResultStr string
    MessageSid string
    WebRoute string
    FullUrl string
    Visits []Visit
}

type JobJson struct {
    Id uint                 `json:"id"`
    CurrentStatus string    `json:"currentStatus"`
    WebStatus string        `json:"webStatus"`
    TimeSent time.Time      `json:"timeSent"`
    FromNum string          `json:"fromNum"`
    ToNum string            `json:"toNum"`
    WebRoute string         `json:"webRoute"`
    FullUrl string          `json:"fullUrl"`
    Visits []VisitJson      `json:"visitJson"`
}

type TwilioMessageJson struct {
    AccountSid string                   `json:"account_sid"`
    ApiVersion string                   `json:"api_version"`
    Body string                         `json:"body"`
    DateCreated string                  `json:"date_created"`
    DateSent string                     `json:"date_sent"`
    DateUpdated string                  `json:"date_updated"`
    Direction string                    `json:"direction"`
    ErrorCode string                    `json:"error_code"`
    ErrorMessage string                 `json:"error_message"`
    From string                         `json:"from"`
    MessagingServiceSid string          `json:"messaging_service_sid"`
    NumMedia string                     `json:"num_media"`
    NumSegments string                  `json:"num_segments"`
    Price string                        `json:"price"`
    PriceUnit string                    `json:"price_unit"`
    Sid string                          `json:"sid"`
    Status string                       `json:"status"`
    SubResourceUri SubResourceUriJson   `json:"subresource_uris"`
    To string                           `json:"to"`
    Uri string                          `json:"uri"`
}

type SubResourceUriJson struct {
    Media string                `json:"media"`
}

func (tapit *Tapit) handleCampaign(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.getCampaigns(w, r)
    } else if strings.ToUpper(r.Method) == "POST" {
        tapit.createCampaign(w, r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) getCampaigns(w http.ResponseWriter, r *http.Request) {
    var campaigns []Campaign
    tapit.db.Find(&campaigns)
    jsonResults, err := json.Marshal(campaignsToJson(campaigns))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}

func campaignsToJson(campaigns []Campaign) []CampaignJson {
    var results []CampaignJson
    for _, currCampaign := range campaigns {
        var currJson CampaignJson
        currJson.Id = currCampaign.ID
        currJson.Name = currCampaign.Name
        currJson.Size = currCampaign.Size
        currJson.FromNumber = currCampaign.FromNumber
        currJson.CurrentStatus = currCampaign.CurrentStatus
        currJson.CreateDate = currCampaign.CreatedAt
        currJson.PhonebookId = currCampaign.PhonebookId
        currJson.TextTemplateId = currCampaign.TextTemplateId
        currJson.WebTemplateId = currCampaign.WebTemplateId
        currJson.ProviderTag = currCampaign.ProviderTag

        results = append(results, currJson)
    }

    return results
}

func (tapit *Tapit) createCampaign(w http.ResponseWriter, r *http.Request) {
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    var newCampaignJson CampaignJson
    err = json.Unmarshal(requestBody, &newCampaignJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if newCampaignJson.Name != "" && newCampaignJson.FromNumber != "" && newCampaignJson.PhonebookId != 0 && newCampaignJson.TextTemplateId != 0 {
        var newCampaign Campaign

        // populate details to be used later
        var newRecords []PhoneRecord
        var newTextTemplateBody string
        var newAccSID string
        var newAuthToken string
        newRecords = tapit.getSpecificPhonebook(newCampaignJson.PhonebookId).Records
        newTextTemplateBody = tapit.getSpecificTextBody(newCampaignJson.TextTemplateId)
        if newCampaignJson.ProviderTag == "twilio" {
            var twilioProvider TwilioProvider
            tapit.db.Last(&twilioProvider)

            newAccSID = twilioProvider.AccountSID
            newAuthToken = twilioProvider.AuthToken
        }

        // update static details
        newCampaign.Name = newCampaignJson.Name
        newCampaign.Size = len(newRecords)
        newCampaign.CurrentStatus = "Not Started"

        newCampaign.FromNumber = newCampaignJson.FromNumber
        newCampaign.PhonebookId = newCampaignJson.PhonebookId
        newCampaign.TextTemplateId = newCampaignJson.TextTemplateId
        newCampaign.WebTemplateId = newCampaignJson.WebTemplateId
        newCampaign.ProviderTag = newCampaignJson.ProviderTag

        // save campaign first
        tapit.db.NewRecord(&newCampaign)
        tapit.db.Create(&newCampaign)
        if newCampaign.ID == 0 {
            notifyPopup(w, r, "failure", "Failed to create campaign", nil)
            return
        }
        // update records
        for _, record := range newRecords {
            var newJob Job
            newJob.CurrentStatus = "Not Started"
            newJob.ProviderTag = newCampaign.ProviderTag
            newJob.AccSID = newAccSID
            newJob.AuthToken = newAuthToken
            newJob.FromNum = newCampaign.FromNumber

            // handle web template only if given
            if newCampaign.WebTemplateId != 0 {
                newJob.WebStatus = "Not Visited"
                newJob.WebRoute = tapit.generateWebTemplateRoute()
                newJob.FullUrl = tapit.globalSettings.WebTemplatePrefix + newJob.WebRoute
            }

            // interpreting records
            var newBodyText string
            newJob.ToNum = record.PhoneNumber
            newBodyText = newTextTemplateBody
            newBodyText = strings.Replace(newBodyText, "{firstName}", record.FirstName, -1)
            newBodyText = strings.Replace(newBodyText, "{lastName}", record.LastName, -1)
            newBodyText = strings.Replace(newBodyText, "{alias}", record.Alias, -1)
            newBodyText = strings.Replace(newBodyText, "{phoneNumber}", record.PhoneNumber, -1)
            newBodyText = strings.Replace(newBodyText, "{url}", newJob.FullUrl, -1)

            newJob.BodyText = newBodyText

            // saving it
            newCampaign.Jobs = append(newCampaign.Jobs, newJob)

            // update campaign
            tapit.db.Save(&newCampaign)
        }

        newCampaignJson.Id = newCampaign.ID
        newCampaignJson.CreateDate = newCampaign.CreatedAt
        newCampaignJson.Size = newCampaign.Size
        newCampaignJson.CurrentStatus = newCampaign.CurrentStatus

        notifyPopup(w, r, "success", "Successfully added new campaign", newCampaignJson)
        return
    } else {
        notifyPopup(w, r, "failure", "Please enter campaign details", nil)
        return
    }
}

func (tapit *Tapit) handleSpecificCampaign(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "PUT" {
        // not implmented -- complexity in changing campaign perimeters
        // tapit.updateCampaign(w, r)
        http.Error(w, "HTTP method not implemented", 400)
        return
    } else if strings.ToUpper(r.Method) == "DELETE" {
        tapit.deleteCampaign(w,r)
        return
    } else if strings.ToUpper(r.Method) == "GET" {
        tapit.getCampaign(w,r)
        return
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) getSpecificCampaign(id uint) Campaign {
    var campaign Campaign
    var jobs []Job

    var dbSearchCampaign Campaign
    dbSearchCampaign.ID = id
    tapit.db.Where(&dbSearchCampaign).First(&campaign)

    var dbSearchJob Job
    dbSearchJob.CampaignId = id
    tapit.db.Where(&dbSearchJob).Find(&jobs)

    campaign.Jobs = jobs
    return campaign
}

func (tapit *Tapit) getCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    phonebook := tapit.getSpecificCampaign(uint(tempID))

    jsonResults, err := json.Marshal(campaignToJson(phonebook))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}

func campaignToJson(campaign Campaign) CampaignJson {
    var cJson CampaignJson
    cJson.Id = campaign.ID
    cJson.Name = campaign.Name
    cJson.FromNumber = campaign.FromNumber
    cJson.Size = campaign.Size
    cJson.CurrentStatus = campaign.CurrentStatus
    cJson.PhonebookId = campaign.PhonebookId
    cJson.TextTemplateId = campaign.TextTemplateId
    cJson.WebTemplateId = campaign.WebTemplateId
    cJson.ProviderTag = campaign.ProviderTag

    // iterating jobs
    for _, job := range campaign.Jobs {
        var currJson JobJson
        currJson.Id = job.ID
        currJson.CurrentStatus = job.CurrentStatus
        currJson.WebStatus = job.WebStatus
        currJson.WebRoute = job.WebRoute
        currJson.FullUrl = job.FullUrl
        currJson.TimeSent = job.TimeSent
        currJson.FromNum = job.FromNum
        currJson.ToNum = job.ToNum

        cJson.Jobs = append(cJson.Jobs, currJson)
    }
    return cJson
}

func (tapit *Tapit) deleteCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var campaign Campaign

    // get phonebook
    var dbSearchCampaign Campaign
    dbSearchCampaign.ID = uint(tempID)
    tapit.db.Where(&dbSearchCampaign).First(&campaign)

    if campaign.ID == uint(tempID) {
        // finally delete it
        tapit.db.Delete(&campaign)
        notifyPopup(w, r, "success", "Successfully deleted campaign", nil)
        return
    } else {
        http.Error(w, "Bad request", 400)
        return
    }
}

func (tapit *Tapit) handleStartCampaign(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.startCampaign(w,r)
        return
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) handleStopCampaign(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.stopCampaign(w,r)
        return
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}


func (tapit *Tapit) startCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var campaign Campaign

    campaign = tapit.getSpecificCampaign(uint(tempID))

    if campaign.ID == uint(tempID) && campaign.CurrentStatus != "Running" && campaign.CurrentStatus != "Completed" {
        // finally start new thread and start working
        go tapit.workerCampaign(campaign)
        campaign.CurrentStatus = "Running"
        tapit.db.Save(&campaign)
        jsonResults := campaignToJson(campaign)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        } else {
            notifyPopup(w, r, "success", "Started campaign", jsonResults)
            return
        }
    } else {
        http.Error(w, "Bad request", 400)
        return
    }
}

func (tapit *Tapit) stopCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var campaign Campaign

    campaign = tapit.getSpecificCampaign(uint(tempID))

    if campaign.ID == uint(tempID) && campaign.CurrentStatus == "Running" {
        var campaignComms CampaignComms
        campaignComms.Action = "stop"
        campaignComms.Campaign = campaign
        tapit.campaignChan <- campaignComms

        // notify
        notifyPopup(w, r, "success", "Paused campaign", nil)
        return
    } else {
        http.Error(w, "Bad request", 400)
        return
    }
}

func (tapit *Tapit) workerCampaign(campaign Campaign) {
    var campaignComms CampaignComms
    var jobChan chan JobComms
    var wg sync.WaitGroup

    jobChan = make(chan JobComms, 1)
    for i:=0; i<tapit.globalSettings.ThreadsPerCampaign; i++ {
        wg.Add(1)
        go tapit.workerJob(jobChan, &wg)
    }

    for _, job := range campaign.Jobs {
        select {
            case campaignComms = <-tapit.campaignChan:
                if campaignComms.Campaign.ID == campaign.ID {
                    if campaignComms.Action == "stop" {
                        // kill all
                        for i:=0; i<tapit.globalSettings.ThreadsPerCampaign; i++ {
                            var stopComms JobComms
                            stopComms.Action = "stop"
                            jobChan <- stopComms
                        }
                        // wait to end
                        wg.Wait()

                        // get updated campaign
                        var newCampaign Campaign
                        var searchCampaign Campaign
                        searchCampaign.ID = campaign.ID
                        tapit.db.Where(&searchCampaign).First(&newCampaign)

                        // update campaign
                        newCampaign.CurrentStatus = "Paused"
                        tapit.db.Save(&newCampaign)
                        return
                    }
                } else {
                    // not mine -- throw it back
                    tapit.campaignChan<- campaignComms
                }
            default:
                if job.CurrentStatus == "Not Started" {
                    var workComms JobComms
                    workComms.Action = "run"
                    workComms.Job = job
                    jobChan <- workComms
                }
        }
    }
    for i:=0; i<tapit.globalSettings.ThreadsPerCampaign; i++ {
        var stopComms JobComms
        stopComms.Action = "stop"
        jobChan <- stopComms
    }

    // wait to end
    wg.Wait()

    // get updated campaign
    var newCampaign Campaign
    var searchCampaign Campaign
    searchCampaign.ID = campaign.ID
    tapit.db.Where(&searchCampaign).First(&newCampaign)

    // update campaign
    newCampaign.CurrentStatus = "Completed"
    tapit.db.Save(&newCampaign)
}

func (tapit *Tapit) workerJob(jobChan chan JobComms, wg *sync.WaitGroup) {
    var currentJob JobComms
    exitCode := false

    for !exitCode {
        currentJob = <-jobChan
        if currentJob.Action != "stop" {
            if currentJob.Job.ProviderTag == "twilio" {

                var resultJson []byte
                resultJson = tapit.twilioSend(currentJob.Job.AccSID, currentJob.Job.AuthToken, currentJob.Job.BodyText, currentJob.Job.FromNum, currentJob.Job.ToNum)
                currentJob.Job.ResultStr = string(resultJson)

                var twilioResult TwilioMessageJson
                err := json.Unmarshal(resultJson, &twilioResult)
                if err != nil {
                    log.Println(err)
                    currentJob.Job.CurrentStatus = "Failed"
                } else if twilioResult.Status == "queued" {
                    currentJob.Job.MessageSid = twilioResult.Sid
                    currentJob.Job.CurrentStatus = "Queued"
                } else if twilioResult.Status == "delivered" {
                    currentJob.Job.MessageSid = twilioResult.Sid
                    currentJob.Job.CurrentStatus = "Delivered"
                    currentJob.Job.TimeSent = time.Now()
                } else {
                    currentJob.Job.CurrentStatus = "Failed"
                }

                // redo until done
                tapit.db.Save(&currentJob.Job)
            }
        } else {
            exitCode = true
        }
    }
    wg.Done()
}

func (tapit *Tapit) clearRunningCampaigns() {
    var campaigns []Campaign
    var searchCampaign Campaign
    searchCampaign.CurrentStatus = "Running"
    tapit.db.Where(&searchCampaign).Find(&campaigns)

    for _, campaign := range campaigns {
        campaign.CurrentStatus = "Paused"
        tapit.db.Save(&campaign)
    }
}

