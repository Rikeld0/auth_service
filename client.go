package main

import (
	"auth/service/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const host = "http://localhost:8080/"

func get() string {
	user := &structs.User{
		Password: "123456",
		Email:    "ivanbarbarian@gmail.com",
	}
	userJ, _ := json.Marshal(user)
	res, _ := http.Post(host+"auth", "application/json", bytes.NewReader(userJ))
	tokens, _ := ioutil.ReadAll(res.Body)
	t := structs.JWT{}
	err := json.Unmarshal(tokens, &t)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return t.AccessToken
}

func get1() {
	res, _ := http.PostForm("https://test-api.k6.io/auth/token/login/", url.Values{
		"username": {"TestUser"},
		"password": {"SuperCroc2020"},
	})
	tokens, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(tokens))
}

func hi(token string) {
	req, _ := http.NewRequest("GET", host+"hi", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	m, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(m))
}

func main() {
	t := get()
	hi(t)
}
