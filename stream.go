package main

import (
	. "github.com/GoSpider"
	"strings"
)

type Stream struct {
	Website  string
	Replynum int
	Name     string
}

//確認網址是否正確
func (stream *Stream) getType() int {
	for i, value := range webs {
		for j, _ := range value {
			if strings.Contains(stream.Website, webs[i][j]) {
				return i
			}
		}
	}
	return 99
}

//取得名稱與回應數
func (stream *Stream) get() {
	html := Gethtml(stream.Website)
	switch stream.getType() {
	case 0:
		test := catch(html, "<div class=\"thread\"", "<div class=\"post reply\"")
		stream.Name = catch(test, "<span class=\"title\">", "</span")      //將討論串標題放入陣列
		stream.Replynum = strings.Count(html, "<div class=\"post reply\"") //將討論串回覆數放入陣列
	case 1:
		test := catch(html, "<div class=\"threadpost\"", "<div class=\"reply\"")
		stream.Name = catch(test, "<span class=\"title\">", "</span>") //將討論串標題放入陣列
		stream.Replynum = strings.Count(html, "<div class=\"reply\"")  //將討論串回覆數放入陣列
	}
}

//Tool

//抓出中間文字
func catch(input, str, end string) string {
	strn := strings.Index(input, str) + strings.Count(str, "") - 1
	endn := strings.Index(input, end)
	if strn < endn {
		catchtml := string(input[strn:endn])
		return catchtml
	} else {
		return ""
	}
}
