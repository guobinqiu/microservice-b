package transports

import (
	"context"
	"encoding/json"
	"microservice/models"
	"net/http"
	"strconv"
)

//自定义请求对象转http请求
func EncodeUserRequest(c context.Context, request *http.Request, r interface{}) error {
	userRequest := r.(models.UserRequest)
	request.URL.Path += "/users/" + strconv.Itoa(userRequest.Uid)
	return nil
}

//收到服务端响应进行解码(json转对象)
func DecodeGetUserNameResponse(c context.Context, response *http.Response) (interface{}, error) {
	var userResponse models.UserResponse
	err := json.NewDecoder(response.Body).Decode(&userResponse)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}

func DecodeDelUserResponse(c context.Context, response *http.Response) (interface{}, error) {
	if response.StatusCode != 204 {
		var e models.Error
		err := json.NewDecoder(response.Body).Decode(&e)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
	var userResponse models.UserResponse
	err := json.NewDecoder(response.Body).Decode(&userResponse)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}
