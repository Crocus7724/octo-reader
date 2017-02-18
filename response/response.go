package response

const (
  // event type
  IssueEvent        = "IssueEvent"
  IssueCommentEvent = "IssueCommentEvent"
  CreateEvent       = "CreateEvent"
  PushEvent         = "PushEvent"
  ReleaseEvent      = "ReleaseEvent"
  GollumEvent       = "GollumEvent"
)

// /authorizations response
type Authorization struct {
  Token string `json:"token"`
}

// events response
type Event struct {
  Type      string `json:"type"`
  Actor     User `json:"actor"`
  Repo      Repo `json:"repo"`
  Payload   Payload `json:"payload"`
  CreatedAt string `json:"created_at"`
  Body      string `json:"body"`
}

type User struct {
  Name string `json:"login"`
  Url  string `json:"url"`
}

type Repo struct {
  Name string `json:"name"`
  Url  string `json:"url"`
}

type Payload struct {
  Action   string `json:"action"`
  Issue    Issue `json:"issue"`
  Pages    []Page `json:"pages"`
  Comments []Comment `json:"comments"`
  Commits  []Commit `json:"commits"`
  Release  Release `json:"release"`

  Description string `json:"description"`
}

type Issue struct {
  Url   string `json:"url"`
  Title string `json:"title"`
  User  User
}

type Comment struct {
  User User `json:"user"`
  Body string `json:"body"`
}

type Page struct {
  PageName string `json:"page_name"`
  Title    string `json:"title"`
  Action   string `json:"action"`
  HtmlUrl  string `json:"html_url"`
}

type Commit struct {
  CommitAuthor
  Author  CommitAuthor `json:"author"`
  Message string `json:"message"`
}

type CommitAuthor struct {
  Name string `json:"name"`
}

type Release struct {
  HtmlUrl string `json:"html_url"`
  TagName string `json:"tag_name"`
  Name    string `json:"name"`
  Author  User `json:"author"`
  Assets  []Asset `json:"assets"`
  Body    string `json:"body"`
}

type Asset struct {
  Name        string `json:"name"`
  Uploader    User `json:"uploader"`
  ContentType string `json:"content_type"`
  State       string `json:"state"`
  BrowserDownloadUrl string `json:"browser_download_url"`
}

type ErrorMessage struct {
  Message string `json:"message"`
  DocumentUrl string `json:"document_url"`
}

// notifications response
type Notification struct {
  Repository Repo `json:"repository"`
  Subject Subject `json:"subject"`
  Reason string `json:"reason"`
  UpdatedAt string `json:"updated_at"`
}

type Subject struct {
  Title string `json:"title"`
  Url string `json:"url"`
  Type string `json:"type"`
}