package service

import (
	"errors"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Benny66/tally-server-cloud/db"
	"github.com/Benny66/tally-server-cloud/db/models"
	"github.com/Benny66/tally-server-cloud/schemas"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/util"
)

func Auth(r *http.Request) (*models.UserModel, error) {
	token := r.Header.Get("token")
	if token == "" {
		return nil, errors.New("token is empty")
	}
	userInfo, err := models.NewUserDao().FindOneWhere("token", token)
	if err != nil {
		return nil, errors.New("token is not found")
	}
	if userInfo.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &userInfo, nil
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	var req = schemas.UserAuthLoginApiReq{
		Code: r.FormValue("code"),
	}

	if req.Code == "" {
		NewResponseJson(w).Error(2)
		return
	}

	wx := wechat.NewWechat()
	memory := cache.NewMemcache()
	cfg := &miniConfig.Config{
		AppID:     os.Getenv("APP_ID"),
		AppSecret: os.Getenv("APP_SECRET"),
		Cache:     memory,
	}
	mini := wx.GetMiniProgram(cfg)
	resCode2Session, err := mini.GetAuth().Code2Session(req.Code)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	token := util.RandomStr(32)
	count, err := models.NewUserDao().FindCountWhere("openid", resCode2Session.OpenID)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	tx := db.Get()
	if count == 0 {
		var userInfo models.UserModel
		userInfo.OpenId = resCode2Session.OpenID
		userInfo.Token = token
		userInfo.NickName = RandomNickname()
		userInfo.AvatarUrl = "http://ru52sfqn4.bkt.clouddn.com/1.png"
		userInfo.Sex = 0
		userInfo.Job = ""
		_, err := models.NewUserDao().Create(tx, &userInfo)
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	} else {
		userInfo, err := models.NewUserDao().FindOneWhere("openid", resCode2Session.OpenID)
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
		_, err = models.NewUserDao().Update(tx, userInfo.ID, map[string]interface{}{
			"token": token,
		})
		if err != nil {
			NewResponseJson(w).Error(1, err.Error())
			return
		}
	}
	NewResponseJson(w).Success(schemas.UserAuthLoginApiRes{
		Token: token,
	})
	return
}

func RandomNickname() string {
	var adjectives = []string{"Happy", "Silly", "Funny", "Crazy", "Smart", "Brave", "Gentle", "Lucky", "Sleepy", "Charming"}
	var nouns = []string{"Penguin", "Kangaroo", "Elephant", "Giraffe", "Tiger", "Monkey", "Lion", "Koala", "Panda", "Zebra"}

	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return adj + noun
}

func Benediction(w http.ResponseWriter, r *http.Request) {
	id := rand.Intn(5) + 1
	phraseInfo, err := models.NewPhraseDao().FindOneWhere("id", id)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	NewResponseJson(w).Success(phraseInfo)
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
	NewResponseJson(w).Success("http://ru52sfqn4.bkt.clouddn.com/clear_day.png")
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	NewResponseJson(w).Success(userInfo)
}

func SetUserInfo(w http.ResponseWriter, r *http.Request) {
	userInfo, err := Auth(r)
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}

	var req = schemas.SetUserInfoApiReq{
		Nickname:  r.FormValue("nick_name"),
		AvatarUrl: r.FormValue("avatar_url"),
		Job:       r.FormValue("job"),
	}
	sexStr := r.FormValue("sex")
	req.Sex, _ = strconv.Atoi(sexStr)
	tx := db.Get()

	_, err = models.UserDao.Update(tx, userInfo.ID, map[string]interface{}{
		"nick_name":  req.Nickname,
		"avatar_url": req.AvatarUrl,
		"job":        req.Job,
		"sex":        req.Sex,
	})
	if err != nil {
		NewResponseJson(w).Error(1, err.Error())
		return
	}
	userInfo.NickName = req.Nickname
	userInfo.AvatarUrl = req.AvatarUrl
	userInfo.Job = req.Job
	userInfo.Sex = req.Sex
	NewResponseJson(w).Success(userInfo)
}
