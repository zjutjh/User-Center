package main

import (
	"github.com/spf13/cobra"
	"github.com/zjutjh/mygo/foundation/command"
	"github.com/zjutjh/mygo/ndb"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/zjutjh/User-Center/register"
)

var tables = []string{
	"user",
	"student",
	"college",
	"mini_program_user",
}

func main() {
	command.Execute(
		register.Boot,
		func(c *cobra.Command) {},
		func(cmd *cobra.Command, args []string) error { return nil },
	)

	g := gen.NewGenerator(gen.Config{
		OutPath: "./biz/dao/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(ndb.Pick())

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
