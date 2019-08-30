package utils

// Article json structure of a article
type Article struct {
    Title string `json:"Title"`
    Desc string `json:"desc"`
    Content string `json:"content"`
}
// ErrorMsg asdasd
type ErrorMsg struct {
    Code int `json:"Code"`
    Msg string `json:"Msg"`
}