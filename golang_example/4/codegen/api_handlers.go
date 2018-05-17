package main

import (
		"encoding/json"
		"net/http"
		"strconv"
	)

func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		h.handlerProfile(w, r)
	case "/user/create":
		h.handlerCreate(w, r)
	default: 
		mp := make(map[string]string)
		mp["error"] = "unknown method"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 404)
		return;
	}
}

func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/create":
		h.handlerCreate(w, r)
	default: 
		mp := make(map[string]string)
		mp["error"] = "unknown method"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 404)
		return;
	}
}

func (h *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	var params ProfileParams
	var err error
// getting and validation
	var raw string
	if r.Method == http.MethodPost{
		raw = r.PostFormValue("login")
		if raw == ""{

		mp := make(map[string]string)
		mp["error"] = "login must be not empty"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Login = raw
	}else if r.Method == http.MethodGet{
		raw = r.URL.Query().Get("login")
		if raw == ""{

		mp := make(map[string]string)
		mp["error"] = "login must be not empty"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Login = raw
	}
res, err := h.Profile(nil, params)
if err != nil {
				mp := make(map[string]string)
				aerr, ok := err.(ApiError)
				if ok {
					mp["error"] = aerr.Err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), err.(ApiError).HTTPStatus)
				} else {
					mp["error"] = err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), http.StatusInternalServerError)
				}
				return
			}
			mp := make(map[string]interface{})
			mp["response"] = (*res)
			mp["error"] = ""
			response, _ := json.Marshal(mp)
			w.Write(response)
}

func (h *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		mp := make(map[string]string)
		mp["error"] = "bad method"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 406)
		return;
}

	if r.Header.Get("X-Auth") != "100500" {

		mp := make(map[string]string)
		mp["error"] = "unauthorized"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 403)
		return;

	}
	var params CreateParams
	var err error
// getting and validation
	var raw string
	if r.Method == http.MethodPost{
		raw = r.PostFormValue("login")
		if raw == ""{

		mp := make(map[string]string)
		mp["error"] = "login must be not empty"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if len(raw) < 10{

		mp := make(map[string]string)
		mp["error"] = "login len must be >= 10"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Login = raw
		raw = r.PostFormValue("full_name")
		params.Name = raw
		raw = r.PostFormValue("status")
		if raw == ""{
raw = "user"
}
flag := false
statusSl := []string{ "user", "moderator", "admin",}
	for _, i := range statusSl{
		if i == raw{flag = true}
	}
	if !flag{

		mp := make(map[string]string)
		mp["error"] = "status must be one of [user, moderator, admin]"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
	}
		params.Status = raw
		raw = r.PostFormValue("age")
tempInt, err := strconv.Atoi(raw)
		if err != nil{
		mp := make(map[string]string)
		mp["error"] = "age must be int"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;

		}
	if tempInt > 128{

		mp := make(map[string]string)
		mp["error"] = "age must be <= 128"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if tempInt < 0{

		mp := make(map[string]string)
		mp["error"] = "age must be >= 0"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Age = tempInt
	}else if r.Method == http.MethodGet{
		raw = r.URL.Query().Get("login")
		if raw == ""{

		mp := make(map[string]string)
		mp["error"] = "login must be not empty"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if len(raw) < 10{

		mp := make(map[string]string)
		mp["error"] = "login len must be >= 10"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Login = raw
		raw = r.URL.Query().Get("full_name")
		params.Name = raw
		raw = r.URL.Query().Get("status")
		if raw == ""{
raw = "user"
}
flag := false
statusSl := []string{ "user", "moderator", "admin",}
	for _, i := range statusSl{
		if i == raw{flag = true}
	}
	if !flag{

		mp := make(map[string]string)
		mp["error"] = "status must be one of [user, moderator, admin]"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
	}
		params.Status = raw
		raw = r.URL.Query().Get("age")
tempInt, err := strconv.Atoi(raw)
		if err != nil{
		mp := make(map[string]string)
		mp["error"] = "age must be int"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;

		}
	if tempInt > 128{

		mp := make(map[string]string)
		mp["error"] = "age must be <= 128"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if tempInt < 0{

		mp := make(map[string]string)
		mp["error"] = "age must be >= 0"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Age = tempInt
	}
res, err := h.Create(nil, params)
if err != nil {
				mp := make(map[string]string)
				aerr, ok := err.(ApiError)
				if ok {
					mp["error"] = aerr.Err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), err.(ApiError).HTTPStatus)
				} else {
					mp["error"] = err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), http.StatusInternalServerError)
				}
				return
			}
			mp := make(map[string]interface{})
			mp["response"] = (*res)
			mp["error"] = ""
			response, _ := json.Marshal(mp)
			w.Write(response)
}

func (h *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		mp := make(map[string]string)
		mp["error"] = "bad method"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 406)
		return;
}

	if r.Header.Get("X-Auth") != "100500" {

		mp := make(map[string]string)
		mp["error"] = "unauthorized"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 403)
		return;

	}
	var params OtherCreateParams
	var err error
// getting and validation
	var raw string
	if r.Method == http.MethodPost{
		raw = r.PostFormValue("username")
		if raw == ""{

		mp := make(map[string]string)
		mp["error"] = "username must be not empty"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if len(raw) < 3{

		mp := make(map[string]string)
		mp["error"] = "username len must be >= 3"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Username = raw
		raw = r.PostFormValue("account_name")
		params.Name = raw
		raw = r.PostFormValue("class")
		if raw == ""{
raw = "warrior"
}
flag := false
classSl := []string{ "warrior", "sorcerer", "rouge",}
	for _, i := range classSl{
		if i == raw{flag = true}
	}
	if !flag{

		mp := make(map[string]string)
		mp["error"] = "class must be one of [warrior, sorcerer, rouge]"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
	}
		params.Class = raw
		raw = r.PostFormValue("level")
tempInt, err := strconv.Atoi(raw)
		if err != nil{
		mp := make(map[string]string)
		mp["error"] = "age must be int"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;

		}
	if tempInt > 50{

		mp := make(map[string]string)
		mp["error"] = "level must be <= 50"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if tempInt < 1{

		mp := make(map[string]string)
		mp["error"] = "level must be >= 1"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Level = tempInt
	}else if r.Method == http.MethodGet{
		raw = r.URL.Query().Get("username")
		if raw == ""{

		mp := make(map[string]string)
		mp["error"] = "username must be not empty"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if len(raw) < 3{

		mp := make(map[string]string)
		mp["error"] = "username len must be >= 3"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Username = raw
		raw = r.URL.Query().Get("account_name")
		params.Name = raw
		raw = r.URL.Query().Get("class")
		if raw == ""{
raw = "warrior"
}
flag := false
classSl := []string{ "warrior", "sorcerer", "rouge",}
	for _, i := range classSl{
		if i == raw{flag = true}
	}
	if !flag{

		mp := make(map[string]string)
		mp["error"] = "class must be one of [warrior, sorcerer, rouge]"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
	}
		params.Class = raw
		raw = r.URL.Query().Get("level")
tempInt, err := strconv.Atoi(raw)
		if err != nil{
		mp := make(map[string]string)
		mp["error"] = "age must be int"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;

		}
	if tempInt > 50{

		mp := make(map[string]string)
		mp["error"] = "level must be <= 50"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
	if tempInt < 1{

		mp := make(map[string]string)
		mp["error"] = "level must be >= 1"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), 400)
		return;
}
		params.Level = tempInt
	}
res, err := h.Create(nil, params)
if err != nil {
				mp := make(map[string]string)
				aerr, ok := err.(ApiError)
				if ok {
					mp["error"] = aerr.Err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), err.(ApiError).HTTPStatus)
				} else {
					mp["error"] = err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), http.StatusInternalServerError)
				}
				return
			}
			mp := make(map[string]interface{})
			mp["response"] = (*res)
			mp["error"] = ""
			response, _ := json.Marshal(mp)
			w.Write(response)
}

