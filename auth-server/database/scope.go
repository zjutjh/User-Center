package database

import "auth-server/util"

type Scope struct {
	ID  		uint			`gorm:"primaryKey"`
	UserID		string
	ClientID	string
	Permission  string
}

func TrimAuthorizedScope(scopes []string, clientId string, userId string) (result []string) {
	var ss []Scope
	db.Where(&Scope{ClientID: clientId, UserID: userId}).Find(&ss)
	for i := 0; i < len(scopes); i++ {
		ok := false
		for j := 0; j < len(ss); j++ {
			if util.IsSubPermission(scopes[i], ss[j].Permission) {
				ok = true
				break
			}
		}
		if !ok {
			result = append(result, scopes[i])
		}
	}
	return
}

func ValidationTokenScope(scopes []string, clientId string, userId string) bool {
	var ss []Scope
	db.Where(&Scope{ClientID: clientId, UserID: userId}).Find(&ss)
	for i := 0; i < len(scopes); i++ {
		ok := false
		for j := 0; j < len(ss); j++ {
			if util.IsSubPermission(scopes[i], ss[j].Permission) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func NewScopeRecord(scopes []string, clientId string, userId string) {
	for i := 0; i < len(scopes); i++ {
		scope := Scope{UserID: userId, ClientID: clientId, Permission: scopes[i]}
		db.Create(&scope)
	}
}