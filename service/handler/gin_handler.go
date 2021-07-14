package handler

import (
	"auth/service/structs"
	"auth/service/usecase"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

type ginHandler struct {
	userS usecase.User
}

func NewGinHandler(userS usecase.User) *ginHandler {
	return &ginHandler{
		userS: userS,
	}
}

//TODO: сделать ошибки для возрата клиенту

func (h *ginHandler) Login(c *gin.Context) {
	usr := structs.NewUser()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(c.Request.Body)
	bodyJ, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	err = json.Unmarshal(bodyJ, &usr)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.WithValue(c.Request.Context(), "req", c.Request)
	tokens, err := h.userS.SignIn(ctx, usr)
	if err != nil {
		fmt.Println(err)
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *ginHandler) Logout(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), "req", c.Request)
	if err := h.userS.SignOut(ctx); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Redirect(http.StatusOK, "/login")
}

func (h *ginHandler) Registration(c *gin.Context) {
	ctx := context.Background()
	usr := structs.NewUser()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(c.Request.Body)
	bodyJ, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	err = json.Unmarshal(bodyJ, &usr)
	if err != nil {
		fmt.Println(err)
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	usr.Rights = "user"
	ctx = context.WithValue(c.Request.Context(), "req", c.Request)
	err = h.userS.SignUp(ctx, usr)
	if err != nil {
		return
	}
	c.Redirect(http.StatusAccepted, "/login")
}

func (h *ginHandler) Hello(c *gin.Context) {
	msg, err := h.userS.GetMsg(c.Request.Context())
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, msg)
}

func (h *ginHandler) HelloA(c *gin.Context) {
	msg, err := h.userS.GetMsg(c.Request.Context())
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, msg+" admin")
}
