package http_helper

import (
  "encoding/json"
  "net/http"
  "ucenter/model"
)

func RawHttpError(w http.ResponseWriter, err string, code int, status int) {
  res, _ := json.Marshal(model.NewResult(nil, code, err))
  http.Error(w, string(res[:]), status)
}
