package handler

import (
	"auth/service/usecase"
	"context"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type middleware struct {
	userS usecase.User
}

func NewMiddleware(userS usecase.User) *middleware {
	return &middleware{userS: userS}
}

func (h *middleware) AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "req", c.Request)
		rctx, err := h.userS.Auth(ctx)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Request = c.Request.WithContext(rctx)
		c.Next()
	}
}

func (h *middleware) Authorize(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := usecase.NewUserValue(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ok, err := enforce(user.Name, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(500, "error occurred when authorizing user")
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, "forbidden")
			return
		}
		c.Next()
	}
}

func enforce(sub string, obj string, act string) (bool, error) {
	enforcer, err := casbin.NewEnforcerSafe("pkg/config/rbac_model.conf", "pkg/config/rbac_policy.csv")
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	ok := enforcer.Enforce(sub, obj, act)
	return ok, nil
}
