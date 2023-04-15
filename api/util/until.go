package util

import (
	v1 "StockMe/api"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/logs"
)

type Util struct {
	v1.API
}

func FormatDateUnix(t time.Time) string {
	if (t == time.Time{}) {
		return ""
	}
	return strconv.Itoa(int(t.UnixNano()))
}

func CallPostHttpWithJsonRespond(url string, bodyValue, headerValue interface{}, method string) (bool, interface{}, map[string][]string) {

	logs.Info("url = ", url)

	bodyRequestByte, _ := json.Marshal(bodyValue)
	headerRequestByte, _ := json.Marshal(headerValue)
	var headerPost map[string]string

	json.Unmarshal(headerRequestByte, &headerPost)

	client := &http.Client{}

	var request *http.Request
	var err1 error

	switch method {
	case "GET":
		request, err1 = http.NewRequest("GET", url, nil)

		var urlQuery map[string]string

		err := json.Unmarshal(bodyRequestByte, &urlQuery)
		if err != nil || request == nil {
			return false, nil, nil
		}

		q := request.URL.Query()

		for paramsKey, paramsValue := range urlQuery {
			q.Add(paramsKey, paramsValue)
		}

		request.URL.RawQuery = q.Encode()
		url = url + "?" + q.Encode()

	case "POST":
		request, err1 = http.NewRequest("POST", url, bytes.NewBuffer(bodyRequestByte))
	}

	for headerKey, headerValue := range headerPost {

		request.Header[headerKey] = []string{headerValue}

	}

	respond, err2 := client.Do(request)

	logs.Info("method = ", method)

	if respond == nil || request == nil {

		return false, nil, nil

	}

	defer respond.Body.Close()
	reBytes, err3 := ioutil.ReadAll(respond.Body)
	interfaceHeader := respond.Header
	_, err4 := json.Marshal(interfaceHeader)

	var respondBody interface{}

	err5 := json.Unmarshal(reBytes, &respondBody)

	descriptionError := ""
	if err1 != nil {
		descriptionError += "err1 = " + err1.Error()
	}
	if err2 != nil {
		descriptionError += "err2 = " + err2.Error()
	}
	if err3 != nil {
		descriptionError += "err3 = " + err3.Error()
	}
	if err4 != nil {
		descriptionError += "err4 = " + err4.Error()
	}
	if err5 != nil {
		descriptionError += "err5 = " + err5.Error()
	}

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {

		return false, nil, nil

	} else {

		return true, respondBody, respond.Header

	}

}

func ToDateTime(s string) (t time.Time) {
	if s == "" {
		return time.Time{}
	}
	ss := strings.Split(s, " ")
	sections := strings.Split(s, ":")
	// hh:mm
	if len(ss) == 2 && len(ss[1]) == 5 {
		logs.Debug("case 0")
		s += ":00"
	} else if sections == nil || len(sections) == 1 {
		logs.Debug("case 1")
		s += " 00:00:00"
	}
	isDash := (strings.Replace(s, "-", "x", -1)) != s
	var err error
	if isDash {
		t, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Now().Location())
	} else {
		t, err = time.ParseInLocation("02/01/2006 15:04:05", s, time.Now().Location())
	}
	if err != nil {
		logs.Error("err parse date", err)
		return time.Time{}
	}
	return t
}
