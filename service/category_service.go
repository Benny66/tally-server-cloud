package service

import (
	"net/http"
	"strconv"

	"github.com/Benny66/tally-server-cloud/db"
	"github.com/Benny66/tally-server-cloud/db/models"
	"github.com/Benny66/tally-server-cloud/schemas"
)

func GetCategoryList(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}
	typeStr := r.PostFormValue("type")
	// type, _ := strconv.Atoi(typeStr)
	categories, err := models.NewCategoryDao().FindAllWhere("type = ? and is_del = ? and (user_id = 0 or user_id = ?)", typeStr, 0, userInfo.ID)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	NewResponseJson(w).Success(categories)
}

func SetCategoryEdit(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}

	var req = schemas.SetCategoryApiReq{
		Name: r.FormValue("name"),
	}
	typeStr := r.FormValue("type")
	req.Type, _ = strconv.Atoi(typeStr)
	if req.Name == "" {
		NewResponseJson(w).Error(2)
		return
	}
	idStr := r.FormValue("id")
	req.Id, _ = strconv.Atoi(idStr)
	tx := db.Get()
	if req.Id == 0 {
		category := models.CategoryModel{
			UserId:  userInfo.ID,
			Name:    req.Name,
			CType:   req.Type,
			IconUrl: "http://ru52sfqn4.bkt.clouddn.com/mot1.png",
		}
		_, err := models.CategoryDao.Create(tx, &category)
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	} else {
		_, err := models.NewCategoryDao().Update(tx, uint(req.Id), map[string]interface{}{
			"name": req.Name,
			"type": req.Type,
		})
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	}
	NewResponseJson(w).Success("")
}
