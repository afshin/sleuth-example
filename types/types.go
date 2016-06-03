package types

// Article holds the metadata for an article.
type Article struct {
	GUID      string     `json:"guid"`
	Byline    string     `json:"byline"`
	Comments  []*Comment `json:"comments,omitempty"`
	Headline  string     `json:"headline"`
	URL       string     `json:"url"`
	Timestamp int64      `json:"time"`
}

// ArticleResponse is the format of article-service responses.
type ArticleResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Data    *Article `json:"data,omitempty"`
}

// Comment holds the text and metadata for an article comment.
type Comment struct {
	GUID      string `json:"guid"`
	Article   string `json:"article"`
	Text      string `json:"text"`
	Timestamp int64  `json:"time"`
}

// CommentResponse is the format of comment-service responses.
type CommentResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    []*Comment `json:"data,omitempty"`
}
