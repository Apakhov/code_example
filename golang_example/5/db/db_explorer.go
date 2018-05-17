package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type dbHandler struct {
	DB     *sql.DB
	tables map[string][]Coloumn
}

type Coloumn struct {
	Field      string
	Type       string
	Collation  interface{}
	Null       string
	Key        string
	Default    interface{}
	Extra      string
	Privilages string
	Comment    string
}

func NewDbExplorer(DB *sql.DB) (h dbHandler, err error) {
	x := 666
	defer fmt.Println(x)
	x = 13
	rows, err := DB.Query("SHOW TABLES")
	if err != nil {
		return
	}

	h.tables = make(map[string][]Coloumn)
	for rows.Next() {
		var r string
		rows.Scan(&r)
		h.tables[r] = nil
	}
	rows.Close()
	for r := range h.tables {
		ww, err1 := DB.Query("SHOW FULL COLUMNS FROM " + r)
		if err1 != nil {
			err = err1
			return
		}
		h.tables[r] = make([]Coloumn, 0)
		for ww.Next() {
			var coloumn Coloumn
			err1 := ww.Scan(&coloumn.Field, &coloumn.Type, &coloumn.Collation, &coloumn.Null, &coloumn.Key, &coloumn.Default, &coloumn.Extra, &coloumn.Privilages, &coloumn.Comment)
			if err1 != nil {
				err = err1
				return
			}
			h.tables[r] = append(h.tables[r], coloumn)
		}
		ww.Close()
	}
	h.DB = DB
	return
}

func (h dbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("new request: method = ", r.Method, "url = ", r.URL.String())
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/" {
			h.getList(w, r)
			return
		}
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) == 2 {
			h.getSome(w, r)
			return
		}
		if len(parts) == 3 {
			h.getRecord(w, r)
			return
		}
	case http.MethodPut:
		h.createRecord(w, r)
		return
	case http.MethodPost:
		h.updateRecord(w, r)
		return
	case http.MethodDelete:
		h.deleteRecord(w, r)
		return
	default:
		returnError(w, "bad method", http.StatusTeapot)
	}

	http.Error(w, "500 InternalServerError", http.StatusInternalServerError)
}

func (h dbHandler) getList(w http.ResponseWriter, r *http.Request) {
	list := make([]interface{}, 0)
	for name := range h.tables {
		list = append(list, name)
	}
	prepreres := make(map[string]interface{})
	preres := make(map[string]interface{})
	prepreres["tables"] = list
	preres["response"] = prepreres
	res, _ := json.Marshal(preres)
	w.Write(res)
}

func (h dbHandler) getSome(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Path[1:]
	table, ok := h.tables[tableName]
	if !ok {
		returnError(w, "unknown table", http.StatusNotFound)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		log.Println("limit error", r.URL.Query().Get("limit"))
		limit = 5
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		log.Println("offset error", r.URL.Query().Get("offset"))
		offset = 0
	}

	prepreres := make([]map[string]interface{}, 0)

	query := "SELECT "
	for i, coloumn := range h.tables[tableName] {
		query += coloumn.Field
		if i != len(h.tables[tableName])-1 {
			query += ", "
		} else {
			query += " FROM " + tableName
		}
	}
	rows, err := h.DB.Query(query)
	if err != nil {
		returnError(w, "get some error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	j := 0
	scanned := 0
	for rows.Next() {
		if j < offset {
			j++
			continue
		}
		getter := make([]interface{}, len(table))
		reciever := make([]interface{}, len(table))
		for i := 0; i < len(table); i++ {
			getter[i] = &reciever[i]
		}
		rows.Scan(getter...)
		record := make(map[string]interface{})
		for i, coloumn := range h.tables[tableName] {
			field, ok := (reciever[i]).([]uint8)
			if !ok {
				record[coloumn.Field] = nil
				continue
			}
			if strings.Contains(coloumn.Type, "int") {
				record[coloumn.Field], _ = strconv.Atoi(string([]byte(field)))
			} else if strings.Contains(coloumn.Type, "text") {
				record[coloumn.Field] = string([]byte(field))
			} else if strings.Contains(coloumn.Type, "varchar") {
				record[coloumn.Field] = string([]byte(field))
			} else {
				record[coloumn.Field] = field
			}
		}
		prepreres = append(prepreres, record)
		scanned++
		if scanned == limit {
			break
		}
	}

	preres := make(map[string]interface{})
	preres["response"] = make(map[string]interface{})
	preres["response"].(map[string]interface{})["records"] = prepreres
	res, _ := json.Marshal(preres)
	w.Write(res)

}

func (h dbHandler) getRecord(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	tableName := parts[1]
	table, ok := h.tables[tableName]
	if !ok {
		returnError(w, "unknown table", http.StatusNotFound)
		return
	}
	var idName string
	for _, coloumn := range table {
		if coloumn.Key == "PRI" {
			idName = coloumn.Field
			break
		}
	}

	query := "SELECT "
	for i, coloumn := range h.tables[tableName] {
		query += coloumn.Field
		if i != len(h.tables[tableName])-1 {
			query += ", "
		} else {
			query += " FROM " + tableName + " WHERE " + idName + " = ?"
		}
	}
	rows, err := h.DB.Query(query, parts[2])
	if err != nil {
		returnError(w, "get some error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	record := make(map[string]interface{})
	found := false
	if rows.Next() {
		found = true
		getter := make([]interface{}, len(table))
		reciever := make([]interface{}, len(table))
		for i := 0; i < len(table); i++ {
			getter[i] = &reciever[i]
		}
		rows.Scan(getter...)
		for i, coloumn := range h.tables[tableName] {
			if coloumn.Field == idName {
				record[idName] = reciever[i]
				continue
			}
			field, ok := (reciever[i]).([]uint8)
			if !ok {
				record[coloumn.Field] = nil
				continue
			}
			if strings.Contains(coloumn.Type, "int") {
				record[coloumn.Field], _ = strconv.Atoi(string([]byte(field)))
			} else if strings.Contains(coloumn.Type, "text") {
				record[coloumn.Field] = string([]byte(field))
			} else if strings.Contains(coloumn.Type, "varchar") {
				record[coloumn.Field] = string([]byte(field))
			} else {
				record[coloumn.Field] = field
			}
		}
	}
	if !found {
		returnError(w, "record not found", http.StatusNotFound)
		return
	}
	preres := make(map[string]interface{})
	preres["response"] = make(map[string]interface{})
	preres["response"].(map[string]interface{})["record"] = record
	res, _ := json.Marshal(preres)
	w.Write(res)

}

type inEdQuery struct {
	fields []string
	vals   []interface{}
}

func (h dbHandler) createRecord(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		returnError(w, "wrong URL in PUT method", http.StatusInternalServerError)
		return
	}
	tableName := parts[1]
	_, ok := h.tables[tableName]
	if !ok {
		returnError(w, "unknown table", http.StatusNotFound)
		return
	}

	body := make([]byte, 1024)
	i, err := r.Body.Read(body)
	if err != io.EOF {
		returnError(w, "smth gone wrong with unpacking body", http.StatusTeapot)
		return
	}

	var request map[string]interface{}
	err = json.Unmarshal(body[:i], &request)
	if err != nil {
		returnError(w, "cant unpack json due to "+err.Error(), http.StatusBadRequest)
		return
	}

	requiredColoumns := make([]Coloumn, 0)
	for _, coloumn := range h.tables[tableName] {
		if coloumn.Null == "NO" && coloumn.Default == nil {
			requiredColoumns = append(requiredColoumns, coloumn)
		}
	}
	missing := ""
	for _, coloumn := range requiredColoumns {
		_, ok := request[coloumn.Field]
		if !ok {
			missing = coloumn.Field
			break
		}
	}
	if missing != "" {
		returnError(w, missing+" required", http.StatusBadRequest)
		return
	}

	var cq inEdQuery
	var idName string
	for _, coloumn := range h.tables[tableName] {
		if coloumn.Key == "PRI" {
			idName = coloumn.Field
		}
		temp, ok := request[coloumn.Field]
		if ok {
			if coloumn.Key == "PRI" && strings.Contains(coloumn.Extra, "auto_increment") {
				continue
			}
			cq.fields = append(cq.fields, coloumn.Field)
			if strings.Contains(coloumn.Type, "int") {
				float, ok := temp.(float64)
				if !ok {
					returnError(w, coloumn.Field+" wrong type", http.StatusBadRequest)
					return
				}
				cq.vals = append(cq.vals, int(float))
			} else if strings.Contains(coloumn.Type, "text") {
				str, ok := temp.(string)
				if !ok {
					returnError(w, coloumn.Field+" wrong type", http.StatusBadRequest)
					return
				}
				cq.vals = append(cq.vals, str)
			} else if strings.Contains(coloumn.Type, "varchar") {
				str, ok := temp.(string)
				if !ok {
					returnError(w, coloumn.Field+" wrong type", http.StatusBadRequest)
					return
				}
				cq.vals = append(cq.vals, str)
			} else {
				returnError(w, "putting type error", http.StatusTeapot)
			}
		}
	}
	query := bytes.NewBuffer(make([]byte, 0))
	query.WriteString("INSERT INTO ")
	query.WriteString(tableName)
	query.WriteString(" (")
	for i, field := range cq.fields {
		query.WriteString("`")
		query.WriteString(field)
		query.WriteString("`")
		if i != len(cq.fields)-1 {
			query.WriteString(", ")
		} else {
			query.WriteString(") ")
		}
	}
	query.WriteString("VALUES (")
	for i := range cq.fields {
		query.WriteString("?")
		if i != len(cq.fields)-1 {
			query.WriteString(", ")
		} else {
			query.WriteString(") ")
		}
	}
	result, err := h.DB.Exec(query.String(), cq.vals...)
	if err != nil {
		returnError(w, "putting db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	preres := make(map[string]interface{})
	preres["response"] = make(map[string]interface{})
	preres["response"].(map[string]interface{})[idName], _ = result.LastInsertId()
	res, _ := json.Marshal(preres)
	w.Write(res)

}

func (h dbHandler) updateRecord(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		returnError(w, "wrong URL in POST method", http.StatusInternalServerError)
		return
	}
	tableName := parts[1]
	_, ok := h.tables[tableName]
	if !ok {
		returnError(w, "unknown table", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		returnError(w, "id must be int", http.StatusBadRequest)
		return
	}
	body := make([]byte, 1024)
	i, err := r.Body.Read(body)
	if err != io.EOF {
		returnError(w, "smth gone wrong with unpacking body", http.StatusTeapot)
		return
	}

	var request map[string]interface{}
	err = json.Unmarshal(body[:i], &request)
	if err != nil {
		returnError(w, "cant unpack json due to "+err.Error(), http.StatusBadRequest)
		return
	}

	var idName string
	var cq inEdQuery
	possibleToEdit := make(map[string]Coloumn)
	for _, coloumn := range h.tables[tableName] {
		if coloumn.Key != "PRI" {
			possibleToEdit[coloumn.Field] = coloumn
		} else {
			idName = coloumn.Field
		}
	}

	for i, val := range request {
		coloumn, ok := possibleToEdit[i]
		if !ok {
			returnError(w, "field "+i+" have invalid type", http.StatusBadRequest)
			return
		}
		if val == nil {
			if coloumn.Null == "NO" {
				returnError(w, "field "+i+" have invalid type", http.StatusBadRequest)
				return
			}
			cq.fields = append(cq.fields, coloumn.Field)
			cq.vals = append(cq.vals, nil)
			continue
		}
		cq.fields = append(cq.fields, coloumn.Field)
		if strings.Contains(coloumn.Type, "int") {
			float, ok := val.(float64)
			if !ok {
				returnError(w, "field "+i+" have invalid type", http.StatusBadRequest)
				return
			}
			cq.vals = append(cq.vals, int(float))
		} else if strings.Contains(coloumn.Type, "text") {
			str, ok := val.(string)
			if !ok {
				returnError(w, "field "+i+" have invalid type", http.StatusBadRequest)
				return
			}
			cq.vals = append(cq.vals, str)
		} else if strings.Contains(coloumn.Type, "varchar") {
			str, ok := val.(string)
			if !ok {
				returnError(w, "field "+i+" have invalid type", http.StatusBadRequest)
				return
			}
			cq.vals = append(cq.vals, str)
		} else {
			returnError(w, "putting type error", http.StatusTeapot)
		}
	}
	query := bytes.NewBuffer(make([]byte, 0))
	query.WriteString("UPDATE ")
	query.WriteString(tableName)
	query.WriteString(" SET ")
	for i, field := range cq.fields {
		query.WriteString("`")
		query.WriteString(field)
		query.WriteString("`")
		query.WriteString(" = ?")
		if i != len(cq.fields)-1 {
			query.WriteString(", ")
		} else {
			query.WriteString(" WHERE ")
			query.WriteString(idName)
			query.WriteString(" = ")
			query.WriteString(strconv.Itoa(id))
		}
	}
	result, err := h.DB.Exec(query.String(), cq.vals...)
	if err != nil {
		returnError(w, "editing db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	preres := make(map[string]interface{})
	preres["response"] = make(map[string]interface{})
	preres["response"].(map[string]interface{})["updated"], _ = result.RowsAffected()
	res, _ := json.Marshal(preres)
	w.Write(res)
}

func (h dbHandler) deleteRecord(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		returnError(w, "wrong URL in POST method", http.StatusInternalServerError)
		return
	}
	tableName := parts[1]
	_, ok := h.tables[tableName]
	if !ok {
		returnError(w, "unknown table", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		returnError(w, "id must be int", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("DELETE FROM "+tableName+" WHERE id = ?", id)
	if err != nil {
		returnError(w, "editing db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	preres := make(map[string]interface{})
	preres["response"] = make(map[string]interface{})
	preres["response"].(map[string]interface{})["deleted"], _ = result.RowsAffected()
	res, _ := json.Marshal(preres)
	w.Write(res)
}

func returnError(w http.ResponseWriter, msg string, status int) {
	log.Println(msg, "status:", status)
	preres := make(map[string]interface{})
	preres["error"] = msg
	res, _ := json.Marshal(preres)
	http.Error(w, string(res), status)
}
