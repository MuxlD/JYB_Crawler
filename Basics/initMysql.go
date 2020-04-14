package Basics

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDB *gorm.DB


func MysqlInit(name, passwd, host, port, datebase string) {
	fmt.Println("Init mysql...")
	//dbConn := "root:332214@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	var dbConn = name + ":" + passwd + "@(" + host + ":" + port + ")/" + datebase + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open("mysql", dbConn)
	if err != nil {
		fmt.Println("数据库连接失败：", err)
	}
	//db.SingularTable(true)
	//// 启用gger，显示详细日志
	db.LogMode(true)

	MysqlDB = db
}

//获取MySql连接
func GetDB() *gorm.DB {
	return MysqlDB
}

func CreateTable(){
	GetDB().AutoMigrate(
		Type{},
		TsUrl{},
	)
}

func StartMySql() {
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
