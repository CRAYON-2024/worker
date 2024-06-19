package api

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CRAYON-2024/worker/common"
	"github.com/CRAYON-2024/worker/internal/entity"
)

const (
	baseUrl = "https://dummyapi.io/data/v1"
	getMethod = "GET"
)

var (
	client = &http.Client{Timeout: 10 * time.Second}
)

func createRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("app-id", os.Getenv("APP_ID"))

	return req, nil
}

func FetchUsers(page int) (entity.UserResponse, error) {
	var userResponse entity.UserResponse
	req, err := createRequest(getMethod, baseUrl+"/user"+"?page="+strconv.Itoa(page))

	if err != nil {
		return userResponse, err
	}

	response, err := client.Do(req)

	if err != nil {
		return userResponse, err
	}

	defer response.Body.Close()

	return common.UnmarshalResponse[entity.UserResponse](response)
}

func FetchUserDetail(userID string) (entity.User, error) {
	var user entity.User

	req, err := createRequest(getMethod, baseUrl+"/user/"+userID)

	if err != nil {
		return user, err
	}

	response, err := client.Do(req)

	if err != nil {
		return user, err
	}

	defer response.Body.Close()

	return common.UnmarshalResponse[entity.User](response)
}


func FetchPosts(page int) ( entity.PostResponse, error) {
	var postResponse entity.PostResponse

	req, err := createRequest(getMethod, baseUrl+"/post"+"?page="+strconv.Itoa(page))

	if err != nil {
		return postResponse, err
	}

	response, err := client.Do(req)

	if err != nil {
		return postResponse, err
	}

	defer response.Body.Close()

	return common.UnmarshalResponse[entity.PostResponse](response)
}