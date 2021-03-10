package main

import (
    "github.com/jinzhu/gorm"
    "github.com/gorilla/mux"
    "time"
    "log"
    "math/rand"
    "net/http"
    "net/http/httputil"
    "strings"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "encoding/csv"
    "bytes"
)

// WebTemplate is the persistent object within Postgres
type WebTemplate struct {
    gorm.Model
    Name string
    TemplateType string         // enum redirect, harvester
    RedirectAgent string
    RedirectNegAgent string
    RedirectPlaceholderHtml string
    RedirectUrl string
    HarvesterBeforeHtml string
    HarvesterAfterHtml string
}

// WebTemplateJson is the temporary object for JSON data passing
type WebTemplateJson struct {
    Id int                          `json:"id"`
    Name string                     `json:"name"`
    TemplateType string             `json:"templateType"`
    RedirectAgent string            `json:"redirectAgent"`
    RedirectNegAgent string         `json:"redirectNegAgent"`
    RedirectPlaceholderHtml string  `json:"redirectPlaceholderHtml"`
    RedirectUrl string              `json:"redirectUrl"`
    HarvesterBeforeHtml string      `json:"harvesterBeforeHtml"`
    HarvesterAfterHtml string       `json:"harvesterAfterHtml"`
    CreateDate time.Time            `json:"createDate"`
}

type Visit struct {
    gorm.Model
    JobId uint
    SourceIp string
    UserAgent string
    Method string
    BodyContent string
    RawRequest string
}

type VisitJson struct {
    Id uint                     `json:"id"`
    JobId uint                  `json:"jobId"`
    SourceIP string             `json:"sourceIp"`
    UserAgent string            `json:"userAgent"`
    Method string               `json:"method"`
    BodyContent string          `json:"bodyContent"`
    RawRequest string           `json:"rawRequest"`
    CreateDate time.Time        `json:"createDate"`
}

func (tapit *Tapit) handleWebTemplate(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.getWebTemplates(w, r)
    } else if strings.ToUpper(r.Method) == "POST" {
        tapit.createWebTemplate(w, r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) getWebTemplates(w http.ResponseWriter, r *http.Request) {
    webTemplates := []WebTemplate{}
    tapit.db.Find(&webTemplates)
    jsonResults, err := json.Marshal(webTemplatesToJson(webTemplates))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}

func webTemplatesToJson(webTemplates []WebTemplate) []WebTemplateJson {
    webTemplateJson := make([]WebTemplateJson, 0)
    for _, webTemplate := range webTemplates {
        var currentWebTemplateJson WebTemplateJson
        currentWebTemplateJson.Id = int(webTemplate.ID)
        currentWebTemplateJson.Name = webTemplate.Name
        currentWebTemplateJson.TemplateType = webTemplate.TemplateType
        currentWebTemplateJson.RedirectAgent = webTemplate.RedirectAgent
        currentWebTemplateJson.RedirectNegAgent = webTemplate.RedirectNegAgent
        currentWebTemplateJson.RedirectPlaceholderHtml = webTemplate.RedirectPlaceholderHtml
        currentWebTemplateJson.RedirectUrl = webTemplate.RedirectUrl
        currentWebTemplateJson.HarvesterBeforeHtml = webTemplate.HarvesterBeforeHtml
        currentWebTemplateJson.HarvesterAfterHtml = webTemplate.HarvesterAfterHtml
        currentWebTemplateJson.CreateDate = webTemplate.CreatedAt

        webTemplateJson = append(webTemplateJson, currentWebTemplateJson)
    }
    return webTemplateJson
}

func (tapit *Tapit) createWebTemplate(w http.ResponseWriter, r *http.Request) {
    // start doing work
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    newWebTemplateJson := WebTemplateJson{}
    err = json.Unmarshal(requestBody, &newWebTemplateJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if newWebTemplateJson.Name != "" && newWebTemplateJson.TemplateType != "" {
        // check that not both user agents are filled
        if newWebTemplateJson.RedirectAgent != "" && newWebTemplateJson.RedirectNegAgent != "" {
            notifyPopup(w, r, "failure", "Please fill in only either positive or negative redirect user agent.", nil)
            return
        }
        newWebTemplate := jsonToWebTemplate(newWebTemplateJson)
        tapit.db.NewRecord(&newWebTemplate)
        tapit.db.Create(&newWebTemplate)
        if newWebTemplate.ID == 0 {
            notifyPopup(w, r, "failure", "Failed to create text template", nil)
            return
        }
        newWebTemplateJson.Id = int(newWebTemplate.ID)
        newWebTemplateJson.CreateDate = newWebTemplate.CreatedAt

        notifyPopup(w, r, "success", "Successfully added new text template", newWebTemplateJson)
        return
    } else {
        notifyPopup(w, r, "failure", "Please fill in all details", nil)
        return
    }
}

func jsonToWebTemplate(currentWebTemplateJson WebTemplateJson) WebTemplate {
    var webTemplate WebTemplate

    webTemplate.Name = currentWebTemplateJson.Name
    webTemplate.TemplateType = currentWebTemplateJson.TemplateType
    webTemplate.RedirectAgent = currentWebTemplateJson.RedirectAgent
    webTemplate.RedirectNegAgent = currentWebTemplateJson.RedirectNegAgent
    webTemplate.RedirectPlaceholderHtml = currentWebTemplateJson.RedirectPlaceholderHtml
    webTemplate.RedirectUrl = currentWebTemplateJson.RedirectUrl
    webTemplate.HarvesterBeforeHtml = currentWebTemplateJson.HarvesterBeforeHtml
    webTemplate.HarvesterAfterHtml = currentWebTemplateJson.HarvesterAfterHtml

    return webTemplate
}

func (tapit *Tapit) handleSpecificWebTemplate(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "PUT" {
        tapit.updateWebTemplate(w, r)
    } else if strings.ToUpper(r.Method) == "DELETE" {
        tapit.deleteWebTemplate(w,r)
    } else if strings.ToUpper(r.Method) == "GET" {
        tapit.getWebTemplate(w,r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) updateWebTemplate(w http.ResponseWriter, r *http.Request) {
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    var newWebTemplateJson WebTemplateJson
    err = json.Unmarshal(requestBody, &newWebTemplateJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if newWebTemplateJson.Name != "" && newWebTemplateJson.TemplateType != "" {
        var newWebTemplate WebTemplate

        // get current phonebook
        var dbSearchWT WebTemplate
        dbSearchWT.ID = uint(newWebTemplateJson.Id)
        tapit.db.Where(&dbSearchWT).First(&newWebTemplate)

        if newWebTemplate.ID == uint(newWebTemplateJson.Id) {
            // update name & template
            newWebTemplate.Name = newWebTemplateJson.Name
            newWebTemplate.TemplateType = newWebTemplateJson.TemplateType
            newWebTemplate.RedirectAgent = newWebTemplateJson.RedirectAgent
            newWebTemplate.RedirectNegAgent = newWebTemplateJson.RedirectNegAgent
            newWebTemplate.RedirectPlaceholderHtml = newWebTemplateJson.RedirectPlaceholderHtml
            newWebTemplate.RedirectUrl = newWebTemplateJson.RedirectUrl
            newWebTemplate.HarvesterBeforeHtml = newWebTemplateJson.HarvesterBeforeHtml
            newWebTemplate.HarvesterAfterHtml = newWebTemplateJson.HarvesterAfterHtml

            // update database
            tapit.db.Save(&newWebTemplate)
            if newWebTemplate.ID == 0 {
                notifyPopup(w, r, "failure", "Failed to update phonebook", nil)
                return
            }
            newWebTemplateJson.Id = int(newWebTemplate.ID)
            newWebTemplateJson.CreateDate = newWebTemplate.CreatedAt

            notifyPopup(w, r, "success", "Successfully updated web template", newWebTemplateJson)
            return
        } else {
            notifyPopup(w, r, "failure", "Failed to update web template", nil)
            return
        }
    } else {
        notifyPopup(w, r, "failure", "Please enter all details", nil)
        return
    }
}

func (tapit *Tapit) deleteWebTemplate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var webTemplate WebTemplate

    // get tt
    var dbSearchWT WebTemplate
    dbSearchWT.ID = uint(tempID)
    tapit.db.Where(dbSearchWT).First(&webTemplate)

    if webTemplate.ID == uint(tempID) {
        // finally delete it
        tapit.db.Delete(&webTemplate)
        notifyPopup(w, r, "success", "Successfully deleted phonebook", nil)
        return
    } else {
        http.Error(w, "Bad request", 400)
        return
    }
}

func (tapit *Tapit) getWebTemplate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var webTemplate WebTemplate

    // get tt
    var dbSearchWT WebTemplate
    dbSearchWT.ID = uint(tempID)
    tapit.db.Where(dbSearchWT).First(&webTemplate)

    if webTemplate.ID == uint(tempID) {
        jsonResults, err := json.Marshal(webTemplateToJson(webTemplate))
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

func (tapit *Tapit) generateWebTemplateRoute() string {
    charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
    // generate 5 char
    var newRoute string
    var successRoute bool
    successRoute = false

    for !successRoute {
        newRoute = ""
        for i:=0; i<5; i++ {
            num := rand.Int() % len(charset)
            newRoute = newRoute + string(charset[num])

            // search if route already exists
            var dbSearchJob Job
            var jobs []Job
            dbSearchJob.WebRoute = newRoute
            tapit.db.Where(&dbSearchJob).Find(&jobs)
            if len(jobs) == 0 {
                successRoute = true
            }
        }
    }
    return newRoute
}

func webTemplateToJson(webTemplate WebTemplate) WebTemplateJson {
    var currentWebTemplateJson WebTemplateJson
    currentWebTemplateJson.Id = int(webTemplate.ID)
    currentWebTemplateJson.Name = webTemplate.Name
    currentWebTemplateJson.TemplateType = webTemplate.TemplateType
    currentWebTemplateJson.RedirectAgent = webTemplate.RedirectAgent
    currentWebTemplateJson.RedirectNegAgent = webTemplate.RedirectNegAgent
    currentWebTemplateJson.RedirectPlaceholderHtml = webTemplate.RedirectPlaceholderHtml
    currentWebTemplateJson.RedirectUrl = webTemplate.RedirectUrl
    currentWebTemplateJson.HarvesterBeforeHtml = webTemplate.HarvesterBeforeHtml
    currentWebTemplateJson.HarvesterAfterHtml = webTemplate.HarvesterAfterHtml
    currentWebTemplateJson.CreateDate = webTemplate.CreatedAt
    return currentWebTemplateJson
}

func (tapit *Tapit) webTemplateRouteHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    vars := mux.Vars(r)
    currRoute := vars["route"]

    currJob := Job{}
    err = tapit.db.Where(&Job{WebRoute:currRoute}).First(&currJob).Error
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }

    currCampaign := Campaign{}
    err = tapit.db.Where(&Campaign{Model: gorm.Model{ID:currJob.CampaignId}}).First(&currCampaign).Error
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }

    currWebTemplate := WebTemplate{}
    err = tapit.db.Where(&WebTemplate{Model: gorm.Model{ID:currCampaign.WebTemplateId}}).First(&currWebTemplate).Error
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }

    // check type for "redirect" or "harvester"
    if currWebTemplate.TemplateType == "redirect" {
        if currWebTemplate.RedirectAgent != "" {
            listOfUA := strings.Split(currWebTemplate.RedirectAgent, ",")
            currCheck := false
            for _, currUA := range listOfUA {
                // check if user agent matches
                if strings.Contains(r.UserAgent(), currUA) {
                    currCheck = true
                }
            }

            // if matches at least once, redirect, otherwise placeholder
            if currCheck == true {
                http.Redirect(w, r, currWebTemplate.RedirectUrl, 302)
            } else {
                w.Write([]byte(currWebTemplate.RedirectPlaceholderHtml))
            }
        } else {
            listOfUA := strings.Split(currWebTemplate.RedirectNegAgent, ",")
            currCheck := true
            for _, currUA := range listOfUA {
                // check if user agent matches
                if strings.Contains(r.UserAgent(), currUA) {
                    currCheck = false
                }
            }

            // if matches at least once, redirect, otherwise placeholder
            if currCheck == true {
                http.Redirect(w, r, currWebTemplate.RedirectUrl, 302)
            } else {
                w.Write([]byte(currWebTemplate.RedirectPlaceholderHtml))
            }
        }
    } else if currWebTemplate.TemplateType == "harvester" {
        // if get show before, if post show after
        if strings.ToUpper(r.Method) == "GET"{
            w.Write([]byte(currWebTemplate.HarvesterBeforeHtml))
        } else if strings.ToUpper(r.Method) == "POST"{
            w.Write([]byte(currWebTemplate.HarvesterAfterHtml))
        } else {
            http.Error(w, "Bad request", 400)
        }
    }

    // saving records
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    var newVisit Visit

    newVisit = Visit{}
    newVisit.JobId = currJob.ID
    if r.Header.Get("X-Forwarded-For") == "" {
        newVisit.SourceIp = r.RemoteAddr
    } else {
        newVisit.SourceIp = r.Header.Get("X-Forwarded-For")
    }
    newVisit.UserAgent = r.UserAgent()
    newVisit.Method = r.Method
    newVisit.BodyContent = string(requestBody)
    rawReqBytes, err := httputil.DumpRequest(r, true)
    if err == nil {
        newVisit.RawRequest = string(rawReqBytes) + string(requestBody)
    }

    // Update visited status
    var visits []Visit
    tapit.db.Where(Visit{JobId: uint(currJob.ID)}).Find(&visits)
    currJob.WebStatus = strconv.Itoa(len(visits) + 1) + " visits"

    tapit.db.Save(&currJob)

    tapit.db.NewRecord(&newVisit)
    tapit.db.Create(&newVisit)

    return
}

func (tapit *Tapit) handleWebFront(w http.ResponseWriter, r *http.Request) {
    var globalSettings GlobalSettings
    err := tapit.db.Last(&globalSettings).Error
    if err != nil {
        w.Write([]byte(""))
        return
    }

    w.Write([]byte(globalSettings.WebFrontPlaceholder))
    return
}

func (tapit *Tapit) handleDownloadView(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        var csvBuffer bytes.Buffer
        vars := mux.Vars(r)
        tempID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Bad request", 400)
            return
        }

        var visits []Visit
        tapit.db.Where(Visit{JobId: uint(tempID)}).Find(&visits)

        // generate csv
        csvWriter := csv.NewWriter(&csvBuffer)
        csvWriter.Write([]string{"ID", "Time", "Source IP", "User Agent", "Method",  "Body Content", "Raw Request"})
        for _, visit := range visits {
            csvWriter.Write([]string{strconv.Itoa(int(visit.ID)), visit.CreatedAt.String(), visit.SourceIp, visit.UserAgent, visit.Method, visit.BodyContent, visit.RawRequest})
        }
        csvWriter.Flush()
        w.Header().Set("Content-Disposition", "attachment; filename=\"results.csv\"")
        w.Write(csvBuffer.Bytes())
        return
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}
