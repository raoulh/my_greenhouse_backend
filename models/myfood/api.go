package myfood

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	MyFoodApiHost = "hub.myfood.eu"
)

type apiOptions struct {
	timeout   time.Duration
	insecure  bool
	proxy     func(*http.Request) (*url.URL, error)
	transport func(*http.Transport)
}

// Optional parameter for NewDuoApi, used to configure timeouts on API calls.
func SetTimeout(timeout time.Duration) func(*apiOptions) {
	return func(opts *apiOptions) {
		opts.timeout = timeout
		return
	}
}

// Optional parameter for testing only.  Bypasses all TLS certificate validation.
func SetInsecure() func(*apiOptions) {
	return func(opts *apiOptions) {
		opts.insecure = true
	}
}

// Optional parameter for NewDuoApi, used to configure an HTTP Connect proxy
// server for all outbound communications.
func SetProxy(proxy func(*http.Request) (*url.URL, error)) func(*apiOptions) {
	return func(opts *apiOptions) {
		opts.proxy = proxy
	}
}

// SetTransport enables additional control over the HTTP transport used to connect to the API.
func SetTransport(transport func(*http.Transport)) func(*apiOptions) {
	return func(opts *apiOptions) {
		opts.transport = transport
	}
}

type MyFoodApi struct {
	host      string
	userAgent string
	apiClient *http.Client
}

// SetCustomHTTPClient allows one to set a completely custom http client that
// will be used to make network calls to the duo api
func (mfapi *MyFoodApi) SetCustomHTTPClient(c *http.Client) {
	mfapi.apiClient = c
}

// Build an return a MyFoodApi struct
func NewMyFoodApi(host string, options ...func(*apiOptions)) *MyFoodApi {
	opts := apiOptions{
		proxy: http.ProxyFromEnvironment,
	}
	for _, o := range options {
		o(&opts)
	}

	tr := &http.Transport{
		Proxy: opts.proxy,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.insecure,
		},
	}
	if opts.transport != nil {
		opts.transport(tr)
	}

	return &MyFoodApi{
		host:      host,
		userAgent: "mygreenhouse_api/1.0",
		apiClient: &http.Client{
			Timeout:   opts.timeout,
			Transport: tr,
		},
	}
}

type requestOptions struct{}

type MyFoodApiOption func(*requestOptions)

func (duoapi *MyFoodApi) buildOptions(options ...MyFoodApiOption) *requestOptions {
	opts := &requestOptions{}
	for _, o := range options {
		o(opts)
	}
	return opts
}

// Make a MyFood Rest API call
// Example: myfood.Call("POST", "/api/identity/token", nil)
func (mfapi *MyFoodApi) call(method string, uri string, params url.Values, body interface{}) (*http.Response, []byte, error) {

	url := url.URL{
		Scheme:   "https",
		Host:     mfapi.host,
		Path:     uri,
		RawQuery: params.Encode(),
	}
	headers := make(map[string]string)
	headers["User-Agent"] = mfapi.userAgent

	var requestBody io.ReadCloser = nil
	if body != nil {
		headers["Content-Type"] = "application/json"

		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		requestBody = ioutil.NopCloser(bytes.NewReader(b))
	}

	return mfapi.makeHttpCall(method, url, headers, requestBody)
}

// Make a MyFood Rest API call using token
// Example: myfood.CallWithToken("GET", "/api/v1/Measures/GetPHMeasureForUser", token, params)
func (mfapi *MyFoodApi) callWithToken(method string, uri string, token string, params url.Values, body interface{}) (*http.Response, []byte, error) {

	url := url.URL{
		Scheme:   "https",
		Host:     mfapi.host,
		Path:     uri,
		RawQuery: params.Encode(),
	}
	headers := make(map[string]string)
	headers["User-Agent"] = mfapi.userAgent
	headers["Authorization"] = "Bearer " + token

	var requestBody io.ReadCloser = nil
	if body != nil {
		headers["Content-Type"] = "application/json"

		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		requestBody = ioutil.NopCloser(bytes.NewReader(b))
	}

	return mfapi.makeHttpCall(method, url, headers, requestBody)
}

func (mfapi *MyFoodApi) makeHttpCall(
	method string,
	url url.URL,
	headers map[string]string,
	body io.ReadCloser) (*http.Response, []byte, error) {

	request, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	//set headers
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	if body != nil {
		request.Body = body
	}

	resp, err := mfapi.apiClient.Do(request)
	var bodyBytes []byte
	if err != nil {
		return resp, bodyBytes, err
	}

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return resp, bodyBytes, err
}

//GetToken authenticate a user and retrieve a token
func (mfapi *MyFoodApi) GetToken(username, password string) (t *TokenData, err error) {

	d := &AuthData{
		Username: username,
		Password: password,
	}

	_, body, err := mfapi.call("POST", "/api/identity/token", url.Values{}, d)
	if err != nil {
		return
	}

	var tres *TokenResultData
	if err = json.Unmarshal(body, tres); err != nil {
		return
	}

	if tres.Failed || !tres.Succeeded {
		return nil, fmt.Errorf("GetToken failed: %v", tres.Messages)
	}

	return &tres.Data.TokenData, nil
}

//RefreshToken refresh an existing token
func (mfapi *MyFoodApi) RefreshToken(token *TokenData) (t *TokenData, err error) {

	_, body, err := mfapi.call("POST", "/api/identity/token/refresh", url.Values{}, token)
	if err != nil {
		return
	}

	var tres *TokenResultData
	if err = json.Unmarshal(body, tres); err != nil {
		return
	}

	if tres.Failed || !tres.Succeeded {
		return nil, fmt.Errorf("RefreshToken failed: %v", tres.Messages)
	}

	return &tres.Data.TokenData, nil
}

//GetProductionUnitDetailForUser get info for the greenhouse id
// timerange option:
// 		LastDay = 0,
//		LastWeek = 1
//		LastThreeMonths =2
func (mfapi *MyFoodApi) GetProductionUnitDetailForUser(token string, id uint, timerange uint) (p *ProdUnitDetailData, err error) {

	_, body, err := mfapi.callWithToken("GET", "/api/v1/ProductionUnit/GetProductionUnitDetailForUser", token, url.Values{
		"id":    []string{strconv.Itoa(int(id))},
		"range": []string{strconv.Itoa(int(timerange))},
	}, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, p); err != nil {
		return
	}

	if p.Failed || !p.Succeeded {
		return nil, fmt.Errorf("GetProductionUnitDetailForUser failed: %v", p.Messages)
	}

	return
}

//GetAllProductionUnitIdsForCurrentUser gets all greenhouse for a user
func (mfapi *MyFoodApi) GetAllProductionUnitIdsForCurrentUser(token string) (p *ProdUnitsData, err error) {

	_, body, err := mfapi.callWithToken("GET", "/api/v1/ProductionUnit/GetAllProductionUnitIdsForCurrentUser", token, url.Values{}, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, p); err != nil {
		return
	}

	if p.Failed || !p.Succeeded {
		return nil, fmt.Errorf("GetAllProductionUnitIdsForCurrentUser failed: %v", p.Messages)
	}

	return
}

//GetPHMeasureForUser get pH measurements for specified greenhouse
// timerange option:
// 		LastDay = 0,
//		LastWeek = 1
//		LastThreeMonths =2
func (mfapi *MyFoodApi) GetPHMeasureForUser(token string, id uint, timerange uint) (p *ResultData, err error) {

	_, body, err := mfapi.callWithToken("GET", "/api/v1/Measures/GetPHMeasureForUser", token, url.Values{
		"id":    []string{strconv.Itoa(int(id))},
		"range": []string{strconv.Itoa(int(timerange))},
	}, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, p); err != nil {
		return
	}

	if p.Failed || !p.Succeeded {
		return nil, fmt.Errorf("GetPHMeasureForUser failed: %v", p.Messages)
	}

	return
}

//GetWaterTemperatureForUser get water temp measurements for specified greenhouse
// timerange option:
// 		LastDay = 0,
//		LastWeek = 1
//		LastThreeMonths =2
func (mfapi *MyFoodApi) GetWaterTemperatureForUser(token string, id uint, timerange uint) (p *ResultData, err error) {

	_, body, err := mfapi.callWithToken("GET", "/api/v1/Measures/GetWaterTemperatureForUser", token, url.Values{
		"id":    []string{strconv.Itoa(int(id))},
		"range": []string{strconv.Itoa(int(timerange))},
	}, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, p); err != nil {
		return
	}

	if p.Failed || !p.Succeeded {
		return nil, fmt.Errorf("GetWaterTemperatureForUser failed: %v", p.Messages)
	}

	return
}

//GetAirTemperatureMeasureForUser get air temp measurements for specified greenhouse
// timerange option:
// 		LastDay = 0,
//		LastWeek = 1
//		LastThreeMonths =2
func (mfapi *MyFoodApi) GetAirTemperatureMeasureForUser(token string, id uint, timerange uint) (p *ResultData, err error) {

	_, body, err := mfapi.callWithToken("GET", "/api/v1/Measures/GetAirTemperatureMeasureForUser", token, url.Values{
		"id":    []string{strconv.Itoa(int(id))},
		"range": []string{strconv.Itoa(int(timerange))},
	}, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, p); err != nil {
		return
	}

	if p.Failed || !p.Succeeded {
		return nil, fmt.Errorf("GetAirTemperatureMeasureForUser failed: %v", p.Messages)
	}

	return
}

//GetHumidityMeasureForUser get humidity measurements for specified greenhouse
// timerange option:
// 		LastDay = 0,
//		LastWeek = 1
//		LastThreeMonths =2
func (mfapi *MyFoodApi) GetHumidityMeasureForUser(token string, id uint, timerange uint) (p *ResultData, err error) {

	_, body, err := mfapi.callWithToken("GET", "/api/v1/Measures/GetHumidityMeasureForUser", token, url.Values{
		"id":    []string{strconv.Itoa(int(id))},
		"range": []string{strconv.Itoa(int(timerange))},
	}, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, p); err != nil {
		return
	}

	if p.Failed || !p.Succeeded {
		return nil, fmt.Errorf("GetHumidityMeasureForUser failed: %v", p.Messages)
	}

	return
}
