package handler

import (
	"auth/service/structs"
	"auth/service/usecase"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type handler struct {
	userS usecase.User
}

func NewHandler(userS usecase.User) *handler {
	return &handler{
		userS: userS,
	}
}

func (h *handler) AuthorizationMidlaware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "req", r)
		rctx, err := h.userS.Auth(ctx)
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusFound)
		}
		next.ServeHTTP(w, r.WithContext(rctx))
	})
}

func (h *handler) AuthHandle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	usr := structs.NewUser()
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
	err = json.Unmarshal(bodyJ, &usr)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx = context.WithValue(r.Context(), "req", r)
	tokens, err := h.userS.SignIn(ctx, usr)
	if err != nil {
		fmt.Println(err)
		return
	}
	tokensJ, _ := json.Marshal(tokens)
	_, _ = w.Write(tokensJ)
}

func (h *handler) RegHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	usr := structs.NewUser()
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
	err = json.Unmarshal(bodyJ, &usr)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx = context.WithValue(r.Context(), "req", r)
	err = h.userS.SignUp(ctx, usr)
	if err != nil {
		return
	}
	http.Redirect(w, r, "/auth", http.StatusAccepted)
}

func (h *handler) Hello(w http.ResponseWriter, r *http.Request) {
	msg, err := h.userS.GetMsg(r.Context())
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = w.Write([]byte(msg))
}
