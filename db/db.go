package db

import (
  "fmt"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "log"
  "ucenter/config"
)

func initDB(cfg config.Db) *gorm.DB {
  if cfg.Type != "postgres" {
    log.Fatalln("unsupported database type")
  }
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", cfg.Host, cfg.User, cfg.Password,
    cfg.DbName, cfg.Port)
  var err error
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    log.Fatalln(err)
  }
  return db
}

func hasColumn(db *gorm.DB, tableName string, columnName string, rawType string) bool {
  var dataType string
  db.Raw("select data_type " +
    "from information_schema.columns " +
    "where table_name = ? and column_name = ?", tableName, columnName).Scan(&dataType)
  if dataType == "" {
    return false
  }
  if dataType != rawType {
    log.Fatalln(fmt.Sprintf("type conflicts: there is already a column in %s named %s with another type %s",
      tableName, columnName, dataType))
  }
  return true
}

func newColumn(db *gorm.DB, tableName string, columnName string, rawType string) {
  result := db.Exec(fmt.Sprintf("alter table %s add \"%s\" %s", tableName, columnName, rawType))
  if result.Error != nil {
    log.Fatalln(result.Error)
  }
}

var authDB, resourceDB *gorm.DB

func init() {
  authDB = initDB(config.GetConfig().Db.Auth)
  err := authDB.AutoMigrate(&User{}, &Scope{})
  if err != nil {
    log.Fatalln(err)
  }
  resourceDB = initDB(config.GetConfig().Db.Resource)
  initResourceDataTable()
}
