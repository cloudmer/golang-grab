package ssc

import (
	"fmt"
	"xmn/core/model"
	"sync"
	"time"
	"log"
	"xmn/core/logger"
	"strings"
	"strconv"
	"xmn/core/mail"
)

//数据包
var contain_multiple_datapackage []*model.DoubleContinuity

//重庆开奖数据
var contain_multiple_cq_data []*model.Cqssc

//新疆开奖数据
var contain_multiple_xj_data []*model.Xjssc

//天津开奖数据
var contain_multiple_tj_data []*model.Tjssc

/*
//台湾开奖数据
var contain_multiple_tw_data []*model.Twssc
*/

//彩票类型
var contain_multiple_ssc_type map[int]string

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
var multipleNewCodes *multipleCode

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
type multipleCode struct {
	codes map[int]string //彩票类型 => 该彩种的最新开奖号码 每个数据包对应的是一个彩种
	lock sync.RWMutex
}

type multipleData struct {
	packageA   map[string]string
	packageB   map[string]string
	cpType     int
	cpTypeName string
	code       []string
	position   string
	packet     *model.DoubleContinuity
}

type cpTypeNewsCodes struct {
	Cq_q3s []string
	Cq_z3s []string
	Cq_h3s []string

	Xj_q3s []string
	Xj_z3s []string
	Xj_h3s []string

	Tj_q3s []string
	Tj_z3s []string
	Tj_h3s []string

	/*
	Tw_q3s []string
	Tw_z3s []string
	Tw_h3s []string
	*/
}

func init()  {
	multipleNewCodes = new(multipleCode)
	multipleNewCodes.codes = make(map[int]string)
}

func ContailMultiple()  {

	fmt.Println("时时彩 - AB - 连续算法报警")

	contain_multiple_ssc_type = make(map[int]string)
	contain_multiple_ssc_type[1] = "重庆时时彩"
	contain_multiple_ssc_type[2] = "天津时时彩"
	contain_multiple_ssc_type[3] = "新疆时时彩"
	contain_multiple_ssc_type[4] = "台湾五分彩"

	double := new(model.DoubleContinuity)
	contain_multiple_datapackage = double.Query()

	cqssc := new(model.Cqssc)
	contain_multiple_cq_data = cqssc.Query("200")

	xjssc := new(model.Xjssc)
	contain_multiple_xj_data = xjssc.Query("200")

	tjssc := new(model.Tjssc)
	contain_multiple_tj_data = tjssc.Query("200")

	/*
	twssc := new(model.Twssc)
	contain_multiple_tw_data = twssc.Query("200")
	*/

	containMultipAnalysis()
}

func containMultipAnalysis()  {

	cq_q3s , cq_z3s , cq_h3s := getSsccodes(CqsscType)
	xj_q3s , xj_z3s , xj_h3s := getSsccodes(XjsscType)
	tj_q3s , tj_z3s , tj_h3s := getSsccodes(TjsscType)
	//tw_q3s , tw_z3s , tw_h3s := getSsccodes(TwsscType)

	cpTypenewsCodes := &cpTypeNewsCodes{
		Cq_q3s: cq_q3s,
		Cq_z3s: cq_z3s,
		Cq_h3s: cq_h3s,

		Xj_q3s: xj_q3s,
		Xj_z3s: xj_z3s,
		Xj_h3s: xj_h3s,

		Tj_q3s: tj_q3s,
		Tj_z3s: tj_z3s,
		Tj_h3s: tj_h3s,

		/*
		Tw_q3s: tw_q3s,
		Tw_z3s: tw_z3s,
		Tw_h3s: tw_h3s,
		*/
	}

	for i := range contain_multiple_datapackage {
		go containMultipAnalysisCodes(contain_multiple_datapackage[i], cpTypenewsCodes)
	}
}

func containMultipAnalysisCodes(packet *model.DoubleContinuity, cpTypenewsCodes *cpTypeNewsCodes)  {
	//检查是否在报警时间段以内
	if (packet.Start >0 && packet.End >0) && (time.Now().Hour() < packet.Start || time.Now().Hour() > packet.End)  {
		log.Println("AB 包含包 自定义周期 - 数据包别名:", packet.Alias, "报警通知非接受时间段内")
		logger.Log("AB 包含包 自定义周期 - 数据包别名: " + packet.Alias + "报警通知非接受时间段内")
		return
	}

	//检查会否属于重复分析


	slice_dataTxt_package_a := strings.Split(packet.Package_a, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackageA := make(map[string]string)
	for i := range slice_dataTxt_package_a {
		dataTxtMapPackageA[slice_dataTxt_package_a[i]] = slice_dataTxt_package_a[i]
	}

	slice_dataTxt_package_b := strings.Split(packet.Package_b, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackageB := make(map[string]string)
	for i := range slice_dataTxt_package_b {
		dataTxtMapPackageB[slice_dataTxt_package_b[i]] = slice_dataTxt_package_b[i]
	}

	/*
	cq_q3s , cq_z3s , cq_h3s := getSsccodes(CqsscType)
	tj_q3s , tj_z3s , tj_h3s := getSsccodes(TjsscType)
	xj_q3s , xj_z3s , xj_h3s := getSsccodes(XjsscType)
	tw_q3s , tw_z3s , tw_h3s := getSsccodes(TwsscType)
	*/

	cq_q3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Cq_q3s,
		cpType: CqsscType,
		cpTypeName: CpTypeName[CqsscType],
		position: "前3",
		packet: packet,
	}

	cq_z3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Cq_z3s,
		cpType: CqsscType,
		cpTypeName: CpTypeName[CqsscType],
		position: "中3",
		packet: packet,
	}

	cq_h3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Cq_h3s,
		cpType: CqsscType,
		cpTypeName: CpTypeName[CqsscType],
		position: "后3",
		packet: packet,
	}

	xj_q3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Xj_q3s,
		cpType: XjsscType,
		cpTypeName: CpTypeName[XjsscType],
		position: "前3",
		packet: packet,
	}

	xj_z3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Xj_z3s,
		cpType: XjsscType,
		cpTypeName: CpTypeName[XjsscType],
		position: "中3",
		packet: packet,
	}

	xj_h3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Xj_h3s,
		cpType: XjsscType,
		cpTypeName: CpTypeName[XjsscType],
		position: "后3",
		packet: packet,
	}

	tj_q3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Tj_q3s,
		cpType: TjsscType,
		cpTypeName: CpTypeName[TjsscType],
		position: "前3",
		packet: packet,
	}

	tj_z3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Tj_z3s,
		cpType: TjsscType,
		cpTypeName: CpTypeName[TjsscType],
		position: "中3",
		packet: packet,
	}

	tj_h3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Tj_h3s,
		cpType: TjsscType,
		cpTypeName: CpTypeName[TjsscType],
		position: "后3",
		packet: packet,
	}

	/*
	tw_q3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Tw_q3s,
		cpType: TwsscType,
		cpTypeName: CpTypeName[TwsscType],
		position: "前3",
		packet: packet,
	}

	tw_z3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Tw_z3s,
		cpType: TwsscType,
		cpTypeName: CpTypeName[TwsscType],
		position: "中3",
		packet: packet,
	}

	tw_h3 := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: cpTypenewsCodes.Tw_h3s,
		cpType: TwsscType,
		cpTypeName: CpTypeName[TwsscType],
		position: "后3",
		packet: packet,
	}
	*/

	go cq_q3.calculate()
	go cq_z3.calculate()
	go cq_h3.calculate()

	go xj_q3.calculate()
	go xj_z3.calculate()
	go xj_h3.calculate()

	go tj_q3.calculate()
	go tj_z3.calculate()
	go tj_h3.calculate()

	/*
	go tw_q3.calculate()
	go tw_z3.calculate()
	go tw_h3.calculate()
	*/
}

//获取时时彩 前中后 开奖号码
func getSsccodes(cyType int) ([]string, []string, []string)  {
	q3codes := make([]string, 0)
	z3codes := make([]string, 0)
	h3codes := make([]string, 0)

	isRepeat := isRepeat(cyType)
	if !isRepeat {
		//fmt.Println(CpTypeName[cyType], "等待出现最新的号码")
		return q3codes, z3codes, h3codes
	}

	//重庆时时彩
	if cyType == CqsscType {
		for i := range contain_multiple_cq_data {
			q3s := contain_multiple_cq_data[i].One + contain_multiple_cq_data[i].Two + contain_multiple_cq_data[i].Three
			z3s := contain_multiple_cq_data[i].Two + contain_multiple_cq_data[i].Three + contain_multiple_cq_data[i].Four
			h3s := contain_multiple_cq_data[i].Three + contain_multiple_cq_data[i].Four + contain_multiple_cq_data[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}

	//新疆时时彩
	if cyType == XjsscType {
		for i := range contain_multiple_xj_data {
			q3s := contain_multiple_xj_data[i].One + contain_multiple_xj_data[i].Two + contain_multiple_xj_data[i].Three
			z3s := contain_multiple_xj_data[i].Two + contain_multiple_xj_data[i].Three + contain_multiple_xj_data[i].Four
			h3s := contain_multiple_xj_data[i].Three + contain_multiple_xj_data[i].Four + contain_multiple_xj_data[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}

	//天津时时彩
	if cyType == TjsscType {
		for i := range contain_multiple_tj_data {
			q3s := contain_multiple_tj_data[i].One + contain_multiple_tj_data[i].Two + contain_multiple_tj_data[i].Three
			z3s := contain_multiple_tj_data[i].Two + contain_multiple_tj_data[i].Three + contain_multiple_tj_data[i].Four
			h3s := contain_multiple_tj_data[i].Three + contain_multiple_tj_data[i].Four + contain_multiple_tj_data[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}

	/*
	//台湾时时彩
	if cyType == TwsscType {
		for i := range contain_multiple_tw_data {
			q3s := contain_multiple_tw_data[i].One + contain_multiple_tw_data[i].Two + contain_multiple_tw_data[i].Three
			z3s := contain_multiple_tw_data[i].Two + contain_multiple_tw_data[i].Three + contain_multiple_tw_data[i].Four
			h3s := contain_multiple_tw_data[i].Three + contain_multiple_tw_data[i].Four + contain_multiple_tw_data[i].Five
			q3codes = append(q3codes, q3s)
			z3codes = append(z3codes, z3s)
			h3codes = append(h3codes, h3s)
		}
	}
	*/

	return q3codes, z3codes, h3codes
}

//是否属于重复分析
func isRepeat(cyType int) bool {
	//重庆
	if cyType == CqsscType {
		//获取本次查询的最新号码
		if len(contain_multiple_cq_data) == 0 {
			return false
		}
		index := len(contain_multiple_cq_data) - 1
		newCode := contain_multiple_cq_data[index].One + contain_multiple_cq_data[index].Two + contain_multiple_cq_data[index].Three + contain_multiple_cq_data[index].Four + contain_multiple_cq_data[index].Five
		//获取内存中最新的重新开奖号码
		new_code := multipleNewCodes.Get(cyType)
		if new_code == newCode {
			return false
		}
		//刷新最新开奖号码
		multipleNewCodes.Set(cyType, newCode)
		return true
	}

	//新疆
	if cyType == XjsscType {
		//获取本次查询的最新号码
		if len(contain_multiple_xj_data) == 0 {
			return false
		}
		index := len(contain_multiple_xj_data) - 1
		newCode := contain_multiple_xj_data[index].One + contain_multiple_xj_data[index].Two + contain_multiple_xj_data[index].Three + contain_multiple_xj_data[index].Four + contain_multiple_xj_data[index].Five
		//获取内存中最新的重新开奖号码
		new_code := multipleNewCodes.Get(cyType)
		if new_code == newCode {
			return false
		}
		//刷新最新开奖号码
		multipleNewCodes.Set(cyType, newCode)
		return true
	}

	//天津
	if cyType == TjsscType {
		//获取本次查询的最新号码
		if len(contain_multiple_tj_data) == 0 {
			return false
		}
		index := len(contain_multiple_tj_data) - 1
		newCode := contain_multiple_tj_data[index].One + contain_multiple_tj_data[index].Two + contain_multiple_tj_data[index].Three + contain_multiple_tj_data[index].Four + contain_multiple_tj_data[index].Five
		//获取内存中最新的重新开奖号码
		new_code := multipleNewCodes.Get(cyType)
		if new_code == newCode {
			return false
		}
		//刷新最新开奖号码
		multipleNewCodes.Set(cyType, newCode)
		return true
	}

	/*
	//台湾
	if cyType == TwsscType {
		//获取本次查询的最新号码
		if len(contain_multiple_tw_data) == 0 {
			return false
		}
		index := len(contain_multiple_tw_data) - 1
		newCode := contain_multiple_tw_data[index].One + contain_multiple_tw_data[index].Two + contain_multiple_tw_data[index].Three + contain_multiple_tw_data[index].Four + contain_multiple_tw_data[index].Five
		//获取内存中最新的重新开奖号码
		new_code := multipleNewCodes.Get(cyType)
		if new_code == newCode {
			return false
		}
		//刷新最新开奖号码
		multipleNewCodes.Set(cyType, newCode)
		return true
	}
	*/

	return false
}

func (md *multipleData) calculate() {

	if len(md.code) == 0 {
		return
	}

	//连续包含A包内 初始化=0
	continuity_a_num := 0
	number := 0
	strLogHtml := ""
	status := false
	for i := range md.code {
		//检查上一期是否 包含A包
		pre_code := ""
		if i != 0 {
			pre_code = md.code[i - 1]
		}

		var pre_in_a bool = false
		//上一期是否包含A包
		if pre_code != "" {
			_, pre_in_a = md.packageA[pre_code]
		}

		_, in_a := md.packageA[md.code[i]]
		_, in_b := md.packageB[md.code[i]]

		strLogHtml += "<div> 开奖号: "+ md.code[i] +"</div>"

		//没有连续出现A 并且 本次出现A包
		if continuity_a_num == 0 && in_a {
			status = true
			//报警计数 +1
			number += 1
			//A 连续出现 +1
			continuity_a_num += 1
			strLogHtml += "<div>本期包含A包 报警+1 = "+ strconv.Itoa(number) +"</div>"
			continue
		}

		//A包连续出现 并且 本次出现A包
		if continuity_a_num >0 && in_a {
			status = false
			//A 连续出现 +1
			continuity_a_num += 1
			strLogHtml += "<div>本期包含A包 [不管]</div>"
			continue
		}

		//A B 包都没出现
		if !in_a && !in_b {
			status = false
			//A 连续出现清零
			continuity_a_num = 0
			strLogHtml += "<div>本期AB都不包含 [不管]</div>"
			continue
		}

		//B包出现
		if in_b {
			// a a a a b 属于这种情况 必须 > 1
			//A 包连续出现
			if continuity_a_num > 1 {
				status = false
				//A 连续出现清零
				continuity_a_num = 0
				strLogHtml += "<div>包含B包,之前连续开A包 [不管]</div>"
				continue
			}

			//A包未连续出现 并且 上一期包含A包
			if continuity_a_num <= 1 && pre_in_a {
				status = false
				//报警计数 清零
				number = 0
				//A 连续出现清零
				continuity_a_num = 0
				strLogHtml += "<div>包含B包 上一期包含A包 报警清零=0</div>"
				continue
			}

			//A包未连续出现 并且 上一期未包含A包
			if continuity_a_num <= 1 && !pre_in_a {
				status = false
				//A 连续出现清零
				continuity_a_num = 0
				strLogHtml += "<div>包含B包 不管 上一期不包含A包</div>"
				continue
			}

		}
	}

	if status && number >= md.packet.Number {
		emialBody := "<div>AB连续 - 彩种: "+ md.cpTypeName + " 位置: " + md.position + "数据包别名: "+ md.packet.Alias + " 报警期数: " + strconv.Itoa(number) + "</div><br/><br/>"
		//emialBody += strLogHtml
		go mail.SendMail(md.cpTypeName + "ab连续报警", emialBody)
	}

}

func (c *multipleCode) Get(k int) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.codes[k]
}

func (c *multipleCode) Set(k int, v string)  {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.codes[k] = v
}