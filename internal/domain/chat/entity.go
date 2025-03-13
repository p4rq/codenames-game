package chat

type Message struct {
    ID      string `json:"id"`
    UserID  string `json:"user_id"`
    Content string `json:"content"`
    Time    string `json:"time"`
}

type Chat struct {
    ID      string    `json:"id"`
    Messages []Message `json:"messages"`
}