package play2

import (
	"xmn/core/model"
	"sync"
	"time"
	"log"
	"xmn/core/logger"
	"strconv"
	"bytes"
	"xmn/core/mail"
	"fmt"
)

var consecutive []*model.Play2

//重庆开奖数据
var consecutive_cq_data []*model.Cqssc

//新疆开奖数据
var consecutive_xj_data []*model.Xjssc

//天津开奖数据
//var consecutive_tj_data []*model.Tjssc

/*
//台湾开奖数据
var consecutive_tw_data []*model.Twssc
*/

//间隔几连号
var consecutiveNumbers map[string]string = make(map[string]string)

//时时彩最新开奖号码
var consecutiveSscNewCodes *consecutive_ssc_new_codes

//时时彩最新开奖号码 加读写锁 防止并行写入map 导致程序崩溃
type consecutive_ssc_new_codes struct {
	codes map[int]string //彩种类型 => 该彩种的最新开奖号码
	lock sync.RWMutex
}

func init()  {
	//初始化连号
	consecutiveNumbers["01"] = "01"
	consecutiveNumbers["12"] = "12"
	consecutiveNumbers["23"] = "23"
	consecutiveNumbers["34"] = "34"
	consecutiveNumbers["45"] = "45"
	consecutiveNumbers["56"] = "56"
	consecutiveNumbers["67"] = "67"
	consecutiveNumbers["78"] = "78"
	consecutiveNumbers["89"] = "89"
	consecutiveNumbers["90"] = "90"

	//初始化 时时彩 最新开奖号码
	consecutiveSscNewCodes = new(consecutive_ssc_new_codes)
	consecutiveSscNewCodes.codes = make(map[int]string)
}

// 时时彩 间隔几连号 算法 连号 01 12 23 34 45 56 78 89 90 视为连号
func Consecutive()  {
	fmt.Println("时时彩 - 间隔几连号 算法")
	play2 := new(model.Play2)
	consecutive = play2.Query(model.Play2Consecutive)

	cqssc := new(model.Cqssc)
	consecutive_cq_data = cqssc.Query("100")

	xjssc := new(model.Xjssc)
	consecutive_xj_data = xjssc.Query("100")

	/*
	tjssc := new(model.Tjssc)
	consecutive_tj_data = tjssc.Query("100")
	*/

	/*
	twssc := new(model.Twssc)
	consecutive_tw_data = twssc.Query("100")
	*/

	consecutiveAnalysis()
}

func consecutiveAnalysis()  {
	for i := range consecutive {
		go consecutiveAnalysisCodes(consecutive[i])
	}
}

func consecutiveAnalysisCodes(config *model.Play2)  {
	//检查是否在报警时间段以内
	if (config.Start >0 && config.End >0) && (time.Now().Hour() < config.Start || time.Now().Hour() > config.End)  {
		log.Println("时时彩-间隔几连号 报警通知非接受时间段内")
		logger.Log("时时彩-连续 报警通知非接受时间段内")
		return
	}

	cq_q3s, cq_z3s, cq_h3s := getCqCodes()
	xj_q3s, xj_z3s, xj_h3s := getXjCodes()
	//tj_q3s, tj_z3s, tj_h3s := getTjCodes()
	/*
	tw_q3s, tw_z3s, tw_h3s := getTwCodes()
	*/

	go func(config *model.Play2) {
		//重庆报警
		var body string
		q3_log_html, q3_num := consecutiveCodesAnalyse(config, cq_q3s, "前三", CpTypeName[CqsscType])
		z3_log_html, z3_num := consecutiveCodesAnalyse(config, cq_z3s, "中三", CpTypeName[CqsscType])
		h3_log_html, h3_num := consecutiveCodesAnalyse(config, cq_h3s, "后三", CpTypeName[CqsscType])
		//_, q3_num := consecutiveCodesAnalyse(config, cq_q3s, "前三", CpTypeName[CqsscType])
		//_, z3_num := consecutiveCodesAnalyse(config, cq_z3s, "中三", CpTypeName[CqsscType])
		//_, h3_num := consecutiveCodesAnalyse(config, cq_h3s, "后三", CpTypeName[CqsscType])

		name := " 连号: " + strconv.Itoa(config.Number) + "期, 周期: "+strconv.Itoa(config.Cycle)+ " 期 "
		if q3_num == config.Cycle {
			body += "<div> 彩种: " + CpTypeName[CqsscType] + " 间隔几连号 "+name+" 报警提示 位置: 前三 周期数: "+ strconv.Itoa(q3_num) + "</div>"
		}
		if z3_num == config.Cycle {
			body += "<div> 彩种: " + CpTypeName[CqsscType] + " 间隔几连号 "+name+" 报警提示 位置: 中三 周期数: "+ strconv.Itoa(z3_num) + "</div>"
		}
		if h3_num == config.Cycle {
			body += "<div> 彩种: " + CpTypeName[CqsscType] + " 间隔几连号 "+name+" 报警提示 位置: 后三 周期数: "+ strconv.Itoa(h3_num) + "</div>"
		}
		body += q3_log_html
		body += z3_log_html
		body += h3_log_html

		if q3_num == config.Cycle || z3_num == config.Cycle || h3_num == config.Cycle {
			//发送邮件
			mail.SendMail(CpTypeName[CqsscType] + " 间隔几连号", body)
		}
	}(config)

	/*
	go func(config *model.Play2) {
		//天津报警
		var body string
		_, q3_num := consecutiveCodesAnalyse(tj_q3s, "前三", CpTypeName[TjsscType])
		_, z3_num := consecutiveCodesAnalyse(tj_z3s, "中三", CpTypeName[TjsscType])
		_, h3_num := consecutiveCodesAnalyse(tj_h3s, "后三", CpTypeName[TjsscType])
		if q3_num == config.Number {
			body += "<div> 彩种: " + CpTypeName[TjsscType] + " 连号报警提示 位置: 前三 期数: "+ strconv.Itoa(q3_num) + "</div>"
		}
		if z3_num == config.Number {
			body += "<div> 彩种: " + CpTypeName[TjsscType] + " 连号报警提示 位置: 中三 期数: "+ strconv.Itoa(z3_num) + "</div>"
		}
		if h3_num == config.Number {
			body += "<div> 彩种: " + CpTypeName[TjsscType] + " 连号报警提示 位置: 后三 期数: "+ strconv.Itoa(h3_num) + "</div>"
		}
		//body += q3_log_html
		//body += z3_log_html
		//body += h3_log_html

		if q3_num == config.Number || z3_num == config.Number || h3_num == config.Number {
			//发送邮件
			mail.SendMail(CpTypeName[TjsscType] + " 间隔几连号", body)
		}
	}(config)
	*/

	go func(config *model.Play2) {
		//新疆报警
		var body string
		q3_log_html, q3_num := consecutiveCodesAnalyse(config, xj_q3s, "前三", CpTypeName[XjsscType])
		z3_log_html, z3_num := consecutiveCodesAnalyse(config, xj_z3s, "中三", CpTypeName[XjsscType])
		h3_log_html, h3_num := consecutiveCodesAnalyse(config, xj_h3s, "后三", CpTypeName[XjsscType])
		//_, q3_num := consecutiveCodesAnalyse(config, xj_q3s, "前三", CpTypeName[XjsscType])
		//_, z3_num := consecutiveCodesAnalyse(config, xj_z3s, "中三", CpTypeName[XjsscType])
		//_, h3_num := consecutiveCodesAnalyse(config, xj_h3s, "后三", CpTypeName[XjsscType])

		name := " 连号: " + strconv.Itoa(config.Number) + "期, 周期: "+strconv.Itoa(config.Cycle)+ " 期 "
		if q3_num == config.Cycle {
			body += "<div> 彩种: " + CpTypeName[XjsscType] + " 间隔几连号 "+name+" 报警提示 位置: 前三 周期数: "+ strconv.Itoa(q3_num) + "</div>"
		}
		if z3_num == config.Cycle {
			body += "<div> 彩种: " + CpTypeName[XjsscType] + " 间隔几连号 "+name+" 报警提示 位置: 中三 周期数: "+ strconv.Itoa(z3_num) + "</div>"
		}
		if h3_num == config.Cycle {
			body += "<div> 彩种: " + CpTypeName[XjsscType] + " 间隔几连号 "+name+" 报警提示 位置: 后三 周期数: "+ strconv.Itoa(h3_num) + "</div>"
		}
		body += q3_log_html
		body += z3_log_html
		body += h3_log_html

		if q3_num == config.Cycle || z3_num == config.Cycle || h3_num == config.Cycle {
			//发送邮件
			mail.SendMail(CpTypeName[XjsscType] + " 间隔几连号", body)
		}
	}(config)

	/*
	go func(config *model.Play2) {
		//台湾报警
		var body string
		//q3_log_html, q3_num := consecutiveCodesAnalyse(tw_q3s, "前三", CpTypeName[TwsscType])
		//z3_log_html, z3_num := consecutiveCodesAnalyse(tw_z3s, "中三", CpTypeName[TwsscType])
		//h3_log_html, h3_num := consecutiveCodesAnalyse(tw_h3s, "后三", CpTypeName[TwsscType])
		_, q3_num := consecutiveCodesAnalyse(tw_q3s, "前三", CpTypeName[TwsscType])
		_, z3_num := consecutiveCodesAnalyse(tw_z3s, "中三", CpTypeName[TwsscType])
		_, h3_num := consecutiveCodesAnalyse(tw_h3s, "后三", CpTypeName[TwsscType])
		if q3_num >= config.Number {
			body += "<div> 彩种: " + CpTypeName[TwsscType] + " 连号报警提示 位置: 前三 期数: "+ strconv.Itoa(q3_num) + "</div>"
		}
		if z3_num >= config.Number {
			body += "<div> 彩种: " + CpTypeName[TwsscType] + " 连号报警提示 位置: 中三 期数: "+ strconv.Itoa(z3_num) + "</div>"
		}
		if h3_num >= config.Number {
			body += "<div> 彩种: " + CpTypeName[TwsscType] + " 连号报警提示 位置: 后三 期数: "+ strconv.Itoa(h3_num) + "</div>"
		}
		//body += q3_log_html
		//body += z3_log_html
		//body += h3_log_html

		if q3_num >= config.Number || z3_num >= config.Number || h3_num >= config.Number {
			//发送邮件
			mail.SendMail(CpTypeName[TwsscType] + " 间隔几连号", body)
		}
	}(config)
	*/

}

//获取重庆 前中后的 开奖号码
func getCqCodes() ([]string, []string, []string) {
	q3s := make([]string, 0)
	z3s := make([]string, 0)
	h3s := make([]string, 0)
	//检查是否重复分析 开奖号码 防止1期号码重复分析与报警
	//重庆检查
	if len := len(consecutive_cq_data); len >0 {
		index := len - 1
		//该彩种到最新的一期 开奖号码
		new_code := consecutive_cq_data[index].One + consecutive_cq_data[index].Two + consecutive_cq_data[index].Three + consecutive_cq_data[index].Four + consecutive_cq_data[index].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := consecutiveSscNewCodes.Get(CqsscType)
		if new_code == newcode {
			log.Println("时时彩-间隔几连号 彩票类型:", CpTypeName[CqsscType], "已经分析过了,等待新的开奖号码出现...")
			return q3s, z3s, h3s
		}

		for i:= range consecutive_cq_data {
			q3s = append(q3s, consecutive_cq_data[i].One + consecutive_cq_data[i].Two + consecutive_cq_data[i].Three)
			z3s = append(z3s, consecutive_cq_data[i].Two + consecutive_cq_data[i].Three + consecutive_cq_data[i].Four)
			h3s = append(h3s, consecutive_cq_data[i].Three + consecutive_cq_data[i].Four + consecutive_cq_data[i].Five)
		}
		//刷新该彩种最新 开奖号码
		consecutiveSscNewCodes.Set(CqsscType, new_code)
	}
	return q3s, z3s, h3s
}

/*
//获取天津 前中后的 开奖号码
func getTjCodes() ([]string, []string, []string) {
	q3s := make([]string, 0)
	z3s := make([]string, 0)
	h3s := make([]string, 0)

	//天津检查
	if len := len(consecutive_tj_data); len >0 {
		index := len - 1
		//该彩种到最新的一期 开奖号码
		new_code := consecutive_tj_data[index].One + consecutive_tj_data[index].Two + consecutive_tj_data[index].Three + consecutive_tj_data[index].Four + consecutive_tj_data[index].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := consecutiveSscNewCodes.Get(TjsscType)
		if new_code == newcode {
			log.Println("时时彩-间隔几连号 彩票类型:", CpTypeName[TjsscType], "已经分析过了,等待新的开奖号码出现...")
			return q3s, z3s, h3s
		}

		for i:= range consecutive_tj_data {
			q3s = append(q3s, consecutive_tj_data[i].One + consecutive_tj_data[i].Two + consecutive_tj_data[i].Three)
			z3s = append(z3s, consecutive_tj_data[i].Two + consecutive_tj_data[i].Three + consecutive_tj_data[i].Four)
			h3s = append(h3s, consecutive_tj_data[i].Three + consecutive_tj_data[i].Four + consecutive_tj_data[i].Five)
		}
		//刷新该彩种最新 开奖号码
		consecutiveSscNewCodes.Set(TjsscType, new_code)
	}
	return q3s, z3s, h3s
}
*/

//获取新疆 前中后的 开奖号码
func getXjCodes() ([]string, []string, []string) {
	q3s := make([]string, 0)
	z3s := make([]string, 0)
	h3s := make([]string, 0)

	//新疆检查
	if len := len(consecutive_xj_data); len >0 {
		index := len - 1
		//该彩种到最新的一期 开奖号码
		new_code := consecutive_xj_data[index].One + consecutive_xj_data[index].Two + consecutive_xj_data[index].Three + consecutive_xj_data[index].Four + consecutive_xj_data[index].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := consecutiveSscNewCodes.Get(XjsscType)
		if new_code == newcode {
			log.Println("时时彩-间隔几连号 彩票类型:", CpTypeName[XjsscType], "已经分析过了,等待新的开奖号码出现...")
			return q3s, z3s, h3s
		}

		for i:= range consecutive_xj_data {
			q3s = append(q3s, consecutive_xj_data[i].One + consecutive_xj_data[i].Two + consecutive_xj_data[i].Three)
			z3s = append(z3s, consecutive_xj_data[i].Two + consecutive_xj_data[i].Three + consecutive_xj_data[i].Four)
			h3s = append(h3s, consecutive_xj_data[i].Three + consecutive_xj_data[i].Four + consecutive_xj_data[i].Five)
		}
		//刷新该彩种最新 开奖号码
		consecutiveSscNewCodes.Set(XjsscType, new_code)
	}
	return q3s, z3s, h3s
}

/*
//获取台湾 前中后的 开奖号码
func getTwCodes() ([]string, []string, []string) {
	q3s := make([]string, 0)
	z3s := make([]string, 0)
	h3s := make([]string, 0)

	//台湾检查
	if len := len(consecutive_tw_data); len >0 {
		index := len - 1
		//该彩种到最新的一期 开奖号码
		new_code := consecutive_tw_data[index].One + consecutive_tw_data[index].Two + consecutive_tw_data[index].Three + consecutive_tw_data[index].Four + consecutive_tw_data[index].Five
		//读取该数据吧 所属的 彩种类型的最新开奖号码
		newcode := consecutiveSscNewCodes.Get(TwsscType)
		if new_code == newcode {
			log.Println("时时彩-间隔几连号 彩票类型:", CpTypeName[TwsscType], "已经分析过了,等待新的开奖号码出现...")
			return q3s, z3s, h3s
		}

		for i:= range consecutive_tw_data {
			q3s = append(q3s, consecutive_tw_data[i].One + consecutive_tw_data[i].Two + consecutive_tw_data[i].Three)
			z3s = append(z3s, consecutive_tw_data[i].Two + consecutive_tw_data[i].Three + consecutive_tw_data[i].Four)
			h3s = append(h3s, consecutive_tw_data[i].Three + consecutive_tw_data[i].Four + consecutive_tw_data[i].Five)
		}
		//刷新该彩种最新 开奖号码
		consecutiveSscNewCodes.Set(TwsscType, new_code)
	}
	return q3s, z3s, h3s
}
*/

func consecutiveCodesAnalyse(config *model.Play2, codes []string, position string, cpName string) (string, int) {
	return continuityModo(config, codes, position, cpName)

	// 计算方式是 连续模式
	if config.Type == 1 {
		return continuityModo(config, codes, position, cpName)
	}else {
		// 计算方式是 非连续模式
		return notContinuityModo(config, codes, position, cpName)
	}
	return "", 0
}

// 连续模式
func continuityModo(config *model.Play2, codes []string, position string, cpName string) (string, int) {
	log_html := ""

	// 连续模式 周期计数器
	continuity_modo_number := 0

	// 连续模式 计算状态
	continuity_modo_status := true

	// 满1周期 出现的位置
	cycle_position := 0
	// 满1周期 时的开奖号
	cycle_position_code := ""

	//参考对象
	var reference string = ""
	var number int = 0
	for i := range codes {

		//该号码是否是组六
		isSix := IsSix(codes[i])

		//排序
		code := CodeSort(codes[i], "asc")

		//检查本期号码 是否有连号
		reference_current_obj := isConsecutiveNumber(code)

		//fmt.Println("开奖号:", codes[i], "排序后:", code, position, cpName, "是否是组6", isSix)

		if isSix == false {
			log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 不是组6 [不管] 期数 = " + strconv.Itoa(number)
			log_html += "<div> [连续模式] 周期计数:  "+ strconv.Itoa(continuity_modo_number) + " 周期计算状态: " + strconv.FormatBool(continuity_modo_status) +"</div>"
			//不是组6 跳出本次循环

			//当前轮循完 刷新下一期的 参考对象
			reference = reference_current_obj
			continue
		}

		//检查上一期是否有参考对象
		if reference != "" {
			//上一期有参考对象
			//上一期出现了 连号 检查本期号码 是否包含上期的连号 的 其中一位
			isContain := consecutiveContainNumber(code, reference)

			//上一期出现了 连号 并且 本期号码 包含上一期连号 其中的1位 清零 并且本期也出现了连号 + 1
			if isContain == 1 && reference_current_obj != "" {
				// 周期累计 状态 继续
				continuity_modo_status = true
				number = 0
				number += 1
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期包含上期连号其中1位 并且 本期出现连号 清0 再 +1 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//上一期出现了 连号 并且 本期号码 包含上一期连号 其中的1位 清零 并且本期未出现连号
			if isContain == 1 && reference_current_obj == "" {
				// 周期累计 状态 继续
				continuity_modo_status = true
				number = 0
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期包含上期连号其中1位 并且 本期未出现连号 清0 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//上一期出现了 连号 并且 本期号码 未包含上一期连号 其中的1位 不管 并且本期出现连号 + 1
			if isContain != 1 && reference_current_obj != "" {
				number += 1
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期未包含上期连号其中1位 并且 本期出现连号 +1 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//上一期出现了 连号 并且 本期号码 未包含上一期连号 其中的1位 不管 并且本期未出现连号 不管
			if isContain != 1 && reference_current_obj == "" {
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期未包含上期连号其中1位 并且 本期未出现连号 [不管] 期数 = " + strconv.Itoa(number) + "</div>"
			}
		} else {
			//上一期没有参考对象
			//本期 号码 未出现连号 不管
			if reference_current_obj == "" {
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期未出现连号 并且 本期未包含上期连号其中1位 并且 本期未出现连号 [不管] 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//本期 号码 出现连号
			if reference_current_obj != "" {
				number += 1
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期未出现连号 并且 本期未包含上期连号其中1位 并且 本期出现连号 +1 期数 = " + strconv.Itoa(number) + "</div>"
			}
		}

		// 周期计算
		if config.Number == number {
			// 1周期出现的 位置
			cycle_position = i
			// 1周期时的开奖号
			cycle_position_code = codes[i]
		}

		// 1周期后的一期开奖
		if i == cycle_position + 1 {
			// 本期出现的是组三 或者是 与上一期 同连号 则周期+1
			if continuity_modo_status == true && isSix == false || (codes[i] == cycle_position_code && reference_current_obj != "") {
				// 1周期累计完毕
				continuity_modo_number += 1
				// 不清参考对象
				// 周期累计暂停
				continuity_modo_status = false
			}
		}

		// 累计的周期数 等于设置的 周期数 就清零了 如果是最新到一期到话
		if continuity_modo_number == config.Cycle && (i + 1) != len(codes) {
			continuity_modo_number = 0
		}

		log_html += "<div> [连续模式] 周期计数:  "+ strconv.Itoa(continuity_modo_number) + " 周期计算状态: " + strconv.FormatBool(continuity_modo_status) +"</div>"

		//当前轮循完 刷新下一期的 参考对象
		reference = reference_current_obj
	}
	log_html += "<br/>"

	return log_html, continuity_modo_number
}

// 非连续模式
func notContinuityModo(config *model.Play2, codes []string, position string, cpName string) (string, int) {
	log_html := ""

	// 非连续模式 周期计数器
	not_continuity_modo_number := 0

	// 非连续模式 计算状态
	not_continuity_modo_status := true

	//参考对象
	var reference string = ""
	var number int = 0
	for i := range codes {

		//该号码是否是组六
		isSix := IsSix(codes[i])

		//排序
		code := CodeSort(codes[i], "asc")

		//检查本期号码 是否有连号
		reference_current_obj := isConsecutiveNumber(code)

		//fmt.Println("开奖号:", codes[i], "排序后:", code, position, cpName, "是否是组6", isSix)

		if isSix == false {
			log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 不是组6 [不管] 期数 = " + strconv.Itoa(number)
			log_html += "<div> [非连续模式] 周期计数:  "+ strconv.Itoa(not_continuity_modo_number) + " 周期计算状态: " + strconv.FormatBool(not_continuity_modo_status) +"</div>"
			//不是组6 跳出本次循环

			//当前轮循完 刷新下一期的 参考对象
			reference = reference_current_obj
			continue
		}

		//检查上一期是否有参考对象
		if reference != "" {
			//上一期有参考对象
			//上一期出现了 连号 检查本期号码 是否包含上期的连号 的 其中一位
			isContain := consecutiveContainNumber(code, reference)

			//上一期出现了 连号 并且 本期号码 包含上一期连号 其中的1位 清零 并且本期也出现了连号 + 1
			if isContain == 1 && reference_current_obj != "" {
				// 周期累计 状态 继续
				not_continuity_modo_status = true
				//number = 0
				//number += 1

				// 连续就清零
				number = 0
				// 周期累计 状态 继续
				not_continuity_modo_status = true

				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期包含上期连号其中1位 并且 本期出现连号 清0 再 +1 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//上一期出现了 连号 并且 本期号码 包含上一期连号 其中的1位 清零 并且本期未出现连号
			if isContain == 1 && reference_current_obj == "" {
				// 周期累计 状态 继续
				not_continuity_modo_status = true
				//number = 0

				// 没开就 + 1
				number += 1

				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期包含上期连号其中1位 并且 本期未出现连号 清0 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//上一期出现了 连号 并且 本期号码 未包含上一期连号 其中的1位 不管 并且本期出现连号 + 1
			if isContain != 1 && reference_current_obj != "" {
				//number += 1

				// 连续就清零
				number = 0
				// 周期累计 状态 继续
				not_continuity_modo_status = true

				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期未包含上期连号其中1位 并且 本期出现连号 +1 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//上一期出现了 连号 并且 本期号码 未包含上一期连号 其中的1位 不管 并且本期未出现连号 不管
			if isContain != 1 && reference_current_obj == "" {
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期出现连号 并且 本期未包含上期连号其中1位 并且 本期未出现连号 [不管] 期数 = " + strconv.Itoa(number) + "</div>"
			}
		} else {
			//上一期没有参考对象
			//本期 号码 未出现连号 不管
			if reference_current_obj == "" {
				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期未出现连号 并且 本期未包含上期连号其中1位 并且 本期未出现连号 [不管] 期数 = " + strconv.Itoa(number) + "</div>"
			}

			//本期 号码 出现连号
			if reference_current_obj != "" {
				//number += 1

				// 连续就清零
				number = 0
				// 周期累计 状态 继续
				not_continuity_modo_status = true

				log_html += "<div> 彩种:"+ cpName +" 开奖号: " + codes[i] + " 排序后 " + code + " 位置: " + position + " 上期参考对象: " + reference + " 上期未出现连号 并且 本期未包含上期连号其中1位 并且 本期出现连号 +1 期数 = " + strconv.Itoa(number) + "</div>"
			}
		}

		// 周期计算
		if config.Number == number && not_continuity_modo_status == true {
			// 1周期累计完毕
			not_continuity_modo_number += 1
			// 不清参考对象
			// 周期累计暂停
			not_continuity_modo_status = false

			// 累计的周期数 等于设置的 周期数 就清零了 如果是最新到一期到话
			if not_continuity_modo_number == config.Cycle && (i + 1) != len(codes) {
				not_continuity_modo_number = 0
			}
		}

		log_html += "<div> [非连续模式] 周期计数:  "+ strconv.Itoa(not_continuity_modo_number) + " 周期计算状态: " + strconv.FormatBool(not_continuity_modo_status) +"</div>"

		//当前轮循完 刷新下一期的 参考对象
		reference = reference_current_obj
	}
	log_html += "<br/>"

	return log_html, not_continuity_modo_number
}

//是否是连续的号码 并返回 最小的连号 例如: 123 是连号 返回最新的2个连号 将 返回 12
func isConsecutiveNumber(code string) string {
	by := []byte(code)
	center_bool := false
	tail_bool := false

	first_int, _ := strconv.Atoi(string(by[0]))
	center_int, _ := strconv.Atoi(string(by[1]))
	tail_int, _ := strconv.Atoi(string(by[2]))

	//检查 下标 0 与 下标 1 是否是连号
	if first_int + 1 == center_int {
		center_bool = true
	}

	//检查 下标 1 与 下标 2 是否是连号
	if center_int + 1 == tail_int {
		tail_bool = true
	}

	//下标0 与下标1 是连号 并且 下标1 与下标2 是连号 将返回最小的2个连号
	if center_bool == true && tail_bool == true {
		front, _ := strconv.Atoi(string(by[0]) + string(by[1]))
		after, _ := strconv.Atoi(string(by[1]) + string(by[2]))

		if front < after {
			return string(by[0]) + string(by[1])
		}
		return string(by[1]) + string(by[2])
	}

	//下标0 与下标1 是连号
	if center_bool == true {
		return string(by[0]) + string(by[1])
	}

	//下标1 与下标2 是连号
	if tail_bool == true {
		return string(by[1]) + string(by[2])
	}
	return ""
}

//是否包含其中一位号码
func consecutiveContainNumber(code string, val string) int {
	code_by := []byte(code)
	val_by := []byte(val)

	show_num := 0
	for i := range val_by {
		if bytes.IndexAny(code_by, string(val_by[i])) >= 0 {
			show_num += 1
		}
	}
	return show_num
}

func (c *consecutive_ssc_new_codes) Get(k int) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.codes[k]
}

func (c *consecutive_ssc_new_codes) Set(k int, v string)  {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.codes[k] = v
}