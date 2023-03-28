package dao

type MemosData struct {
	Content string `json:"content"`
}
type MoodReader struct {
	Highlights []MoodReaderHighlights `json:"highlights"`
}
type MoodReaderHighlights struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
	Note   string `json:"note"`
}
type Daily struct {
	Data       string `json:"data"`
	MdShowData string `json:"md_show_data"`
	Date       string `json:"date"`
	ServerTime string `json:"serverTime"`
}

type UrlStruct struct {
	Url string `json:"url"`
}

type PostJson struct {
	Content string `json:"content"`
}

// Couch Json struct

type FileDoc struct {
	ID       string   `json:"_id"`
	Rev      string   `json:"_rev,omitempty"`
	Children []string `json:"children"`
	Ctime    int64    `json:"ctime"`
	Mtime    int64    `json:"mtime"`
	Size     int      `json:"size"`
	Type     string   `json:"type"`
	Deleted  bool     `json:"deleted"`
}

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
