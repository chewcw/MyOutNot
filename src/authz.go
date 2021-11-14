package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthzService struct {
	r            *gin.Engine
	accessToken  string
	refreshToken string
	userOid      string
	LoginURL     string
	TokenURL     string
}

func NewAuthzService() *AuthzService {
	authzService := AuthzService{
		r:        gin.Default(),
		LoginURL: "https://login.microsoft.com/" + tenantID + "/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoft.com/" + tenantID + "/oauth2/v2.0/token",
	}

	// redirect to MSFT login
	authzService.r.GET("/login", func(c *gin.Context) {
		u := url.URL{
			Path: authzService.LoginURL,
		}
		q, _ := url.ParseQuery(u.RawQuery)
		q.Add("client_id", clientID)
		q.Add("response_type", "code")
		q.Add("redirect_uri", redirectURL)
		q.Add("scope", "user.read offline_access")
		q.Add("response_mode", "query")

		u.RawQuery = q.Encode()
		c.Redirect(http.StatusFound, u.RequestURI())
	})

	// callback url after successfully authorized
	authzService.r.GET("/evelyn", func(c *gin.Context) {
		authzCode, _ := c.GetQuery("code")
		data := url.Values{
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"scope":         {"user.read offline_access"},
			"code":          {authzCode},
			"redirect_uri":  {redirectURL},
			"grant_type":    {"authorization_code"},
		}
		resp, err := http.PostForm(authzService.TokenURL, data)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var res map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Fatal(err)
		}

		authzService.accessToken = res["access_token"].(string)
		authzService.refreshToken = res["refresh_token"].(string)

		// parse the jwt
		oid, err := parseJwt(authzService.accessToken, "oid")
		if err != nil {
			log.Fatal(err)
		}
		authzService.userOid = oid
	})

	return &authzService
}

// Refresh the access token
func (a *AuthzService) Refresh() {
	if a.refreshToken == "" {
		log.Println("No refresh token, not going to refresh access token for now")
		return
	}

	log.Println("Refreshing access token")

	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"refresh_token": {a.refreshToken},
		"scope":         {"user.read offline_access"},
		"grant_type":    {"refresh_token"},
	}
	resp, err := http.PostForm(a.TokenURL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}

	a.accessToken = res["access_token"].(string)
	a.refreshToken = res["refresh_token"].(string)
}

func parseJwt(accessToken string, claim string) (string, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		log.Println(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims[claim].(string), nil
	} else {
		errorMsg := fmt.Sprintf("`%s` claim not found in the jwt token", claim)
		return "", errors.New(errorMsg)
	}
}
