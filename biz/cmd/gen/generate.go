package main

import (
	"github.com/zjutjh/User-Center/biz/database"
	"github.com/zjutjh/User-Center/biz/viper"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

var tables = []string{
	"user",
	"student",
	"college",
	"mini_program_user",
}

func main() {
	viper.InitViper()
	g := gen.NewGenerator(gen.Config{
		OutPath: "./dao/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	database.NewRunOptions().Init()
	g.UseDB(database.DB)

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
