package code

import (
	"net/http"

	"ditto.co.jp/awsman/api"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/labstack/echo/v4"
)

// codecommitListRepositories - リポジトリ一覧取得
// @Summary リポジトリ一覧取得
// @Description リポジトリ一覧取得
// @Tags CodeCommit
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /code/codecommit/repositories [get]
func (s *Service) codecommitListRepositories(c echo.Context) error {
	resp := api.Response{
		Code: 20000,
		Data: "success",
	}
	token := c.QueryParam("token")
	var awstoken *string
	if token != "" {
		awstoken = aws.String(token)
	}
	commit := s.Aws().CodeCommit()

	input := codecommit.ListRepositoriesInput{
		NextToken: awstoken,
	}
	result, err := commit.ListRepositories(&input)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}

	response := &Response{
		NextToken:    aws.StringValue(result.NextToken),
		Repositories: make([]*Repository, 0),
	}
	for _, r := range result.Repositories {
		rp := Repository{
			ID:      aws.StringValue(r.RepositoryId),
			Name:    aws.StringValue(r.RepositoryName),
			Branchs: make([]string, 0),
		}

		bi := codecommit.ListBranchesInput{
			RepositoryName: r.RepositoryName,
		}
		branchs, err := commit.ListBranches(&bi)
		if err != nil {
			return err
		}
		for _, b := range branchs.Branches {
			rp.Branchs = append(rp.Branchs, aws.StringValue(b))
		}

		response.Repositories = append(response.Repositories, &rp)
	}

	resp.Data = response

	return c.JSON(http.StatusOK, resp)
}
