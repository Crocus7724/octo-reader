package main

import (
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  "github.com/crocus7724/octo-reader/clients"
  "github.com/crocus7724/octo-reader/response"
)

func main() {
  Write("Listening Github Notification.")
  Write("Exit Ctrl+C...\n")
  config := clients.CreateConfig()
  httpClient := new(http.Client)
  client, err := clients.NewNotificationClient(httpClient, config)

  CheckError(err)

  resultChannel := make(chan []response.Notification, 1)
  errChannel := make(chan error, 1)
  finishChannel := make(chan struct{}, 1)

  go client.LoopRequest(resultChannel, errChannel, finishChannel)

  go func() {
    for {
      select {
      case notification := <-resultChannel:
        writeNotifications(&notification)
      case err := <-errChannel:
        WriteError(err)
      }
    }
  }()

  sig := make(chan os.Signal, 1)

  // SIGINT(ctrl+c押されたときとか)にチャンネルに通知
  signal.Notify(sig, syscall.SIGINT)

  // 終了通知待ち
  <-sig
  close(finishChannel)
  fmt.Println("\nEnd Github Notification.")
}

func writeNotifications(notifications *[]response.Notification) {
  if notifications == nil {
    return
  }

  for _, notification := range *notifications {
    Write("Reason: " + notification.Reason + " " + notification.UpdatedAt)
    Write(notification.Repository.Name + "(" + notification.Repository.Url + ")")
    Write(notification.Subject.Type + ": " + notification.Subject.Title + "(" + notification.Subject.Url + ")")
  }
}

// Eventsレスポンスを書き出す
func writeEvents(events *[]response.Event) {
  if events == nil {
    return
  }

  for _, event := range *events {
    Write("Type: " + event.Type)

    Write("Author: " + event.Actor.Name + "    Date: " + event.CreatedAt)
    Write("Repository: " + event.Repo.Name + "(" + event.Repo.Url + ")")

    payload := event.Payload
    switch event.Type {
    case response.IssueCommentEvent:
      fallthrough
    case response.IssueEvent:
      issue := payload.Issue
      Write("Issue Title: " + issue.Title + "(" + issue.Url + ")")
      Write("Action: " + event.Payload.Action)
    case response.CreateEvent:
      Write("Description: " + payload.Description)
    case response.PushEvent:
      commits := payload.Commits

      for _, commit := range commits {
        Write("  " + commit.Name + ": " + commit.Message)
      }
    case response.GollumEvent:
      pages := payload.Pages

      for _, page := range pages {
        Write("Page Action: " + page.Action)
        Write("Page Name: " + page.PageName + "(" + page.HtmlUrl + ")")
      }
    case response.ReleaseEvent:
      Write("Release Action: " + payload.Action)

      release := payload.Release
      Write(release.Name + "(" + release.HtmlUrl + ") " + release.TagName)
      Write(release.Body)
      for _, asset := range release.Assets {
        Write("[Asset] " + asset.Name + "(" + asset.BrowserDownloadUrl + ")")
      }
    }

    fmt.Println()
  }
}
