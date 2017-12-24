package play1

import (
	"xmn/core/model"
	"sync"
)

//重庆开奖数据
var cqCodes []*model.Cqssc

//新疆开奖数据
var xjCodes []*model.Xjssc

//天津开奖数据
var tjCodes []*model.Tjssc

//重庆时时彩 类型
var cqsscType int = 1

//天津时时彩 类型
var tjsscType int = 2

//新疆时时彩 类型
var xjsscType int = 3

//台湾时时彩 类型
var twsscType int = 4

//彩票名集合
var cpTypeName map[int]string = make(map[int]string)

//所有彩票类型的开奖号
type allCpCodes struct {
	cq_q3s []string
	cq_z3s []string
	cq_h3s []string

	tj_q3s []string
	tj_z3s []string
	tj_h3s []string

	xj_q3s []string
	xj_z3s []string
	xj_h3s []string
}

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
var newsCode *newsCodes

//时时彩 各个数据包 对应的 各个彩种的最新开奖号码
type newsCodes struct {
	codes map[int]string //彩票类型 => 该彩种的最新开奖号码 每个数据包对应的是一个彩种
	lock sync.RWMutex
}

func init()  {
	cpTypeName[cqsscType] = "重庆时时彩"
	cpTypeName[tjsscType] = "天津时时彩"
	cpTypeName[xjsscType] = "新疆时时彩"
	cpTypeName[twsscType] = "台湾时时彩"

	newsCode = new(newsCodes)
	newsCode.codes = make(map[int]string)
}

//获取
func (c *newsCodes) Get(k int) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.codes[k]
}

//更新
func (c *newsCodes) Set(k int, v string)  {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.codes[k] = v
}