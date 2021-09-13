/*
	麻将的显示， 画出来的

 */
package view

import (
	"fmt"
	"strings"
	"sync"
)

var (
	// 牌花色的格式
	// 1
	t1 = map[int][]int{
		3:{5},
	}
	// 2
	t2 = map[int][]int{
		2:{5},
		4:{5},
	}
	// 3
	t3 = map[int][]int{
		2:{4},
		3:{5},
		4:{6},
	}
	// 4
	t4 = map[int][]int{
		2:{3, 7},
		4:{3, 7},
	}
	// 5
	t5 = map[int][]int{
		2:{3, 7},
		3:{5},
		4:{3, 7},
	}
	// 6
	t6 = map[int][]int{
		2:{4, 6},
		3:{4, 6},
		4:{4, 6},
	}
	// 7
	t7 = map[int][]int{
		1:{3},
		2:{5},
		3:{7,3,4},
		4:{3,4},
	}
	// 8
	t8 = map[int][]int{
		1:{4, 6},
		2:{4, 6},
		3:{4, 6},
		4:{4, 6},
	}
	// 9
	t9 = map[int][]int{
		2:{3,5,7},
		3:{3,5,7},
		4:{3,5,7},
	}
	// 7条用
	tt7 = map[int][]int{
		2:[]int{5},
		3:[]int{3,5,7},
		4:[]int{3,5,7},
	}

	// 筒
	Tong1 = PaiShow(t1, "0")
	Tong2 = PaiShow(t2, "0")
	Tong3 = PaiShow(t3, "0")
	Tong4 = PaiShow(t4, "0")
	Tong5 = PaiShow(t5, "0")
	Tong6 = PaiShow(t6, "0")
	Tong7 = PaiShow(t7, "0")
	Tong8 = PaiShow(t8, "0")
	Tong9 = PaiShow(t9, "0")

	// 条
	Tiao1 = PaiShow(t1, "I")
	Tiao2 = PaiShow(t2, "I")
	Tiao3 = PaiShow(t3, "I")
	Tiao4 = PaiShow(t4, "I")
	Tiao5 = PaiShow(t5, "I")
	Tiao6 = PaiShow(t6, "I")
	Tiao7 = PaiShow(tt7, "I")
	Tiao8 = Tiao8Txt()
	Tiao9 = PaiShow(t9, "I")

	// 万
	Wan1 = Wan("一")
	Wan2 = Wan("二")
	Wan3 = Wan("三")
	Wan4 = Wan("四")
	Wan5 = Wan("五")
	Wan6 = Wan("六")
	Wan7 = Wan("七")
	Wan8 = Wan("八")
	Wan9 = Wan("九")

	PaiWan = []string{Wan1, Wan2, Wan3, Wan4, Wan5, Wan6, Wan7, Wan8, Wan9}
	PaiTong = []string{Tong1, Tong2, Tong3, Tong4, Tong5, Tong6, Tong7, Tong8, Tong9}
	PaiTiao = []string{Tiao1, Tiao2, Tiao3, Tiao4, Tiao5, Tiao6, Tiao7, Tiao8, Tiao9}
	PaiKey = []string{}

	PaiValue = map[string]int{
		Wan1:11,  Wan2:12,  Wan3:13,  Wan4:14,  Wan5:15,  Wan6:16,  Wan7:17,  Wan8:18,  Wan9:19,
		Tong1:21, Tong2:22, Tong3:23, Tong4:24, Tong5:25, Tong6:26, Tong7:27, Tong8:28, Tong9:29,
		Tiao1:31, Tiao2:32, Tiao3:33, Tiao4:34, Tiao5:35, Tiao6:36, Tiao7:37, Tiao8:38, Tiao9:39,
	}

	PaiValue2 = map[int]string{
		11:Wan1,  12:Wan2,  13:Wan3,  14:Wan4,  15:Wan5,  16:Wan6,  17:Wan7,  18:Wan8,  19:Wan9,
		21:Tong1, 22:Tong2, 23:Tong3, 24:Tong4, 25:Tong5, 26:Tong6, 27:Tong7, 28:Tong8, 29:Tong9,
		31:Tiao1, 32:Tiao2, 33:Tiao3, 34:Tiao4, 35:Tiao5, 36:Tiao6, 37:Tiao7, 38:Tiao8, 39:Tiao9,
	}

	PaiValue3 = map[int]string{
		11:"1万",  12:"2万",  13:"3万",  14:"4万",  15:"5万",  16:"6万",  17:"7万",  18:"8万",  19:"9万",
		21:"1筒",  22:"2筒",  23:"3筒",  24:"4筒",  25:"5筒",  26:"6筒",  27:"7筒",  28:"8筒",  29:"9筒",
		31:"1条",  32:"2条",  33:"3条",  34:"4条",  35:"5条",  36:"6条",  37:"7条",  38:"8条",  39:"9条",
	}

	Pai []int

	gameChan = make(chan string)
	endChan = make(chan string)

	// 用户当前的牌 - 用于操作
	nowPai = []int{}
	nowPaiLock = sync.Mutex{}

	nowPaiVal = []int{}
	nowPaiMap = map[string]int{}
	nowPaiText = map[string]string{}

)

var pay = [][]string{
	{" ", "_", "_", "_", "_", "_", "_", "_", "_", "_", " "},
	{" ", "|", " ", " ", " ", " ", " ", " ", " ", "|", " "},
	{" ", "|", " ", " ", " ", " ", " ", " ", " ", "|", " "},
	{" ", "|", " ", " ", " ", " ", " ", " ", " ", "|", " "},
	{" ", "|", " ", " ", " ", " ", " ", " ", " ", "|", " "},
	{" ", "|", "_", "_", "_", "_", "_", "_", "_", "|", " "},
}

func PaiShow(sit map[int][]int, hua string) string {
	// 深copy二维数组
	dst:=make([][]string,len(pay))
	for i,_:=range pay{
		dst[i]=make([]string,len(pay[0]))
		copy(dst[i],pay[i])
	}
	// 牌
	rList := []string{}
	for k, v := range dst{
		if n, ok := sit[k]; ok {
			for _,m := range n{
				v[m] = hua
			}
		}
		rList = append(rList, strings.Join(v,""))
	}
	pShow := strings.Join(rList, "\n")
	//fmt.Println(pShow)
	return pShow
}

func Tiao8Txt() string {
	return ` _________ 
 |       | 
 |  |/\| | 
 |       | 
 |  |\/| | 
 |_______| `
}

func Wan(s string) string {
	return fmt.Sprintf(` _________ 
 |       | 
 |   %s  | 
 |       | 
 |   万  | 
 |_______| `, s)
}