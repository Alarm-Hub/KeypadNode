package communicator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Phill93/DoorManager/config"
	"github.com/Phill93/DoorManager/log"
	jwt "github.com/dgrijalva/jwt-go"
	"io"
	"net/http"
	"os"
	"time"
)

// Handles the communication with the controller app

var cfg = config.Config()

var apiClient = &http.Client{
	Timeout: time.Second * cfg.GetDuration("http_timeout"),
}

type communicator struct {
	access      string
	refresh     string
	lastRefresh time.Time
	baseUrl     string
}

func NewCommunicator(access string, refresh string, baseUrl string) *communicator {
	return &communicator{access: access, refresh: refresh, baseUrl: baseUrl}
}

func (c communicator) VerifyAccess() (bool, error) {
	if c.baseUrl != "" && c.access != "" {
		p, err := json.Marshal(map[string]string{
			"token": c.access,
		})
		if err != nil {
			return false, err
		}
		res, err := apiClient.Post(fmt.Sprintf("%s/api/token/verify/", c.baseUrl), "application/json;charset=utf-8", bytes.NewBuffer(p))
		if err != nil {
			return false, err
		}
		if res.StatusCode == 200 {
			return true, nil
		}
		return false, err
	} else {
		return false, errors.New("Base Url or access token empty")
	}
}

func (c *communicator) Refresh() error {
	if c.baseUrl != "" && c.access != "" {
		p, err := json.Marshal(map[string]string{
			"refresh": c.refresh,
		})
		if err != nil {
			return err
		}
		res, err := apiClient.Post(fmt.Sprintf("%s/api/token/refresh/", c.baseUrl), "application/json;charset=utf-8", bytes.NewBuffer(p))
		if err != nil {
			return err
		}
		type response struct {
			access  string
			refresh string
		}
		var r response
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return err
		}
		c.access = r.access
		c.refresh = r.refresh
		c.lastRefresh = time.Now()
		_ = c.VerifyAccess()
		return nil
	} else {
		return errors.New("Base Url or access token empty")
	}
}

func (c communicator) ValidateCode(code string) (bool, error) {
	if c.baseUrl != "" && c.access != "" {
		p, err := json.Marshal(map[string]string{
			"code": code,
		})
		if err != nil {
			return false, err
		}
		res, err := apiClient.Post(fmt.Sprintf("%s/api/codes/verify/", c.baseUrl), "application/json;charset=utf-8", bytes.NewBuffer(p))
		if err != nil {
			return false, err
		}
		type response struct {
			valid bool
		}
		var r response
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return false, err
		}
		if r.valid == true {
			return true, nil
		}
		return false, nil
	} else {
		return false, errors.New("Base Url or access token empty")
	}
}

func parseToken(tokenString string) interface{} {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Fatal(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Fatal("Can't convert claims")
	}
	return claims
}

func (c *communicator) TokenWatcher() {
	for {
		ok, err := c.VerifyAccess()
		if err == nil {
			log.Fatal(err)
		}
		if ok {
			claims := parseToken(c.access)
			var tm time.Time
			switch exp := claims["exp"].(type) {
			case float64:
				tm = time.Unix(int64(exp), 0)
			case json.Number:
				v, _ := exp.Int64()
				tm = time.Unix(v, 0)
			}
			diff := time.Now().Sub(tm).Seconds()
			if diff < cfg.GetFloat64("token_refresh_threshold") {
				err = c.Refresh()
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal("Token invalid!")
		}
		time.Sleep(time.Second * 10)
	}
}
