package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
)

type TemplateData struct {
	Resp string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := r.URL.Query().Get("resp")
		data := &TemplateData{Resp: resp}
		tmpl, _ := template.ParseFiles("template.html")
		tmpl.Execute(w, data)
	})
	http.HandleFunc("/settings", GetSettings)
	http.HandleFunc("/state_instance", GetStateInstance)
	http.HandleFunc("/send_message", SendMessage)
	http.HandleFunc("/send_file", SendFile)
	http.ListenAndServe(":8080", nil)
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idInstance := r.FormValue("idInstance")
	token := r.FormValue("apiToken")
	link := fmt.Sprintf("https://7103.api.greenapi.com/waInstance%s/getSettings/%s", idInstance, token)
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
	fmt.Println(string(resp))

	encodedResp := url.QueryEscape(string(resp))
	http.Redirect(w, r, "/?resp="+encodedResp, http.StatusSeeOther)
}

func GetStateInstance(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	idInstance := r.FormValue("idInstance")
	token := r.FormValue("apiToken")
	link := fmt.Sprintf("https://7103.api.greenapi.com/waInstance%s/getStateInstance/%s", idInstance, token)
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
	fmt.Println(string(resp))

	encodedResp := url.QueryEscape(string(resp))
	http.Redirect(w, r, "/?resp="+encodedResp, http.StatusSeeOther)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {

}

func SendFile(w http.ResponseWriter, r *http.Request) {
	// Ваша логика для SendFile
}
