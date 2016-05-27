package plugin

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// AuthService for authenticating against a backend
type AuthService interface {
	// Login user to the given URL with provided username and password
	Login(url, user, password string) (string, error)
	// Logout user from the given URL
	//Logout(url string) error
}

const defaultServerURL = "someDefaultServer"

// Login command
func Login(c *cli.Context) {
	// TODO: Implement auth service
	authService := dummyAuthService{"token", nil}
	doLogin(c, authService)
}

func doLogin(c *cli.Context, authService AuthService) {

	user := c.String("username")
	password := c.String("password")

	if user == "" || password == "" {
		fmt.Println("Username and password must be provided.")
		return
	}

	serverURL := c.Args().First()
	if serverURL == "" {
		serverURL = defaultServerURL
	}

	// Authenticate against backend
	authToken, err := authService.Login(serverURL, user, password)
	if err != nil {
		fmt.Println("Error authenticating with server.", err)
		return
	}

	// Store token in config
	cfg, err := readConfig()
	if err != nil {
		fmt.Println("Error reading config file.", err)
		return
	}

	cfg.Auths[serverURL] = BackendAuth{authToken}

	err = writeConfig(cfg)
	if err != nil {
		fmt.Println("Error saving config file.", err)
		return
	}

	fmt.Println("Login successful. Authentication token saved.")
}

type dummyAuthService struct {
	token string
	err   error
}

func (das dummyAuthService) Login(url, user, password string) (string, error) {
	if das.err != nil {
		return "", das.err
	}
	return das.token, nil
}
