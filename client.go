package main

import (
	"auth/service/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const host = "http://localhost:8080/v1/"

func logUser() {
	user := &structs.User{
		Name:     "bob",
		Password: "123456",
		Email:    "test@test.test",
		Rights:   "user",
	}
	userJ, _ := json.Marshal(user)
	res, err := http.Post(host+"registration", "application/json", bytes.NewReader(userJ))
	if err != nil {
		fmt.Println(err)
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	return
}

func get() (string, error) {
	user := &structs.User{
		Password: "123456",
		Email:    "test@test.test",
	}
	userJ, _ := json.Marshal(user)
	res, err := http.Post(host+"login", "application/json", bytes.NewReader(userJ))
	if err != nil {
		return "", err
	}
	tokens, _ := ioutil.ReadAll(res.Body)
	t := structs.JWT{}
	err = json.Unmarshal(tokens, &t)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return t.AccessToken, nil
}

func hi(token string) {
	req, _ := http.NewRequest("GET", host+"admin/hi", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	m, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(m))
}

func exit(token string) {
	req, _ := http.NewRequest("GET", host+"logout", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	m, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(m))
}

func main() {
	//logUser()
	t, err := get()
	if err != nil {
		log.Println(err)
		return
	}
	hi(t)
	//exit(t)

}
