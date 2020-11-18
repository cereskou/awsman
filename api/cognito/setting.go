package cognito

import (
	"net/http"

	"ditto.co.jp/awsman/api"
	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/model"
	"github.com/labstack/echo/v4"
)

// cognitoSettings - 設定情報取得
// @Summary 設定情報取得
// @Description 設定情報取得
// @Tags Cognito
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Security ApiKeyAuth
// @Router /cognito/setting [get]
func (s *Service) cognitoSettings(c echo.Context) error {
	conf, _ := config.Load()
	settings, err := s.DB().GetCognitoSettings()
	if err != nil {
		return err
	}
	haskey := false
	for _, kv := range settings {
		if kv.Key == "userpool" {
			haskey = true
		}
	}
	if !haskey {
		settings = append(settings, &model.Setting{
			Category: "cognito",
			Key:      "userpool",
			Value:    conf.Cognito.UserPoolID,
		})
	}

	resp := api.Response{
		Code: 20000,
		Data: settings,
	}
	return c.JSON(http.StatusOK, &resp)
}

// cognitoSaveSettings - 設定情報保存
// @Summary 設定情報保存
// @Description 設定情報保存
// @Tags Cognito
// @Produce json
// @Param data body array false "data"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Security ApiKeyAuth
// @Router /cognito/setting [post]
func (s *Service) cognitoSaveSettings(c echo.Context) error {
	settings := make([]*model.Setting, 0)
	if err := c.Bind(&settings); err != nil {
		return err
	}

	conf, _ := config.Load()
	for _, sv := range settings {
		err := s.DB().CreateUpdateSetting(sv)
		if err != nil {
			return err
		}
		if sv.Key == "userpool" {
			conf.Cognito.UserPoolID = sv.Value
		}
	}

	resp := api.Response{
		Code: 20000,
		Data: "success",
	}
	return c.JSON(http.StatusOK, &resp)
}
