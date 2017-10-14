package CustomPackage

import (
	"xmn/core/model"
	"time"
	"log"
	"xmn/core/logger"
	"strings"
	"strconv"
	"xmn/core/mail"
	"fmt"
)

//时时彩 自定包 算法
//A包(包含)连续N期(自定义) 为一个周期, N个周期(报警期数) - 算法

//计算分析结构体
type computing struct {
	packet_map   map[string]string
	cpType     int
	cpTypeName string
	code       []string
	position   string
	packet     *model.CustomPackage
}

//开始计算
func Calculation()  {
	fmt.Println("时时彩 - AB 自定义A包 - 周期报警算法")

	//获取开奖号
	cqssc := new(model.Cqssc)
	cqCodes = cqssc.Query("200")

	tjssc := new(model.Tjssc)
	tjCodes = tjssc.Query("200")

	xjscc := new(model.Xjssc)
	xjCodes = xjscc.Query("200")

	//获取数据包
	cPackage := new(model.CustomPackage)
	configPackage := cPackage.Query()

	//先获取 前中后的开奖号码 再遍历 就可以判断是否属于重复分析
	cq_q3s, cq_z3s, cq_h3s := getFrontCenterAfterCodes(cqsscType)
	tj_q3s, tj_z3s, tj_h3s := getFrontCenterAfterCodes(tjsscType)
	xj_q3s, xj_z3s, xj_h3s := getFrontCenterAfterCodes(xjsscType)

	allCodes := &allCpCodes{
		cq_q3s: cq_q3s,
		cq_z3s: cq_z3s,
		cq_h3s: cq_h3s,

		tj_q3s: tj_q3s,
		tj_z3s: tj_z3s,
		tj_h3s: tj_h3s,

		xj_q3s: xj_q3s,
		xj_z3s: xj_z3s,
		xj_h3s: xj_h3s,
	}

	for i := range configPackage {
		go analysis(configPackage[i], allCodes)
	}

}

//解析数据包
func analysis(packet *model.CustomPackage, allCodes *allCpCodes)  {
	//检查是否在报警时间段以内
	if (packet.Start >0 && packet.End >0) && (time.Now().Hour() < packet.Start || time.Now().Hour() > packet.End)  {
		log.Println("AB 包含包 自定义周期 - 数据包别名:", packet.Alias, "报警通知非接受时间段内")
		logger.Log("AB 包含包 自定义周期 - 数据包别名: " + packet.Alias + "报警通知非接受时间段内")
		return
	}

	//数据包包解析成map
	slice_dataTxt_package := strings.Split(packet.Package, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackage := make(map[string]string)
	for i := range slice_dataTxt_package {
		dataTxtMapPackage[slice_dataTxt_package[i]] = slice_dataTxt_package[i]
	}

	//重庆前3
	cq_q3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.cq_q3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "前3",
		packet: packet,
	}

	//重庆中3
	cq_z3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.cq_z3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "中3",
		packet: packet,
	}

	//重庆后3
	cq_h3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.cq_h3s,
		cpType: cqsscType,
		cpTypeName: cpTypeName[cqsscType],
		position: "后3",
		packet: packet,
	}

	//天津前3
	tj_q3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.tj_q3s,
		cpType: tjsscType,
		cpTypeName: cpTypeName[tjsscType],
		position: "前3",
		packet: packet,
	}

	//天津中3
	tj_z3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.tj_z3s,
		cpType: tjsscType,
		cpTypeName: cpTypeName[tjsscType],
		position: "中3",
		packet: packet,
	}

	//天津后3
	tj_h3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.tj_h3s,
		cpType: tjsscType,
		cpTypeName: cpTypeName[tjsscType],
		position: "后3",
		packet: packet,
	}

	//新疆前3
	xj_q3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.xj_q3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "前3",
		packet: packet,
	}

	//新疆中3
	xj_z3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.xj_z3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "中3",
		packet: packet,
	}

	//新疆后3
	xj_h3 := &computing{
		packet_map: dataTxtMapPackage,
		code: allCodes.xj_h3s,
		cpType: xjsscType,
		cpTypeName: cpTypeName[xjsscType],
		position: "后3",
		packet: packet,
	}

	go cq_q3.calculate()
	go cq_z3.calculate()
	go cq_h3.calculate()

	go tj_q3.calculate()
	go tj_z3.calculate()
	go tj_h3.calculate()

	go xj_q3.calculate()
	go xj_z3.calculate()
	go xj_h3.calculate()
}

//计算分析
func (c *computing) calculate()  {
	/*
	c.code = make([]string, 0)
	c.code = append(c.code, "001")
	c.code = append(c.code, "001")
	c.code = append(c.code, "001")
	c.code = append(c.code, "001")
	c.code = append(c.code, "001")
	c.code = append(c.code, "012")
	c.code = append(c.code, "001")
	*/

	//连续包含A包 计数 连续包含A包 大于 自定义周期数 就算开 要清零 重新计算了 , 当值 等于 自定义周期的时候 周期计数就要累加(报警期数)
	continuity_num := 0

	//A包连续计算 状态 默认打开
	//清零后 要等B包出现后 再计算几A连续
	//列子
	//设置3A
	//a
	//a
	//a  +1
	//a  清零
	//a  [ 暂停连续A包计算 不管 ]
	//b  [ b包出现 暂停状态 设置为 允许计算 ]
	//a  [ 1a ]
	continuity_num_status := true

	//周期计数 也是报警期数计数
	cycle_count := 0

	strHtmlContents := "<div>自定义连续 "+ strconv.Itoa(c.packet.Continuity) +"A 报警提示</div>"
	strHtmlContents += "<div>彩票类型: "+ c.cpTypeName+ " 数据包别名: "+ c.packet.Alias + " 计算位置: "+ c.position +"</div>"
	strHtmlContents += "<div>当前设置 "+ strconv.Itoa(c.packet.Continuity)+ " A 阀值, 报警期数设置为: "+ strconv.Itoa(c.packet.Number) +"期</div><br/><br/>"

	for i := range c.code {
		_, in_package := c.packet_map[c.code[i]]
		//不考虑AB包是否有重复的值

		strHtmlContents += "<div>开奖号: "+ c.code[i] + "</div>"

		//不包含A包的意思是 包含B包 A包连续中断 重新计数
		if !in_package {
			//A包连续清零
			continuity_num = 0
			//A包连续状态 设置 允许计算连续值
			continuity_num_status = true
			strHtmlContents += "<div>包含_B包, A包连续值累加 开启, A包连续: "+ strconv.Itoa(continuity_num)+ " 报警累计: "+ strconv.Itoa(cycle_count) +"</div><br/><br/>"
			continue
		}

		//暂停连续A包计算
		if !continuity_num_status {
			strHtmlContents += "<div>等待B包出现, 才计算A包连续</div><br/><br/>"
			continue
		}

		//包含A包 A包连续累加
		if in_package {
			continuity_num += 1
			strHtmlContents += "<div>包含_A包, A包连续 +1,  A包连续: "+ strconv.Itoa(continuity_num)+ " 报警累计: "+ strconv.Itoa(cycle_count) + "</div>"
		}

		//当A包连续 到达 自定义A包连续阀值, 报警期数 累加
		if continuity_num == c.packet.Continuity {
			//A包连续 累加
			cycle_count += 1
			strHtmlContents += "<div>A包连续值 = 等于 A包阀值, 报警期数 +1,  A包连续: "+ strconv.Itoa(continuity_num)+ " 报警累计: "+ strconv.Itoa(cycle_count) + "</div>"
		}

		//当A包连续 大于 自定义A包连续阀值, 清零
		if continuity_num > c.packet.Continuity {
			//清零 A包连续值
			continuity_num = 0
			//暂停 A包连续
			continuity_num_status = false
			//清零 报警期数
			cycle_count = 0
			strHtmlContents += "<div>A包连续值 > 大于 A包阀值, A包连续值 清零, A包连续值累加 暂停, 报警期数 清零,  A包连续: "+ strconv.Itoa(continuity_num)+ " 报警累计: "+ strconv.Itoa(cycle_count) + "</div>"
		}

		strHtmlContents += "<br/><br/>"
	}

	//最新的一期 是否包含A包 最新的一期包含A包 才报警
	last_in_a := false
	if len(c.code) > 0 {
		last_code := c.code[len(c.code) - 1]
		_, in_a := c.packet_map[last_code]
		if in_a {
			last_in_a = true
		}
	}

	//到达报警条件
	if last_in_a && cycle_count >= c.packet.Number {
		emailTitle := "<div>自定义连续 "+ strconv.Itoa(c.packet.Continuity) +"A 报警"+ " 彩种: "+ c.cpTypeName + " 位置: "+ c.position + " 包别名: "+ c.packet.Alias +" 报警 ["+ strconv.Itoa(cycle_count) +"]期 提示</div> <br/><br/>"
		go mail.SendMail(c.cpTypeName + "AB包自定义报警", emailTitle + strHtmlContents)
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
		//fmt.Println(CpTypeName[cyType], "等待出现最新的号码")
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

	//天津时时彩
	if cpType == tjsscType {
		for i := range tjCodes {
			q3s := tjCodes[i].One + tjCodes[i].Two + tjCodes[i].Three
			z3s := tjCodes[i].Two + tjCodes[i].Three + tjCodes[i].Four
			h3s := tjCodes[i].Three + tjCodes[i].Four + tjCodes[i].Five
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

	//天津时时彩
	if cpType == tjsscType {
		//获取本次查询的最新号码
		if len(tjCodes) == 0 {
			return false
		}
		index := len(tjCodes) - 1
		newCode = tjCodes[index].One + tjCodes[index].Two + tjCodes[index].Three + tjCodes[index].Four + tjCodes[index].Five
	}

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