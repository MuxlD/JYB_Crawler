package Basics

import "github.com/jinzhu/gorm"

const (
	JYBFL = "https://cs.jiaoyubao.cn/edu/"
	JYB = "https://cs.jiaoyubao.cn"
)

func InitConf()  {
	ConfSql=MySql{
		"root",
		"332214",
		"127.0.0.1",
		"3306",
		"test",
	}
}

var ConfSql MySql

type MySql struct {
	Name     string
	Password string
	Host     string
	Port     string
	Database string

}
func StartMySql()  {
	InitConf()
	MysqlInit(
		ConfSql.Name,
		ConfSql.Password,
		ConfSql.Host,
		ConfSql.Port,
		ConfSql.Database,
	)
	CreateTable()
}
type TrainingSchool struct {
	ID          string   `json:"id"`
	TypeName    string   `json:"type_name"`
	TypeUrl     string   `json:"type_url"`
	TypeId      uint     `json:"type_id"`
	Name        string   `json:"name"`
	Url         string   `json:"url"`
	BrightSpot  []string `json:"bright_spot"`             //亮点，特色
	Info        string   `gorm:"size:1000" json:"info"`   //简介
	Course      []string `gorm:"size:2000" json:"course"` //课程
	Campus      string   `gorm:"size:1000" json:"campus"` //校区
	PhoneNumber string   `json:"phone_number"`            //联系电话
}

type Type struct {
	gorm.Model
	Name    string `json:"name"`
	TypeUrl string `json:"type_url"`
	MaxPage int    `json:"max_page"`
	Count   int    `json:"count"`
}

type TsUrl struct {
	gorm.Model
	TypeID uint   `json:"type_id"`
	Url    string `json:"url"`
}