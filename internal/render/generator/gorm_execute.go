/*
 * @Author: lida lidaemail@qq.com
 * @LastEditors: lida lidaemail@qq.com
 */
package generator_handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/imlida/go-gin-api/configs"
	"github.com/imlida/go-gin-api/internal/pkg/core"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type gormExecuteRequest struct {
	Db     string `form:"db" binding:"required"`
	Tables string `form:"tables" binding:"required"`
}

func (h *handler) GormExecute() core.HandlerFunc {
	dir, _ := os.Getwd()
	projectPath := strings.Replace(dir, "\\", "/", -1)
	outPath := projectPath + "/internal/dal/query"

	return func(c core.Context) {
		req := new(gormExecuteRequest)
		if err := c.ShouldBindPostForm(req); err != nil {
			c.Payload("参数传递有误")
			return
		}

		dbReq := req.Db
		dbCfg := strings.Split(dbReq, "|")
		mysqlConf := configs.Get().MySQL[strings.TrimSpace(dbCfg[0])]
		dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local", mysqlConf.User, mysqlConf.Pass, mysqlConf.Addr, mysqlConf.Name, mysqlConf.Charset)

		// specify the output directory (default: "./query")
		// ### if you want to query without context constrain, set mode gen.WithoutContext ###
		g := gen.NewGenerator(gen.Config{
			OutPath: outPath,
			Mode:    gen.WithQueryInterface,
			/* Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface*/
			//if you want the nullable field generation property to be pointer type, set FieldNullable true
			/* FieldNullable: true,*/
			//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
			/* FieldCoverable: true,*/
			// if you want generate field with unsigned integer type, set FieldSignable true
			/* FieldSignable: true,*/
			//if you want to generate index tags from database, set FieldWithIndexTag true
			/* FieldWithIndexTag: true,*/
			//if you want to generate type tags from database, set FieldWithTypeTag true
			/* FieldWithTypeTag: true,*/
			//if you need unit tests for query code, set WithUnitTest true
			/* WithUnitTest: true, */
		})

		// reuse the database connection in Project or create a connection here
		// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
		db, _ := gorm.Open(mysql.Open(dsn))
		g.UseDB(db)

		// apply basic crud api on structs or table models which is specified by table name with function
		// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
		// 想对已有的model生成crud等基础方法可以直接指定model struct ，例如model.User{}
		// 如果是想直接生成表的model和crud方法，则可以指定表的名称，例如g.GenerateModel("company")
		// 想自定义某个表生成特性，比如struct的名称/字段类型/tag等，可以指定opt，例如g.GenerateModel("company",gen.FieldIgnore("address")), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address"))
		// g.ApplyBasic(model.User{}, g.GenerateModel("company"), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address")))

		tables := strings.Split(req.Tables, ",")

		for _, table := range tables {
			g.ApplyBasic(g.GenerateModel(table))
		}

		// g.ApplyBasic(g.GenerateModel("user"))

		// generate all tables, ex: g.ApplyBasic(g.GenerateAllTable()...)
		// g.GenerateAllTable()

		// generate a model struct map to table `people` in database
		// g.GenerateModel("people")

		// generate a struct and specify struct's name
		// g.GenerateModelAs("people", "People")

		// apply diy interfaces on structs or table models
		// 如果想给某些表或者model生成自定义方法，可以用ApplyInterface，第一个参数是方法接口，可以参考DIY部分文档定义
		// g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

		// execute the action of code generation
		g.Execute()

		c.Payload("生成成功")

	}
}
