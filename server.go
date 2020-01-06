package main

import (
    "io"
    "os"
    "fmt"
    "log"
    "path"
    "net/http"
    "encoding/json"
    "./utils"
)

// ArticleDB global array to simulate a database
var (
    ArticleDB []utils.Article
    Configs utils.Configurations
)

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!\n")
    log.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
    w.Header().Add("Content-Type", "application/json")
    log.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(ArticleDB)
}

func uploadFile(w http.ResponseWriter, r*http.Request){
    // curl -F 'file=@notes.txt' localhost:5678/upload
    log.Println("method:", r.Method)
    r.ParseMultipartForm(32 << 20)
    file, handler, err := r.FormFile("file")
    if err != nil {
        log.Println(err)
        w.WriteHeader(404)
        w.Header().Add("Content-Type", "application/json")
        json.NewEncoder(w).Encode(utils.ErrorMsg{Code: 404, Msg: "File not found"})
        return
    }
    defer file.Close()
    fmt.Fprintf(w, "%v", handler.Header)
    f, err := os.OpenFile(path.Join(Configs.UploadFolder, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        log.Println(err)
        return
    }
    defer f.Close()
    io.Copy(f, file)
}

func handleRequests() {
    url := fmt.Sprintf("%s:%d", Configs.ServerURL, Configs.ServerPort)
    http.HandleFunc("/", homePage)
    // add our articles route and map it to our returnAllArticles function like so
    http.HandleFunc("/articles", returnAllArticles)
    // upload
    http.HandleFunc("/upload", uploadFile)
    // /public -> serve static files from static folder
    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./static"))))
    log.Fatal(http.ListenAndServe(url, nil))
}

func main() {
    ArticleDB = append(ArticleDB, utils.Article{Title: "Title: 1", Desc: "Article Description", Content: "Article Content"})
    ArticleDB = append(ArticleDB, utils.Article{Title: "Title: 2", Desc: "Article Description", Content: "Article Content"})

    Configs = utils.GetConfigs()
    // check if upload folder exist
    if _,err := os.Stat(path.Join(".", Configs.UploadFolder)); os.IsNotExist(err) {
        os.Mkdir(path.Join(".", Configs.UploadFolder), os.ModePerm)
    }

    log.Printf("Serving API on %[1]s port %[2]d (http://%[1]s:%[2]d/)\n", Configs.ServerURL, Configs.ServerPort)
    handleRequests()
}