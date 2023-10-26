package utility

import "regexp"

func IsValidEmail(email string) bool {
	// 正则表达式用于验证邮箱格式
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`

	// 编译正则表达式
	re := regexp.MustCompile(regexPattern)

	// 使用正则表达式匹配邮箱
	return re.MatchString(email)
}
