package goffive

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Auth struct {
	User              string `json:"username"`
	Pass              string `json:"password"`
	LoginProviderName string `json:"loginProviderName"`
}

type AuthResponse struct {
	Token struct {
		Token            string `json:"token"`
		Name             string `json:"name"`
		Timeout          int64  `json:"timeout"`
		StartTime        string `json:"startTime"`
		ExpirationMicros uint64 `json:"expirationMicros"`
		LastUpdateMicros uint64 `json:"lastUpdateMicros"`
	} `json:"token"`
}

type Client struct {
	user     string
	password string
	url      string
	conn     *http.Client
	token    string

	ASM *asm
	LTM *ltm
}

func authorization(user string, pass string, url string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},

		},
	}
	path := fmt.Sprintf("%s/mgmt/shared/authn/login", url)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(&Auth{
		User:              user,
		Pass:              pass,
		LoginProviderName: "tmos",
	})

	req, err := http.NewRequest("POST", path, b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var authResp AuthResponse
	dec := json.NewDecoder(bytes.NewBuffer(body))
	err = dec.Decode(&authResp)
	if err != nil {
		return "", err
	}

	return authResp.Token.Token, nil
}

func New(user string, pass string, url string) (*Client, error) {

	if user == "" || pass == "" || url == "" {
		return nil, errors.New("didn't find credentials")
	}

	token, err := authorization(user, pass, url)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	fmt.Println(token)
	return &Client{
		conn:     client,
		user:     user,
		password: pass,
		url:      url,
		token:    token,
		ASM: &asm{
			conn: client,
			token: token,
			url: url,
		},
		LTM: &ltm{
			conn: client,
			token: token,
			url: url,
		},
	}, nil
}
