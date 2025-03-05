package dao

type MemosData struct {
	Content string `json:"content"`
}

type UrlStruct struct {
	Url string `json:"url"`
}

type PostJson struct {
	Content string `json:"content"`
}

// CouchDB Json struct

type FileDoc struct {
	ID       string   `json:"_id"`
	Rev      string   `json:"_rev,omitempty"`
	Path     string   `json:"path"`
	Children []string `json:"children"`
	Ctime    int64    `json:"ctime"`
	Mtime    int64    `json:"mtime"`
	Size     int      `json:"size"`
	Type     string   `json:"type"`
	Deleted  bool     `json:"deleted"`
	Eden     struct{} `json:"eden"`
}

// CouchDB Json struct
type LeftDoc struct {
	ID   string `json:"_id"`
	Rev  string `json:"_rev,omitempty"`
	Data string `json:"data"`
	Type string `json:"type"`
}

type ErrorJson struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}
