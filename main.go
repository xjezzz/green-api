package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	APIBaseURL = "https://7103.api.greenapi.com/waInstance%s/%s/%s"
)

type TemplateData struct {
	Resp string
}
type Message struct {
	ChatId  string `json:"chatId"`
	Message string `json:"message"`
}
type File struct {
	ChatId          string `json:"chatId"`
	Url             string `json:"urlFile"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption"`
	QuotedMessageId string `json:"quotedMessageId"`
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/settings", GetSettings)
	http.HandleFunc("/state_instance", GetStateInstance)
	http.HandleFunc("/send_message", SendMessage)
	http.HandleFunc("/send_file", SendFileByUrl)
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	resp := r.URL.Query().Get("resp")
	data := &TemplateData{Resp: resp}
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, "getSettings")

}

func GetStateInstance(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, "getStateInstance")
}
func SendMessage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idInstance := r.FormValue("idInstance")
	token := r.FormValue("apiToken")
	phone := r.FormValue("phoneNumber1")
	text := r.FormValue("message")

	if idInstance == "" || token == "" || phone == "" || text == "" {
		http.Error(w, "All fields must be filled", http.StatusBadRequest)
		return
	}
	message := Message{ChatId: phone + "@c.us", Message: text}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	link := fmt.Sprintf("https://7103.api.greenapi.com/waInstance%s/sendMessage/%s", idInstance, token)
	body, err := http.Post(link, "application/json", bytes.NewReader(jsonMessage))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer body.Body.Close()
	resp, err := io.ReadAll(body.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	encodedResp := url.QueryEscape(string(resp))
	http.Redirect(w, r, "/?resp="+encodedResp, http.StatusSeeOther)
}

func SendFileByUrl(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idInstance := r.FormValue("idInstance")
	token := r.FormValue("apiToken")
	phone := r.FormValue("phoneNumber2")
	fileUrl := r.FormValue("fileUrl")

	if idInstance == "" || token == "" || phone == "" || fileUrl == "" {
		http.Error(w, "All fields must be filled", http.StatusBadRequest)
		return
	}
	ext, err := GetFileContentType(fileUrl)
	if err != nil {
		fmt.Println(err)
		return

	}
	message := File{ChatId: phone + "@c.us", Url: fileUrl, FileName: ext}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	link := fmt.Sprintf("https://7103.api.greenapi.com/waInstance%s/sendFileByUrl/%s", idInstance, token)
	body, err := http.Post(link, "application/json", bytes.NewReader(jsonMessage))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer body.Body.Close()
	resp, err := io.ReadAll(body.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	encodedResp := url.QueryEscape(string(resp))
	http.Redirect(w, r, "/?resp="+encodedResp, http.StatusSeeOther)
}

func handleRequest(w http.ResponseWriter, r *http.Request, endpoint string) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idInstance := r.FormValue("idInstance")
	token := r.FormValue("apiToken")

	if idInstance == "" || token == "" {
		http.Error(w, "All fields must be filled", http.StatusBadRequest)
		return
	}

	link := fmt.Sprintf(APIBaseURL, idInstance, endpoint, token)
	body, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer body.Body.Close()
	resp, err := io.ReadAll(body.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	encodedResp := url.QueryEscape(string(resp))
	http.Redirect(w, r, "/?resp="+encodedResp, http.StatusSeeOther)
}

func GetFileContentType(url string) (string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return "", fmt.Errorf("No Content-Type header found")
	}

	mimeType := strings.Split(contentType, ";")[0]

	return mimeType, nil
}
