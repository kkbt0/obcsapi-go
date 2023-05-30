package tools

type RJson struct {
	Code    int            `json:"code"`
	Msg     string         `json:"msg"`
	Data    map[string]any `json:"data"`
	Success bool           `json:"success"`
}

type H map[string]any
