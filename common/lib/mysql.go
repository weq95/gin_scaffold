package lib

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBMySQL  *gorm.DB

//初始化数据库
func init()  {
	db,err:= gorm.Open(mysql.New(mysql.Config{
		DSN: "root:123456@tcp(127.0.0.1:3306)/microservice?charset=utf8&parseTime=True&loc=Local",
		DefaultStringSize: 256,//string类型字段的默认长度
		DisableDatetimePrecision: true,//禁用datetime精度,MySQL5.6之前的数据库不支持
		DontSupportRenameIndex: true,//重命名索引时采用删除并新建的方式, MySQL5.7 之前的数据库和MariaDB 不支持重命名索引
		DontSupportRenameColumn: true,//用`change`重命名列, MySQL8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,//根据当前MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DBMySQL = db
}

func GetMysql() *gorm.DB {
	return DBMySQL
}