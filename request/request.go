package request

type Authorization struct {
  Scopes []string `json:"scopes"`
  Note   string `json:"note"`
}
