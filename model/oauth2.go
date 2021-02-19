package model

type OAuth2TokenResult struct {
	AccessToken 	string		`json:"access_token"`
	RefreshToken	string		`json:"refresh_token"`
	Scope			string		`json:"scope"`
	TokenType		string		`json:"token_type"`
	Error			string		`json:"error"`
	ErrorDesc		string		`json:"error_description"`
}

type AuthRequiredError struct {}

func (a *AuthRequiredError) Error() string {
	return "auth_required"
}