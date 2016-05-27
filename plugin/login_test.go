package plugin

import (
	"flag"
	"fmt"
	"testing"

	"github.com/codegangsta/cli"
)

var tests = []struct {
	username      string
	password      string
	server        string
	authToken     string
	authErr       error
	tokenExpected bool
}{
	{
		username:      "user",
		password:      "password",
		server:        "someServer",
		authToken:     "someToken",
		authErr:       nil,
		tokenExpected: true,
	},
	{
		username:      "user",
		password:      "",
		server:        "someServer",
		authToken:     "someToken",
		authErr:       nil,
		tokenExpected: false,
	},
	{
		username:      "",
		password:      "password",
		server:        "someServer",
		authToken:     "someToken",
		authErr:       nil,
		tokenExpected: false,
	},
	{
		username:      "user",
		password:      "password",
		server:        "",
		authToken:     "someToken",
		authErr:       nil,
		tokenExpected: true,
	},
	{
		username:      "user",
		password:      "password",
		server:        "",
		authToken:     "someToken",
		authErr:       fmt.Errorf("Invalid credentials provided"),
		tokenExpected: false,
	},
}

func TestLoginCmd(t *testing.T) {
	for _, test := range tests {
		deleteConfig()
		context := buildContext(test.username, test.password, test.server)
		authService := buildAuthService(test.authToken, test.authErr)

		doLogin(context, authService)

		cfg, err := readConfig()
		if err != nil {
			t.Error(err)
		}
		server := test.server
		if server == "" {
			server = defaultServerURL
		}
		token := cfg.Auths[server].AuthToken

		if test.tokenExpected && token != test.authToken {
			t.Error("Expected the token in the config file, but was not there")
		}

		if token == test.authToken && !test.tokenExpected {
			t.Error("Token was set in config file, but was not expected.")
		}
	}
}

func TestAuthTokenUpdated(t *testing.T) {
	deleteConfig()
	context := buildContext("user", "password", "someServer")
	authService := buildAuthService("token1", nil)

	doLogin(context, authService)

	cfg, err := readConfig()
	if err != nil {
		t.Error("Error reading config", err)
	}
	if cfg.Auths["someServer"].AuthToken != "token1" {
		t.Error("Token was not set in config file")
	}

	authService2 := buildAuthService("token2", nil)

	doLogin(context, authService2)

	cfg, err = readConfig()
	if err != nil {
		t.Error("Error reading config", err)
	}
	if cfg.Auths["someServer"].AuthToken != "token2" {
		t.Error("Token was not updated in config file")
	}

}

func buildContext(username, password, server string) *cli.Context {
	set := flag.NewFlagSet("test", 0)
	set.String("username", "", "")
	set.String("password", "", "")

	set.Parse([]string{"--username", username, "--password", password, server})

	context := cli.NewContext(nil, set, nil)
	return context
}

func buildAuthService(token string, err error) dummyAuthService {
	return dummyAuthService{token, err}
}
