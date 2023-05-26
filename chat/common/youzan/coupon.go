package youzan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	// 缓存token
	CacheToken *TokenResponse
)

type TokenResponse struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Data    struct {
		Expires     int64  `json:"expires"`
		Scope       string `json:"scope"`
		AccessToken string `json:"access_token"`
		AuthorityID string `json:"authority_id"`
	} `json:data`
	Message string `json:"message"`
	MToken  string `json:"m_token"`
}

func RefreshAccessKey(clientID, clientSecret, grantID string) (*TokenResponse, error) {

	if CacheToken != nil && CacheToken.Data.Expires > (time.Now().UnixMilli()+60000) {
		return CacheToken, nil
	}
	url := "https://open.youzanyun.com/auth/token"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
  "client_id": "%s",
  "client_secret": "%s",
  "authorize_type": "silent",
  "grant_id": "%s",
  "refresh": "false"
    }`, clientID, clientSecret, grantID))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rsp := TokenResponse{}
	err = json.Unmarshal(body, &rsp)

	if err != nil {
		return nil, err
	}

	if !rsp.Success {
		return nil, err
	}
	rsp.getMToken()
	CacheToken = &rsp
	return &rsp, nil
}

// 内部掉用
func (s *TokenResponse) getMToken() {

	url := "https://open.youzanyun.com/api/youzan.mei.dept.bind/3.0.0?access_token=%s"
	url = fmt.Sprintf(url, s.Data.AccessToken)
	method := "POST"

	payload := strings.NewReader(`{
  "dept_id": 1
    }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	type Rsp struct {
		Data    string `json:"data"`
		Success bool   `json:"success"`
	}
	rsp := Rsp{}
	err = json.Unmarshal(body, &rsp)

	if err != nil {
		return
	}

	if !rsp.Success {
		return
	}
	s.MToken = rsp.Data
}

func ComponCancelAfterWrite(ticketNo, token, mToken string) (string, error) {

	url := "https://open.youzanyun.com/api/youzan.mei.verification.do/1.0.0?access_token=" + token
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
  "m_token": "%s",
  "ticket_no": "%s",
  "verify_ticket_type": 1
    }`, mToken, ticketNo))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	rsp := TokenResponse{}
	err = json.Unmarshal(body, &rsp)

	if err != nil {
		return "", err
	}

	if !rsp.Success {
		return rsp.Message, errors.New("核销失败" + string(body))
	}
	return rsp.Message, nil
}

func ParseCouponNo(picString string) (string, error) {

	splits := strings.Split(picString, "\n")
	reg, _ := regexp.Compile("[^0-9]")
	for _, s := range splits {
		s = reg.ReplaceAllString(s, "")
		if len(s) >= 12 {
			return s, nil
		}
	}
	return "", errors.New("parse failed")
}
