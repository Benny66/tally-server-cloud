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
	http.HandleFunc("auth-login", service.AuthLogin)
	http.HandleFunc("benediction", service.Benediction)
	// http.HandleFunc("upload-file", service.UploadFile)
	http.HandleFunc("get-weather", service.GetWeather)

	//
	http.HandleFunc("get-user-info", service.GetUserInfo)
	http.HandleFunc("set-user-info", service.SetUserInfo)

	http.HandleFunc("get-user-book", service.GetUserBooks)
	http.HandleFunc("get-book-info", service.GetUserBookInfo)
	http.HandleFunc("set-book-edit", service.SetBookEdit)

	http.HandleFunc("/api/get-category-list", service.GetCategoryList)
	http.HandleFunc("set-category-edit", service.SetCategoryEdit)

	http.HandleFunc("get-tally-main-total", service.GetTallyMainTotal)
	http.HandleFunc("get-tally-main-list", service.GetTallyMainList)

	http.HandleFunc("set-tally-edit", service.SetTallyMainEdit)
	http.HandleFunc("get-tally-info", service.GetTallyMainInfo)
	http.HandleFunc("get-tally-sta", service.GetTallySta)

	log.Fatal(http.ListenAndServe(":80", nil))
}


