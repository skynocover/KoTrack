/*Copyright (C) 2019 EricWu <skynocover@gmail.com>
See COPYING for license details.*/
package main

import (
	"github.com/zserge/lorca"
	"log"
	"net/url"

	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strconv"
)

var (
	vieweight  int = 1000
	viewheight int = 820

	ui, err = lorca.New("", "", vieweight, viewheight)

	file = "file/strdat"
	webs = [][]string{{"aqua.komica.org/04/pixmicat", "luna.komica.org/12/pixmicat.php", "sora.komica.org/00/pixmicat", "aqua.komica.org/cs/pixmicat."},
		{"gzone-anime.info/UnitedSites/academic/pixmicat", "alleyneblade.mymoe.moe/queensblade/", "2cat.komica.org/~tedc21thc/new/pixmicat", "2cat.komica.org/~tedc21thc/live/pixmicat", "2cat.komica.org/~kirur/orz/pixmicat"}}
)

func main() {
	var follow Follow
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	ui.Load("data:text/html," + url.PathEscape(read("./file/index.html")))

	follow.load()
	//點選增加按鈕
	ui.Bind("add", func(url string) {
		follow.input.Website = url
		if len(follow.Streams) == 6 { //若此刻討論串數量已達到6
			ui.Eval(`window.alert("討論串數量超過上限") `)
			ui.Eval(`document.getElementById("input").value = ""`)
		} else if follow.input.getType() == 99 { //確認網址是否是符合的版面
			ui.Eval(`window.alert("網址無法使用,非正確版面") `)
			ui.Eval(`document.getElementById("input").value = ""`)
		} else {
			follow.add() //增加跟串數量
			ui.Eval(`document.getElementById("input").value = ""`)
		}
	})
	//點選刪除按鈕
	ui.Bind("del", func(selected string) {
		sel, err := strconv.Atoi(selected) //選擇的index
		if err != nil {
			ui.Eval(`window.alert("未選擇討論串") `) //接起err防止程式崩潰
		} else {
			follow.del(sel)
			initial(len(follow.Streams), 6) //從現有陣列長度開始初始化
		}
	})
	//點選前往按鈕
	ui.Bind("go", func(selected string) {
		sel, err := strconv.Atoi(selected) //選擇的index
		if err != nil {                    //通常發生於沒有討論串的時候
			ui.Eval(`window.alert("未選擇討論串") `) //接起err防止程式崩潰
		} else {
			openbrowser(follow.Streams[sel].Website)
		}
	})
	<-ui.Done()
}

//初始化streamname以及replynum
func initial(str int, end int) {
	//從陣列長度 (若長度現在是5則從title5(第6個)開始初始化)
	for i := str; i < end; i++ {
		ui.Eval(`document.getElementById("title` + strconv.Itoa(i) + `").textContent = "追串` + strconv.Itoa(i+1) + `"`)
		ui.Eval(`document.getElementById("reply` + strconv.Itoa(i) + `").textContent = "回應數"`)
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
