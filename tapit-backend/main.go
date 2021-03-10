package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "log"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "path/filepath"
    "math/rand"
)

// Tapit is the general struct with shared objects
type Tapit struct {
    db *gorm.DB
    globalSettings GlobalSettings
    campaignChan chan CampaignComms
}

func generateFileHandler(path string) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        r.Header.Add("Cache-Control", "private, max-age=604800") // 7 days
        //r.Header.Add("Cache-Control", "private, max-age=1") // 1 sec -- debug
        http.ServeFile(w, r, path)
    }
}

func iterateStatic(r *mux.Router, path string, startWebPath string) {
    files, err := ioutil.ReadDir(path)
    if err!=nil {
        log.Fatal(err)
    }

    for _, f := range files {
        if !f.IsDir() && f.Name()[0] != '.' {
            r.HandleFunc(startWebPath + f.Name(), generateFileHandler(path+"/"+f.Name()))
            log.Println(startWebPath + f.Name()+" added to path")
        } else if f.IsDir() && f.Name()[0] != '.' {
            iterateStatic(r, path + "/" + string(f.Name()), startWebPath + string(f.Name() + "/"))
        }
    }
}

func generateRoutes(r *mux.Router, indexPath string, routes []string) {
    for _, route := range routes {
        r.HandleFunc(route, generateFileHandler(indexPath))
        log.Println(route+" added as route")
    }
}

func main() {
    // Setting up DB
    host := "postgres-tapit"
    db, err := gorm.Open("postgres", "sslmode=disable host=" + host + " port=5432 user=tapit dbname=tapit password=secret-tapit-password")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // DB Migrations
    db.AutoMigrate(&GlobalSettings{})
    db.AutoMigrate(&Session{})
    db.AutoMigrate(&User{})
    db.AutoMigrate(&TextTemplate{})
    db.AutoMigrate(&WebTemplate{})
    db.AutoMigrate(&TwilioProvider{})
    db.AutoMigrate(&Phonebook{})
    db.AutoMigrate(&PhoneRecord{})
    db.AutoMigrate(&Campaign{})
    db.AutoMigrate(&Job{})
    db.AutoMigrate(&Visit{})

    // Setting up Tapit app
    var tapit Tapit
    tapit.db = db

    var globalSettings GlobalSettings

    // handle global settings
    err = tapit.db.Last(&globalSettings).Error
    if err != nil {
        globalSettings.SecretRegistrationCode = "Super-Secret-Code"
        globalSettings.ThreadsPerCampaign = 2
        globalSettings.BcryptCost = 12
        globalSettings.MaxRequestRetries = 5
        globalSettings.WaitBeforeRetry = 1000
        globalSettings.WebTemplatePrefix = "https://www.attacker.com/"
        globalSettings.WebFrontPlaceholder = ""

        tapit.db.NewRecord(&globalSettings)
        tapit.db.Create(&globalSettings)
    }

    tapit.globalSettings = globalSettings

    // Seeding random
    rand.Seed(time.Now().UnixNano())

    // Clear running campaigns & starting background jobs
    tapit.clearRunningCampaigns()
    go tapit.workerTwilioChecker()
    tapit.campaignChan = make(chan CampaignComms, 10)

    // Setting up mux
    r := mux.NewRouter()

    // Get current dir
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }

    // Setting up static routes (frontend)
    iterateStatic(r, dir + "/static/", "/")
    routes := []string{
        "/",
        "/login",
        "/register",
        "/profile",
        "/campaign",
        "/campaign/new",
        "/campaign/{id}/view",
        "/phonebook",
        "/phonebook/new",
        "/phonebook/{id}/edit",
        "/text-template",
        "/text-template/new",
        "/text-template/{id}/edit",
        "/web-template",
        "/web-template/new",
        "/web-template/{id}/edit",
        "/global-settings",
        "/provider",
    }
    indexPath := dir + "/static/index.html"
    generateRoutes(r, indexPath, routes)

    // Setting up API routes
    r.HandleFunc("/api/login", tapit.login)
    r.HandleFunc("/api/logout", tapit.logout)
    r.HandleFunc("/api/register", tapit.register)
    r.HandleFunc("/api/myself", tapit.authenticationHandler(tapit.myselfHandler))

    r.Handle("/api/text-template",tapit.authenticationHandler(tapit.handleTextTemplate))
    r.Handle("/api/text-template/{id}",tapit.authenticationHandler(tapit.handleSpecificTextTemplate))
    r.Handle("/api/web-template",tapit.authenticationHandler(tapit.handleWebTemplate))
    r.Handle("/api/web-template/{id}",tapit.authenticationHandler(tapit.handleSpecificWebTemplate))
    r.Handle("/api/provider/twilio",tapit.authenticationHandler(tapit.handleTwilioProvider))
    r.Handle("/api/phonebook",tapit.authenticationHandler(tapit.handlePhonebook))
    r.Handle("/api/phonebook/{id}",tapit.authenticationHandler(tapit.handleSpecificPhonebook))
    r.Handle("/api/import-phonebook",tapit.authenticationHandler(tapit.importPhonebook))
    r.Handle("/api/campaign",tapit.authenticationHandler(tapit.handleCampaign))
    r.Handle("/api/campaign/{id}",tapit.authenticationHandler(tapit.handleSpecificCampaign))
    r.Handle("/api/campaign/{id}/start",tapit.authenticationHandler(tapit.handleStartCampaign))
    r.Handle("/api/campaign/{id}/pause",tapit.authenticationHandler(tapit.handleStopCampaign))
    r.Handle("/api/globalsettings",tapit.authenticationHandler(tapit.handleGlobalSettings))
    r.Handle("/api/jobs/{id}/visits",tapit.authenticationHandler(tapit.handleDownloadView))

    // Starting management web server
    r.Handle("/", r)
    log.Println("Starting management web server on port 8000...")
    go http.ListenAndServe(":8000", r)

    // Handle WebTemplate Routes
    webTemplateRouter := mux.NewRouter()
    webTemplateRouter.HandleFunc("/", tapit.handleWebFront)
    webTemplateRouter.HandleFunc("/{route}", tapit.webTemplateRouteHandler)

    // Starting victim route web server
    webTemplateRouter.Handle("/", webTemplateRouter)
    log.Println("Starting victim routes on port 8001...")
    http.ListenAndServe(":8001", webTemplateRouter)

}
