package main

import (
    "net/http"
    "strings"
    "io/ioutil"
    "encoding/json"
    "github.com/jinzhu/gorm"
    "math/rand"
    "golang.org/x/crypto/bcrypt"
)

type UserJson struct {
    Username string             `json:"username"`
    Password string             `json:"password"`
    Name string                 `json:"name"`
    Email string                `json:"email"`
    SecretCode string           `json:"secretCode"`
}

type User struct {
    gorm.Model
    Username string
    PasswordHash string
    Name string
    Email string
}


type Session struct {
    gorm.Model
    SessionID string
    UserID uint
}

func (tapit *Tapit) login(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "POST" {
        // start doing work
        requestBody, err:= ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Bad request", 400)
            return
        }
        userJson := UserJson{}
        err = json.Unmarshal(requestBody, &userJson)
        if err != nil {
            http.Error(w, "Bad request", 400)
            return
        }
        currUser := User{}
        tapit.db.Where(&User{Username:userJson.Username}).First(&currUser)
        // user exists
        if currUser.Username == userJson.Username {
            // checking hash...
            if checkPasswordHash(currUser.PasswordHash, userJson.Password) {
                userJson.Password = ""
                userJson.Name = currUser.Name
                userJson.Email = currUser.Email
                messageOutput := NotificationJson{
                    Text: "Successfully logged in!",
                    ResultType:   "success",
                    Payload: userJson,
                    }
                jsonResults, err := json.Marshal(messageOutput)
                if err!=nil {
                    http.Error(w, "Internal server error", 500)
                    return
                }
                w.Header().Set("Content-Type", "application/json")
                authCookie := tapit.generateCookie(currUser)
                http.SetCookie(w, &authCookie)
                w.Write(jsonResults)
                return
            } else {
                notifyPopup(w, r, "failure", "Username or password is incorrect", nil)
                return
            }
        } else {
            tapit.hashPassword("nothing-to-do-waste-time")
            notifyPopup(w, r, "failure", "Username or password is incorrect", nil)
            return
        }
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) register(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "POST" {
        // start doing work
        requestBody, err:= ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Bad request", 400)
            return
        }

        userJson := UserJson{}
        err = json.Unmarshal(requestBody, &userJson)
        if err != nil {
            http.Error(w, "Bad request", 400)
            return
        }

        // checks if secret code is correct
        if userJson.SecretCode != tapit.globalSettings.SecretRegistrationCode {
            messageOutput := NotificationJson{
            Text: "Your secret code is incorrect. Please try again.",
            ResultType:   "failure",
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

        //check if user exists
        currUser := User{}
        tapit.db.Where(&User{Username: userJson.Username}).First(&currUser)
        if currUser.Username != "" {
            messageOutput := NotificationJson{
            Text: "Username exists. Please choose another one.",
            ResultType:   "failure",
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

        //input validation that all are filled
        if userJson.Username == "" || userJson.Name == "" || userJson.Email == "" || userJson.Password == "" {
            messageOutput := NotificationJson{
            Text: "Please fill up all the information",
            ResultType:   "failure",
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

        // creates user...
        currUser.Username = userJson.Username
        currUser.Name = userJson.Name
        currUser.Email = userJson.Email
        currUser.PasswordHash, _ = tapit.hashPassword(userJson.Password)
        var jsonResults []byte
        if (tapit.db.NewRecord(&currUser)) {
            tapit.db.Create(&currUser)
            userJson.Password = ""
            messageOutput := NotificationJson{
                Text: "Successfully registered!",
                ResultType:   "success",
                Payload: userJson,
                }
            jsonResults, err = json.Marshal(messageOutput)
            if err!=nil {
                http.Error(w, "Internal server error", 500)
                return
            }
        } else {
            http.Error(w, "Internal server error", 500)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        authCookie := tapit.generateCookie(currUser)
        http.SetCookie(w, &authCookie)
        w.Write(jsonResults)
        return
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) logout(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "POST" {
        // start doing work
        var currSession Session
        authCookie, err := r.Cookie("tapitsession")
        if err!=nil {
            http.Error(w, "Not authorised", 401)
            return
        }
        authCookieStr := authCookie.String()[13:]
        tapit.db.Where(&Session{SessionID: authCookieStr}).First(&currSession)
        if currSession.SessionID != authCookieStr {
            http.Error(w, "Not authorised", 401)
            return
        } else {
            tapit.db.Delete(&currSession)
            messageOutput := NotificationJson{
            Text: "Successfully logged out",
            ResultType:   "success",
            Payload: "",
            }
            jsonResults, err := json.Marshal(messageOutput)
            if err!=nil {
                http.Error(w, "Internal server error", 500)
                return
            }
            delCookie := tapit.deleteCookie()
            http.SetCookie(w, &delCookie)
            w.Header().Set("Content-Type", "application/json")
            w.Write(jsonResults)
        }
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) authenticationHandler(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var currSession Session
        authCookie, err := r.Cookie("tapitsession")
        if err!=nil {
            http.Error(w, "Not authorised", 401)
            return
        }
        authCookieStr := authCookie.String()[13:]
        tapit.db.Where(&Session{SessionID: authCookieStr}).First(&currSession)
        if currSession.SessionID != authCookieStr {
            http.Error(w, "Not authorised", 401)
            return
        } else {
            next.ServeHTTP(w, r)
            return
        }
    }
}

func (tapit *Tapit) generateCookie(user User) http.Cookie {
    newToken := generateToken()
    newSession := Session{}
    tapit.db.Where(&Session{SessionID: newToken}).First(&newSession)
    for newToken == newSession.SessionID {
        newToken = generateToken()
        tapit.db.Where(&Session{SessionID: newToken}).First(&newSession)
    }
    newSession.UserID = user.ID
    newSession.SessionID = newToken
    tapit.db.NewRecord(&newSession)
    tapit.db.Create(&newSession)

    newCookie := http.Cookie {
        Name: "tapitsession",
        Value: newToken,
        Path: "/",
        MaxAge: 60*60*24*365*10,
        HttpOnly: true,
    }
    return newCookie
}

func (tapit *Tapit) deleteCookie() http.Cookie {
    newCookie := http.Cookie {
        Name: "tapitsession",
        Value: "",
        Path: "/",
        MaxAge: 0,
        HttpOnly: true,
    }
    return newCookie
}

func generateToken() string {
    var tokenResult strings.Builder
    var r int
    tokenCharset := "abcdefghijklmnopqrstuvwxyz0123456789"
    for i:=0; i<16; i++ {
        r = rand.Int() % len(tokenCharset)
        tokenResult.WriteRune(rune(tokenCharset[r]))
    }
    return tokenResult.String()
}

func (tapit *Tapit) hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), tapit.globalSettings.BcryptCost)
    return string(bytes), err
}

func checkPasswordHash(hash string, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (tapit *Tapit) myselfHandler(w http.ResponseWriter, r *http.Request) {
    if strings.ToUpper(r.Method) == "GET" {
        tapit.checkUser(w, r)
    } else if strings.ToUpper(r.Method) == "PUT" {
        tapit.updateUser(w, r)
    } else {
        http.Error(w, "HTTP method not implemented", 400)
        return
    }
}

func (tapit *Tapit) checkUser(w http.ResponseWriter, r *http.Request) {
    var currSession Session
    authCookie, err := r.Cookie("tapitsession")
    if err!=nil {
        http.Error(w, "Not authorised", 401)
        return
    }
    authCookieStr := authCookie.String()[13:]
    tapit.db.Where(&Session{SessionID: authCookieStr}).First(&currSession)
    if currSession.SessionID != authCookieStr {
        http.Error(w, "Not authorised", 401)
        return
    } else {
        currUser := User{}
        searchUser := User{}
        searchUser.ID = currSession.UserID
        tapit.db.Where(searchUser).First(&currUser)
        currentUserJson := UserJson{}
        currentUserJson.Username = currUser.Username
        currentUserJson.Name = currUser.Name
        currentUserJson.Email = currUser.Email
        jsonResults, err := json.Marshal(currentUserJson)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.Write(jsonResults)
            return
        }
    }
}

func (tapit *Tapit) updateUser(w http.ResponseWriter, r *http.Request) {
    var currSession Session
    authCookie, err := r.Cookie("tapitsession")
    if err!=nil {
        http.Error(w, "Not authorised", 401)
        return
    }
    authCookieStr := authCookie.String()[13:]
    tapit.db.Where(&Session{SessionID: authCookieStr}).First(&currSession)
    if currSession.SessionID != authCookieStr {
        http.Error(w, "Not authorised", 401)
        return
    } else {
        requestBody, err:= ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Bad request", 400)
            return
        }
        userJson := UserJson{}
        err = json.Unmarshal(requestBody, &userJson)
        if err != nil {
            http.Error(w, "Bad request", 400)
            return
        }
        currUser := User{}
        searchUser := User{}
        searchUser.ID = currSession.UserID
        tapit.db.Where(searchUser).First(&currUser)
        if currUser.ID == currSession.UserID && currUser.Username == userJson.Username {
            currUser.Name = userJson.Name
            currUser.Email = userJson.Email
            currUser.PasswordHash, _ = tapit.hashPassword(userJson.Password)
            tapit.db.Save(&currUser)
            userJson.Password = ""
            // writing output
            notifyPopup(w, r, "success", "Successfully changed profile!", userJson)
            return
        } else {
            http.Error(w, "Not authorised", 401)
            return
        }
    }
}
