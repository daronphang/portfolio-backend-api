package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"portfolio_golang/src/zaplog"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/yaml.v2"
)

type EnvConf struct {
	MSSQLUSERNAME 			string `yaml:"MSSQLUSERNAME"`
	MSSQLPASSWORD 			string `yaml:"MSSQLPASSWORD"`
	MSSQLSERVER   			string `yaml:"MSSQLSERVER"`
	MSSQLPORT     			int    `yaml:"MSSQLPORT"`
	MYSQLUSERNAME 			string `yaml:"MYSQLUSERNAME"`
	MYSQLPASSWORD 			string `yaml:"MYSQLPASSWORD"`
	MYSQLSERVER   			string `yaml:"MYSQLSERVER"`
	MYSQLPORT     			int    `yaml:"MYSQLPORT"`
	MYSQLPROTOCOL			string `yaml:"MYSQLPROTOCOL"`
}

var AppCfg = &EnvConf{}
var GmailClient = &http.Client{}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		zaplog.Logger.Fatal(fmt.Sprintf("unable to get cwd: %s", err.Error()))
	}

	// read appCfg env
	AppCfg.Read(cwd)

	// generate Google OAuth HTTP client for Gmail
	getGmailClient(cwd)
}

func (c *EnvConf) Read(cwd string) {
	filepath := path.Join(cwd, "config", "env.yaml")
	envYaml, err := os.ReadFile(filepath)
	if err != nil {
		zaplog.Logger.Fatal(err.Error())
	}
	err = yaml.Unmarshal(envYaml, c)
	if err != nil {
		zaplog.Logger.Fatal(err.Error())
	}
	zaplog.Logger.Info("succecssfully read environment variables")
}

// returns generated client for gmail
// auto refreshes as necessary
func getGmailClient(cwd string) {
	// get credentials from file
	filepath := path.Join(cwd, "config", "gmail.credentials.json")
	b, err := os.ReadFile(filepath)
	if err != nil {
		zaplog.Logger.Fatal(err.Error())
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to parse gmail credentials to config: %v", err)
		zaplog.Logger.Fatal(errMsg)
	}

	// get access/refresh tokens from file
	filepath = path.Join(cwd, "config", "gmail.token.json")
	tok, err := tokenFromFile(filepath)
	if err != nil {
		zaplog.Logger.Fatal(err.Error())
	}

	GmailClient = config.Client(context.Background(), tok)
	zaplog.Logger.Info("succecssfully established gmail client")
}

// retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
			return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}