package spaceshipsuper

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/michaelmcallister/spaceshipsuper/http"
)

var authVersion = map[string]string{
	"version": "5",
	"status":  "-1",
}

type authPayload struct {
	SignupID  string `json:"signup_id"`
	Source    string
	Decoupled bool
	Oauth2    struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string
		Scope        string
	}
}

type authRequest struct {
	username    string
	password    string
	fingerprint string
}

// basicAuth concatenates username and password and returns a base64 encoded
// string.
// Initial login and authentication is gated by basicAuth.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// refreshAuth checks the validity of the current token, by default it will not
// refresh if the current token is valid. Set force to true to override this
// behaviour.
func (c *Client) refreshAuth(force bool) error {
	// If our token hasn't expired and we aren't forcing a new one, just re-use.
	if !force && time.Now().Before(c.expiry) {
		return nil
	}

	return c.auth()
}

// auth handles the initial authentication dance, which involves sending the
// username and password (and setting the basicauth header).
func (c *Client) auth() error {
	const authEndpoint = apiEndpointBase + "/v1/shim/login"
	fprint, _ := json.Marshal(authVersion)

	authRequest, _ := json.Marshal(authRequest{
		username:    c.Username,
		password:    c.Password,
		fingerprint: string(fprint),
	})

	resp, err := http.DoPost(authEndpoint, basicAuth(c.Username, c.Password), authRequest)
	if err != nil {
		return err
	}

	authp := &authPayload{}
	if err := json.Unmarshal(resp, authp); err != nil {
		return err
	}

	c.accessToken = authp.Oauth2.AccessToken
	c.expiry = time.Now().Add(time.Second * time.Duration(authp.Oauth2.ExpiresIn))
	return nil
}
