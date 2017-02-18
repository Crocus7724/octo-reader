package main

import (
  "encoding/json"
  "fmt"
  "os"
  "github.com/crocus7724/octo-reader/clients"
)

// 単なるPrintlnのラッパー
func Write(a interface{}) {
  write(a)
}

func WriteDebug(a interface{}) {
  if clients.Debug {
    Write(a)
  }
}

// 単なるFprint(error)のラッパー
func WriteError(a interface{}) {
  fmt.Fprint(os.Stderr, a)
}

// エラーがないかチェックする
// エラーが有った場合、そこでプログラム終了
func CheckError(err error) {
  if err != nil {
    WriteError(err)
    os.Exit(-1)
  }
}

func write(a interface{}) {
  switch v := a.(type) {
  case string:
    fmt.Println(v)
  case int:
    fmt.Println(v)
  case byte:
    fmt.Println(v)
  case bool:
    fmt.Println(v)
  default:
    j, err := json.MarshalIndent(v, "", "  ")

    if err != nil {
      WriteError(err)
    }
    fmt.Println(string(j))
  }
}
