package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/joaovicdsantos/whosbest-api/app/models"
	"github.com/joaovicdsantos/whosbest-api/app/services"
)

func GetUrlParam(url string, re *regexp.Regexp) interface{} {
	res := re.FindAllStringSubmatch(url, 1)
	return res[0][1]
}

func ParseBodyToStruct(r *http.Request, model interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return fmt.Errorf("Invalid requisition body")
	}

	if err = json.Unmarshal(body, model); err != nil {
		return fmt.Errorf("Invalid requisition body")
	}
	return nil
}

func GetCurrentUser(payload map[string]interface{}, db *sql.DB) models.User {
	username := fmt.Sprintf("%s", payload["username"])

	var userService services.UserService
	userService.DB = db

	user := userService.GetOneByUsername(username)

	return user
}
