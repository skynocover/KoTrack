package main

import (
	"strconv"
)

type Follow struct {
	input   Stream   //從頁面獲得的輸入
	Streams []Stream //跟串列表
}

//載入跟串
func (follow *Follow) load() {
	Load(file, follow) //載入檔案
	follow.setHtml()   //放置html
	follow.refresh()   //更新跟串
	Save(file, follow)
}

//放入頁面內
func (follow *Follow) setHtml() {
	ui.Eval(`var op = document.getElementById("select")`)
	ui.Eval(`op.options.length=0`) //將select初始化
	for i := 0; i < len(follow.Streams); i++ {
		ui.Eval(`document.getElementById("title` + strconv.Itoa(i) + `").innerHTML = "` + follow.Streams[i].Name + `"`)                       //放入標題
		ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").innerHTML =  "回應數` + strconv.Itoa(follow.Streams[i].Replynum) + `"`) //放入回覆數
		ui.Eval(`op[` + strconv.Itoa(i) + `] = new Option("` + follow.Streams[i].Name + `","` + strconv.Itoa(i) + `")`)                       //放入選擇討論串
	}
}

//重新確認每個討論串
func (follow *Follow) refresh() {
	for i := 0; i < len(follow.Streams); i++ {
		temp := follow.Streams[i]
		temp.get()
		if temp.Replynum > follow.Streams[i].Replynum {
			ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").innerHTML =  "有` + strconv.Itoa(temp.Replynum-follow.Streams[i].Replynum) + `個新回應"`)
			follow.Streams[i] = temp
		} else if temp.Replynum == 0 {
			ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").innerHTML =  "討論串被刪除"`)
			follow.Streams[i] = temp
		}
	}
}

//增加跟串的數量
func (follow *Follow) add() {
	follow.Streams = append(follow.Streams, follow.input)
	follow.Streams[len(follow.Streams)-1].get()
	follow.setHtml()
	Save(file, &follow)
}

//刪除跟串
func (follow *Follow) del(sel int) {
	for i := sel; i < len(follow.Streams)-1; i++ { //從選擇的那一個index開始每一個往前一個移動
		follow.Streams[i] = follow.Streams[i+1]
	}
	//將所有陣列的長度-1
	arr := follow.Streams[0 : len(follow.Streams)-1]
	follow.Streams = arr
	follow.setHtml()
	Save(file, &follow)
}
