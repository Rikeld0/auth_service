package usecase

import (
	"auth/pkg/gen_key"
	"auth/pkg/myJwt"
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

var tokenUser = map[string]string{}
var userIP = map[string]string{}

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
	repoU  repo.Users
	repoUK repo.UsersKey
}

func NewUserService(repoU repo.Users, repoUK repo.UsersKey) User {
	return &userService{
		repoU:  repoU,
		repoUK: repoUK,
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
	uuid, ok := userIP[r.RemoteAddr]
	if !ok {
		return nil, errors.New("error ip")
	}
	keys, err := u.repoUK.Get(ctx, uuid)
	c := structs.Cl{}
	token, err := getToken(r)
	if err != nil {
		return nil, err
	}
	claims, err := myJwt.VerefyJWT(token, keys)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(claims), &c)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	usr, err := u.repoU.GetUserId(ctx, c.Iss)
	return context.WithValue(r.Context(), "user", usr), nil
}

func (u *userService) SignUp(ctx context.Context, user *structs.User) error {
	r, err := ReqValue(ctx)
	if err != nil {
		return err
	}
	uuid, err := u.repoU.InsertUser(ctx, user)
	if err != nil {
		return err
	}
	userIP[r.RemoteAddr] = uuid
	priv, pub := gen_key.GenKey()
	err = u.repoUK.Put(ctx, uuid, string(priv), string(pub))
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) SignIn(ctx context.Context, user *structs.User) (*myJwt.JWT, error) {
	r, err := ReqValue(ctx)
	if err != nil {
		return nil, err
	}
	usr, err := u.repoU.Get(ctx, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	userIP[r.RemoteAddr] = usr.Uuid
	keys, err := u.repoUK.Get(ctx, usr.Uuid)
	jwtToken, err := myJwt.GenerateJWT(usr.Uuid, keys)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//tokenUser[usr.Uuid] = jwtToken.AccessToken
	return jwtToken, nil
}

func (u *userService) GetMsg(ctx context.Context) (string, error) {
	usr, err := NewUserValue(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_ = usr
	//if _, ok := tokenUser[usr.Uuid]; !ok {
	//	fmt.Println("no")
	//	return "", errors.New("")
	//}
	return "hello", nil
}
