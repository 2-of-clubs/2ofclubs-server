package handler

import (
	"../model"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

/*
	Common methods shared amongst the different models
*/

func GenerateSaltedPass(password string) (string, bool) {
	saltedHashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}
	return string(saltedHashPass), true
}
func RecordExists(db *gorm.DB, column string, val string, t interface{}) bool {
	return !db.Where(column+"= ?", val).First(t).RecordNotFound()
}

func ExtractPersonInfo(r *http.Request) model.Person {
	decoder := json.NewDecoder(r.Body)
	p := model.NewPerson()
	decoder.Decode(&p)
	return p
}

func ParseJSON(response interface{}) string {
	data, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func WriteData(data string, code int, w http.ResponseWriter) int {
	w.WriteHeader(code)
	n, err := fmt.Fprint(w, data)
	if err != nil {
		return -1
	}
	return n
}
