package main

import (
    "github.com/jinzhu/gorm"
    "github.com/gorilla/mux"
    "time"
    "net/http"
    "strings"
    "encoding/json"
    "io/ioutil"
    "strconv"
)

// TextTemplate is the persistent object within Postgres
type TextTemplate struct {
    gorm.Model
    Name string
    TemplateStr string
}

// TextTemplateJson is the temporary object for JSON data passing
type TextTemplateJson struct {
    Id int                          `json:"id"`
    Name string                     `json:"name"`
    TemplateStr string              `json:"templateStr"`
    CreateDate time.Time            `json:"createDate"`
}

func (tapit *Tapit) handleTextTemplate(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.getTextTemplates(w, r)
    } else if strings.ToUpper(r.Method) == "POST" {
        tapit.createTextTemplate(w, r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) getTextTemplates(w http.ResponseWriter, r *http.Request) {
    textTemplates := []TextTemplate{}
    tapit.db.Find(&textTemplates)
    jsonResults, err := json.Marshal(textTemplatesToJson(textTemplates))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}

func textTemplatesToJson(textTemplates []TextTemplate) []TextTemplateJson {
    textTemplateJson := make([]TextTemplateJson, 0)
    for _, textTemplate := range textTemplates {
        var currentTextTemplateJson TextTemplateJson
        currentTextTemplateJson.Id = int(textTemplate.ID)
        currentTextTemplateJson.Name = textTemplate.Name
        currentTextTemplateJson.TemplateStr = textTemplate.TemplateStr
        currentTextTemplateJson.CreateDate = textTemplate.CreatedAt

        textTemplateJson = append(textTemplateJson, currentTextTemplateJson)
    }
    return textTemplateJson
}

func jsonToTextTemplate(textTemplateJson TextTemplateJson) TextTemplate {
    var resultTextTemplate TextTemplate
    resultTextTemplate.Name = textTemplateJson.Name
    resultTextTemplate.TemplateStr = textTemplateJson.TemplateStr
    return resultTextTemplate
}

func (tapit *Tapit) createTextTemplate(w http.ResponseWriter, r *http.Request) {
    // start doing work
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    newTextTemplateJson := TextTemplateJson{}
    err = json.Unmarshal(requestBody, &newTextTemplateJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if newTextTemplateJson.Name != "" && newTextTemplateJson.TemplateStr != "" {
        newTextTemplate := jsonToTextTemplate(newTextTemplateJson)
        tapit.db.NewRecord(&newTextTemplate)
        tapit.db.Create(&newTextTemplate)
        if newTextTemplate.ID == 0 {
            notifyPopup(w, r, "failure", "Failed to create text template", nil)
            return
        }
        newTextTemplateJson.Id = int(newTextTemplate.ID)
        newTextTemplateJson.CreateDate = newTextTemplate.CreatedAt

        notifyPopup(w, r, "success", "Successfully added new text template", newTextTemplate)
        return
    } else {
        notifyPopup(w, r, "failure", "Please fill in all details", nil)
        return
    }
}

func (tapit *Tapit) getSpecificTextBody(id uint) string {
    var textTemplate TextTemplate

    var dbSearchTT TextTemplate
    dbSearchTT.ID = id
    tapit.db.Where(&dbSearchTT).First(&textTemplate)

    return textTemplate.TemplateStr
}

func (tapit *Tapit) handleSpecificTextTemplate(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "PUT" {
        tapit.updateTextTemplate(w, r)
    } else if strings.ToUpper(r.Method) == "DELETE" {
        tapit.deleteTextTemplate(w,r)
    } else if strings.ToUpper(r.Method) == "GET" {
        tapit.getTextTemplate(w,r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) updateTextTemplate(w http.ResponseWriter, r *http.Request) {
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    var newTextTemplateJson TextTemplateJson
    err = json.Unmarshal(requestBody, &newTextTemplateJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if newTextTemplateJson.Name != "" {
        var newTextTemplate TextTemplate

        // get current phonebook
        var dbSearchTT TextTemplate
        dbSearchTT.ID = uint(newTextTemplateJson.Id)
        tapit.db.Where(&dbSearchTT).First(&newTextTemplate)

        if newTextTemplate.ID == uint(newTextTemplateJson.Id) {
            // update name & template
            newTextTemplate.Name = newTextTemplateJson.Name
            newTextTemplate.TemplateStr = newTextTemplateJson.TemplateStr

            // update database
            tapit.db.Save(&newTextTemplate)
            if newTextTemplate.ID == 0 {
                notifyPopup(w, r, "failure", "Failed to update phonebook", nil)
                return
            }
            newTextTemplateJson.Id = int(newTextTemplate.ID)
            newTextTemplateJson.CreateDate = newTextTemplate.CreatedAt

            notifyPopup(w, r, "success", "Successfully updated text template", newTextTemplateJson)
            return
        } else {
            notifyPopup(w, r, "failure", "Failed to update text template", nil)
            return
        }
    } else {
        notifyPopup(w, r, "failure", "Please enter the phonebook name", nil)
        return
    }
}

func (tapit *Tapit) deleteTextTemplate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var textTemplate TextTemplate

    // get tt
    var dbSearchTT TextTemplate
    dbSearchTT.ID = uint(tempID)
    tapit.db.Where(dbSearchTT).First(&textTemplate)

    if textTemplate.ID == uint(tempID) {
        // finally delete it
        tapit.db.Delete(&textTemplate)
        notifyPopup(w, r, "success", "Successfully deleted phonebook", nil)
        return
    } else {
        http.Error(w, "Bad request", 400)
        return
    }
}

func (tapit *Tapit) getTextTemplate(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var textTemplate TextTemplate

    // get tt
    var dbSearchTT TextTemplate
    dbSearchTT.ID = uint(tempID)
    tapit.db.Where(dbSearchTT).First(&textTemplate)

    if textTemplate.ID == uint(tempID) {
        jsonResults, err := json.Marshal(textTemplateToJson(textTemplate))
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

func textTemplateToJson(textTemplate TextTemplate) TextTemplateJson {
    var result TextTemplateJson
    result.Id = int(textTemplate.ID)
    result.Name = textTemplate.Name
    result.TemplateStr = textTemplate.TemplateStr
    result.CreateDate = textTemplate.CreatedAt
    return result
}

