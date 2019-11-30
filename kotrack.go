/*Copyright (C) 2019 EricWu <skynocover@gmail.com>
See COPYING for license details.*/
package main

import (
	"github.com/zserge/lorca"
	"log"
	"net/url"

	"fmt"
	. "github.com/GoSpider"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (
	viewtitle      = "視窗標題"
	inputlabel     = "標題"
	bottonname     = "按鈕名稱"
	vieweight  int = 1000
	viewheight int = 820

	ui, err             = lorca.New("", "", vieweight, viewheight)
	streamname []string //追串的名稱陣列
	website    []string //追串的網址陣列
	replynum   []int    //追串的回應數陣列

	webs = [][]string{{"aqua.komica.org/04/pixmicat", "luna.komica.org/12/pixmicat.php", "sora.komica.org/00/pixmicat", "aqua.komica.org/cs/pixmicat."},
		{"gzone-anime.info/UnitedSites/academic/pixmicat", "alleyneblade.mymoe.moe/queensblade/", "2cat.komica.org/~tedc21thc/new/pixmicat", "2cat.komica.org/~tedc21thc/live/pixmicat", "2cat.komica.org/~kirur/orz/pixmicat"}}
)

func main() {
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	window := read("./file/index.html")
	ui.Load("data:text/html," + url.PathEscape(window))
	load()      //載入setting內儲存的資料,若無則新建一個
	htmlInput() //將陣列放入至html內
	refresh()   //重新確認每個討論串
	//按鈕_新增
	ui.Bind("add", func(input string) {
		if len(streamname) == 6 { //若此刻討論串數量以達到6
			ui.Eval(`window.alert("討論串數量超過上限") `)
			ui.Eval(`document.getElementById("input").value = ""`)
		} else if websiteCheck(input) == 99 { //確認網址是否是符合的版面
			ui.Eval(`window.alert("網址無法使用,非正確版面") `)
			ui.Eval(`document.getElementById("input").value = ""`)
		} else {
			html := Gethtml(input) //抓出輸入網址的html
			websiteSwitch(input, html)
			htmlInput() //將討論串標題與回覆放入頁面
			save()      //儲存陣列到setting
			ui.Eval(`document.getElementById("input").value = ""`)
		}
	})
	ui.Bind("del", func(selected string) {
		sel, err := strconv.Atoi(selected) //選擇的index
		if err != nil {
			ui.Eval(`window.alert("未選擇討論串") `) //接起err防止程式崩潰
		} else {
			for i := sel; i < len(streamname)-1; i++ { //從選擇的那一個index開始每一個往前一個移動
				streamname[i] = streamname[i+1]
				website[i] = website[i+1]
				replynum[i] = replynum[i+1]
			}
			//將所有陣列的長度-1
			arr := streamname[0 : len(streamname)-1]
			streamname = arr
			arr = website[0 : len(website)-1]
			website = arr
			arrtemp := replynum[0 : len(replynum)-1]
			replynum = arrtemp
			htmlInput() //將陣列放入html
			//selectInput() //將標題放入到select
			initial(len(streamname), 6) //從現有陣列長度開始初始化
			save()
		}
	})
	ui.Bind("go", func(selected string) {
		sel, err := strconv.Atoi(selected)
		if err != nil { //通常發生於沒有討論串的時候
			ui.Eval(`window.alert("未選擇討論串") `) //接起err防止程式崩潰
		} else {
			openbrowser(website[sel])
		}
	})
	<-ui.Done()
}
func websiteSwitch(web, html string) {
	k := websiteCheck(web)
	switch k {
	case 0:
		test := catch(html, "<div class=\"thread\"", "<div class=\"post reply\"")
		streamname = append(streamname, catch(test, "<span class=\"title\">", "</span")) //將討論串標題放入陣列
		replynum = append(replynum, strings.Count(html, "<div class=\"post reply\""))    //將討論串回覆數放入陣列
		website = append(website, web)                                                   //將網址放入陣列
	case 1:
		test := catch(html, "<div class=\"threadpost\"", "<div class=\"reply\"")
		streamname = append(streamname, catch(test, "<span class=\"title\">", "</span>")) //將討論串標題放入陣列
		replynum = append(replynum, strings.Count(html, "<div class=\"reply\""))          //將討論串回覆數放入陣列
		website = append(website, web)                                                    //將網址放入陣列
	}
}

//確認網址是否正確
func websiteCheck(web string) int {
	for i, value := range webs {
		for j, _ := range value {
			if strings.Contains(web, webs[i][j]) {
				return i
			}
		}
	}
	return 99
}

//初始化streamname以及replynum
func initial(str int, end int) {
	//從陣列長度 (若長度現在是5則從title5(第6個)開始初始化)
	for i := str; i < end; i++ {
		ui.Eval(`document.getElementById("title` + strconv.Itoa(i) + `").textContent = "追串` + strconv.Itoa(i+1) + `"`)
		ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").textContent = "回應數"`)
	}
}

//將陣列放入頁面內
func htmlInput() {
	ui.Eval(`var op = document.getElementById("select")`)
	ui.Eval(`op.options.length=0`) //將select初始化
	for i := 0; i < len(streamname); i++ {
		ui.Eval(`document.getElementById("title` + strconv.Itoa(i) + `").innerHTML = "` + streamname[i] + `"`)                 //放入標題
		ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").innerHTML =  "回應數` + strconv.Itoa(replynum[i]) + `"`) //放入回覆數

		ui.Eval(`op[` + strconv.Itoa(i) + `] = new Option("` + streamname[i] + `","` + strconv.Itoa(i) + `")`) //放入選擇討論串
	}
}

//重新確認每個討論串
func refresh() {
	for i := 0; i < len(replynum); i++ {
		k := websiteCheck(website[i])
		temp := 0
		switch k {
		case 0:
			temp = strings.Count(Gethtml(website[i]), "<div class=\"post reply\"") //確認現在的回應數
		case 1:
			temp = strings.Count(Gethtml(website[i]), "<div class=\"reply\"") //確認現在的回應數
		}
		if temp > replynum[i] { //若比現在的大則更新
			ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").innerHTML =  "有` + strconv.Itoa(temp-replynum[i]) + `個新回應"`)
			//ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").style.color = "red"`) //放入回覆數
			replynum[i] = temp
		} else if temp == 0 {
			ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").innerHTML =  "討論串被刪除"`)
			replynum[i] = temp
		}
	}
	save()
}

//將三個陣列儲存至setting
func save() {
	file, err := os.Create("./file/setting")
	if err != nil {
		return
	}
	defer file.Close()
	for i := 0; i < len(streamname); i++ {
		file.WriteString(streamname[i] + "\n")
		file.WriteString(website[i] + "\n")
		file.WriteString(strconv.Itoa(replynum[i]) + "\n")
	}
}

//將setting的內容讀出來放到陣列中
func load() {
	bs, err := ioutil.ReadFile("./file/setting") //讀取檔案
	if err != nil {
		file, err := os.Create("./file/setting") //若無則建立檔案
		if err != nil {
			return
		}
		defer file.Close()
	} else {
		streamInput(strings.Split(string(bs), "\n")) //將讀取到的輸入進陣列中
	}
}

//將字串放入streamname以及replynum
func streamInput(input []string) {
	//清空陣列
	streamname = nil
	website = nil
	replynum = nil
	for i := 0; i < len(input)-2; i = i + 3 { //輸入的陣列為三個一循環
		streamname = append(streamname, input[i])
		website = append(website, input[i+1])
		temp, _ := strconv.Atoi(input[i+2])
		replynum = append(replynum, temp)
	}
}

//Tool

//讀取檔案
func read(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

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

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
