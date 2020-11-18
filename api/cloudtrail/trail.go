package cloudtrail

import (
	"fmt"
	"net/http"

	"ditto.co.jp/awsman/api"
	"ditto.co.jp/awsman/model"
	"ditto.co.jp/awsman/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/labstack/echo/v4"
)

// cloudtrailListEvent - イベント取得
// @Summary イベント取得
// @Description イベント取得
// @Tags CloudTrail
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /cloudtrail/events [get]
func (s *Service) cloudtrailListEvent(c echo.Context) error {
	token := c.QueryParam("token")
	name := c.QueryParam("name")

	resp := api.Response{
		Code: 20000,
		Data: "success",
	}
	var awstoken *string
	if token != "" {
		awstoken = aws.String(token)
	}
	trail := s.Aws().CloudTrail()
	input := &cloudtrail.LookupEventsInput{
		// LookupAttributes: make([]*cloudtrail.LookupAttribute, 0),
		NextToken: awstoken,
	}
	if name != "" {
		input.LookupAttributes = make([]*cloudtrail.LookupAttribute, 0)

		attr := cloudtrail.LookupAttribute{
			AttributeKey:   aws.String("EventName"),
			AttributeValue: aws.String(name),
		}
		input.LookupAttributes = append(input.LookupAttributes, &attr)
	}

	result, err := trail.LookupEvents(input)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusOK, resp)
	}

	response := &Response{
		NextToken: aws.StringValue(result.NextToken),
		Events:    make([]*Event, 0),
	}

	for _, event := range result.Events {
		m := model.CloudTrailEvent{
			EventID:         aws.StringValue(event.EventId),
			AccessKeyID:     aws.StringValue(event.AccessKeyId),
			CloudTrailEvent: aws.StringValue(event.CloudTrailEvent),
			EventName:       aws.StringValue(event.EventName),
			EventSource:     aws.StringValue(event.EventSource),
			EventTime:       utils.FormatTime(aws.TimeValue(event.EventTime), ""),
			Username:        aws.StringValue(event.Username),
		}
		err := s.DB().CreateCloudTrailEvent(&m)
		if err != nil {
			return err
		}
	}

	for _, event := range result.Events {
		fmt.Println(aws.StringValue(event.CloudTrailEvent))
		e := &Event{
			EventID:         aws.StringValue(event.EventId),
			AccessKeyID:     aws.StringValue(event.AccessKeyId),
			CloudTrailEvent: aws.StringValue(event.CloudTrailEvent),
			EventName:       aws.StringValue(event.EventName),
			EventSource:     aws.StringValue(event.EventSource),
			EventTime:       utils.FormatTime(aws.TimeValue(event.EventTime), ""),
			Username:        aws.StringValue(event.Username),
			Resources:       make([]*Resource, 0),
		}
		for _, resource := range event.Resources {
			e.Resources = append(e.Resources, &Resource{
				Name: aws.StringValue(resource.ResourceName),
				Type: aws.StringValue(resource.ResourceType),
			})
		}

		response.Events = append(response.Events, e)
	}

	resp.Data = response

	return c.JSON(http.StatusOK, resp)
}
