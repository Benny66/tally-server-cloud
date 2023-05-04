package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NewResponseJson(w http.ResponseWriter) *responseJson {
	return &responseJson{
		w: w,
	}
}

type responseJson struct {
	w http.ResponseWriter
}

/*
* description: 设置响应头
* author: shahao
* created on: 19-11-19 下午3:12
* param param_1:
* param param_2:
* return return_1:
 */
func (r *responseJson) SetHeader(key, value string) *responseJson {
	r.w.Header().Set(key, value)
	return r
}

/*
* description: 成功返回数据构造
* author: shahao
* created on: 19-11-19 下午2:17
* param data: 返回的数据
* param message: 返回提示信息
* return :
 */
func (r *responseJson) Success(data interface{}) {
	var res = ResultData{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.w.Write(msg)
}

/*
* description: 错误返回数据构造
* author: shahao
* created on: 19-11-19 下午2:17
* param data: 返回的数据
* param message: 返回提示信息
* return :
 */
func (r *responseJson) Error(errorCode int, params ...interface{}) {
	res := ResultData{
		Code: errorCode,
		Msg:  "error",
		Data: "",
	}
	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	httpCode := http.StatusInternalServerError
	switch errorCode {
	case 8:
		httpCode = http.StatusUnauthorized
	}
	r.w.WriteHeader(httpCode)
	r.w.Write(msg)
}


type ResultData struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}
