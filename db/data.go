package db

import (
  "fmt"
  "log"
  "strings"
  "ucenter/config"
)

type Data struct {
  ID          uint64        `gorm:"primaryKey"`
  Username    string        `gorm:"unique"`
}

func initResourceDataTable() {
  if err := resourceDB.AutoMigrate(&Data{}); err != nil {
    log.Println(err)
  }
  dataFields := config.GetDataFields()
  for k, v := range *dataFields {
    rawType := typeResolver(strings.Split(v.Type, ":")[0])
    if !hasColumn(resourceDB, "data", k, rawType) {
      newColumn(resourceDB, "data", k, rawType)
    }
  }
}

func typeResolver(rawType string) string {
  if rawType == "text" {
    return "text"
  } else if rawType == "int" {
    return "bigint"
  }
  log.Fatalln("type resolver: unknown type")
  return ""
}

func GetDataRecordByUser(userId string, path string) interface{} {
  var result string
  resourceDB.Raw(fmt.Sprintf("select \"%s\" from data where username = '%s'", path, userId)).Scan(&result)
  return result
}