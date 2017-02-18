package clients

// 間違ってつくちゃったやつ

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

type EventClient struct {
  Config      GithubConfig
  httpRequest *http.Request
  httpClient  *http.Client
}

func NewEventClient(httpClient *http.Client, config GithubConfig) (*EventClient, error) {
  var err error
  client := EventClient{
    Config: config,
  }

  client.httpRequest, err = createNewEventsRequest(config)

  if err != nil {
    return nil, err
  }

  client.httpClient = httpClient

  return &client, nil
}

func (c *EventClient) LoopRequest(result chan []response.Event, err chan error, cancel chan struct{}) {
  var lastModified string = ""
  var interval int = 0

  go func() {
    for {
      event, httpErr := requestEvents(c, &lastModified, &interval)

      if httpErr != nil {
        err <- httpErr
      }
      if event != nil {
        result <- *event
      }

      time.Sleep(time.Duration(interval) * time.Second)
    }
  }()

  <-cancel
}

func (c *EventClient) do() (*http.Response, error) {
  return c.httpClient.Do(c.httpRequest)
}

// notificationsにリクエストぶん投げる
func requestEvents(client *EventClient, lastModified *string, interval *int) (*[]response.Event, error) {
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

  var events []response.Event

  err = json.Unmarshal(byteArray, &events)

  if err != nil {
    var errorMessage response.ErrorMessage
    err := json.Unmarshal(byteArray, &errorMessage)

    if err != nil {
      return nil, err
    }
    return nil, errors.New(errorMessage.Message)
  }

  return &events, nil
}

// users/{username}/events用のhttp.Requestを返す
func createNewEventsRequest(config GithubConfig) (*http.Request, error) {
  host, err := config.GetApiHost()

  if err != nil {
    return nil, err
  }

  url := host + "users/" + config.Name + "/events"
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
