package helpers

import (
	"regexp"
)

func GetUrlParam(url string, re *regexp.Regexp) interface{} {
	res := re.FindAllStringSubmatch(url, 1)
	return res[0][1]
}
