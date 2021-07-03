package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	defs map[string]interface{}
)

// Init initializes the layer
func Init(_defs interface{}) interface{} {
	defs = _defs.(map[string]interface{})
	return nil
}

// InitDebug inits debug
func InitDebug(_defs interface{}) interface{} {
	// save args
	defs = _defs.(map[string]interface{})
	return nil
}

// Update updates the sheet
func Update(args interface{}) interface{} {
	_url := defs["url"].(string)
	sheetName := defs["sheetName"].(string)

	_args := args.([]interface{})
	v := url.Values{}
	v.Set("sheet", sheetName)
	v.Add("col1", fmt.Sprintf("%v", time.Now()))
	v.Add("col2", fmt.Sprintf("%v", _args[0].(float64)))
	v.Add("col3", fmt.Sprintf("%v", _args[1].(float64)))
	v.Add("col4", fmt.Sprintf("%v", _args[5].(bool)))
	v.Add("col5", fmt.Sprintf("%v", _args[3].(bool)))
	v.Add("col6", fmt.Sprintf("%v", _args[4].(bool)))
	v.Add("col7", fmt.Sprintf("%v", _args[2].(float64)))
	v.Add("col8", "FALSE")
	_url = _url + "?" + v.Encode()
	// fmt.Println("SHEETS", _url, args, v.Encode())
	_, err := http.Get(_url)
	if err != nil {
		fmt.Println("Unable to update sheet, err:", err)
	}
	return nil
}
