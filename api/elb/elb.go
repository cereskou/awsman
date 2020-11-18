package elb

import (
	"net/http"

	"ditto.co.jp/awsman/api"
	"ditto.co.jp/awsman/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/labstack/echo/v4"
)

// elbv2ListLoadBalancers - ELB一覧を取得します
// @Summary ELB一覧を取得します
// @Description ELB一覧を取得します
// @Tags ELBV2
// @Accept json
// @Produce json
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /elbv2/listelb [get]
func (s *Service) elbv2ListLoadBalancers(c echo.Context) error {
	// token := c.QueryParam("token")
	// name := c.QueryParam("name")

	resp := api.Response{
		Code: 20000,
		Data: "success",
	}

	elb := s.Aws().ELBV2()

	input := elbv2.DescribeLoadBalancersInput{}

	result, err := elb.DescribeLoadBalancers(&input)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusInternalServerError, resp)
	}
	lbs := make([]*LoadBalancer, 0)
	for _, lb := range result.LoadBalancers {
		b := &LoadBalancer{
			Arn:                   aws.StringValue(lb.LoadBalancerArn),
			Name:                  aws.StringValue(lb.LoadBalancerName),
			IPAddressType:         aws.StringValue(lb.IpAddressType),
			VpcID:                 aws.StringValue(lb.VpcId),
			Type:                  aws.StringValue(lb.Type),
			Scheme:                aws.StringValue(lb.Scheme),
			DNSName:               aws.StringValue(lb.DNSName),
			CreatedTime:           utils.FormatTime(aws.TimeValue(lb.CreatedTime), ""),
			CustomerOwnedIpv4Pool: aws.StringValue(lb.CustomerOwnedIpv4Pool),
			CanonicalHostedZoneID: aws.StringValue(lb.CanonicalHostedZoneId),
			AvailabilityZones:     make([]*AvailabilityZone, 0),
			SecurityGroups:        make([]string, 0),
			State: &State{
				Code:   aws.StringValue(lb.State.Code),
				Reason: aws.StringValue(lb.State.Reason),
			},
		}
		//SecurityGroups
		for _, g := range lb.SecurityGroups {
			b.SecurityGroups = append(b.SecurityGroups, aws.StringValue(g))
		}
		//AvailabilityZones
		for _, z := range lb.AvailabilityZones {
			az := &AvailabilityZone{
				ZoneName:  aws.StringValue(z.ZoneName),
				SubnetID:  aws.StringValue(z.SubnetId),
				OutpostID: aws.StringValue(z.OutpostId),
				Addresses: make([]*Address, 0),
			}
			for _, addr := range z.LoadBalancerAddresses {
				ad := &Address{
					AllocationID:       aws.StringValue(addr.AllocationId),
					IPAddress:          aws.StringValue(addr.IpAddress),
					PrivateIPv4Address: aws.StringValue(addr.PrivateIPv4Address),
				}
				az.Addresses = append(az.Addresses, ad)
			}

			b.AvailabilityZones = append(b.AvailabilityZones, az)

		}
		lbs = append(lbs, b)
	}

	resp.Data = lbs

	return c.JSON(http.StatusOK, resp)
}

// elbv2ListLoadBalancers - リーセンナー一覧取得
// @Summary リーセンナー一覧取得
// @Description リーセンナー一覧取得
// @Tags ELBV2
// @Accept json
// @Produce json
// @Param arn query string true "LoadBalancerArn"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /elbv2/listeners [get]
func (s *Service) elbv2GetListeners(c echo.Context) error {
	resp := api.Response{
		Code: 20000,
		Data: "success",
	}
	arn := c.QueryParam("arn")

	elb := s.Aws().ELBV2()

	input := elbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(arn),
	}

	result, err := elb.DescribeListeners(&input)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusInternalServerError, resp)
	}

	listeners := make([]*Listener, 0)
	for _, l := range result.Listeners {
		listener := &Listener{
			Arn:             aws.StringValue(l.ListenerArn),
			LoadBalancerArn: aws.StringValue(l.LoadBalancerArn),
			Port:            aws.Int64Value(l.Port),
			Protocol:        aws.StringValue(l.Protocol),
			SslPolicy:       aws.StringValue(l.SslPolicy),
		}

		listeners = append(listeners, listener)
	}

	resp.Data = listeners

	return c.JSON(http.StatusOK, resp)
}

// elbv2GetRules - ルール一覧取得
// @Summary ルール一覧取得
// @Description ルール一覧取得
// @Tags ELBV2
// @Accept json
// @Produce json
// @Param arn query string true "ListenerArn"
// @Success 200 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.HTTPError
// @Security ApiKeyAuth
// @Router /elbv2/rules [get]
func (s *Service) elbv2GetRules(c echo.Context) error {
	resp := api.Response{
		Code: 20000,
		Data: "success",
	}
	arn := c.QueryParam("arn")

	elb := s.Aws().ELBV2()
	input := elbv2.DescribeRulesInput{
		ListenerArn: aws.String(arn),
	}
	result, err := elb.DescribeRules(&input)
	if err != nil {
		resp.Code = 50000
		resp.Data = err.Error()

		return c.JSON(http.StatusInternalServerError, resp)
	}

	rules := make([]*Rule, 0)
	for _, r := range result.Rules {
		n := &Rule{
			Arn:        aws.StringValue(r.RuleArn),
			Priority:   aws.StringValue(r.Priority),
			IsDefault:  aws.BoolValue(r.IsDefault),
			Conditions: make([]*RuleCondition, 0),
			Actions:    make([]*Action, 0),
		}
		//Conditions -
		for _, c := range r.Conditions {
			nc := &RuleCondition{
				Field:  aws.StringValue(c.Field),
				Values: make([]string, 0),
			}
			if c.HostHeaderConfig != nil {
				nc.HostHeaderConfig = &HostHeaderConditionConfig{
					Values: make([]string, 0),
				}
				for _, v := range c.HostHeaderConfig.Values {
					nc.HostHeaderConfig.Values = append(nc.HostHeaderConfig.Values, aws.StringValue(v))
				}
			}
			if c.HttpHeaderConfig != nil {
				nc.HTTPHeaderConfig = &HTTPHeaderConditionConfig{
					HTTPHeaderName: aws.StringValue(c.HttpHeaderConfig.HttpHeaderName),
					Values:         make([]string, 0),
				}
				for _, v := range c.HttpHeaderConfig.Values {
					nc.HTTPHeaderConfig.Values = append(nc.HTTPHeaderConfig.Values, aws.StringValue(v))
				}
			}
			if c.HttpRequestMethodConfig != nil {
				nc.HTTPRequestMethodConfig = &HTTPRequestMethodConditionConfig{
					Values: make([]string, 0),
				}
				for _, v := range c.HttpRequestMethodConfig.Values {
					nc.HTTPRequestMethodConfig.Values = append(nc.HTTPRequestMethodConfig.Values, aws.StringValue(v))
				}
			}
			if c.PathPatternConfig != nil {
				nc.PathPatternConfig = &PathPatternConditionConfig{
					Values: make([]string, 0),
				}
				for _, v := range c.PathPatternConfig.Values {
					nc.PathPatternConfig.Values = append(nc.PathPatternConfig.Values, aws.StringValue(v))
				}
			}
			if c.SourceIpConfig != nil {
				nc.SourceIPConfig = &SourceIPConditionConfig{
					Values: make([]string, 0),
				}
				for _, v := range c.SourceIpConfig.Values {
					nc.SourceIPConfig.Values = append(nc.SourceIPConfig.Values, aws.StringValue(v))
				}
			}
			if c.QueryStringConfig != nil {
				nc.QueryStringConfig = &QueryStringConditionConfig{
					Values: make([]*QueryStringKeyValuePair, 0),
				}
				for _, v := range c.QueryStringConfig.Values {
					nv := &QueryStringKeyValuePair{
						Key:   aws.StringValue(v.Key),
						Value: aws.StringValue(v.Value),
					}
					nc.QueryStringConfig.Values = append(nc.QueryStringConfig.Values, nv)
				}
			}

			for _, rcval := range c.Values {
				nc.Values = append(nc.Values, aws.StringValue(rcval))
			}
			n.Conditions = append(n.Conditions, nc)
		}
		//Actions -
		for _, a := range r.Actions {
			na := &Action{
				Type:           aws.StringValue(a.Type),
				TargetGroupArn: aws.StringValue(a.TargetGroupArn),
				Order:          aws.Int64Value(a.Order),
			}

			if a.AuthenticateCognitoConfig != nil {
				na.AuthenticateCognitoConfig = &AuthenticateCognitoActionConfig{
					AuthenticationRequestExtraParams: make(map[string]string),
					OnUnauthenticatedRequest:         aws.StringValue(a.AuthenticateCognitoConfig.OnUnauthenticatedRequest),
					Scope:                            aws.StringValue(a.AuthenticateCognitoConfig.Scope),
					SessionCookieName:                aws.StringValue(a.AuthenticateCognitoConfig.SessionCookieName),
					SessionTimeout:                   aws.Int64Value(a.AuthenticateCognitoConfig.SessionTimeout),
					UserPoolArn:                      aws.StringValue(a.AuthenticateCognitoConfig.UserPoolArn),
					UserPoolClientID:                 aws.StringValue(a.AuthenticateCognitoConfig.UserPoolClientId),
					UserPoolDomain:                   aws.StringValue(a.AuthenticateCognitoConfig.UserPoolDomain),
				}
			}
			if a.AuthenticateOidcConfig != nil {
				na.AuthenticateOidcConfig = &AuthenticateOidcActionConfig{
					AuthenticationRequestExtraParams: make(map[string]string),
					AuthorizationEndpoint:            aws.StringValue(a.AuthenticateOidcConfig.AuthorizationEndpoint),
					ClientID:                         aws.StringValue(a.AuthenticateOidcConfig.ClientId),
					ClientSecret:                     aws.StringValue(a.AuthenticateOidcConfig.ClientSecret),
					Issuer:                           aws.StringValue(a.AuthenticateOidcConfig.Issuer),
					OnUnauthenticatedRequest:         aws.StringValue(a.AuthenticateOidcConfig.OnUnauthenticatedRequest),
					Scope:                            aws.StringValue(a.AuthenticateOidcConfig.Scope),
					SessionCookieName:                aws.StringValue(a.AuthenticateOidcConfig.SessionCookieName),
					SessionTimeout:                   aws.Int64Value(a.AuthenticateOidcConfig.SessionTimeout),
					TokenEndpoint:                    aws.StringValue(a.AuthenticateOidcConfig.TokenEndpoint),
					UseExistingClientSecret:          aws.BoolValue(a.AuthenticateOidcConfig.UseExistingClientSecret),
					UserInfoEndpoint:                 aws.StringValue(a.AuthenticateOidcConfig.UserInfoEndpoint),
				}
			}
			if a.FixedResponseConfig != nil {
				na.FixedResponseConfig = &FixedResponseActionConfig{
					ContentType: aws.StringValue(a.FixedResponseConfig.ContentType),
					MessageBody: aws.StringValue(a.FixedResponseConfig.MessageBody),
					StatusCode:  aws.StringValue(a.FixedResponseConfig.StatusCode),
				}
			}
			if a.ForwardConfig != nil {
				na.ForwardConfig = &ForwardActionConfig{
					TargetGroupStickinessConfig: &TargetGroupStickinessConfig{
						DurationSeconds: aws.Int64Value(a.ForwardConfig.TargetGroupStickinessConfig.DurationSeconds),
						Enabled:         aws.BoolValue(a.ForwardConfig.TargetGroupStickinessConfig.Enabled),
					},
					TargetGroups: make([]*TargetGroupTuple, 0),
				}
				//ForwardConfig.TargetGroups -
				for _, tg := range a.ForwardConfig.TargetGroups {
					gt := &TargetGroupTuple{
						TargetGroupArn: aws.StringValue(tg.TargetGroupArn),
						Weight:         aws.Int64Value(tg.Weight),
					}
					na.ForwardConfig.TargetGroups = append(na.ForwardConfig.TargetGroups, gt)
				}
			}
			if a.RedirectConfig != nil {
				na.RedirectConfig = &RedirectActionConfig{
					Host:       aws.StringValue(a.RedirectConfig.Host),
					Path:       aws.StringValue(a.RedirectConfig.Path),
					Port:       aws.StringValue(a.RedirectConfig.Port),
					Protocol:   aws.StringValue(a.RedirectConfig.Protocol),
					Query:      aws.StringValue(a.RedirectConfig.Query),
					StatusCode: aws.StringValue(a.RedirectConfig.StatusCode),
				}
			}
			n.Actions = append(n.Actions, na)
		}

		rules = append(rules, n)
	}

	resp.Data = rules

	return c.JSON(http.StatusOK, resp)
}
