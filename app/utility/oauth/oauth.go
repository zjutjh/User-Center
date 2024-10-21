package oauth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func CheckByOauth(username string, password string) (string, error) {
	f := Fetch{}
	f.Init()
	log.Println(username)
	log.Println(password)
	loginHome, err := f.GetRaw(OauthLoginHome())
	if err != nil {
		return "", err
	}
	if len(f.Cookie) < 1 {
		return "", err
	}
	doc, _ := goquery.NewDocumentFromReader(loginHome.Body)
	hiddenInput := doc.Find("input[type=hidden][name=execution]")

	execution := hiddenInput.AttrOr("value", "")
	loginData := genOauthLoginData(username, password, execution, &f)

	postRedirectUrl, err := f.PostFormRedirect(OauthLoginHome(), loginData)
	log.Println(postRedirectUrl)
	if err != nil {
		return "", err
	}
	f.Cookie = []*http.Cookie{}
	_, err = f.Get(postRedirectUrl.String())
	if err != nil {
		return "", err
	}
	log.Println(f.Cookie)
	s, err := f.Get("http://www.me.zjut.edu.cn/api/basic/info")
	if err != nil {
		return "", err
	}
	var m Model
	err = json.Unmarshal(s, &m)
	return m.Data.Yhm, err
}

func genOauthLoginData(username, password, execution string, f *Fetch) url.Values {
	s, _ := f.Get(OauthLoginGetPublickey())

	encodePassword, _ := GetEncryptPassword(s, password)
	return url.Values{
		"username":   {username},
		"mobileCode": {},
		"password":   {encodePassword},
		"authcode":   {},
		"execution":  {execution},
		"_eventId":   {"submit"}}
}

type Model struct {
	Data struct {
		Bmid              string `json:"bmid"`
		Bmmc              string `json:"bmmc"`
		HeadPortrait      string `json:"headPortrait"`
		IdentityCard      string `json:"identityCard"`
		Jsdm              string `json:"jsdm"`
		Jsmc              string `json:"jsmc"`
		LastAccessTime    string `json:"lastAccessTime"`
		Mbkjj             string `json:"mbkjj"`
		Mbtb              string `json:"mbtb"`
		Mbwb              string `json:"mbwb"`
		ModifyPwdUrl      string `json:"modifyPwdUrl"`
		Nc                string `json:"nc"`
		Qm                string `json:"qm"`
		RowEnd            int    `json:"row_end"`
		RowStart          int    `json:"row_start"`
		SecurityCenterUrl string `json:"securityCenterUrl"`
		Sfgly             string `json:"sfgly"`
		Txdz              string `json:"txdz"`
		UserType          string `json:"userType"`
		Username          string `json:"username"`
		Xszyzt            string `json:"xszyzt"`
		Xyid              string `json:"xyid"`
		Yhm               string `json:"yhm"`
		Ysjkzzdtb1        string `json:"ysjkzzdtb1"`
		Ysjkzzdtb2        string `json:"ysjkzzdtb2"`
	} `json:"data"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func OauthLoginGetPublickey() string {
	return "https://oauth.zjut.edu.cn/cas/v2/getPubKey"
}

func OauthLoginHome() string {
	return "https://oauth.zjut.edu.cn/cas/login?service=http%3A%2F%2Fwww.me.zjut.edu.cn%2F"
}
