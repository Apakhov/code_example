package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type dataUser struct {
	ID        int    `json:"id" xml:"id"`
	Age       int    `json:"age" xml:"age"`
	FirstName string `json:"first_name" xml:"first_name"`
	LastName  string `json:"last_name" xml:"last_name"`
	Gender    string `json:"gender" xml:"gender"`
	About     string `json:"about" xml:"about"`
}

type dataUsers struct {
	List []dataUser `xml:"row"`
}

const rightAccessToken = "wrongAccessToken"

func (du dataUser) convertToUser() User {
	return User{
		Id:     du.ID,
		Name:   du.FirstName + du.LastName,
		Age:    du.Age,
		Gender: du.Gender,
		About:  du.About,
	}
}

func (dus dataUsers) convertToUsers() []User {
	users := make([]User, 0)
	for _, du := range dus.List {
		users = append(users, du.convertToUser())
	}
	return users
}

func filter(users []User, query string, limit int, offset int) []User {
	res := make([]User, 0)
	for _, user := range users {
		if offset > 0 {
			offset--
			continue
		}
		if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
			if limit == 0 {
				return res
			}
			res = append(res, user)
			limit--
		}
	}
	return res
}

type Server struct {
	rightAccessToken string
	datasetPath      string
}

func (s Server) SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		res, _ := json.Marshal(SearchErrorResponse{"bad method"})
		http.Error(w, string(res), http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Accesstoken") != s.rightAccessToken {
		res, _ := json.Marshal(SearchErrorResponse{"bad Accesstoken"})
		http.Error(w, string(res), http.StatusUnauthorized)
		return
	}

	xmlData, err := ioutil.ReadFile(s.datasetPath) //"99_hw/coverage/dataset.xml"
	if err != nil {
		log.Println("cant get dataset.xml", err.Error())
		res, _ := json.Marshal(SearchErrorResponse{err.Error()})
		http.Error(w, string(res), http.StatusInternalServerError)
		return
	}

	data := new(dataUsers)
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		res, _ := json.Marshal(SearchErrorResponse{err.Error()})
		http.Error(w, string(res), http.StatusInternalServerError)
		return
	}
	query := r.URL.Query().Get("query")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		res, _ := json.Marshal(SearchErrorResponse{err.Error()})
		http.Error(w, string(res), http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		res, _ := json.Marshal(SearchErrorResponse{err.Error()})
		http.Error(w, string(res), http.StatusBadRequest)
		return
	}
	users := filter(data.convertToUsers(), query, limit, offset)

	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		res, _ := json.Marshal(SearchErrorResponse{err.Error()})
		http.Error(w, string(res), http.StatusBadRequest)
		return
	}
	if orderBy != OrderByAsc && orderBy != OrderByAsIs && orderBy != OrderByDesc {
		res, _ := json.Marshal(SearchErrorResponse{"ErrorBadOrderField"})
		http.Error(w, string(res), http.StatusBadRequest)
		return
	}
	var usersToSort []User
	if len(users) == limit {
		usersToSort = users[0 : len(users)-1]
	} else {
		usersToSort = users[0:len(users)]
	}
	if orderBy != OrderByAsIs {
		orderField := r.URL.Query().Get("order_field")
		switch orderField {
		case "Id":
			sort.Slice(usersToSort, func(i int, j int) bool { return usersToSort[i].Id < usersToSort[j].Id })
		case "Age":
			sort.Slice(usersToSort, func(i int, j int) bool { return usersToSort[i].Age < usersToSort[j].Age })
		case "Name":
			sort.Slice(usersToSort, func(i int, j int) bool { return usersToSort[i].Name < usersToSort[j].Name })
		case "":
			sort.Slice(usersToSort, func(i int, j int) bool { return usersToSort[i].Name < usersToSort[j].Name })
		default:
			res, _ := json.Marshal(SearchErrorResponse{"ErrorBadOrderField"})
			http.Error(w, string(res), http.StatusBadRequest)
			return
		}
		if orderBy == OrderByAsc {
			for i := 0; i < len(usersToSort)/2; i++ {
				usersToSort[i], usersToSort[len(usersToSort)-1-i] = usersToSort[len(usersToSort)-1-i], usersToSort[i]
			}
		}
	}

	res, _ := json.Marshal(users)
	w.Write(res)
	return
}
