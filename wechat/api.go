package wechat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	appid = "wx782c26e4c19acffb"

	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36"
	accept    = "application/json, text/plain, */*"

	endpointUUID      = "https://login.weixin.qq.com/jslogin"
	endpointLoginPoll = "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"
)

var defaultClient = &http.Client{}

// UUID fetches UUID for subsequent QR code image request
func UUID() (string, error) {
	re, err := regexp.Compile("window.QRLogin.code = (\\d+); window.QRLogin.uuid = \"(\\S+?)\"")
	if err != nil {
		return "", err
	}

	queryParams := url.Values{}
	queryParams.Add("appid", appid)
	queryParams.Add("fun", "new")
	queryParams.Add("lang", "en_US")
	queryParams.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	queryString := queryParams.Encode()

	client := &http.Client{}

	request, err := http.NewRequest("GET", endpointUUID+"?"+queryString, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("User-Agent", userAgent)
	request.Header.Add("Accept", accept)
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	matches := re.FindAllStringSubmatch(string(body), -1)
	log.Printf("Get regex match result: %v", matches)
	if matches != nil && len(matches[0]) == 3 && matches[0][1] == "200" {
		return matches[0][2], nil
	}
	return "", fmt.Errorf("Received code %s", matches[0][1])
}
