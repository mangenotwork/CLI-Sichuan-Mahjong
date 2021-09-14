package game

import (
	"fmt"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/Client/view"
	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/utils"
	"os"
	"strings"
	"text/template"
	"time"
)

func Game(){
	nowPai := []int{11,12,13,14,15,16,17,18,19,21,12,13,14,15}
	// 渲染发牌动画
	for j := 1; j<len(nowPai); j++{
		GameView(nowPai[:j], false)
		time.Sleep(60 * time.Millisecond) // dlay
	}
	GameView(nowPai, true)
	for{}
}

var nowPaiMap = map[string]int{}
var nowPaiText = map[string]string{}


// 打牌区数据结构
type DpqData struct {
	Syp string // 剩余牌    XX张
	GameTime string // 游戏时长   XX秒
	NowDoingUser string // 玩家4
	User1Name string //非自己的玩家1   玩家1 (操作倒计时 10s)
	User1SP string  // 玩家1-手牌   XX张
	User1PG string  // 玩家1-碰杠   碰X筒, 杠X万
	User1DP string  // 玩家1-打的牌   X万, X万, X条
	User2Name string //非自己的玩家2   玩家2 (操作倒计时 10s)
	User2SP string  // 玩家2-手牌   XX张
	User2PG string  // 玩家2-碰杠   碰X筒, 杠X万
	User2DP string  // 玩家2-打的牌   X万, X万, X条
	User3Name string //非自己的玩家3   玩家3 (操作倒计时 10s)
	User3SP string  // 玩家3-手牌   XX张
	User3PG string  // 玩家3-碰杠   碰X筒, 杠X万
	User3DP string  // 玩家3-打的牌   X万, X万, X条
	MyName string	// 自己的名称  玩家4 你自己 (操作倒计时 10s)
	MySP string  // 自己的-手牌   XX张
	MyPG string  // 自己的-碰杠   碰X筒, 杠X万
	MyDP string  // 自己的-打的牌   X万, X万, X条
}

func GameView(mypai []int, isDp bool){
	utils.Cle()

	fmt.Println(`===============================麻将 - 血战到底===============================
         `)
	mb := `
剩余牌: {{.Syp}}    游戏时长: {{.GameTime}}     当前操作玩家: {{.NowDoingUser}}

------------------------------- {{.User1Name}} ---------------------------------------
手牌: {{.User1SP}}  
碰杠: {{.User1PG}}
已打: {{.User1DP}}

------------------------------- {{.User2Name}} ---------------------------------------
手牌: {{.User2SP}}
碰杠: {{.User2PG}}
已打: {{.User2DP}}

------------------------------- {{.User3Name}} ---------------------------------------
手牌: {{.User3SP}}
碰杠: {{.User3PG}}
已打: {{.User3DP}}

------------------------------- 自己: {{.MyName}} ---------------------------------------
手牌: {{.MySP}}
碰杠: {{.MyPG}}
已打: {{.MyDP}}
`

	dpqShow := DpqData{
		Syp:          "88张",
		GameTime:     "50秒",
		NowDoingUser: "玩家4",
		User1Name:    "玩家1",
		User1SP:      "6张",
		User1PG:      "碰3筒, 杠3万",
		User1DP:      "1万, 2万, 3条, ",
		User2Name:    "玩家2",
		User2SP:      "6张",
		User2PG:      "碰3筒, 杠3万",
		User2DP:      "1万, 2万, 3条, ",
		User3Name:    "玩家3",
		User3SP:      "6张",
		User3PG:      "碰3筒, 杠3万",
		User3DP:      "1万, 2万, 3条, ",
		MyName:       "玩家4 你自己 (操作倒计时 10s)",
		MySP:         "6张",
		MyPG:         "碰3筒, 杠3万",
		MyDP:         "1万, 2万, 3条, ",
	}

	tmpl, _ := template.New("test").Parse(mb)
	_ = tmpl.Execute(os.Stdout, dpqShow)

	// 碰牌区
	fmt.Println(`

------------------------------- 你的牌 ---------------------------------------
`)

	// 初始化操作数据
	nowPaiMap = make(map[string]int)
	nowPaiText = make(map[string]string)

	// 显示画出来的牌
	showList := make([]string, 0)

	// 如果还有牌
	if len(mypai) > 0 {
		paiTxt := make([]string, 0)
		// isDp 是否显示打牌序号
		if !isDp {
			for _, v := range mypai {
				paiTxt = append(paiTxt, view.PaiValue2[v])
			}
		} else {
			// 添加操作号
			for j, v := range mypai {
				num := j + 1
				value := view.PaiValue2[v] + fmt.Sprintf("\n   %d打牌   ", num)
				paiTxt = append(paiTxt, value)
				nowPaiText[fmt.Sprintf("%d", num)] = view.PaiValue3[v]
				nowPaiMap[fmt.Sprintf("%d", num)] = v
			}
		}

		showList = strings.Split(paiTxt[0], "\n")
		if len(showList) > 1 {
			for _, v := range paiTxt[1:] {
				//fmt.Println(v)
				showList = ShowAdd(v, showList)
			}
		}
	}else{
		showList = []string{"\n", "\n", "\n"}
	}

	fmt.Println(strings.Join(showList, "\n"))


	if len(mypai) < 1{
		fmt.Sprintf("没牌了....")
		return
	}

	fmt.Println(`
——————————————————————————————————————————————————————————————
游戏玩法(输入相关字母):
	q :退出
	a :换牌
	1~14 : 打牌
——————————————————————————————————————————————————————————————
输入:`)


}

func ShowAdd(t string, tList []string) []string {
	temp := strings.Split(t, "\n")
	for i:=0; i<len(tList);i++{
		tList[i] = tList[i]+temp[i]
	}
	return tList
}