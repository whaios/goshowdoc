package main

import (
	"strconv"

	"github.com/whaios/goshowdoc/datadict"
	"github.com/whaios/goshowdoc/log"
	"github.com/whaios/goshowdoc/runapi"
)

// UpdateDataDict 生成数据字典
func UpdateDataDict(driver, host, user, pwd, db, schema, cat string) {
	var dd datadict.DataDict
	switch driver {
	case datadict.MySQL:
		dd = datadict.NewMySQL(host, user, pwd, db)
	case datadict.PostgreSQL:
		dd = datadict.NewPostgreSQL(host, user, pwd, db, schema)
	case datadict.SQLServer:
		dd = datadict.NewSQLServer(host, user, pwd, db)
	case datadict.SQlite:
		dd = datadict.NewSqlite(host)
	default:
		log.Error("不支持的数据库类型 %s", driver)
		return
	}

	if err := dd.Open(); err != nil {
		log.Error("无法连接到数据库: %s", err.Error())
		return
	}
	defer dd.Close()

	tbs, err := dd.Query()
	if err != nil {
		log.Error("查询数据库出错: %s", err.Error())
		return
	}
	max := len(tbs)
	for i, tb := range tbs {
		if err = runapi.UpdateByApi(cat, tb.Name, strconv.FormatInt(int64(i+1), 10), tb.Markdown()); err != nil {
			log.Error("更新文档[%s/%s]失败: %s", cat, tb.Name, err.Error())
			return
		}
		log.DrawProgressBar("更新文档", i+1, max)
	}
	log.Success("更新完成")
}
