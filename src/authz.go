package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

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
	firstTime    bool
}

func NewAuthzService() *AuthzService {
	authzService := AuthzService{
		r:         gin.Default(),
		LoginURL:  "https://login.microsoft.com/" + tenantID + "/oauth2/v2.0/authorize",
		TokenURL:  "https://login.microsoft.com/" + tenantID + "/oauth2/v2.0/token",
		firstTime: true,
	}

	authzService.r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hello")
	})

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

		if resp.StatusCode != 200 {
			b, err := json.Marshal(res)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(b))
			return
		}

		authzService.accessToken = res["access_token"].(string)
		authzService.refreshToken = res["refresh_token"].(string)

		// parse the jwt
		oid, err := parseJwt(authzService.accessToken, "oid")
		if err != nil {
			log.Fatal(err)
		}
		authzService.userOid = oid

		// save access token and refresh token locally
		storeLocally(authzService.accessToken, authzService.refreshToken)
	})

	return &authzService
}

// Refresh the access token
func (a *AuthzService) Refresh() {
	if a.refreshToken == "" {
		log.Println("getting refresh token locally")
		var err error
		a.accessToken, a.refreshToken, err = readLocally()
		if err != nil {
			log.Println(err)
			return
		}
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

	log.Printf("Got access_token: %s\n", a.accessToken)
	log.Printf("Got refresh_token: %s\n", a.refreshToken)

	a.firstTime = false

	// parse the jwt
	oid, err := parseJwt(a.accessToken, "oid")
	if err != nil {
		log.Fatal(err)
	}
	a.userOid = oid

	// save refresh token to azure table
	storeLocally(a.accessToken, a.refreshToken)
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

func storeLocally(accessToken, refreshToken string) {
	at := []byte(accessToken)
	f1, err := os.Create(path.Join(localFileLocation, "access_token"))
	if err != nil {
		log.Println(err)
		return
	}
	defer f1.Close()
	if _, err = f1.Write(at); err != nil {
		log.Println(err)
		return
	}

	rt := []byte(refreshToken)
	f2, err := os.Create(path.Join(localFileLocation, "refresh_token"))
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = f2.Write(rt); err != nil {
		log.Println(err)
		return
	}
}

func readLocally() (accessToken, refreshToken string, err error) {
	at, err := os.ReadFile(path.Join(localFileLocation, "access_token"))
	if err != nil {
		return "", "", err
	}

	rt, err := os.ReadFile(path.Join(localFileLocation, "refresh_token"))
	if err != nil {
		return "", "", err
	}

	return string(at), string(rt), nil
}
