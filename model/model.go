package model

type Info struct {
	ID   string `gorm:"type:varchar(255)"`
	Name string
}

type Collection struct {
	ID     string `gorm:"type:varchar(255)"`
	InfoID string `gorm:"type:varchar(255)"`
	Items  []Item `json:"item"`
}

type Item struct {
	ID           string `gorm:"type:varchar(255)"`
	Name         string
	CollectionID string    `gorm:"type:varchar(255)" json:"-"`
	IsFolder     bool      `gorm:"default:false"`
	Folders      []Folder  `gorm:"constraint:OnDelete:CASCADE;" json:"item"`
	Examples     []Example `gorm:"constraint:OnDelete:CASCADE;" json:"response"`
}

type Folder struct {
	ID       string `gorm:"type:varchar(255)"`
	Name     string
	ItemID   string    `gorm:"type:varchar(255)" json:"-"`
	Examples []Example `gorm:"constraint:OnDelete:CASCADE;" json:"response"`
}

type Request struct {
	ID       string `gorm:"type:varchar(255)"`
	ParentID string `gorm:"type:varchar(255)"`
	Method   string `gorm:"default:GET"`
	Headers  []Header
	Body     Body
	Url      Url
}

type Url struct {
	ID        string `gorm:"type:varchar(255)"`
	RequestID string `gorm:"type:varchar(255)"`
	Raw       string
	Queries   []Query
}

type Header struct {
	ID          string
	RequestID   string `gorm:"type:varchar(255)"`
	Key         string
	Value       string
	Disabled    bool `gorm:"default:false"`
	Type        string
	Description string
}

type Query struct {
	ID          string
	UrlID       string `gorm:"type:varchar(255)"`
	Key         string `gorm:"null"`
	Value       string `gorm:"null"`
	Disabled    bool   `default:"false"`
	Description string
}

type Body struct {
	ID        string
	RequestID string `gorm:"type:varchar(255)"`
	Mode      string
	Raw       string
	Options   string
	Disabled  bool `default:"false"`
}

type Example struct {
	ID       string `gorm:"type:varchar(255)"`
	Name     string
	ItemID   string `gorm:"type:varchar(255) Null"`
	FolderID string `gorm:"type:varchar(255)"`
}

// type Event struct {
// 	ID       string
// 	Listen   string
// 	ScriptID string
// 	Script   Script
// 	Disabled bool `default:"false"`
// }

// type Variable struct {
// 	ID            string
// 	Key           string
// 	Value         string
// 	Type          string
// 	Name          string
// 	DescriptionID string
// 	Description   Description
// 	System        bool `default:"false"`
// 	Disabled      bool `default:"false"`
// }

// type Description struct {
// 	ID      string
// 	Content string
// }

// type Request struct {
// 	ID string
// 	UrlID         string
// 	Url           Url
// 	ProxyID       string
// 	Proxy         Proxy
// Method string
// DescriptionID string
// Description   Description
// HeaderID      string
// Header        Header
// BodyID        string
// Body          Body
// }

type Response struct {
	ID        string
	ExampleID string `gorm:"type:varchar(255)"`
	Header    string
	Body      string
}

// type Proxy struct {
// 	ID       string
// 	Match    string `default:"http+https://*/*"`
// 	Host     string
// 	Port     uint `default:"8080"`
// 	Tunnel   bool `default:"false"`
// 	Disabled bool `default:"false"`
// }

// type Header struct {
// 	ID            string
// 	Key           string
// 	Value         string
// 	Disabled      string
// 	DescriptionID string
// 	Description   Description
// }

// type Body struct {
// 	ID           string
// 	Mode         string
// 	Raw          string
// 	Graphql      string
// 	UrlEncodedID string
// 	UrlEncoded   UrlEncoded
// 	FormatData   string
// 	Options      string
// 	Disabled     bool `default:"false"`
// }

// type UrlEncoded struct {
// 	ID            string
// 	Key           string
// 	Value         string
// 	Disabled      bool `default:"false"`
// 	DescriptionID string
// 	Description   Description
// }

// type Script struct {
// 	ID   string
// 	Type string
// 	Exec string
// 	Src  string
// 	Url  Url `gorm:"foreignKey:Src"`
// 	Name string
// }
