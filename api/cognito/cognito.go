package cognito

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ditto.co.jp/awsman/api"
	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/cx"
	"ditto.co.jp/awsman/logger"
	"ditto.co.jp/awsman/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/labstack/echo/v4"
)

//cognitoImportUsers -
func (s *Service) cognitoImportUsers() (int64, error) {
	count := int64(0)
	conf, _ := config.Load()
	cognito := s.Aws().Cognito()

	//get groups
	groupinput := cognitoidentityprovider.ListGroupsInput{
		UserPoolId: aws.String(conf.Cognito.UserPoolID),
	}
	groups, err := cognito.ListGroups(&groupinput)
	if err != nil {
		return 0, err
	}

	err = s.DB().DeleteCognitoGroup("")
	if err != nil {
		logger.Error(err)
	}
	err = s.DB().DeleteCognitoGroup("")
	if err != nil {
		logger.Error(err)
	}

	modelgroups := make([]*model.CognitoGroup, 0)
	for _, g := range groups.Groups {
		group := model.CognitoGroup{
			UserPoolID:  *g.UserPoolId,
			Name:        *g.GroupName,
			Description: *g.Description,
		}
		err = s.DB().CreateCognitoGroup(&group)
		if err != nil {
			return 0, err
		}

		modelgroups = append(modelgroups, &group)
	}

	// map[userid] - map[group]
	usergroup := make(map[string]map[string]int)
	for _, g := range modelgroups {
		listuserinput := cognitoidentityprovider.ListUsersInGroupInput{
			GroupName:  aws.String(g.Name),
			UserPoolId: aws.String(g.UserPoolID),
		}
		listuseroutput, err := cognito.ListUsersInGroup(&listuserinput)
		if err != nil {
			return 0, err
		}
		//users
		for _, u := range listuseroutput.Users {
			if gs, ok := usergroup[*u.Username]; !ok {
				gs = make(map[string]int)
				gs[g.Name] = 1

				usergroup[*u.Username] = gs
			} else {
				gs[g.Name] = 1
			}
		}
	}

	input := cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(conf.Cognito.UserPoolID),
	}
	output, err := cognito.ListUsers(&input)
	if err != nil {
		return 0, err
	}
	for _, u := range output.Users {
		user := model.CognitoUser{
			UUID:       *u.Username,
			UserStatus: *u.UserStatus,
			Enabled:    *u.Enabled,
			UserPoolID: conf.Cognito.UserPoolID,
			CreateDate: (*u.UserCreateDate).Unix(),
			UpdateDate: (*u.UserLastModifiedDate).Unix(),
		}

		//attributes
		for _, a := range u.Attributes {
			if *a.Name == "name" {
				user.Name = *a.Value
			}
			if *a.Name == "email" {
				user.Email = *a.Value
			}
			if *a.Name == "sub" {
				user.Sub = *a.Value
			}
		}
		//group
		if gs, ok := usergroup[*u.Username]; ok {
			keys := make([]string, 0)
			for k := range gs {
				keys = append(keys, k)
			}
			user.GroupName = strings.Join(keys, ",")
		}

		err = s.DB().CreateCognitoUser(&user)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

// cognitoListUsers - ユーザー一覧をAWS Cognitoと同期する
// @Summary ユーザー一覧をAWS Cognitoと同期する
// @Description ユーザー一覧をAWS Cognitoと同期する
// @Tags Cognito
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /cognito/user/sync [get]
func (s *Service) cognitoSyncUser(c echo.Context) error {
	resp := api.Response{
		Code: 20000,
		Data: "success",
	}

	//clear
	err := s.DB().DeleteCognitoUser("")
	if err != nil {
		return err
	}
	err = s.DB().DeleteCognitoGroup("")
	if err != nil {
		return err
	}

	//get data from aws cognito
	_, err = s.cognitoImportUsers()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

// cognitoListUsers - ユーザー一覧を取得する
// @Summary ユーザー一覧を取得する
// @Description ユーザー一覧を取得する
// @Tags Cognito
// @Accept json
// @Produce json
// @Param email query string false "Email"
// @Param group query string false "group"
// @Param page query int false "page number" default(1)
// @Param limit query int false "limit" default(10)
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /cognito/user/list [get]
func (s *Service) cognitoListUser(c echo.Context) error {
	param := c.QueryParam("limit")
	limit, _ := strconv.ParseInt(param, 10, 64)
	param = c.QueryParam("page")
	page, _ := strconv.ParseInt(param, 10, 64)
	sort := c.QueryParam("sort")

	offset := limit * (page - 1)
	input := cx.CognitoUserSearch{
		Email:  c.QueryParam("email"),
		Group:  c.QueryParam("group"),
		Offset: int(offset),
		Limit:  int(limit),
		Sort:   sort,
	}

	count, err := s.DB().GetCognitoUserCount(&input)
	if err != nil {
		return err
	}
	if count == 0 {
		count, err = s.cognitoImportUsers()
		if err != nil {
			return err
		}
	}
	conf, _ := config.Load()
	cognito := s.Aws().Cognito()
	coginput := cognitoidentityprovider.DescribeUserPoolInput{
		UserPoolId: aws.String(conf.Cognito.UserPoolID),
	}
	cogoutput, err := cognito.DescribeUserPool(&coginput)
	if err != nil {
		return err
	}
	lmdate := cogoutput.UserPool.LastModifiedDate
	fmt.Println(lmdate)

	groups, err := s.DB().GetCognitoGroup(conf.Cognito.UserPoolID)
	if err != nil {
		return err
	}

	users, err := s.DB().GetCognitoUsers(&input)
	if err != nil {
		return err
	}
	data := Response{
		Total:  count,
		Items:  users,
		Groups: groups,
	}
	resp := api.Response{
		Code: 20000,
		Data: data,
	}
	return c.JSON(http.StatusOK, resp)
}

// cognitoUpdateUser - ユーザー情報を更新する
// @Summary ユーザー情報を更新する
// @Description ユーザー情報を更新する
// @Tags Cognito
// @Produce json
// @Param uuid path string true "uuid"
// @Param data body model.CognitoUser false "data"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Security ApiKeyAuth
// @Router /cognito/user/{uuid} [put]
func (s *Service) cognitoUpdateUser(c echo.Context) error {
	uuid := c.Param("uuid")
	fmt.Println(uuid)

	var input = model.CognitoUser{}
	if err := c.Bind(&input); err != nil {
		return err
	}
	resp := api.Response{
		Code: 20000,
		Data: "success",
	}

	curuser, err := s.DB().GetCognitoUser(uuid)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}
	err = s.DB().CreateCognitoUser(&input)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}
	ua := make([]*cognitoidentityprovider.AttributeType, 0)
	if curuser.Name != input.Name {
		ua = append(ua, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("name"),
			Value: aws.String(input.Name),
		})
	}
	if curuser.Email != input.Email {
		ua = append(ua, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("email"),
			Value: aws.String(input.Email),
		})
	}
	//update
	conf, _ := config.Load()
	cognito := s.Aws().Cognito()
	if len(ua) > 0 {
		attinput := cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserPoolId:     aws.String(conf.Cognito.UserPoolID),
			Username:       aws.String(uuid),
			UserAttributes: ua,
		}
		_, err = cognito.AdminUpdateUserAttributes(&attinput)
		if err != nil {
			resp.Code = 50000
			resp.Data = err.Error()

			return c.JSON(http.StatusOK, resp)
		}
	}
	if curuser.GroupName != input.GroupName {
		if curuser.GroupName != "" {
			rminput := cognitoidentityprovider.AdminRemoveUserFromGroupInput{
				UserPoolId: aws.String(conf.Cognito.UserPoolID),
				Username:   aws.String(uuid),
				GroupName:  aws.String(curuser.GroupName),
			}
			_, err = cognito.AdminRemoveUserFromGroup(&rminput)
			if err != nil {
				resp.Code = 50000
				resp.Data = err.Error()

				return c.JSON(http.StatusOK, resp)
			}
		}
		if input.GroupName != "" {
			grpinput := cognitoidentityprovider.AdminAddUserToGroupInput{
				UserPoolId: aws.String(conf.Cognito.UserPoolID),
				Username:   aws.String(uuid),
				GroupName:  aws.String(input.GroupName),
			}
			_, err = cognito.AdminAddUserToGroup(&grpinput)
			if err != nil {
				resp.Code = 50000
				resp.Data = err.Error()

				return c.JSON(http.StatusOK, resp)
			}

		}
	}
	//有効無効
	if curuser.Enabled != input.Enabled {
		if input.Enabled {
			input := cognitoidentityprovider.AdminEnableUserInput{
				UserPoolId: aws.String(conf.Cognito.UserPoolID),
				Username:   aws.String(uuid),
			}
			_, err = cognito.AdminEnableUser(&input)
		} else {
			input := cognitoidentityprovider.AdminDisableUserInput{
				UserPoolId: aws.String(conf.Cognito.UserPoolID),
				Username:   aws.String(uuid),
			}
			_, err = cognito.AdminDisableUser(&input)
		}
		if err != nil {
			resp.Code = 50000
			resp.Data = err.Error()

			return c.JSON(http.StatusOK, resp)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// cognitoDeleteUser - ユーザー情報を削除する
// @Summary ユーザー情報を削除する
// @Description ユーザー情報を削除する
// @Tags Cognito
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Security ApiKeyAuth
// @Router /cognito/user/{uuid} [delete]]
func (s *Service) cognitoDeleteUser(c echo.Context) error {
	uuid := c.Param("uuid")
	fmt.Println(uuid)

	resp := api.Response{
		Code: 20000,
		Data: "success",
	}

	err := s.DB().DeleteCognitoUser(uuid)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}
	conf, _ := config.Load()
	cognito := s.Aws().Cognito()
	input := cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(conf.Cognito.UserPoolID),
		Username:   aws.String(uuid),
	}
	_, err = cognito.AdminDeleteUser(&input)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()
		fmt.Println(err)

		return c.JSON(http.StatusOK, resp)
	}

	return c.JSON(http.StatusOK, resp)
}

// cognitoEnableUser - ユーザーの有効無効
// @Summary ユーザーの有効無効
// @Description ユーザーの有効無効
// @Tags Cognito
// @Produce json
// @Param uuid path string true "uuid"
// @Param status query string true "string enums" Enums(enabled, disabled)
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Security ApiKeyAuth
// @Router /cognito/user/enable/{uuid} [post]
func (s *Service) cognitoEnableUser(c echo.Context) error {
	uuid := c.Param("uuid")
	fmt.Println(uuid)
	status := c.QueryParam("status")

	resp := api.Response{
		Code: 20000,
		Data: "success",
	}
	err := s.DB().EnableCognitoUser(uuid, status)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}
	conf, _ := config.Load()
	cognito := s.Aws().Cognito()
	switch strings.ToLower(status) {
	case "enabled":
		input := cognitoidentityprovider.AdminEnableUserInput{
			UserPoolId: aws.String(conf.Cognito.UserPoolID),
			Username:   aws.String(uuid),
		}
		_, err = cognito.AdminEnableUser(&input)
	case "disabled":
		input := cognitoidentityprovider.AdminDisableUserInput{
			UserPoolId: aws.String(conf.Cognito.UserPoolID),
			Username:   aws.String(uuid),
		}
		_, err = cognito.AdminDisableUser(&input)
	}
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}

	return c.JSON(http.StatusOK, resp)
}

// cognitoCreateUser - ユーザー情報の新規作成
// @Summary ユーザー情報の新規作成
// @Description ユーザー情報の新規作成
// @Tags Cognito
// @Produce json
// @Param data body model.CognitoUser true "data"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Security ApiKeyAuth
// @Router /cognito/user [POST]
func (s *Service) cognitoCreateUser(c echo.Context) error {
	var input = model.CognitoUser{}
	if err := c.Bind(&input); err != nil {
		return err
	}
	resp := api.Response{
		Code: 20000,
		Data: "success",
	}

	conf, _ := config.Load()
	cognito := s.Aws().Cognito()
	userinput := cognitoidentityprovider.AdminCreateUserInput{
		DesiredDeliveryMediums: []*string{
			aws.String("EMAIL"),
		},
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.Email),
			},
		},
		UserPoolId: aws.String(conf.Cognito.UserPoolID),
		Username:   aws.String(input.Email),
	}
	if input.Name != "" {
		userinput.UserAttributes = append(userinput.UserAttributes, &cognitoidentityprovider.AttributeType{
			Name:  aws.String("name"),
			Value: aws.String(input.Name),
		})
	}
	useroutput, err := cognito.AdminCreateUser(&userinput)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()
		fmt.Println(err)

		return c.JSON(http.StatusOK, resp)
	}
	if input.GroupName != "" {
		grpinput := cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(conf.Cognito.UserPoolID),
			Username:   useroutput.User.Username,
			GroupName:  aws.String(input.GroupName),
		}

		_, err = cognito.AdminAddUserToGroup(&grpinput)
		if err != nil {
			resp.Code = 50000
			resp.Data = err.Error()

			return c.JSON(http.StatusOK, resp)
		}
	}
	user := model.CognitoUser{
		UUID:       *useroutput.User.Username,
		UserStatus: *useroutput.User.UserStatus,
		Enabled:    *useroutput.User.Enabled,
		GroupName:  input.GroupName,
		UserPoolID: conf.Cognito.UserPoolID,
		CreateDate: useroutput.User.UserCreateDate.Unix(),
		UpdateDate: useroutput.User.UserLastModifiedDate.Unix(),
	}
	//attributes
	for _, a := range useroutput.User.Attributes {
		if *a.Name == "name" {
			user.Name = *a.Value
		}
		if *a.Name == "email" {
			user.Email = *a.Value
		}
		if *a.Name == "sub" {
			user.Sub = *a.Value
		}
	}
	err = s.DB().CreateCognitoUser(&user)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}

	return c.JSON(http.StatusOK, resp)
}

// cognitoSetPassword - パスワードリセット通知
// @Summary パスワードリセット通知
// @Description パスワードリセット通知
// @Tags Cognito
// @Produce json
// @Param uuid path string true "uuid"
// @Param data body model.CognitoUser true "data"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Security ApiKeyAuth
// @Router /cognito/setpassword/{uuid} [post]]
func (s *Service) cognitoSetPassword(c echo.Context) error {
	uuid := c.Param("uuid")

	var input = ChangePasswordRequest{}
	if err := c.Bind(&input); err != nil {
		return err
	}
	resp := api.Response{
		Code: 20000,
		Data: "success",
	}
	permanent := false
	if input.Permanent == "1" {
		permanent = true
	}
	conf, _ := config.Load()
	cognito := s.Aws().Cognito()
	passinput := cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(conf.Cognito.UserPoolID),
		Username:   aws.String(uuid),
		Password:   aws.String(input.Password),
		Permanent:  aws.Bool(permanent),
	}
	_, err := cognito.AdminSetUserPassword(&passinput)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}

	return c.JSON(http.StatusOK, resp)
}
