package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Benny66/tally-server-cloud/db"
	"github.com/Benny66/tally-server-cloud/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	// http.HandleFunc("/", service.IndexHandler)
	// http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/api/auth-login", service.AuthLogin)
	http.HandleFunc("/api/benediction", service.Benediction)
	// http.HandleFunc("upload-file", service.UploadFile)
	http.HandleFunc("/api/get-weather", service.GetWeather)

	//
	http.HandleFunc("/api/get-user-info", service.GetUserInfo)
	http.HandleFunc("/api/set-user-info", service.SetUserInfo)

	http.HandleFunc("/api/get-user-book", service.GetUserBooks)
	http.HandleFunc("/api/get-book-info", service.GetUserBookInfo)
	http.HandleFunc("/api/set-book-edit", service.SetBookEdit)

	http.HandleFunc("/api/get-category-list", service.GetCategoryList)
	http.HandleFunc("/api/set-category-edit", service.SetCategoryEdit)

	http.HandleFunc("/api/get-tally-main-total", service.GetTallyMainTotal)
	http.HandleFunc("/api/get-tally-main-list", service.GetTallyMainList)

	http.HandleFunc("/api/set-tally-edit", service.SetTallyMainEdit)
	http.HandleFunc("/api/get-tally-info", service.GetTallyMainInfo)
	http.HandleFunc("/api/get-tally-sta", service.GetTallySta)

	log.Fatal(http.ListenAndServe(":80", nil))
}


