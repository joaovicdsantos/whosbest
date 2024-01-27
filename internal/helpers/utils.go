package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/joaovicdsantos/whosbest-api/internal/models"
	"github.com/joaovicdsantos/whosbest-api/internal/services"
)

func GetUrlParam(url string, re *regexp.Regexp) interface{} {
	res := re.FindAllStringSubmatch(url, 1)
	return res[0][1]
}

func ParseBodyToStruct(r *http.Request, model interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return fmt.Errorf("invalid requisition body")
	}

	if err = json.Unmarshal(body, model); err != nil {
		return fmt.Errorf("invalid requisition body")
	}
	return nil
}

func ParseMapToStruct(genericMap map[string]interface{}, model interface{}) error {
	jsonBody, err := json.Marshal(genericMap)
	if err != nil {
		return fmt.Errorf("invalid generic map")
	}
	if err = json.Unmarshal(jsonBody, model); err != nil {
		return fmt.Errorf("invalid generic map")
	}
	return nil
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap)
	return
}

func GetCurrentUser(payload map[string]interface{}, db *sql.DB) models.User {
	username := fmt.Sprintf("%s", payload["username"])

	var userService services.UserService
	userService.DB = db

	user := userService.GetOneByUsername(username)

	return user
}
