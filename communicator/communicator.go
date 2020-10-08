package communicator

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/Phill93/DoorManager/config"
  "github.com/Phill93/DoorManager/log"
  jwt "github.com/dgrijalva/jwt-go"
  "github.com/spf13/viper"
  "io/ioutil"
  "net/http"
  "os"
  "time"
)

// Handles the communication with the controller app

var cfg = config.Config()

var apiClient = &http.Client{
	Timeout: time.Second * cfg.GetDuration("http_timeout"),
}

type Communicator struct {
	access      string
	refresh     string
	lastRefresh time.Time
	baseUrl     string
}

type tokens struct {
  Access      string  `json:"access"`
  Refresh     string  `json:"refresh"`
}

func NewCommunicator(baseUrl string) Communicator {
  var com Communicator
  com.baseUrl = baseUrl
  t := com.loadTokens()
  com.access = t.Access
  com.refresh = t.Refresh
  ok, err := com.VerifyAccess()
  if err != nil {
    log.Panic(err)
  }
  if !ok {
    log.Panic("Tokens are not valid!")
  }
	return com
}

func (c Communicator) loadTokens() tokens {
  jsonFile, err := os.Open("tokens.json")
  if err != nil {
    log.Panic(err)
  }
  defer jsonFile.Close()
  jsonBytes, err := ioutil.ReadAll(jsonFile)
  if err != nil {
    log.Panic(err)
  }
  var tokens tokens
  err = json.Unmarshal(jsonBytes, &tokens)
  if err != nil {
    log.Panic(err)
  }
  return tokens
}

func (c Communicator) saveTokens(tokens tokens) {
  file, err := json.Marshal(&tokens)
  if err != nil {
    log.Panic(err)
  }
  err = ioutil.WriteFile("tokens.json",file,0600)
  if err != nil {
    log.Panic(err)
  }
}

func (c Communicator) VerifyAccess() (bool, error) {
	if c.baseUrl != "" && c.access != "" {
	  log.Debug("Try to verify token")
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
		  log.Debug("Token is valid")
			return true, nil
		}
		log.Debug("Token is invalid")
		return false, err
	} else {
		return false, errors.New("Base Url or access token empty")
	}
}

func (c *Communicator) Refresh() error {
	if c.baseUrl != "" && c.access != "" {
	  log.Debug("Try to refresh tokens")
		p, err := json.Marshal(map[string]string{
			"refresh": c.refresh,
		})
		if err != nil {
			return err
		}
		res, err := apiClient.Post(fmt.Sprintf("%s/api/token/refresh/", c.baseUrl), "application/json;charset=utf-8", bytes.NewBuffer(p))
		defer func() {
		  if res != nil {
		    _ = res.Body.Close()
      }
    }()
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
		  log.Error("Failed to refresh token")
		  return nil
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      log.Error(err)
    }
    var r tokens
		err = json.Unmarshal(body, &r)
		if err != nil {
			return err
		}
		c.access = r.Access
		c.refresh = r.Refresh
		c.saveTokens(r)
		if cfg.ConfigFileUsed() != "" {
		  _ = viper.WriteConfig()
    }
		c.lastRefresh = time.Now()
		_, _ = c.VerifyAccess()
		return nil
	} else {
		return errors.New("Base Url or access token empty")
	}
}

func (c Communicator) ValidateCode(code string) (bool, error) {
	if c.baseUrl != "" && c.access != "" {
		p, err := json.Marshal(map[string]string{
			"code": code,
		})
		if err != nil {
			return false, err
		}
		req, err := http.NewRequest("post", fmt.Sprintf("%s/api/codes/verify/", c.baseUrl), bytes.NewBuffer(p))
		if err != nil {
		  return false, err
    }
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.access))
		res, err := apiClient.Do(req)
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

func (c *Communicator) TokenWatcher() {
	for {
		ok, err := c.VerifyAccess()
		if err == nil {
			log.Fatal(err)
		}
		if ok {
			claims := parseToken(c.access)
			var tm time.Time
			fmt.Print(claims)
			//switch exp := claims["exp"].(type) {
      //  case float64:
      //    tm = time.Unix(int64(exp), 0)
      //  case json.Number:
      //    v, _ := exp.Int64()
      //    tm = time.Unix(v, 0)
			//}
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
