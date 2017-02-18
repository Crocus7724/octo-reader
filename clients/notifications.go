package clients

import (
  "encoding/base64"
  "encoding/json"
  "errors"
  "io/ioutil"
  "net/http"
  "strconv"
  "time"

  "github.com/crocus7724/octo-reader/response"
)

type NotificationClient struct {
  Config      GithubConfig
  httpRequest *http.Request
  httpClient  *http.Client
}

func NewNotificationClient(httpClient *http.Client, config GithubConfig) (*NotificationClient, error) {
  client := NotificationClient{
    Config: config,
  }

  var err error
  client.httpRequest, err = createNewNotificationRequest(config)

  if err != nil {
    return nil, err
  }

  client.httpClient = httpClient

  return &client, nil
}

func (c *NotificationClient) LoopRequest(result chan []response.Notification, err chan error, cancel chan struct{}) {
  var lastModified = ""
  var interval int = 0

  go func() {
    for {
      notifications, httpErr := requestNotifications(c, &lastModified, &interval)

      if httpErr != nil {
        err <- httpErr
      }
      if notifications != nil {
        result <- *notifications
      }

      time.Sleep(time.Duration(interval) * time.Second)
    }
  }()

  <-cancel
}

func (c *NotificationClient) do() (*http.Response, error) {
  return c.httpClient.Do(c.httpRequest)
}

func requestNotifications(client *NotificationClient, lastModified *string, interval *int) (*[]response.Notification, error) {
  if *lastModified != "" {
    client.httpRequest.Header.Set("If-Modified-Since", *lastModified)
  }

  resp, err := client.do()

  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  // Not Modified
  // lastModifiedから変更なし
  if resp.StatusCode == 304 {
    return nil, nil
  }

  // レスポンスヘッダのLast-Modifiedで書き換え
  *lastModified = resp.Header.Get("Last-Modified")
  *interval, err = strconv.Atoi(resp.Header.Get("X-Poll-Interval"))

  if err != nil {
    return nil, err
  }

  byteArray, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    return nil, err
  }

  var notifications []response.Notification

  err = json.Unmarshal(byteArray, &notifications)

  if err != nil {
    var errorMessage response.ErrorMessage
    err := json.Unmarshal(byteArray, &errorMessage)

    if err != nil {
      return nil, err
    }
    return nil, errors.New(errorMessage.Message)
  }

  return &notifications, nil
}

func createNewNotificationRequest(config GithubConfig) (*http.Request, error) {
  host, err := config.GetApiHost()

  if err != nil {
    return nil, err
  }

  url := host + "notifications"
  req, err := http.NewRequest("GET", url, nil)

  if err != nil {
    return nil, err
  }

  if config.Token == "" {
    if config.Password == "" {
      return nil, errors.New("passwordかtokenを設定して下さい")
    }
    encoding := base64.StdEncoding.EncodeToString([]byte(config.Name + ":" + config.Password))
    req.Header.Set("Authorization", "Basic "+encoding)
  } else {
    req.Header.Set("Authorization", "token "+config.Token)
  }

  req.Header.Set("Accept", "application/json")

  return req, nil
}
