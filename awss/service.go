package awss

import (
	"net/http"
	"net/url"

	"ditto.co.jp/awsman/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/codecommit"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

//Service -
type Service struct {
	_session    *session.Session
	_cip        *cognitoidentityprovider.CognitoIdentityProvider
	_cloudtrail *cloudtrail.CloudTrail
	_codecommit *codecommit.CodeCommit
	_elbv2      *elbv2.ELBV2
}

//New -
func New() *Service {
	svc := &Service{}
	sess := svc.getSession()
	svc._session = sess

	return svc
}

//Close -
func (s *Service) Close() {
	if s._cip != nil {
	}
}

func (s *Service) getSession() *session.Session {
	cfg, _ := config.Load()
	//Proxy
	var httpClient *http.Client
	if len(cfg.Proxy) > 0 {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return url.Parse(cfg.Proxy)
				},
			},
		}
	}

	//認証情報を作成します。
	cred := credentials.NewStaticCredentials(
		cfg.Aws.AccessKey,
		cfg.Aws.SecretKey,
		"")

	//セッション作成します
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Aws.Region),
		Credentials: cred,
		HTTPClient:  httpClient,
		MaxRetries:  aws.Int(cfg.Aws.Retry),
	}))

	return sess
}

//Cognito -
func (s *Service) Cognito() *cognitoidentityprovider.CognitoIdentityProvider {
	if s._cip == nil {
		s._cip = cognitoidentityprovider.New(s._session)
	}

	return s._cip
}

//CloudTrail -
func (s *Service) CloudTrail() *cloudtrail.CloudTrail {
	if s._cloudtrail == nil {
		s._cloudtrail = cloudtrail.New(s._session)
	}

	return s._cloudtrail
}

//CodeCommit -
func (s *Service) CodeCommit() *codecommit.CodeCommit {
	if s._codecommit == nil {
		s._codecommit = codecommit.New(s._session)
	}

	return s._codecommit
}

//ELBV2 -
func (s *Service) ELBV2() *elbv2.ELBV2 {
	if s._elbv2 == nil {
		s._elbv2 = elbv2.New(s._session)
	}
	return s._elbv2
}
