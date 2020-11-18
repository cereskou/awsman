package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ditto.co.jp/awsman/logger"
	"ditto.co.jp/awsman/utils"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/kardianos/osext"
)

//
var (
	configfile = "awsman.json"
	setting    *Config
)

//AwsConfig -
type AwsConfig struct {
	AccessKey string `json:"access_key"`        //AWSアクセスキー
	SecretKey string `json:"secret_access_key"` //シークレットキー
	Region    string `json:"region"`            //リージョン
	Retry     int    `json:"retry"`             //リトライ回数
}

//AccountConfig -
type AccountConfig struct {
	Name     string `json:"name"`     // username
	Password string `json:"password"` //password
	Avatar   string `json:"avatar"`   //avatar png file
	Expires  int64  `json:"expires"`  //単位：時間
}

//RsaConfig -
type RsaConfig struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

//DbConfig for sqlite3
type DbConfig struct {
	Name       string `json:"Name"`
	AutoVacuum bool   `json:"AutoVacuum"`
}

//CognitoConfig -
type CognitoConfig struct {
	UserPoolID string `json:"userpoolid"`
}

//Config -
type Config struct {
	Version   string        `json:"-"`         //Version
	ExeName   string        `json:"-"`         //モジュールパス
	FileName  string        `json:"-"`         //s3service.jsonパス
	Dir       string        `json:"-"`         //work directory
	LogFile   string        `json:"-"`         //logfile fullpath
	Host      string        `json:"Host"`      //サーバーIP
	Port      int           `json:"Port"`      //ポート
	LogDir    string        `json:"LogDir"`    //ログファイルの出力先
	LogOutput string        `json:"LogOutput"` //出力先[screen, file]
	LogFormat string        `json:"LogFormat"` //ログファイル名フォーマット
	Proxy     string        `json:"proxy"`     //proxy server for aws client
	Avatar    string        `json:"Avatar"`    //default avatar image
	Account   AccountConfig `json:"account"`   //Default root account
	Rsa       RsaConfig     `json:"rsa"`       //Rsa key
	Aws       AwsConfig     `json:"aws"`       //aws config
	Db        DbConfig      `json:"db"`        //DB設定
	Cognito   CognitoConfig `json:"cognito"`   //Cognito
}

//Load -
func Load() (*Config, error) {
	if setting != nil {
		return setting, nil
	}

	//initialize
	setting = &Config{
		Host: "localhost",
		Port: 9898,
	}

	exename, _ := osext.Executable()
	setting.ExeName = exename

	dir := filepath.Dir(exename)
	name := filepath.Join(dir, configfile)
	setting.FileName = name
	setting.Dir = dir
	fmt.Println(name)

	fr, err := os.Open(name)
	if err != nil {
		return setting, err
	}
	defer fr.Close()

	//s3service -json
	err = utils.JSON.NewDecoder(fr).Decode(&setting)
	if err != nil {
		return setting, err
	}
	if setting.Dir == "" {
		setting.Dir = dir
	}

	//initializer
	setting.Init()

	//token expires time in hour
	if setting.Account.Expires <= 0 {
		setting.Account.Expires = 12
	}

	//credential
	if setting.Aws.AccessKey == "" || setting.Aws.SecretKey == "" {
		dir, _ := dirWindows()
		dir = filepath.ToSlash(dir)

		credsfile := fmt.Sprintf("%s/.aws/credentials", dir)
		creds := credentials.NewSharedCredentials(credsfile, "default")
		credValue, err := creds.Get()
		if err != nil {
			//環境変数
			setting.Aws.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
			setting.Aws.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
			setting.Aws.Region = os.Getenv("AWS_DEFAULT_REGION")
		} else {
			setting.Aws.AccessKey = credValue.AccessKeyID
			setting.Aws.SecretKey = credValue.SecretAccessKey
		}
	}
	if setting.Aws.Region == "" {
		setting.Aws.Region = "ap-northeast-1"
	}

	return setting, nil
}

//Init -
func (c *Config) Init() {
	if c.LogOutput == "" {
		c.LogOutput = "file"
	}
	var logfile string
	if c.LogOutput == "file" || c.LogOutput == "both" {
		err := utils.MakeDirectory(c.LogDir)
		if err != nil {
			c.LogDir = ""
		}

		logfile = getLogFilename(c.LogFormat)
		logfile = filepath.Join(c.LogDir, logfile)

		c.LogFile = logfile
	}
	logger.SetOutput(c.LogOutput, logfile)
}

//getLogFilename -
func getLogFilename(format string) string {
	logfile := "awsman.log" //20060102150405

	pos := strings.Index(format, "%")
	if pos > -1 {
		lpos := strings.LastIndex(format, "%")
		if lpos > -1 {
			if lpos > pos {
				fstr := format[pos+1 : lpos]
				if len(fstr) > 0 {
					ffmt := format[:pos] + "%v" + format[lpos+1:]
					//20060102150405 "%yyyymmddHHMMSS%"
					fstr = strings.ReplaceAll(fstr, "yyyy", "2006")
					fstr = strings.ReplaceAll(fstr, "mm", "01")
					fstr = strings.ReplaceAll(fstr, "dd", "02")
					fstr = strings.ReplaceAll(fstr, "HH", "15")
					fstr = strings.ReplaceAll(fstr, "MM", "04")
					fstr = strings.ReplaceAll(fstr, "SS", "05")

					logfile = fmt.Sprintf(ffmt, utils.NowJST().Format(fstr)) //20060102150405
				}
			}
		}
	} else {
		if format != "" {
			logfile = format
		}
	}

	return logfile
}

func dirWindows() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// Prefer standard environment variable USERPROFILE
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home, nil
	}

	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, or USERPROFILE are blank")
	}

	return home, nil
}
