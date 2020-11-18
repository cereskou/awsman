package elb

//State -
type State struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

//Address -
type Address struct {
	AllocationID       string `json:"allocationid"`
	IPAddress          string `json:"ipaddress"`
	PrivateIPv4Address string `json:"private_ipv4_address"`
}

//AvailabilityZone -
type AvailabilityZone struct {
	ZoneName  string     `json:"name"`
	SubnetID  string     `json:"subnetid"`
	OutpostID string     `json:"outpostid"`
	Addresses []*Address `json:"addresses"`
}

//LoadBalancer -
type LoadBalancer struct {
	Arn                   string              `json:"arn"`
	Name                  string              `json:"name"`
	IPAddressType         string              `json:"ipaddresstype"`
	VpcID                 string              `json:"vpcid"`
	Type                  string              `json:"type"`
	State                 *State              `json:"state,omitempty"`
	Scheme                string              `json:"scheme"`
	DNSName               string              `json:"dnsname"`
	CreatedTime           string              `json:"created_time"`
	CustomerOwnedIpv4Pool string              `json:"owned_ipv4_pool"`
	CanonicalHostedZoneID string              `json:"hosted_zone_id"`
	AvailabilityZones     []*AvailabilityZone `json:"availability_zone"`
	SecurityGroups        []string            `json:"security_groups"`
}

//Listener -
type Listener struct {
	Arn             string `json:"arn"`
	LoadBalancerArn string `json:"loadbalancerarn"`
	Port            int64  `json:"port"`
	Protocol        string `json:"protocol"`
	SslPolicy       string `json:"sslpolicy"`
}

//AuthenticateCognitoActionConfig -
type AuthenticateCognitoActionConfig struct {
	AuthenticationRequestExtraParams map[string]string `json:"extra_params"`
	OnUnauthenticatedRequest         string            `json:"unauthenticated_request"`
	Scope                            string            `json:"scope"`
	SessionCookieName                string            `json:"session_cookie_name"`
	SessionTimeout                   int64             `json:"session_timeout"`
	UserPoolArn                      string            `json:"user_pool_arn"`
	UserPoolClientID                 string            `json:"user_pool_client_id"`
	UserPoolDomain                   string            `json:"user_pool_domain"`
}

//AuthenticateOidcActionConfig -
type AuthenticateOidcActionConfig struct {
	AuthenticationRequestExtraParams map[string]string `json:"extra_params"`
	AuthorizationEndpoint            string            `json:"endpoint"`
	ClientID                         string            `json:"client_id"`
	ClientSecret                     string            `json:"client_secret"`
	Issuer                           string            `json:"issuer"`
	OnUnauthenticatedRequest         string            `json:"unauthenticated_request"`
	Scope                            string            `json:"scope"`
	SessionCookieName                string            `json:"session_cookie_name"`
	SessionTimeout                   int64             `json:"session_timeout"`
	TokenEndpoint                    string            `json:"token_endpoint"`
	UseExistingClientSecret          bool              `json:"use_existing_client_secret"`
	UserInfoEndpoint                 string            `json:"user_info_endpoint"`
}

//FixedResponseActionConfig -
type FixedResponseActionConfig struct {
	ContentType string `json:"content_type"`
	MessageBody string `json:"body"`
	StatusCode  string `json:"statuscode"`
}

// TargetGroupStickinessConfig -
type TargetGroupStickinessConfig struct {
	DurationSeconds int64 `json:"duration"`
	Enabled         bool  `json:"enabled"`
}

//TargetGroupTuple -
type TargetGroupTuple struct {
	TargetGroupArn string `json:"arn"`
	Weight         int64  `json:"weight"`
}

//ForwardActionConfig -
type ForwardActionConfig struct {
	TargetGroupStickinessConfig *TargetGroupStickinessConfig `json:"stickiness_config"`
	TargetGroups                []*TargetGroupTuple          `json:"target_groups"`
}

//RedirectActionConfig -
type RedirectActionConfig struct {
	Host       string `json:"host"`
	Path       string `json:"path"`
	Port       string `json:"port"`
	Protocol   string `json:"protocol"`
	Query      string `json:"query"`
	StatusCode string `json:"statuscode"`
}

//Action -
type Action struct {
	Type                      string                           `json:"type"`
	TargetGroupArn            string                           `json:"target_group_arn"`
	Order                     int64                            `json:"order"`
	AuthenticateCognitoConfig *AuthenticateCognitoActionConfig `json:"cognito_aconfig,omitempty"`
	AuthenticateOidcConfig    *AuthenticateOidcActionConfig    `json:"oidc_config,omitempty"`
	FixedResponseConfig       *FixedResponseActionConfig       `json:"fixed_response_config,omitempty"`
	ForwardConfig             *ForwardActionConfig             `json:"forward_config,omitempty"`
	RedirectConfig            *RedirectActionConfig            `json:"redirect_config,omitempty"`
}

// HostHeaderConditionConfig -
type HostHeaderConditionConfig struct {
	Values []string `json:"values"`
}

//HTTPHeaderConditionConfig -
type HTTPHeaderConditionConfig struct {
	HTTPHeaderName string   `json:"name"`
	Values         []string `json:"values"`
}

//HTTPRequestMethodConditionConfig -
type HTTPRequestMethodConditionConfig struct {
	Values []string `json:"values"`
}

//PathPatternConditionConfig -
type PathPatternConditionConfig struct {
	Values []string `json:"values"`
}

// QueryStringKeyValuePair -
type QueryStringKeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//QueryStringConditionConfig -
type QueryStringConditionConfig struct {
	Values []*QueryStringKeyValuePair `json:"values"`
}

//SourceIPConditionConfig -
type SourceIPConditionConfig struct {
	Values []string `json:"values"`
}

//RuleCondition -
type RuleCondition struct {
	Field                   string                            `json:"field"`
	HostHeaderConfig        *HostHeaderConditionConfig        `json:"host_header,omitempty"`
	HTTPHeaderConfig        *HTTPHeaderConditionConfig        `json:"http_header,omitempty"`
	HTTPRequestMethodConfig *HTTPRequestMethodConditionConfig `json:"http_request_method,omitempty"`
	PathPatternConfig       *PathPatternConditionConfig       `json:"path_pattern,omitempty"`
	QueryStringConfig       *QueryStringConditionConfig       `json:"query_string,omitempty"`
	SourceIPConfig          *SourceIPConditionConfig          `json:"source_ip,omitempty"`
	Values                  []string                          `json:"values"`
}

//Rule -
type Rule struct {
	Arn        string           `json:"arn"`
	Priority   string           `json:"priority"`
	IsDefault  bool             `json:"default"`
	Conditions []*RuleCondition `json:"conditions"`
	Actions    []*Action        `json:"actions"`
}
