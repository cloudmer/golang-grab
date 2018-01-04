package play1

import (
	"fmt"
	"xmn/core/model"
	"time"
	"log"
	"strings"
	"xmn/core/logger"
	"strconv"
	"xmn/core/mail"
)

//计算分析结构体
type computing struct {
	packet_a_map   map[string]string
	packet_b_map   map[string]string
	cpType     int
	cpTypeName string
	code       []string
	position   string
	packet     *model.Play1
}

//开始计算
func Calculation()  {
	fmt.Println("时时彩 - a出现几期的b")

	//获取开奖号
	cqssc := new(model.Cqssc)
	cqCodes = cqssc.Query("200")

	xjscc := new(model.Xjssc)
	xjCodes = xjscc.Query("200")

	//获取数据包
	cPackage := new(model.Play1)
	configPackage := cPackage.Query()

	cq_q3s, cq_z3s, cq_h3s := getFrontCenterAfterCodes(cqsscType)
	xj_q3s, xj_z3s, xj_h3s := getFrontCenterAfterCodes(xjsscType)

	allCodes := &allCpCodes{
		cq_q3s: cq_q3s,
		cq_z3s: cq_z3s,
		cq_h3s: cq_h3s,

		xj_q3s: xj_q3s,
		xj_z3s: xj_z3s,
		xj_h3s: xj_h3s,
	}

	for i := range configPackage {
		go analysis(configPackage[i], allCodes)
	}
}

//解析数据包
func analysis(packet *model.Play1, allCodes *allCpCodes)  {
	//检查是否在报警时间段以内
	if (packet.Start >0 && packet.End >0) && (time.Now().Hour() < packet.Start || time.Now().Hour() > packet.End)  {
		log.Println("a出现几期的b - 数据包别名:", packet.Alias, "报警通知非接受时间段内")
		logger.Log("a出现几期的b - 数据包别名: " + packet.Alias + "报警通知非接受时间段内")
		return
	}

	//数据包 a包 解析成map
	slice_dataTxt_package_a := strings.Split(packet.Package_a, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackageA := make(map[string]string)
	for i := range slice_dataTxt_package_a {
		dataTxtMapPackageA[slice_dataTxt_package_a[i]] = slice_dataTxt_package_a[i]
	}

	//数据包 b包 解析成map
	slice_dataTxt_package_b := strings.Split(packet.Package_b, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackageB := make(map[string]string)
	for i := range slice_dataTxt_package_b {
		dataTxtMapPackageB[slice_dataTxt_package_b[i]] = slice_dataTxt_package_b[i]
	}

	//重庆前3
	cq_q3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		packet_b_map: dataTxtMapPackageB,
		code: allCodes.cq_q3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "前3",
		packet: packet,
	}

	//重庆中3
	cq_z3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		packet_b_map: dataTxtMapPackageB,
		code: allCodes.cq_z3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "中3",
		packet: packet,
	}

	//重庆后3
	cq_h3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		packet_b_map: dataTxtMapPackageB,
		code: allCodes.cq_h3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "后3",
		packet: packet,
	}

	//新疆前3
	xj_q3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		packet_b_map: dataTxtMapPackageB,
		code: allCodes.xj_q3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "前3",
		packet: packet,
	}

	//新疆中3
	xj_z3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		packet_b_map: dataTxtMapPackageB,
		code: allCodes.xj_z3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "中3",
		packet: packet,
	}

	//新疆后3
	xj_h3 := &computing{
		packet_a_map: dataTxtMapPackageA,
		packet_b_map: dataTxtMapPackageB,
		code: allCodes.xj_h3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "后3",
		packet: packet,
	}

	go cq_q3.calculate()
	go cq_z3.calculate()
	go cq_h3.calculate()

	go xj_q3.calculate()
	go xj_z3.calculate()
	go xj_h3.calculate()
}

//计算分析
func (c *computing) calculate()  {
	if len(c.code) == 0 {
		//fmt.Println("a出现几期的b", c.packet.Alias, c.cpTypeName , c.position , "等待出现最新的号码")
		return
	}
	
	// 循环一次后 是否出现a 包
	pre_a_show := false
	// 循环一次后 是否出现b 包
	pre_b_show := false
	// 循环一次后 是否出现c 包
	pre_c_show := false

	// 连续出现自定义几期b 计数器
	continuity_b_number := 0

	// 报警期数 计数器
	number := 0

	var strHtmlLog string = ""

	for i := range c.code {
		// 是否在 a 包
		_, in_a_package := c.packet_a_map[c.code[i]]
		// 是否在 b 包
		_, in_b_package := c.packet_b_map[c.code[i]]
		// 是否在 c 包 【不在 ab 内就算c包】
		in_c_package := false
		if !in_a_package && !in_b_package {
			in_c_package = true
		}

		// 上一期 出现a包 并且 本期 出现了b包 或 本期出现了c包
		if pre_a_show && (in_b_package || in_c_package) {
			continuity_b_number += 1
		}

		// 上一期出现c包 或者 上一期出现b 包 并且 本期出现c包 或者 本期出现 b包
		if (pre_c_show || pre_b_show) && (in_c_package || in_b_package) {
			continuity_b_number += 1
		}

		// 累加连续b 已经达到 了自定义几b 值 报警期数 + 1
		if continuity_b_number == c.packet.ContinuityNumber {
			number += 1
		}

		// 设置自定义连续几B, 开奖超出自定义连续几B 就算开奖 清零 并且本期 未出现c包 只有b包才能清零
		if continuity_b_number > c.packet.ContinuityNumber && in_b_package {
			// 连续b 清零
			continuity_b_number = 0
			// 报警期数 清零
			number = 0
		}

		strHtmlLog += "<div>开奖号: </div>" + c.code[i] + "<br/>"
		strHtmlLog += "<div>上一期包含A包吗: " + strconv.FormatBool(pre_a_show) + "</div>"
		strHtmlLog += "<div>上一期包含B包吗: " + strconv.FormatBool(pre_b_show) + "</div>"
		strHtmlLog += "<div>上一期包含C包吗: " + strconv.FormatBool(pre_c_show) + "</div>"
		strHtmlLog += "<div>---<br/>"
		strHtmlLog += "<div>本期包含A包吗: " + strconv.FormatBool(in_a_package) + "</div>"
		strHtmlLog += "<div>本期包含B包吗: " + strconv.FormatBool(in_b_package) + "</div>"
		strHtmlLog += "<div>本期包含C包吗: " + strconv.FormatBool(in_c_package) + "</div>"
		strHtmlLog += "<div>设置的连续 " + strconv.Itoa(c.packet.ContinuityNumber) + " B 已经累加到 " + strconv.Itoa(continuity_b_number) + "</div>"
		strHtmlLog += "<div>报警期数累加到 " + strconv.Itoa(number) + "</div><br/>"

		// 本次循环执行完后 再更新 本期是否出现了a 包
		if in_a_package {
			pre_a_show = true
		}

		// 本次循环执行完后 再更新 本期是否出现了b 包
		if in_b_package {
			pre_b_show = true
		}

		// 本次循环执行完后 再更新 本期是否出现了c 包
		if in_c_package {
			pre_c_show = true
		}
	}

	go mail.SendMail(c.cpTypeName + "a出现几期的b", strHtmlLog)


	// 最新的一期 包含 c 包 或者包含 b包 并且 达到报警期数
	if (pre_b_show || pre_c_show) && number >= c.packet.Number {
		emailTitle := "<div>a出现几期的b 设置的连续"+ strconv.Itoa(c.packet.ContinuityNumber) +"B 报警"+ " 彩种: "+ c.cpTypeName + " 位置: "+ c.position + " 包别名: "+ c.packet.Alias +" 报警 ["+ strconv.Itoa(number) +"]期 提示</div> <br/><br/>"
		emailTitle += strHtmlLog
		go mail.SendMail(c.cpTypeName + "a出现几期的b", emailTitle)
	}

}

//获取 前中后3 开奖号
func getFrontCenterAfterCodes(cpType int) ([]string, []string, []string) {
	q3codes := make([]string, 0)
	z3codes := make([]string, 0)
	h3codes := make([]string, 0)

	//是否属于重复分析
	isRepeat := isRepeat(cpType)
	if !isRepeat {
		//fmt.Println(cpTypeName[cpType], "等待出现最新的号码")
		return q3codes, z3codes, h3codes
	}

	//重庆时时彩
	if cpType == cqsscType {
		for i := range cqCodes {
			q3s := cqCodes[i].One + cqCodes[i].Two + cqCodes[i].Three
			z3s := cqCodes[i].Two + cqCodes[i].Three + cqCodes[i].Four
			h3s := cqCodes[i].Three + cqCodes[i].Four + cqCodes[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}

	//新疆时时彩
	if cpType == xjsscType {
		for i:= range xjCodes {
			q3s := xjCodes[i].One + xjCodes[i].Two + xjCodes[i].Three
			z3s := xjCodes[i].Two + xjCodes[i].Three + xjCodes[i].Four
			h3s := xjCodes[i].Three + xjCodes[i].Four + xjCodes[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}

	return q3codes, z3codes, h3codes
}

//是否属于重复分析
func isRepeat(cpType int) bool {

	//数据库最新开奖号
	var newCode string

	//内存中最新开奖号
	var new_code string

	//重庆时时彩
	if cpType == cqsscType {
		//获取本次查询的最新号码
		if len(cqCodes) == 0 {
			return false
		}
		index := len(cqCodes) - 1
		newCode = cqCodes[index].One + cqCodes[index].Two + cqCodes[index].Three + cqCodes[index].Four + cqCodes[index].Five
	}

	// 新疆时时彩
	if cpType == xjsscType {
		//获取本次查询的最新号码
		if len(xjCodes) == 0 {
			return false
		}
		index := len(xjCodes) - 1
		newCode = xjCodes[index].One + xjCodes[index].Two + xjCodes[index].Three + xjCodes[index].Four + xjCodes[index].Five
	}

	//获取内存中最新的重新开奖号码
	new_code = newsCode.Get(cpType)
	if new_code == newCode {
		return false
	}
	//刷新最新开奖号码
	newsCode.Set(cpType, newCode)
	return true
}