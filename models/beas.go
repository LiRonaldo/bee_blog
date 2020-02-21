package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	//username:password@tcp(127.0.0.1:3306)/db_name?charset=utf8
	url := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	orm.RegisterDataBase("default", "mysql", url)
	orm.RegisterModel(new(User), new(Category), new(Post), new(Config), new(Comment))

}

func TableName(str string) string {
	return beego.AppConfig.String("dbprefix") + str
}
