package main

import (
    "github.com/jinzhu/gorm"
    "github.com/gorilla/mux"
    "github.com/tealeg/xlsx"
    "time"
    "net/http"
    "strings"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "log"
    "io"
    "bytes"
)

type Phonebook struct {
    gorm.Model
    Name string
    Size int
    Records []PhoneRecord           `gorm:"foreignkey:PhonebookID"`
}

type PhonebookJson struct {
    Id uint                          `json:"id"`
    Name string                     `json:"name"`
    Size int                        `json:"size"`
    CreateDate time.Time            `json:"createDate"`
    Records []PhoneRecordJson       `json:"records"`
}

type PhoneRecord struct {
    gorm.Model
    PhonebookID uint
    FirstName string
    LastName string
    Alias string
    PhoneNumber string
}

type PhoneRecordJson struct {
    Id uint                         `json:"id"`
    FirstName string                `json:"firstName"`
    LastName string                 `json:"lastName"`
    Alias string                    `json:"alias"`
    PhoneNumber string              `json:"phoneNumber"`
}

func (tapit *Tapit) handlePhonebook(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.getPhonebooks(w, r)
    } else if strings.ToUpper(r.Method) == "POST" {
        tapit.createPhonebook(w, r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) getPhonebooks(w http.ResponseWriter, r *http.Request) {
    var phonebooks []Phonebook
    tapit.db.Find(&phonebooks)
    jsonResults, err := json.Marshal(phonebooksToJson(phonebooks))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}

func phonebooksToJson(pb []Phonebook) []PhonebookJson {
    var pbJson []PhonebookJson
    for _, currObj := range pb {
        var currPbJson PhonebookJson
        currPbJson.Id = currObj.ID
        currPbJson.Name = currObj.Name
        currPbJson.CreateDate = currObj.CreatedAt
        currPbJson.Size = currObj.Size

        pbJson = append(pbJson, currPbJson)
    }
    return pbJson
}

func (tapit *Tapit) createPhonebook(w http.ResponseWriter, r *http.Request) {
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    var newPhonebookJson PhonebookJson
    err = json.Unmarshal(requestBody, &newPhonebookJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if newPhonebookJson.Name != "" {
        var newPhonebook Phonebook

        // update name & size
        newPhonebook.Name = newPhonebookJson.Name
        newPhonebook.Size = len(newPhonebookJson.Records)

        // update records
        for _, record := range newPhonebookJson.Records {
            var newRecord PhoneRecord
            newRecord.FirstName = record.FirstName
            newRecord.LastName = record.LastName
            newRecord.Alias = record.Alias
            newRecord.PhoneNumber = record.PhoneNumber

            newPhonebook.Records = append(newPhonebook.Records, newRecord)
        }

        // update database
        tapit.db.NewRecord(&newPhonebook)
        tapit.db.Create(&newPhonebook)
        if newPhonebook.ID == 0 {
            notifyPopup(w, r, "failure", "Failed to create phonebook", nil)
            return
        }
        newPhonebookJson.Id = newPhonebook.ID
        newPhonebookJson.CreateDate = newPhonebook.CreatedAt
        newPhonebookJson.Size = newPhonebook.Size

        notifyPopup(w, r, "success", "Successfully added new phonebook", newPhonebookJson)
        return
    } else {
        notifyPopup(w, r, "failure", "Please enter the phonebook name", nil)
        return
    }
}

func (tapit *Tapit) getSpecificPhonebook(id uint) Phonebook {
    var phonebook Phonebook
    var records []PhoneRecord

    var dbPhonebookSearch Phonebook
    dbPhonebookSearch.ID = id
    tapit.db.Where(&dbPhonebookSearch).First(&phonebook)

    var dbSearchPhoneRecord PhoneRecord
    dbSearchPhoneRecord.PhonebookID = id
    tapit.db.Where(&dbSearchPhoneRecord).Find(&records)

    phonebook.Records = records
    return phonebook
}

func (tapit *Tapit) getPhonebook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    phonebook := tapit.getSpecificPhonebook(uint(tempID))

    jsonResults, err := json.Marshal(phonebookToJson(phonebook))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}

func phonebookToJson(pb Phonebook) PhonebookJson {
    var pbJson PhonebookJson
    pbJson.Id = pb.ID
    pbJson.Name = pb.Name
    pbJson.CreateDate = pb.CreatedAt
    pbJson.Size = pb.Size
    for _, record := range pb.Records {
        var recordJson PhoneRecordJson
        recordJson.Id = record.ID
        recordJson.FirstName = record.FirstName
        recordJson.LastName = record.LastName
        recordJson.Alias = record.Alias
        recordJson.PhoneNumber = record.PhoneNumber

        pbJson.Records = append(pbJson.Records, recordJson)
    }
    return pbJson
}

func (tapit *Tapit) handleSpecificPhonebook(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "PUT" {
        tapit.updatePhonebook(w, r)
    } else if strings.ToUpper(r.Method) == "DELETE" {
        tapit.deletePhonebook(w,r)
    } else if strings.ToUpper(r.Method) == "GET" {
        tapit.getPhonebook(w,r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) updatePhonebook(w http.ResponseWriter, r *http.Request) {
    requestBody, err:= ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    var newPhonebookJson PhonebookJson
    err = json.Unmarshal(requestBody, &newPhonebookJson)
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }
    if newPhonebookJson.Name != "" {
        var newPhonebook Phonebook

        // get current phonebook
        var dbSearchPhonebook Phonebook
        tapit.db.Where(&dbSearchPhonebook).First(&newPhonebook)

        // update name & size
        newPhonebook.Name = newPhonebookJson.Name
        newPhonebook.Size = len(newPhonebookJson.Records)

        // update records
        for _, record := range newPhonebookJson.Records {
            var newRecord PhoneRecord
            newRecord.FirstName = record.FirstName
            newRecord.LastName = record.LastName
            newRecord.Alias = record.Alias
            newRecord.PhoneNumber = record.PhoneNumber

            newPhonebook.Records = append(newPhonebook.Records, newRecord)
        }

        // update database
        tapit.db.Save(&newPhonebook)
        if newPhonebook.ID == 0 {
            notifyPopup(w, r, "failure", "Failed to create phonebook", nil)
            return
        }
        newPhonebookJson.Id = newPhonebook.ID
        newPhonebookJson.CreateDate = newPhonebook.CreatedAt
        newPhonebookJson.Size = newPhonebook.Size

        notifyPopup(w, r, "success", "Successfully added new phonebook", newPhonebookJson)
        return
    } else {
        notifyPopup(w, r, "failure", "Please enter the phonebook name", nil)
        return
    }
}

func (tapit *Tapit) deletePhonebook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tempID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Bad request", 400)
        return
    }

    // start working
    var phonebook Phonebook

    // get phonebook
    var dbSearchPhonebook Phonebook
    dbSearchPhonebook.ID = uint(tempID)
    tapit.db.Where(&dbSearchPhonebook).First(&phonebook)

    if phonebook.ID == uint(tempID) {
        // finally delete it
        tapit.db.Delete(&phonebook)
        notifyPopup(w, r, "success", "Successfully deleted phonebook", nil)
        return
    } else {
        http.Error(w, "Bad request", 400)
        return
    }
}

func (tapit *Tapit) importPhonebook(w http.ResponseWriter, r *http.Request) {
    var records []PhoneRecordJson
    err := r.ParseForm()
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }

    // 100 M reserved
    err = r.ParseMultipartForm(100000000)
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }

    importFile, _, err := r.FormFile("phonebookFile")
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }

    var buff bytes.Buffer

    // use buffer to bytes
    io.Copy(&buff, importFile)
    fileBytes := buff.Bytes()

    excelFile, err := xlsx.OpenBinary(fileBytes)
    if err != nil {
        log.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }

    for num, row := range excelFile.Sheet["import"].Rows {
        if num != 0 {
            var tempRecord PhoneRecordJson
            tempRecord.FirstName = row.Cells[0].Value
            tempRecord.LastName = row.Cells[1].Value
            tempRecord.Alias = row.Cells[2].Value
            tempRecord.PhoneNumber = row.Cells[3].Value
            records = append(records, tempRecord)
        }
    }
    jsonResults, err := json.Marshal(records)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResults)
        return
    }
}
