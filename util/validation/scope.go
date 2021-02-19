package validation

import (
  "strings"
  "ucenter/config"
)

func IsValidPermission(p string) bool {
  res := strings.Split(p, ".")
  for i := 0; i < len(res); i++ {
    if res[i] == "" {
      return false
    }
  }
  return true
}

func IsSubPermission(current, parent string) bool {
  if !strings.HasPrefix(current, parent) {
    return false
  }
  tail := strings.TrimPrefix(current, parent)
  if tail == "" || strings.HasPrefix(tail, ".") {
    return true
  }
  return false
}

func IsAppliedPermission(current, clientId string) bool {
  cfg := config.GetConfig()
  length := len(cfg.Client)
  for i := 0; i < length; i++ {
    if cfg.Client[i].ID == clientId {
      length = len(cfg.Client[i].Scope)
      ok := false
      for j := 0; j < length; j++ {
        if IsSubPermission(current, cfg.Client[i].Scope[j]) {
          ok = true
          break
        }
      }
      return ok
    }
  }
  return false
}
