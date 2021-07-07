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

// Expects args as an array of interfaces
func isValidInput(args interface{}) bool {
	_args, ok := args.([]interface{})
	if !ok {
		return false
	}
	if _args[0].(float64) > 100.0 {
		return false
	}
	if _args[1].(float64) > 100.0 {
		return false
	}
	// if _args[5].(bool) > 100.0 {
	// 	return false
	// }
	// if _args[3].(bool) > 100.0 {
	// 	return false
	// }
	// if _args[4].(bool) > 100 {
	// 	return false
	// }
	if _args[2].(float64) > 100.0 {
		return false
	}
	return true
}

// Update updates the sheet
func Update(args interface{}) interface{} {
	_url := defs["url"].(string)
	sheetName := defs["sheetName"].(string)

	if !isValidInput(args) {
		fmt.Printf("Invalid input = %v", args)
		return nil
	}
	_args := args.([]interface{})
	v := url.Values{}
	v.Set("sheet", sheetName)
	v.Add("col1", fmt.Sprintf("%v", time.Now().Format(time.RFC3339)))
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
