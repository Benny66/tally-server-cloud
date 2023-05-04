package service

import (
	"net/http"
	"strconv"

	"github.com/Benny66/tally-server-cloud/db"
	"github.com/Benny66/tally-server-cloud/db/models"
	"github.com/Benny66/tally-server-cloud/schemas"
)

func GetUserBooks(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}

	books, err := models.NewBookDao().FindAllWhere("user_id", userInfo.ID)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	NewResponseJson(w).Success(books)
}

func GetUserBookInfo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PostFormValue("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		NewResponseJson(w).Error(2)
		return
	}

	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}

	book, err := models.NewBookDao().FindOneWhere("user_id = ? and id = ?", userInfo.ID, id)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	NewResponseJson(w).Success(book)
}
func SetBookEdit(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}

	var req = schemas.SetBookEditApiReq{
		Name: r.FormValue("name"),
	}
	sortStr := r.FormValue("sort")
	req.Sort, _ = strconv.Atoi(sortStr)
	if req.Name == "" {
		NewResponseJson(w).Error(2)
		return
	}
	idStr := r.FormValue("id")
	req.Id, _ = strconv.Atoi(idStr)
	tx := db.Get()
	if req.Id == 0 {
		book := models.BookModel{
			UserId: userInfo.ID,
			Name:   req.Name,
			Sort:   req.Sort,
		}
		_, err := models.BookDao.Create(tx, &book)
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	} else {
		_, err := models.BookDao.Update(tx, uint(req.Id), map[string]interface{}{
			"name": req.Name,
			"sort": req.Sort,
		})
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	}
	NewResponseJson(w).Success("")
}
