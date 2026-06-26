package dbx

import (
	"fmt"
	"net/url"
)

type MysqlConf struct {
	Host      string `json:",optional"`
	Port      int    `json:",optional"`
	Username  string
	Password  string
	Database  string
	Charset   string `json:",optional"`
	ParseTime bool   `json:",optional"`
	Loc       string `json:",optional"`
}

func (c MysqlConf) DSN() string {
	host := c.Host
	if host == "" {
		host = "127.0.0.1"
	}

	port := c.Port
	if port == 0 {
		port = 3306
	}

	charset := c.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	loc := c.Loc
	if loc == "" {
		loc = "Asia/Shanghai"
	}

	parseTime := c.ParseTime
	if !c.ParseTime {
		parseTime = true
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		host,
		port,
		c.Database,
		charset,
		parseTime,
		url.QueryEscape(loc),
	)
}
