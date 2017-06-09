package wechat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
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

// ClientManager is a manager for HTTP client, including cookie management.
type ClientManager struct {
	client      *http.Client
	redirectURI string
}

var clientManager *ClientManager

func init() {
	jar, _ := cookiejar.New(nil)
	clientManager = &ClientManager{
		client: &http.Client{
			Jar: jar,
		},
	}
}

func withDefaultHeader(request *http.Request) *http.Request {
	request.Header.Add("User-Agent", userAgent)
	request.Header.Add("Accept", accept)
	return request
}

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

	request, err := http.NewRequest("GET", endpointUUID+"?"+queryString, nil)
	if err != nil {
		return "", err
	}

	response, err := clientManager.client.Do(withDefaultHeader(request))
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

// PollLogin fetches login state and update redirect URI once user logins.
func PollLogin(uuid string) (string, error) {
	values := url.Values{}
	values.Add("loginicon", "true")
	values.Add("uuid", uuid)
	values.Add("tip", "0")
	values.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	values.Add("_", strconv.FormatInt(time.Now().Unix(), 10))

	uri := endpointLoginPoll + "?" + values.Encode()
	request, _ := http.NewRequest("GET", uri, nil)
	response, err := clientManager.client.Do(withDefaultHeader(request))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Getting login state: %s", body)

	re := regexp.MustCompile("window.code=(\\d+);")
	matches := re.FindAllStringSubmatch(string(body), -1)
	if matches != nil && len(matches[0]) == 2 {
		code := matches[0][1]
		if code == "200" {
			re = regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
			matches = re.FindAllStringSubmatch(string(body), -1)
			if matches != nil && len(matches[0]) == 2 {
				clientManager.redirectURI = matches[0][1]
				log.Printf("Getting redirect URI %s", matches[0][1])
			}
		}
		return code, nil
	}

	return "", nil
}
