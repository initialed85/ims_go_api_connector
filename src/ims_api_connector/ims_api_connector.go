package ims_api_connector

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"log"
)

//
// main object
//

type IMSAPIConnector struct {
	Username, Password, BaseURL, Key string
	JSONDecoder                      json.Decoder
	Authenticated                    bool
	client                           *http.Client
}

//
// related objects
//

type authenticationResponse struct {
	Key            string `json:"key"`
	NonFieldErrors []string `json:"non_field_errors"`
}

type Asset struct {
	ID                int `json:"id"`
	Name              string `json:"name"`
	IsDeleted         bool `json:""`
	LastUpdated       time.Time `json:"last_updated"`
	Note              string `json:"note"`
	JSONData          interface{} `json:""`
	TypeID            int `json:"type_id"`
	PrimaryIPDeviceID int `json:"primary_ip_device_id"`
	SiteID            int `json:"site_id"`
	Tags              []int `json:"tags"`
}

//
// private methods
//

func (self *IMSAPIConnector) buildResource(resource string) string {
	return self.BaseURL + strings.Trim(resource, "/") + "/"
}

func (self *IMSAPIConnector) buildRequest(method, request_url, body string) (*http.Request, error) {

	var validMethod bool = false

	switch method {
	case http.MethodGet:
		validMethod = true
	case http.MethodPost:
		validMethod = true
	}

	if !validMethod {
		return nil, errors.New(fmt.Sprint("method of %v was invalid", method))
	}

	req, err := http.NewRequest(method, request_url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cache-Control", "no-cache")
	if self.Authenticated {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Token "+self.Key)
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

func (self *IMSAPIConnector) buildUsernameAndPasswordParams() string {
	params_values := url.Values{}

	params_values.Set("username", self.Username)
	params_values.Set("password", self.Password)

	return params_values.Encode()
}

func (self *IMSAPIConnector) handleAuthenticationResponse(body []byte) bool {
	var authentication authenticationResponse

	err := json.Unmarshal(body, &authentication)
	if err != nil {
		log.Fatal(err)
	}

	if authentication.Key != "" {
		self.Key = authentication.Key
		self.Authenticated = true
	} else {
		self.Key = ""
		self.Authenticated = false
	}

	return self.Authenticated
}

func (self *IMSAPIConnector) handleGetAssetsResponse(body []byte) ([]Asset, error) {
	var assets []Asset

	err := json.Unmarshal(body, &assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

//
// public methods
//

func New(username, password, baseurl string, timeout int) *IMSAPIConnector {
	if !(strings.HasPrefix(baseurl, "http://") || strings.HasPrefix(baseurl, "https://")) {
		baseurl = "http://" + baseurl
	}

	if !strings.HasSuffix(baseurl, "/api/") {
		baseurl = strings.TrimRight(baseurl, "/") + "/api/"
	}

	client := http.Client{Timeout: time.Second * time.Duration(timeout)}

	self := &IMSAPIConnector{Username: username, Password: password, BaseURL: baseurl, client: &client}

	return self
}

func (self *IMSAPIConnector) Post(req *http.Request) ([]byte, error) {
	resp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func (self *IMSAPIConnector) Authenticate() (bool, error) {
	req, err := self.buildRequest(http.MethodPost, self.buildResource("auth/login"), self.buildUsernameAndPasswordParams())
	if err != nil {
		return false, err
	}

	self.buildUsernameAndPasswordParams()

	body, err := self.Post(req)
	if err != nil {
		return false, err
	}

	return self.handleAuthenticationResponse(body), nil
}

func (self *IMSAPIConnector) GetAssets() ([]Asset, error) {
	req, err := self.buildRequest(http.MethodGet, self.buildResource("assets"), "")
	if err != nil {
		return nil, err
	}

	body, err := self.Post(req)
	if err != nil {
		return nil, err
	}

	assets, err := self.handleGetAssetsResponse(body)
	if err != nil {
		return nil, err
	}

	return assets, nil
}
