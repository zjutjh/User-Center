package main

import (
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zjutjh/User-Center/common/dbx"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type rootConfig struct {
	UserRPC struct {
		Mysql dbx.MysqlConf
	}
}

var (
	configFile = flag.String("f", "config.yaml", "the config file")
	tables     = []string{"user", "student", "college", "mini_program_user"}
)

func main() {
	flag.Parse()

	var c rootConfig
	conf.MustLoad(*configFile, &c)

	db, err := gorm.Open(mysql.Open(c.UserRPC.Mysql.DSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./apps/user-rpc/internal/dao/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)

	m := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			return "int8"
		},
	}
	g.WithDataTypeMap(m)

	for _, table := range tables {
		tableName := g.GenerateModel(
			table,
			gen.FieldType("deleted_at", "soft_delete.DeletedAt"),
			gen.FieldGORMTag("deleted_at", func(tag field.GormTag) field.GormTag {
				return tag.Set("softDelete", "milli")
			}),
			gen.FieldJSONTag("deleted_at", "-"),
		)
		g.ApplyBasic(tableName)
	}

	g.Execute()
}
