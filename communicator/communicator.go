package communicator

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/Phill93/DoorManager/config"
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

func (c communicator) VerifyAccess() error {
  if c.baseUrl != "" && c.access != "" {
    p, err := json.Marshal(map[string]string{
      "token": c.access,
    })
    if err != nil {
      return err
    }
    res, err := apiClient.Post(fmt.Sprintf("%s/api/token/verify/", c.baseUrl), "application/json;charset=utf-8", bytes.NewBuffer(p))
    if err != nil {
      return err
    }
    _, err = io.Copy(os.Stdout, res.Body)
    if err != nil {
      return err
    }
    return nil
  } else {
    return errors.New("Base Url or access token empty")
  }
}
