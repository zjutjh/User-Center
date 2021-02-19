package db

import (
  "errors"
  "github.com/lib/pq"
  "gorm.io/gorm"
  "ucenter/model"
)

type User struct {
  ID          uint      `gorm:"primaryKey"`
  Username    string    `gorm:"unique"`
  Password    string
}

func VerifyPassword(username, password string) (bool, *model.Result) {
  var user User
  result := authDB.Where(&User{Username: username}).First(&user)
  if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    return false, model.NewResult(nil,  1002, "user not found")
  }
  if user.Password != password {
    return false, model.NewResult(nil, 1001, "incorrect username or password")
  }
  return true, nil
}

func SetPassword(username, password string) {
  var user User
  authDB.Where(&User{Username: username}).First(&user)
  user.Password = password
  authDB.Save(&user)
}

func NewUser(username, password string) (bool, string) {
  user := User{Username: username, Password: password}
  result := authDB.Create(&user)
  if err, ok := result.Error.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
    return false, "user already exists"
  }
  if result.Error != nil {
    return false, "unknown internal error"
  }
  return true, ""
}
