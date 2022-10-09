package top

type ResponseList struct {
	List       []interface{} `json:"list"`
	Pagination Pagination    `json:"pagination"`
}

type Pagination struct {
	Page      int                    `json:"page"`
	Rows      int                    `json:"rows"`
	Size      int                    `json:"size"`
	Limit     int                    `json:"limit"`
	Pages     int                    `json:"pages"`
	URLLayout string                 `json:"urlLayout"`
	Data      map[string]interface{} `json:"data"`
}
