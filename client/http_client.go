package client

import (
	"encoding/json"
	"fmt"
	c "github.com/fomk/docker-hub-limit-exporter/config"
	"github.com/patrickmn/go-cache"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)
var s = cache.New(5*time.Minute,6*time.Minute)
var (
	authUrl = "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull"
	url = "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest"
	httpClient = http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

)

type HubResponse struct {
	Limit float64
	Remaining float64
}

type token struct {
	Token string `json:"token"`
	Expires int `json:"expires_in"`
}

func GetMetrics() (HubResponse, error)  {
	token, tokenErr := getToken()
	if tokenErr != nil {
		return HubResponse{}, tokenErr
	}
	req, respErr := http.NewRequest("HEAD", url, nil)
	if respErr != nil {
		return HubResponse{}, fmt.Errorf("cannot create request  %s", respErr)
	}


	req.Header.Set("User-Agent", "dockerhub_rate_limit_exporter")

	req.Header.Set("Authorization", "Bearer " + token)

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return HubResponse{}, fmt.Errorf("cannot get metrics %s", getErr)
	}

	if res.StatusCode == 401 {
		token, tokenErr := getToken()
		if tokenErr != nil {
			return HubResponse{}, tokenErr
		}
		req.Header.Set("Authorization", "Bearer " + token)
		res, getErr = httpClient.Do(req)
		if getErr != nil {
			return HubResponse{}, fmt.Errorf("cannot get metrics %s", getErr)
		}
	}


	Limit, limErr := strconv.ParseFloat(strings.Split(res.Header.Get("ratelimit-limit"), ";")[0], 64)
	if limErr != nil {
		return HubResponse{}, fmt.Errorf("cannot read limit %s", limErr)
	}
	Remaining, remErr := strconv.ParseFloat(strings.Split(res.Header.Get("ratelimit-remaining"), ";")[0],64)
	if remErr != nil {
		return HubResponse{}, fmt.Errorf("cannot read remaining %s", remErr)
	}

	return HubResponse{
		Limit: Limit,
		Remaining: Remaining,
	}, nil
}

func getToken() (string, error) {

	jwt, found := s.Get("hubToken")

	if found {
		return fmt.Sprintf("%v", jwt), nil
	}

	req, reqErr := http.NewRequest("GET", authUrl, nil)
	if reqErr != nil {
		return "", fmt.Errorf("cannot create request %s", reqErr)
	}

	if *c.HubUser != "" {
		req.SetBasicAuth(*c.HubUser, *c.HubPass)
	}

	req.Header.Set("User-Agent", "dockerhub_rate_limit_exporter")


	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return "", fmt.Errorf("cannot obtain token %s", getErr)
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", fmt.Errorf("cannot read response %s", readErr)
	}
	token := token{}
	jsonErr := json.Unmarshal(body, &token)
	if jsonErr != nil {
		return "", fmt.Errorf("cannot parse response %s", jsonErr)
	}
	s.Set("hubToken",token.Token, time.Duration(token.Expires)*time.Second)

	return token.Token, nil
}