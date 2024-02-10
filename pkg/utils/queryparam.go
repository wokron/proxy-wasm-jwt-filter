package utils

import "strings"

func GetQueryParamValue(path string, key string) (string, bool) {
	params := GetQueryParamsFromPath(path)
	value, ok := params[key]
	if !ok {
		return "", false
	}
	return value, true
}

func GetQueryParamsFromPath(path string) map[string]string {
	params := map[string]string{}
	index := strings.Index(path, "?")
	if index < 0 {
		return params
	}

	paramsString := path[index+1:]
	for _, paramPair := range strings.Split(paramsString, "&") {
		var paramKey, paramValue string
		index := strings.Index(paramPair, "=")
		if index < 0 {
			paramKey = paramPair
			paramValue = ""
		} else {
			paramKey = paramPair[:index]
			paramValue = paramPair[index+1:]
		}
		params[paramKey] = paramValue
	}
	return params
}
