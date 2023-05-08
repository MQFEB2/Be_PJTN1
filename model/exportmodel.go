package model

type EpInfo struct {
	ID     string `gorm:"type:varchar(255)" json:"_postman_id"`
	Name   string `json:"name"`
	Schema string `json:"schema" gorm:"default:https://schema.getpostman.com/json/collection/v2.1.0/collection.json"`
}

type EpCollection struct {
	EpInfo  EpInfo        `gorm:"type:varchar(255)" json:"info"`
	EpItems []interface{} `json:"item"`
}

type EpItemFolder struct {
	Name   string          `json:"name"`
	Folder []EpItemRequest `json:"item"`
}

type EpItemRequest struct {
	Name     string       `json:"name"`
	Protocol interface{}  `json:"protocolProfileBehavior"`
	Request  EpRequest    `json:"request"`
	Response []EpResponse `json:"response"`
}

type EpRequest struct {
	Method    string     `json:"method"`
	EpHeaders []EpHeader `json:"header"`
	EpBody    *EpBody    `json:"body,omitempty"`
	EpUrl     *EpUrl     `json:"url,omitempty"`
}

type EpUrl struct {
	Raw       string        `json:"raw,omitempty"`
	Protocol  string        `json:"protocol,omitempty"`
	Host      []interface{} `json:"host,omitempty"`
	Port      string        `json:"port,omitempty"`
	Path      []interface{} `json:"path,omitempty"`
	EpQueries []EpQuery     `json:"query,omitempty"`
}

type EpHeader struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Disabled    bool   `gorm:"default:false" json:"disable,omitempty"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type EpQuery struct {
	Key         string `gorm:"null" json:"key"`
	Value       string `gorm:"null" json:"value"`
	Disabled    bool   `default:"false" json:"disable,omitempty"`
	Description string `json:"description,omitempty"`
}

type EpBody struct {
	Mode     string      `json:"mode,omitempty"`
	Raw      string      `json:"raw,omitempty"`
	Options  interface{} `json:"options,omitempty"`
	Disabled bool        `default:"false" json:"disable,omitempty"`
}

type EpResponse struct {
	Name     string    `json:"name"`
	Request  EpRequest `json:"originalRequest"`
	Status   string    `json:"status"`
	Code     string    `json:"code"`
	Language string    `json:"_postman_previewlanguage"`
	Header   string    `json:"header"`
	Body     string    `json:"body"`
}
