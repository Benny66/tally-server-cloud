package service

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Benny66/tally-server-cloud/db"
	"github.com/Benny66/tally-server-cloud/db/models"
	"github.com/Benny66/tally-server-cloud/schemas"
)

func GetTallyMainTotal(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}
	var req = schemas.GetTallyMainListApiReq{
		Date: r.FormValue("date"),
	}

	bookIdStr := r.FormValue("book_id")
	req.BookId, _ = strconv.Atoi(bookIdStr)
	var query string = fmt.Sprintf("user_id = %d ", userInfo.ID)
	if req.BookId != 0 {
		query += fmt.Sprintf(" and book_id = %d ", req.BookId)
	}
	if req.Date != "" {
		dateStartTime, err := Parse(fmt.Sprintf("%s 00:00:00", req.Date), "2006-01-02 15:04:05")
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
		dateStartTime = time.Date(dateStartTime.Year(), dateStartTime.Month(), 1, 0, 0, 0, 0, dateStartTime.Location())
		dateEndTime := dateStartTime.AddDate(0, 1, 0)
		query += fmt.Sprintf("and date >= '%s' and date < '%s'", dateStartTime.String(), dateEndTime.String())
	}
	if req.Date == "" && req.StartTime != "" && req.EndTime != "" {
		req.StartTime = fmt.Sprintf("%s 00:00:00", req.StartTime)
		req.EndTime = fmt.Sprintf("%s 23:59:59", req.EndTime)
		query += fmt.Sprintf("and date >= '%s' and date <= '%s'", req.StartTime, req.EndTime)
	}
	tallies, err := models.MainDao.FindAllWhere(query)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	var res = schemas.GetTallyMainTotalApiRes{}
	for _, v := range tallies {
		switch v.MType {
		case 1:
			res.Expend += v.Money
		case 2:
			res.Income += v.Money
		case 3:
			res.Disregard += v.Money
		}
	}
	NewResponseJson(w).Success(res)
}
func GetTallyMainList(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}
	var req = schemas.GetTallyMainListApiReq{}

	bookIdStr := r.FormValue("book_id")
	req.BookId, _ = strconv.Atoi(bookIdStr)
	categoryIdStr := r.FormValue("category_id")
	req.CategoryId, _ = strconv.Atoi(categoryIdStr)
	typeStr := r.FormValue("type")
	req.Type, _ = strconv.Atoi(typeStr)

	req.StartTime = r.PostFormValue("start_time")
	req.EndTime = r.PostFormValue("end_time")
	var query string = fmt.Sprintf("user_id = %d ", userInfo.ID)
	if req.BookId != 0 {
		query += fmt.Sprintf(" and book_id = %d ", req.BookId)
	}
	if req.CategoryId != 0 {
		query += fmt.Sprintf("and category_id = %d ", req.CategoryId)
	}
	if req.Type != 0 {
		query += fmt.Sprintf("and type = %d ", req.Type)
	}
	if req.StartTime != "" && req.EndTime != "" {
		req.StartTime = fmt.Sprintf("%s 00:00:00", req.StartTime)
		req.EndTime = fmt.Sprintf("%s 23:59:59", req.EndTime)
		query += fmt.Sprintf("and date >= '%s' and date <= '%s'", req.StartTime, req.EndTime)
	}

	tallies, err := models.MainDao.FindAllWhere(query)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	bookIds := []uint{}
	bookIdMap := make(map[uint]bool, 0)
	categoryIds := []uint{}
	categoryIdMap := make(map[uint]bool, 0)
	for _, tally := range tallies {
		if _, ok := bookIdMap[tally.BookId]; !ok {
			bookIds = append(bookIds, tally.BookId)
		}
		if _, ok := categoryIdMap[tally.BookId]; !ok {
			categoryIds = append(categoryIds, tally.CategoryId)
		}
	}
	books, err := models.BookDao.FindAllWhere("id in (?)", bookIds)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	bookMap := make(map[uint]models.BookModel, 0)

	for _, book := range books {
		bookMap[book.ID] = book
	}
	categories, err := models.CategoryDao.FindAllWhere("id in (?)", categoryIds)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	categoryMap := make(map[uint]models.CategoryModel, 0)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}
	var tallyList []schemas.GetTallyMainListApiRes
	for _, tally := range tallies {
		tallyRes := schemas.GetTallyMainListApiRes{
			ID:         tally.ID,
			UserId:     tally.UserId,
			BookId:     tally.BookId,
			CategoryId: tally.CategoryId,
			MType:      tally.MType,
			Money:      tally.Money,
			Name:       tally.Name,
			Desc:       tally.Desc,
			Date:       tally.Date.String(),
			IsDel:      tally.IsDel,
			CreatedAt:  tally.CreatedAt.String(),
		}
		tallyRes.Book = bookMap[tallyRes.BookId]
		tallyRes.Category = categoryMap[tallyRes.CategoryId]
		tallyList = append(tallyList, tallyRes)
	}
	NewResponseJson(w).Success(tallyList)
}

func SetTallyMainEdit(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}

	var req = schemas.SetTallyMainEditApiReq{
		Name: r.FormValue("name"),
		Desc: r.FormValue("desc"),
		Date: r.FormValue("date"),
	}
	idStr := r.FormValue("id")
	req.Id, _ = strconv.Atoi(idStr)
	bookIdStr := r.FormValue("book_id")
	req.BookId, _ = strconv.Atoi(bookIdStr)
	categoryIdStr := r.FormValue("category_id")
	req.CategoryId, _ = strconv.Atoi(categoryIdStr)
	typeStr := r.FormValue("type")
	req.Type, _ = strconv.Atoi(typeStr)
	moneySrr := r.FormValue("money")
	req.Money, _ = strconv.ParseFloat(moneySrr, 64)
	if req.BookId == 0 || req.CategoryId == 0 || req.Type == 0 || req.Money == 0 {
		NewResponseJson(w).Error(2)
		return
	}
	tx := db.Get()
	var main models.MainModel
	if req.Id == 0 {
		main = models.MainModel{
			UserId:     userInfo.ID,
			Name:       req.Name,
			BookId:     uint(req.BookId),
			CategoryId: uint(req.CategoryId),
			MType:      req.Type,
			Money:      req.Money,
			Desc:       req.Desc,
		}
		dateTime, err := Parse(req.Date, "2006-01-02")
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
		main.Date = models.ModelTime(dateTime)
		_, err = models.MainDao.Create(tx, &main)
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	} else {
		_, err := models.MainDao.Update(tx, uint(req.Id), map[string]interface{}{
			"name":  req.Name,
			"type":  req.Type,
			"money": req.Money,
			"desc":  req.Desc,
			"date":  req.Date,
		})
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	}
	NewResponseJson(w).Success(main)

}

func GetTallyMainInfo(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(8, err.Error())
		return
	}
	id := r.PostFormValue("id")
	tally, err := models.MainDao.FindOneWhere("user_id = ? and  id = ? ", userInfo.ID, id)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	NewResponseJson(w).Success(tally)
}

func GetTallySta(w http.ResponseWriter, r *http.Request) {
	typeStr := r.FormValue("type")
	date := r.FormValue("date")
	var query string = "1 = 1 "
	if typeStr != "" {
		query += fmt.Sprintf("and type = %s ", typeStr)
	}
	if date != "" {
		dateStartTime, err := Parse(fmt.Sprintf("%s 00:00:00", date), "2006-01-02 15:04:05")
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
		dateStartTime = time.Date(dateStartTime.Year(), dateStartTime.Month(), 1, 0, 0, 0, 0, dateStartTime.Location())
		dateEndTime := dateStartTime.AddDate(0, 1, 0)

		fmt.Println(dateStartTime, dateEndTime)

		query += fmt.Sprintf("and date >= '%s' and date < '%s'", dateStartTime.String(), dateEndTime.String())
	}
	mainSta, err := models.MainDao.StaWhere(query)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	var res []schemas.GetTallyStaRes
	for _, item := range mainSta {
		res = append(res, schemas.GetTallyStaRes{
			NickName:  item.NickName,
			AvatarUrl: item.AvatarUrl,
			Pay:       item.Pay,
			Time:      item.Time.String(),
		})
	}
	NewResponseJson(w).Success(res)
}

func Parse(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}
