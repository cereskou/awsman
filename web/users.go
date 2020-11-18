package web

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"ditto.co.jp/awsman/api"
	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/model"
	"ditto.co.jp/awsman/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// UsersInfo - ユーザー情報取得
// @Summary ユーザー情報取得
// @Description ユーザー情報取得
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /info [get]
func (s *Service) UsersInfo(c echo.Context) error {
	token := c.QueryParam("token")
	fmt.Println(token)
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("Unexpected signing methd")
		}
		return s._rsa.GetPublicKey(), nil
	})
	if err != nil || !t.Valid {
		return err
	}

	// user := c.Get("user").(*jwt.Token)
	claims := t.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	ui, err := s.DB().GetUser(name)
	if err != nil {
		return err
	}
	ui.Roles = strings.Split(ui.Role, ",")
	resp := api.Response{
		Code: 20000,
		Data: ui,
	}
	return c.JSON(http.StatusOK, resp)
}

// CreateUser - ユーザー作成
// @Summary ユーザー作成
// @Description ユーザー作成
// @Tags User
// @Accept json
// @Produce json
// @Param data body model.User false "data"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /user [post]
func (s *Service) CreateUser(c echo.Context) error {
	var input = model.User{}
	if err := c.Bind(&input); err != nil {
		return err
	}
	conf, _ := config.Load()
	input.Password = utils.HashPassword(input.Password)
	if input.Avatar == "" {
		input.Avatar = conf.Avatar
	}
	err := s.DB().CreateUser(&input)
	if err != nil {
		return err
	}

	resp := api.Response{
		Code: 20000,
		Data: input,
	}

	return c.JSON(http.StatusOK, resp)

}

// Login - ログイン
// @Summary ログイン
// @Description ログイン
// @Tags User
// @Accept json
// @Produce json
// @Param data body Login false "data"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Router /user/login [post]
func (s *Service) Login(c echo.Context) error {
	var input = Login{}
	if err := c.Bind(&input); err != nil {
		return err
	}

	user, err := s.DB().GetUser(input.Username)
	if err != nil {
		resp := api.Response{
			Code: 40001,
			Data: "Invalid Username or Password",
		}
		return c.JSON(http.StatusUnauthorized, resp)
	}

	if !utils.CompareHashedPassword(user.Password, input.Password) {
		resp := api.Response{
			Code: 40001,
			Data: "Invalid Username or Password",
		}
		return c.JSON(http.StatusUnauthorized, resp)
	}

	token, err := s.generateToken(user.Name, true)
	if err != nil {
		return err
	}
	resp := api.Response{
		Code: 20000,
		Data: token,
	}

	// c.SetCookie(&http.Cookie{
	// 	Name:  "Authorization",
	// 	Value: token.AccessToken,
	// 	Path:  "/",
	// 	//Domain:   "localhost",
	// 	MaxAge:   86400,
	// 	HttpOnly: true,
	// })

	// return c.JSON(http.StatusForbidden, map[string]string{
	// 	"data": "Wrong password or username",
	// })
	return c.JSON(http.StatusOK, resp)
}

// Logout - ログアウト
// @Summary ログアウト
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /user/logout [post]
func (s *Service) Logout(c echo.Context) error {
	user := User{
		ID:     "id",
		Name:   "mail@company.com",
		Roles:  []string{"admin"},
		Avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
	}
	resp := api.Response{
		Code: 20000,
		Data: user,
	}
	return c.JSON(http.StatusOK, resp)
}

//DeleteUser -
func (s *Service) DeleteUser(c echo.Context) error {
	users := []User{
		{
			ID:   "id0",
			Name: "mail@company.com",
		},
	}

	return c.JSON(http.StatusOK, users)
}
