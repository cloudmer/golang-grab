package ssc

import (
	"strings"
	"strconv"
	"sort"
)

//重庆时时彩 类型
var CqsscType int = 1

//天津时时彩 类型
var TjsscType int = 2

//新疆时时彩 类型
var XjsscType int = 3

//台湾时时彩 类型
var TwsscType int = 4

var CpTypeName map[int]string = make(map[int]string)

func init()  {
	CpTypeName[CqsscType] = "重庆时时彩"
	CpTypeName[TjsscType] = "天津时时彩"
	CpTypeName[XjsscType] = "新疆时时彩"
	CpTypeName[TwsscType] = "台湾时时彩"
}

//是否是组六
func IsSix(str string) bool {
	flag := true
	by := []byte(str)
	for i:=0; i<len(by); i++  {
		i_str := string(by[i])
		if strings.IndexAny(str, i_str) != strings.LastIndex(str, i_str) {
			flag = false
		}
	}
	return flag
}

func CodeSort(code string, mold string) string {
	code_int_slice := make([]int, 0)
	by := []byte(code)
	for i := range by {
		by_i_int, _ := strconv.Atoi(string(by[i]))
		code_int_slice = append(code_int_slice, by_i_int)
	}
	if mold == "desc" {
		sort.Sort(sort.Reverse(sort.IntSlice(code_int_slice)))
	}
	if mold == "asc" {
		sort.Ints(code_int_slice)
	}

	var str string
	for i:= range code_int_slice {
		str += strconv.Itoa(code_int_slice[i])
	}
	return  str
}
