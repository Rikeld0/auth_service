package usecase

import (
	"auth/service/repo"
	"auth/service/structs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const token = "user"

func NewUserValue(ctx context.Context) (*structs.User, error) {
	buf, ok := ctx.Value(token).(*structs.User)
	if !ok {
		return nil, errors.New("error user")
	}
	return buf, nil
}

func ReqValue(ctx context.Context) (*http.Request, error) {
	buf, ok := ctx.Value("req").(*http.Request)
	if !ok {
		return nil, errors.New("error req")
	}
	return buf, nil
}

type userService struct {
	repoU   repo.Users
	repoUK  repo.UsersKey
	repoJwt repo.Jwt
}

func NewUserService(repoU repo.Users, repoUK repo.UsersKey, repoJwt repo.Jwt) User {
	return &userService{
		repoU:   repoU,
		repoUK:  repoUK,
		repoJwt: repoJwt,
	}
}

func getToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("no token")
	}
	splits := strings.SplitN(auth, " ", 2)
	if len(splits) < 2 {
		return "", errors.New("Bad authorization string")
	}
	return splits[1], nil
}

func (u *userService) Auth(ctx context.Context) (context.Context, error) {
	r, err := ReqValue(ctx)
	if err != nil {
		return nil, err
	}
	uuid, err := u.repoU.FindUserId(ctx, r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	keys, err := u.repoUK.Get(ctx, uuid)
	c := structs.Cl{}
	token, err := getToken(r)
	if err != nil {
		return nil, err
	}
	claims, err := u.repoJwt.Validate(token, keys)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(claims, &c)
	usr, err := u.repoU.GetUserId(ctx, c.Iss)
	//TODO: сделать проверку жизни токена
	return context.WithValue(r.Context(), "user", usr), nil
}

func (u *userService) SignUp(ctx context.Context, user *structs.User) error {
	uuid, err := u.repoU.InsertUser(ctx, user)
	if err != nil {
		return err
	}
	err = u.repoUK.Put(ctx, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) SignIn(ctx context.Context, user *structs.User) (*structs.JWT, error) {
	r, err := ReqValue(ctx)
	if err != nil {
		return nil, err
	}
	usr, err := u.repoU.Get(ctx, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	//TODO: подумать над проверкой пароля
	//if ok := structs.CheckPass(usr.Password, user.Password); !ok {
	//	return nil, errors.New("password not valid")
	//}
	if err = u.repoU.SaveUserIDAndIP(ctx, usr.Uuid, r.RemoteAddr); err != nil {
		return nil, err
	}
	keys, err := u.repoUK.Get(ctx, usr.Uuid)
	jwtToken, err := u.repoJwt.Generate(usr.Uuid, keys)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jwtToken, nil
}

func (u *userService) SignOut(ctx context.Context) error {
	usr, err := NewUserValue(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	r, err := ReqValue(ctx)
	if err != nil {
		return err
	}
	if err = u.repoU.DelUserIDAndIP(ctx, usr.Uuid, r.RemoteAddr); err != nil {
		return err
	}
	return nil
}

func (u *userService) GetMsg(ctx context.Context) (string, error) {
	_, err := NewUserValue(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return "hello", nil
}
