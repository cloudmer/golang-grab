package ssc

import (
	"fmt"
	"xmn/core/model"
	"sync"
	"time"
	"log"
	"xmn/core/logger"
	"strings"
	"errors"
	"strconv"
	"xmn/core/mail"
)

//数据包
var contain_multiple_datapackage []*model.Double

//重庆开奖数据
var contain_multiple_cq_data []*model.Cqssc

//天津开奖数据
var contain_multiple_tj_data []*model.Tjssc

//新疆开奖数据
var contain_multiple_xj_data []*model.Xjssc

//彩票类型
var contain_multiple_ssc_type map[int]string

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
var multipleNewCodes *multipleCode

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
type multipleCode struct {
	codes map[int]string //数据包id => 该彩种的最新开奖号码 每个数据包对应的是一个彩种
	lock sync.RWMutex
}

type multipleData struct {
	packageA   map[string]string
	packageB   map[string]string
	cpType     int
	cpTypeName string
	code       interface{}
	packet     *model.Double
}

func init()  {
	multipleNewCodes = new(multipleCode)
	multipleNewCodes.codes = make(map[int]string)
}

func ContailMultiple()  {

	fmt.Println("时时彩 - 包含数据包 - 自定义多期为一周期 算法")

	contain_multiple_ssc_type = make(map[int]string)
	contain_multiple_ssc_type[1] = "重庆时时彩"
	contain_multiple_ssc_type[2] = "天津时时彩"
	contain_multiple_ssc_type[3] = "新疆时时彩"
	contain_multiple_ssc_type[4] = "台湾五分彩"

	double := new(model.Double)
	contain_multiple_datapackage = double.Query()

	cqssc := new(model.Cqssc)
	contain_multiple_cq_data = cqssc.Query("200")

	tjssc := new(model.Tjssc)
	contain_multiple_tj_data = tjssc.Query("200")

	xjssc := new(model.Xjssc)
	contain_multiple_xj_data = xjssc.Query("200")

	containMultipAnalysis()
}

func containMultipAnalysis()  {
	for i := range contain_multiple_datapackage {
		go containMultipAnalysisCodes(contain_multiple_datapackage[i])
	}
}

func containMultipAnalysisCodes(packet *model.Double)  {
	//检查是否在报警时间段以内
	if (packet.Start >0 && packet.End >0) && (time.Now().Hour() < packet.Start || time.Now().Hour() > packet.End)  {
		log.Println("AB 包含包 自定义周期 - 数据包别名:", packet.Alias, "报警通知非接受时间段内")
		logger.Log("AB 包含包 自定义周期 - 数据包别名: " + packet.Alias + "报警通知非接受时间段内")
		return
	}

	slice_dataTxt_package_a := strings.Split(packet.Package_a, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackageA := make(map[string]string)
	for i := range slice_dataTxt_package_a {
		dataTxtMapPackageA[slice_dataTxt_package_a[i]] = slice_dataTxt_package_a[i]
	}

	slice_dataTxt_package_b := strings.Split(packet.Package_a, "\r\n")
	//slice data txt to slice data txt map
	dataTxtMapPackageB := make(map[string]string)
	for i := range slice_dataTxt_package_b {
		dataTxtMapPackageB[slice_dataTxt_package_b[i]] = slice_dataTxt_package_b[i]
	}

	cq := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: contain_multiple_cq_data,
		cpType: CqsscType,
		cpTypeName: CpTypeName[CqsscType],
		packet: packet,
	}

	tj := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: contain_multiple_tj_data,
		cpType: TjsscType,
		cpTypeName: CpTypeName[TjsscType],
		packet: packet,
	}

	xj := &multipleData{
		packageA: dataTxtMapPackageA,
		packageB: dataTxtMapPackageB,
		code: contain_multiple_xj_data,
		cpType: XjsscType,
		cpTypeName: CpTypeName[XjsscType],
		packet: packet,
	}

	go cq.calculate()
	go tj.calculate()
	go xj.calculate()
}

//开奖号是否最新 是否重复统计计算
func multipCodeIsNewest(cp_tpe int, packet *model.Double) (bool, error) {
	//重庆开奖号
	if cp_tpe == CqsscType {
		//获取最新开奖号码
		model := new(model.Cqssc)
		new_code, err := model.GetNesCode()
		if err != nil {
			//没有最新开奖号码 数据库中没有数据
			return false, errors.New("重庆时时彩, 还没有开奖号码")
		}

		//检查是否重复 计算
		newcode := multipleNewCodes.Get(CqsscType)
		if newcode == new_code {
			return false, errors.New("AB 包 自定义周期 重庆时时彩 数据包别名: " + packet.Alias + " 最新的一期 已经分析过了... 等待出现新的开奖号 ")
		}

		//最新开奖号 与 内存中的最新开奖号 不相同 刷新内存最新开奖号值
		multipleNewCodes.Set(CqsscType, new_code)
		return true, nil
	}

	//天津开奖号
	if cp_tpe == CqsscType {
		//获取最新开奖号码
		model := new(model.Tjssc)
		new_code, err := model.GetNesCode()
		if err != nil {
			//没有最新开奖号码 数据库中没有数据
			return false, errors.New("天津时时彩, 还没有开奖号码")
		}
		//检查是否重复 计算
		newcode := multipleNewCodes.Get(TjsscType)
		if newcode == new_code {
			return false, errors.New("AB 包 自定义周期 天津时时彩 数据包别名: " + packet.Alias + " 最新的一期 已经分析过了... 等待出现新的开奖号 ")
		}

		//最新开奖号 与 内存中的最新开奖号 不相同 刷新内存最新开奖号值
		multipleNewCodes.Set(TjsscType, new_code)
		return true, nil
	}

	//新疆开奖号
	if cp_tpe == CqsscType {
		//获取最新开奖号码
		model := new(model.Xjssc)
		new_code, err := model.GetNesCode()
		if err != nil {
			//没有最新开奖号码 数据库中没有数据
			return false, errors.New("新疆时时彩, 还没有开奖号码")
		}
		//检查是否重复 计算
		newcode := multipleNewCodes.Get(XjsscType)
		if newcode == new_code {
			return false, errors.New("AB 包 自定义周期 新疆时时彩 数据包别名: " + packet.Alias + " 最新的一期 已经分析过了... 等待出现新的开奖号 ")
		}

		//最新开奖号 与 内存中的最新开奖号 不相同 刷新内存最新开奖号值
		multipleNewCodes.Set(XjsscType, new_code)
		return true, nil
	}
	return false, nil
}

func (md *multipleData) calculate()  {
	isOk, err := multipCodeIsNewest(CqsscType, md.packet)
	if isOk == false && err != nil {
		fmt.Println(err)
		return
	}

	type codeStruct struct {
		One   string
		Two   string
		Three string
		Four  string
		Five  string
	}

	codes := make([]*codeStruct, 0)

	if md.cpType == CqsscType {
		for i := range md.code.([]*model.Cqssc) {
			codes = append(codes, &codeStruct{
				One: md.code.([]*model.Cqssc)[i].One,
				Two: md.code.([]*model.Cqssc)[i].Two,
				Three: md.code.([]*model.Cqssc)[i].Three,
				Four: md.code.([]*model.Cqssc)[i].Four,
				Five: md.code.([]*model.Cqssc)[i].Five,
			})
		}
	}
	if md.cpType == TjsscType {
		for i := range md.code.([]*model.Tjssc) {
			codes = append(codes, &codeStruct{
				One: md.code.([]*model.Tjssc)[i].One,
				Two: md.code.([]*model.Tjssc)[i].Two,
				Three: md.code.([]*model.Tjssc)[i].Three,
				Four: md.code.([]*model.Tjssc)[i].Four,
				Five: md.code.([]*model.Tjssc)[i].Five,
			})
		}
	}
	if md.cpType == XjsscType {
		for i := range md.code.([]*model.Xjssc) {
			codes = append(codes, &codeStruct{
				One: md.code.([]*model.Xjssc)[i].One,
				Two: md.code.([]*model.Xjssc)[i].Two,
				Three: md.code.([]*model.Xjssc)[i].Three,
				Four: md.code.([]*model.Xjssc)[i].Four,
				Five: md.code.([]*model.Xjssc)[i].Five,
			})
		}
	}

	//连续中 计数器
	continuity_lucky := 0
	//连续中 周期计数器
	continuity_lucky_cycle_number := 0

	//连续未中 计数器
	continuity_regret := 0
	//连续未中 周期计数器
	continuity_regret_cycle_number := 0

	strHtml := ""

	for i := range codes {
		q3 := codes[i].One + codes[i].Two + codes[i].Three
		z3 := codes[i].Two + codes[i].Three + codes[i].Four
		h3 := codes[i].Three + codes[i].Four + codes[i].Five

		strHtml += " <div>开奖号码: " + codes[i].One + codes[i].Two + codes[i].Three + codes[i].Four + codes[i].Five + "</div>"

		_, q3_in := md.packageB[q3]
		_, z3_in := md.packageB[z3]
		_, h3_in := md.packageB[h3]

		if !q3_in || !z3_in || !h3_in {
			continuity_lucky = 0
			continuity_regret += 1

			strHtml += "<div>"
			if !q3_in {
				strHtml += "<span> 前未包含B </span>"
			}

			if !z3_in {
				strHtml += "<span> 中未包含B </span>"
			}

			if !h3_in {
				strHtml += "<span> 后未包含B </span>"
			}
			strHtml += "</div>"
			strHtml += "<div> 连续未中计数="+ strconv.Itoa(continuity_regret) + " 连续中计数= "+ strconv.Itoa(continuity_lucky) +" </div><br/>"

		} else {
			continuity_regret = 0
			continuity_lucky += 1

			strHtml += "<div>"
			if q3_in {
				strHtml += "<span> 前包含B </span>"
			}

			if z3_in {
				strHtml += "<span> 中包含B </span>"
			}

			if h3_in {
				strHtml += "<span> 后包含B </span>"
			}
			strHtml += "</div>"
			strHtml += "<div> 连续未中计数="+ strconv.Itoa(continuity_regret) + " 连续中计数= "+ strconv.Itoa(continuity_lucky) +" </div><br/>"
		}

		if continuity_regret == md.packet.Cycle {
			continuity_regret = 0
			continuity_regret_cycle_number += 1
			strHtml += "<br/>"
			strHtml += " <div>连续未中已到自定义周期数["+ strconv.Itoa(md.packet.Cycle) +"] 连续未中周期计数器="+ strconv.Itoa(continuity_regret_cycle_number) +" </div><br/>"
		}

		if continuity_lucky == md.packet.Cycle {
			continuity_lucky = 0
			continuity_lucky_cycle_number += 1
			strHtml += "<br/>"
			strHtml += " <div>连续中已到自定义周期数["+ strconv.Itoa(md.packet.Cycle) +"] 连续中周期计数器="+ strconv.Itoa(continuity_lucky_cycle_number) +" </div><br/>"
		}

	}


	//连续中周期 报警
	if continuity_lucky_cycle_number >= md.packet.CycleNumber {
		mailContents := "<div>" + md.cpTypeName +
		    " 数据包别名: " + md.packet.Alias +
			" AB包 自定周期报警通知 当前连续" + strconv.Itoa(md.packet.Cycle) + "期为1周期" +
			" 当前周期报警数 " + strconv.Itoa(md.packet.CycleNumber) +
			" 当前已经连续 " + strconv.Itoa(continuity_lucky_cycle_number) + "周期 未中 </div> <br/></br> "
		mailContents += strHtml
		go mail.SendMail(md.cpTypeName + " AB包 自定义周期 ", mailContents)
	}

	//连续中周期 报警
	if continuity_regret_cycle_number >= md.packet.CycleNumber {
		mailContents := "<div>" + md.cpTypeName +
			" 数据包别名: " + md.packet.Alias +
			" AB包 自定周期报警通知 当前连续" + strconv.Itoa(md.packet.Cycle) + "期为1周期" +
			" 当前周期报警数" + strconv.Itoa(md.packet.CycleNumber) +
			" 当前已经连续 " + strconv.Itoa(continuity_lucky_cycle_number) + "周期 中 </div>"
		mailContents += strHtml
		go mail.SendMail(md.cpTypeName + " AB包 自定义周期 ", mailContents)
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