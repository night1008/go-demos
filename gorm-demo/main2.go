package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	endpoint string
	username string
	password string
)

type APIClient struct {
	http.Client
	Endpoint string
	Username string
	Password string
}

type App struct {
	ID    uint
	AppID string
}

type Instance struct {
	ID         uint
	AppID      string
	InstanceID string
}

func (c *APIClient) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "applacation/json")
	req.SetBasicAuth(c.Username, c.Password)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *APIClient) do(method string, uri string, args interface{}, returns interface{}) (int, error) {
	req, err := c.newRequest(method, c.Endpoint+uri, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if returns != nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if err := json.Unmarshal(body, returns); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return resp.StatusCode, nil
}

func init() {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:8080/api/v1", "Input api endpoint")
	flag.StringVar(&username, "username", "admin", "Input basic request auth username")
	flag.StringVar(&password, "password", "secret", "Input basic request auth password")
}

func main() {
	flag.Parse()

	fmt.Println(endpoint, username, password)

	client := APIClient{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}

	list := make([]App, 0)
	limit := 100
	page := 1
	for {
		var listPage struct {
			Total int
			List  []App
		}
		_, err := client.do("GET", "/apps", fmt.Sprintf("?page=%d&limit=%d", page, limit), &list)
		if err != nil {
			panic(err)
		}
		list = append(list, listPage.List...)
		if listPage.Total == len(list) {
			break
		}
		page += 1
	}
	for _, item := range list {
		tag := getAppTag(item.AppID)
		addTag(&client, "apps", item.ID, []string{tag})
	}
}

func addTag(client *APIClient, taggableType string, taggableID uint, tags []string) {
	args := make(map[string]interface{}, 0)
	args["TaggableType"] = taggableType
	args["TaggableID"] = taggableID
	args["Tags"] = tags
	_, err := client.do("POST", "/tags/attach", args, nil)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("==> add tag ", args)
	}
}

func getAppTag(appID string) string {
	if strings.HasSuffix(appID, "-detection") {
		return "behavior"
	} else {
		return ""
	}
}
