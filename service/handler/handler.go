package handler

import (
	"auth/pkg/myJwt"
	"auth/service/structs"
	"auth/service/usecase"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type handler struct {
	userS usecase.User
}

func NewHandler(userS usecase.User) *handler{
	return &handler{
		userS: userS,
	}
}

func (h *handler) Authorization(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		auth := r.Header.Get("Authorization")
		if auth != ""{
			token := strings.Split(auth, " ")
			fmt.Println(token)
			_, err := myJwt.VerefyJWT(token[1])
			next.ServeHTTP(w, r)
			if err != nil{
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}else{
			_, _ = w.Write([]byte("err auth"))
			return
		}
	})
}

func (h *handler) AuthHandle(w http.ResponseWriter, r *http.Request) {
	authData := structs.NewUser()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)
	bodyJ, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte(err.Error()))
	}
	err = json.Unmarshal(bodyJ, &authData)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	token, err := h.userS.SignIn(ctx, authData)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(token)
	_, _ = w.Write([]byte(token))
}

func (h *handler) Hello(w http.ResponseWriter, r *http.Request){
	_, _ = w.Write([]byte("hello"))
}
