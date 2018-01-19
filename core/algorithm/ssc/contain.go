package ssc

import (
	"fmt"
	"xmn/core/model"
	"strings"
	"time"
	"log"
	"xmn/core/logger"
	"xmn/core/mail"
	"strconv"
	"sync"
)

//数据包
var contain_datapackage []*model.Packet

//重庆开奖数据
var contain_cq_data []*model.Cqssc

//新疆开奖数据
var contain_xj_data []*model.Xjssc

//天津开奖数据
//var contain_tj_data []*model.Tjssc

/*
//台湾开奖数据
var contain_tw_data []*model.Twssc
*/

//彩票类型
var contain_ssc_type map[int]string

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
var newCodes *code

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
type code struct {
	codes map[int]string //数据包id => 该彩种的最新开奖号码 每个数据包对应的是一个彩种
	lock sync.RWMutex
}

//多协程 共享 各彩种最新开奖号 支撑并发 读取 写入
func init()  {
	newCodes = new(code)
	newCodes.codes = make(map[int]string)
}

//时时彩
//包含数据包 算法
func Contain()  {
	fmt.Println("时时彩 - 包含数据包 算法")

	contain_ssc_type = make(map[int]string)
	contain_ssc_type[1] = "重庆时时彩"
	contain_ssc_type[2] = "天津时时彩"
	contain_ssc_type[3] = "新疆时时彩"
	contain_ssc_type[4] = "台湾五分彩"

	packet := new(model.Packet)
	contain_datapackage = packet.Query()

	cqssc := new(model.Cqssc)
	contain_cq_data = cqssc.Query("300")

	xjssc := new(model.Xjssc)
	contain_xj_data = xjssc.Query("300")

	/*
	tjssc := new(model.Tjssc)
	contain_tj_data = tjssc.Query("300")
	*/

	/*
	twssc := new(model.Twssc)
	contain_tw_data = twssc.Query("300")
	*/

	containAnalysis()
}

func containAnalysis()  {
	for i := range contain_datapackage {
		go containAnalysisCodes(contain_datapackage[i])
	}
}

func containAnalysisCodes(packet *model.Packet)  {
	//log.Println(contain_ssc_type[packet.Type], "时时彩－包含数据包 正在分析　数据包别名:", packet.Alias)
	slice_dataTxt := strings.Split(packet.DataTxt, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMap := make(map[string]string)
	for i := range slice_dataTxt {
		dataTxtMap[slice_dataTxt[i]] = slice_dataTxt[i]
	}

	//fmt.Println(dataTxtMap)

	//检查是否在报警时间段以内
	if (packet.Start >0 && packet.End >0) && (time.Now().Hour() < packet.Start || time.Now().Hour() > packet.End)  {
		log.Println("彩票类型:", contain_ssc_type[packet.Type], "数据包别名:", packet.Alias, "报警通知非接受时间段内")
		logger.Log("彩票类型: " +  contain_ssc_type[packet.Type] + " 数据包别名: " + packet.Alias + " 报警通知非接受时间段内 ")
		return
	}

	//开奖数据
	codes := make([]string, 0)
	//重庆时时彩
	if packet.Type == 1 && len(contain_cq_data) > 0 {
		//检查 该彩种到最新的一期 是否重复分析
		new_code := contain_cq_data[len(contain_cq_data) - 1].One + contain_cq_data[len(contain_cq_data) - 1].Two + contain_cq_data[len(contain_cq_data) - 1].Three + contain_cq_data[len(contain_cq_data) - 1].Four + contain_cq_data[len(contain_cq_data) - 1].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := newCodes.Get(packet.Id)
		if newcode == new_code {
			log.Println(contain_ssc_type[packet.Type], "数据包别名:", packet.Alias, "最新的一期 已经分析过了... 等待出现新的开奖号")
			return
		} else {
			//最新开奖号 与 内存中的最新开奖号 不相同 刷新内存最新开奖号值
			newCodes.Set(packet.Id, new_code)
		}

		for i := range contain_cq_data {
			code := contain_cq_data[i].One + contain_cq_data[i].Two + contain_cq_data[i].Three + contain_cq_data[i].Four +contain_cq_data[i].Five
			codes = append(codes, code)
		}
	}

	//新疆时时彩
	if packet.Type == 3 && len(contain_xj_data) > 0 {
		//检查 该彩种到最新的一期 是否重复分析
		new_code := contain_xj_data[len(contain_xj_data) - 1].One + contain_xj_data[len(contain_xj_data) - 1].Two + contain_xj_data[len(contain_xj_data) - 1].Three + contain_xj_data[len(contain_xj_data) - 1].Four + contain_xj_data[len(contain_xj_data) - 1].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := newCodes.Get(packet.Id)
		if newcode == new_code {
			log.Println(contain_ssc_type[packet.Type], "数据包别名:", packet.Alias, "最新的一期 已经分析过了... 等待出现新的开奖号")
			return
		} else {
			//最新开奖号 与 内存中的最新开奖号 不相同 刷新内存最新开奖号值
			newCodes.Set(packet.Id, new_code)
		}

		for i := range contain_xj_data {
			code := contain_xj_data[i].One + contain_xj_data[i].Two + contain_xj_data[i].Three + contain_xj_data[i].Four +contain_xj_data[i].Five
			codes = append(codes, code)
		}
	}

	/*
	//天津时时彩
	if packet.Type == 2 && len(contain_tj_data) > 0 {
		//检查 该彩种到最新的一期 是否重复分析
		new_code := contain_tj_data[len(contain_tj_data) - 1].One + contain_tj_data[len(contain_tj_data) - 1].Two + contain_tj_data[len(contain_tj_data) - 1].Three + contain_tj_data[len(contain_tj_data) - 1].Four + contain_tj_data[len(contain_tj_data) - 1].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := newCodes.Get(packet.Id)
		if newcode == new_code {
			log.Println(contain_ssc_type[packet.Type], "数据包别名:", packet.Alias, "最新的一期 已经分析过了... 等待出现新的开奖号")
			return
		} else {
			//最新开奖号 与 内存中的最新开奖号 不相同 刷新内存最新开奖号值
			newCodes.Set(packet.Id, new_code)
		}

		for i := range contain_tj_data {
			code := contain_tj_data[i].One + contain_tj_data[i].Two + contain_tj_data[i].Three + contain_tj_data[i].Four +contain_tj_data[i].Five
			codes = append(codes, code)
		}
	}
	*/

	/*
	//台湾时时彩
	if packet.Type == 4 && len(contain_tw_data) > 0 {
		//检查 该彩种到最新的一期 是否重复分析
		new_code := contain_tw_data[len(contain_tw_data) - 1].One + contain_tw_data[len(contain_tw_data) - 1].Two + contain_tw_data[len(contain_tw_data) - 1].Three + contain_tw_data[len(contain_tw_data) - 1].Four + contain_tw_data[len(contain_tw_data) - 1].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := newCodes.Get(packet.Id)
		if newcode == new_code {
			log.Println(contain_ssc_type[packet.Type], "数据包别名:", packet.Alias, "最新的一期 已经分析过了... 等待出现新的开奖号")
			return
		} else {
			//最新开奖号 与 内存中的最新开奖号 不相同 刷新内存最新开奖号值
			newCodes.Set(packet.Id, new_code)
		}

		for i := range contain_tw_data {
			code := contain_tw_data[i].One + contain_tw_data[i].Two + contain_tw_data[i].Three + contain_tw_data[i].Four +contain_tw_data[i].Five
			codes = append(codes, code)
		}
	}
	*/

	//fmt.Println(contain_ssc_type[packet.Type])
	//fmt.Println(codes)

	//各单位报警期数 初始化
	var q3_number int = 0
	var z3_number int = 0
	var h3_number int = 0

	//各单位报警 是否有上期参考对象 初始化
	var q3_reference bool = false
	var z3_reference bool = false
	var h3_reference bool = false

	var q3_html_log string = ""
	var z3_html_log string = ""
	var h3_html_log string = ""

	//周期计数器 自定义周期 一清零
	var q3_cycle int = 0
	var z3_cycle int = 0
	var h3_cycle int = 0

	//周期报警计数器
	var q3_cycle_number int = 0
	var z3_cycle_number int = 0
	var h3_cycle_number int = 0

	for i := range codes{
		code_byte := []byte(codes[i])
		//前三号码
		q3 := string(code_byte[0]) + string(code_byte[1]) + string(code_byte[2])
		//中三号码
		z3 := string(code_byte[1]) + string(code_byte[2]) + string(code_byte[3])
		//后三号码
		h3 := string(code_byte[2]) + string(code_byte[3]) + string(code_byte[4])

		//各单位是否在 数据包内 初始化
		var q3_in bool = false
		var z3_in bool = false
		var h3_in bool = false

		//前三号码 是否在数据包内
		if _, ok := dataTxtMap[q3]; ok {
			q3_in = true
		}
		//中三号码 是否在数据包内
		if _, ok := dataTxtMap[z3]; ok {
			z3_in = true
		}
		//后三号码 是否在数据包内
		if _, ok := dataTxtMap[h3]; ok {
			h3_in = true
		}

		//前三没有上一期 开奖数据 参考对象 and 前三出现在数据包里
		if !q3_reference && q3_in {
			q3_number = q3_number + 1
			q3_cycle += 1
			q3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 前三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + q3 +"</div>"
			q3_html_log += "<div>没有基准 并且 当前号码在数据包内 +1 = " + strconv.Itoa(q3_number) + "</div>"
			q3_html_log += "<div>当期自定义连续 "+ strconv.Itoa(packet.Cycle) +" 期 未开为1周期,  计算是否满1周的值累加到 = " + strconv.Itoa(q3_cycle) + "</div>"

			//检查是否满1周期
			if q3_cycle == packet.Cycle {
				q3_cycle_number += 1
				q3_cycle = 0
				q3_html_log += "<div>[前三 当期满自定义一周期, 周期报警计数 +1 = "+ strconv.Itoa(q3_cycle_number) +"]</div>"
				q3_html_log += "<div>1周期已满, 重新计算连续未开是否满1周期 累计值当前为0</div>"
			}
			q3_html_log += "<br/>"
			//fmt.Println(contain_ssc_type[packet.Type], "q3", q3, "+1=", q3_number)
		} else if q3_reference && q3_in  {
			//前三有上一期 开奖数据 参考对象 and 前三出现在数据包里
			q3_number = 0
			q3_number = q3_number + 1
			q3_cycle = 0
			q3_cycle = q3_cycle + 1
			q3_cycle_number = 0
			q3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 前三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + q3 +"</div>"
			q3_html_log += "<div>有基准 并且 当前号码在数据包内 清零 再 +1 = " + strconv.Itoa(q3_number) + "</div>"
			q3_html_log += "<div>计算是否满1周期 的 计数器 清零 再 +1 = " + strconv.Itoa(q3_cycle) + "</div>"
			q3_html_log += "<div>周期报警期数 清零 = " + strconv.Itoa(q3_cycle_number) + "</div>"

			//检查是否满1周期
			if q3_cycle == packet.Cycle {
				q3_cycle_number += 1
				q3_cycle = 0
				q3_html_log += "<div>[前三 当期满自定义一周期, 周期报警计数 +1 = "+ strconv.Itoa(q3_cycle_number) +"]</div>"
				q3_html_log += "<div>1周期已满, 重新计算连续未开是否满1周期 累计值当前为0</div>"
			}
			q3_html_log += "<br/>"
			//fmt.Println(contain_ssc_type[packet.Type], "q3", q3, "清0 +1=", q3_number)
		} else if !q3_reference && !q3_in {
			q3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 前三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + q3 +"</div>"
			q3_html_log += "<div>没有基准 并且 当前号码不在数据包内 [不管]"+ "</div></br>"
		} else if q3_reference && !q3_in {
			q3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 前三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + q3 +"</div>"
			q3_html_log += "<div>有基准 并且 当前号码不在数据包内 [不管]"+ "</div></br>"
		}

		//中三没有上一期 开奖数据 参考对象 and 中三出现在数据包里
		if !z3_reference && z3_in {
			z3_number = z3_number + 1
			z3_cycle += 1
			z3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 中三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + z3 +"</div>"
			z3_html_log += "<div>没有基准 并且 当前号码在数据包内 +1 = " + strconv.Itoa(z3_number) + "</div>"
			z3_html_log += "<div>当期自定义连续 "+ strconv.Itoa(packet.Cycle) +" 期 未开为1周期,  计算是否满1周的值累加到 = " + strconv.Itoa(z3_cycle) + "</div>"

			//检查是否满1周期
			if z3_cycle == packet.Cycle {
				z3_cycle_number += 1
				z3_cycle = 0
				z3_html_log += "<div>[中三 当期满自定义一周期, 周期报警计数 +1 = "+ strconv.Itoa(z3_cycle_number) +"]</div>"
				z3_html_log += "<div>1周期已满, 重新计算连续未开是否满1周期 累计值当前为0</div>"
			}
			z3_html_log += "<br/>"
		} else if z3_reference && z3_in  {
			//中三有上一期 开奖数据 参考对象 and 中三出现在数据包里
			z3_number = 0
			z3_number = z3_number + 1
			z3_cycle = 0
			z3_cycle = z3_cycle + 1
			z3_cycle_number = 0
			z3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 中三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + z3 +"</div>"
			z3_html_log += "<div>有基准 并且 当前号码在数据包内 清零 再 +1 = " + strconv.Itoa(z3_number) + "</div>"
			z3_html_log += "<div>计算是否满1周期 的 计数器 清零 再 +1 = " + strconv.Itoa(z3_cycle) + "</div>"
			z3_html_log += "<div>周期报警期数 清零 = " + strconv.Itoa(z3_cycle_number) + "</div>"

			//检查是否满1周期
			if z3_cycle == packet.Cycle {
				z3_cycle_number += 1
				z3_cycle = 0
				z3_html_log += "<div>[中三 当期满自定义一周期, 周期报警计数 +1 = "+ strconv.Itoa(z3_cycle_number) +"]</div>"
				z3_html_log += "<div>1周期已满, 重新计算连续未开是否满1周期 累计值当前为0</div>"
			}
			z3_html_log += "<br/>"
		} else if !z3_reference && !z3_in {
			z3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 中三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + z3 +"</div>"
			z3_html_log += "<div>没有基准 并且 当前号码不在数据包内 [不管]"+ "</div></br>"
		} else if z3_reference && !z3_in {
			z3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 中三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + z3 +"</div>"
			z3_html_log += "<div>有基准 并且 当前号码不在数据包内 [不管]"+ "</div></br>"
		}

		//后三没有上一期 开奖数据 参考对象 and 后三出现在数据包里
		if !h3_reference && h3_in {
			h3_number = h3_number + 1
			h3_cycle += 1
			h3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 后三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + h3 +"</div>"
			h3_html_log += "<div>没有基准 并且 当前号码在数据包内 +1 = " + strconv.Itoa(h3_number) + "</div>"
			h3_html_log += "<div>当期自定义连续 "+ strconv.Itoa(packet.Cycle) +" 期 未开为1周期,  计算是否满1周的值累加到 = " + strconv.Itoa(h3_cycle) + "</div>"

			//检查是否满1周期
			if h3_cycle == packet.Cycle {
				h3_cycle_number += 1
				h3_cycle = 0
				h3_html_log += "<div>[后三 当期满自定义一周期, 周期报警计数 +1 = "+ strconv.Itoa(h3_cycle_number) +"]</div>"
				h3_html_log += "<div>1周期已满, 重新计算连续未开是否满1周期 累计值当前为0</div>"
			}
			z3_html_log += "<br/>"
		} else if h3_reference && h3_in  {
			//后三有上一期 开奖数据 参考对象 and 后三出现在数据包里
			h3_number = 0
			h3_number = h3_number + 1
			h3_cycle = 0
			h3_cycle = h3_cycle + 1
			h3_cycle_number = 0
			h3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 后三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + h3 +"</div>"
			h3_html_log += "<div>有基准 并且 当前号码在数据包内 清零 再 +1 = " + strconv.Itoa(h3_number) + "</div>"
			h3_html_log += "<div>计算是否满1周期 的 计数器 清零 再 +1 = " + strconv.Itoa(h3_cycle) + "</div>"
			h3_html_log += "<div>周期报警期数 清零 = " + strconv.Itoa(h3_cycle_number) + "</div>"

			//检查是否满1周期
			if h3_cycle == packet.Cycle {
				h3_cycle_number += 1
				h3_cycle = 0
				h3_html_log += "<div>[后三 当期满自定义一周期, 周期报警计数 +1 = "+ strconv.Itoa(h3_cycle_number) +"]</div>"
				h3_html_log += "<div>1周期已满, 重新计算连续未开是否满1周期 累计值当前为0</div>"
			}
			h3_html_log += "<br/>"
		} else if !h3_reference && !h3_in {
			h3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 后三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + h3 +"</div>"
			h3_html_log += "<div>没有基准 并且 当前号码不在数据包内 [不管]"+ "</div></br>"
		} else if h3_reference && !h3_in {
			h3_html_log += "<div>数据包别名: "+ packet.Alias + " 位置: 后三 " + " 数据包别名: " + packet.Alias+ " 开奖号: " + h3 +"</div>"
			h3_html_log += "<div>有基准 并且 当前号码不在数据包内 [不管]"+ "</div></br>"
		}

		//前三参考对象
		if q3_in {
			q3_reference = true
		} else {
			q3_reference = false
		}

		//中三参考对象
		if z3_in {
			z3_reference = true
		} else {
			z3_reference = false
		}

		//后三参考对象
		if h3_in {
			h3_reference = true
		} else {
			h3_reference = false
		}
	}

	//最新的一期有数据包里的数据 才报警
	if !q3_reference {
		q3_number = 0
	}
	if !z3_reference {
		z3_number = 0
	}
	if !h3_reference {
		h3_number = 0
	}

	//fmt.Println(contain_ssc_type[packet.Type], "q3 期数", q3_number)
	//fmt.Println(contain_ssc_type[packet.Type], "z3 期数", z3_number)
	//fmt.Println(contain_ssc_type[packet.Type], "h3 期数", h3_number)

	//fmt.Println(contain_ssc_type[packet.Type], "q3 周期数", q3_cycle_number)
	//fmt.Println(contain_ssc_type[packet.Type], "z3 周期数", z3_cycle_number)
	//fmt.Println(contain_ssc_type[packet.Type], "h3 周期数", h3_cycle_number)

	//fmt.Println(q3_html_log)
	//fmt.Println(z3_html_log)
	//fmt.Println(h3_html_log)

	var body string = ""
	var cycle_body string = ""

	//前三报警
	if q3_number >= packet.RegretNumber {
		body += "<div> 彩种: " + contain_ssc_type[packet.Type] + " 数据包别名: " + packet.Alias + " 位置 前三 " + strconv.Itoa(q3_number) + " 期 报警！</div><br/>"
	}

	//中三报警
	if z3_number >= packet.RegretNumber {
		body += "<div> 彩种: " + contain_ssc_type[packet.Type] + " 数据包别名: " + packet.Alias + " 位置 中三 " + strconv.Itoa(z3_number) + " 期 报警！</div><br/>"
	}

	//后三报警
	if h3_number >= packet.RegretNumber {
		body += "<div> 彩种: " + contain_ssc_type[packet.Type] + " 数据包别名: " + packet.Alias + " 位置 后三 " + strconv.Itoa(h3_number) + " 期 报警！</div><br/>"
	}

	//发送邮件
	if body != "" {
		go mail.SendMail(contain_ssc_type[packet.Type] + " 包含数据包", body)
	}

	//前三 自定义周期报警
	if q3_reference && q3_cycle_number > 0 && q3_cycle_number >= packet.CycleNumber {
		cycle_body += "<div> 自定义周期 彩种: "+ contain_ssc_type[packet.Type] + " 数据包别名: "+ packet.Alias + " 位置 前三 "+ strconv.Itoa(q3_cycle_number) +" 周期报警! </div><br/>"
		//cycle_body += q3_html_log
	}

	//中三 自定义周期报警
	if z3_reference && z3_cycle_number > 0 && z3_cycle_number >= packet.CycleNumber {
		cycle_body += "<div> 自定义周期 彩种: "+ contain_ssc_type[packet.Type] + " 数据包别名: "+ packet.Alias + " 位置 中三 "+ strconv.Itoa(z3_cycle_number) +" 周期报警! </div><br/>"
		//cycle_body += z3_html_log
	}

	//后三 自定义周期报警
	if h3_reference && h3_cycle_number > 0 && h3_cycle_number >= packet.CycleNumber {
		cycle_body += "<div> 自定义周期 彩种: "+ contain_ssc_type[packet.Type] + " 数据包别名: "+ packet.Alias + " 位置 后三 "+ strconv.Itoa(h3_cycle_number) +" 周期报警! </div><br/>"
		//cycle_body += h3_html_log
	}

	//自定义周期报警 发送邮件
	if cycle_body != "" && contain_ssc_type[packet.Type] != "台湾五分彩"  {
		//log.Println(contain_ssc_type[packet.Type], "自定义周期 包含数据包, 正在发送邮件")
		//log.Println("前三参考对象:", q3_reference, "前三:", q3_cycle_number)
		//log.Println("中三参考对象:", z3_reference, "中三:", z3_cycle_number)
		//log.Println("后三参考对象:", h3_reference, "后三:", h3_cycle_number)
		go mail.SendMail(contain_ssc_type[packet.Type] + " 自定义周期 包含数据包", cycle_body)
	}
}

func (c *code) Get(k int) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.codes[k]
}

func (c *code) Set(k int, v string)  {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.codes[k] = v
}