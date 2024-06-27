package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template.html")
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
	fmt.Println("DATA", idInstance, token)
	url := fmt.Sprintf("https://7103.api.greenapi.com/waInstance%s/getSettings/%s", idInstance, token)
	body, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := io.ReadAll(body.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resp))

}
func GetStateInstance(w http.ResponseWriter, r *http.Request) {

}

func SendMessage(w http.ResponseWriter, r *http.Request) {

}

func SendFile(w http.ResponseWriter, r *http.Request) {

}
