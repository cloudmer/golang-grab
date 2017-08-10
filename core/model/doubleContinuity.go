package model

import (
	"strconv"
)

//报警状态码 开启
const DC_STATUS_ON int = 1
//报警状态码 关闭
const DC_STATUS_OFF int = 0

type DoubleContinuity struct {
	Id          int
	Alias       string
	Package_a   string
	Package_b   string
	Status      int
	Start       int
	End         int
	Number      int
}

func (model *DoubleContinuity) Query() []*DoubleContinuity {
	//只查询开启了报警的数据
	str_sql := "SELECT * FROM `double_continuity` WHERE status=" + strconv.Itoa(DC_STATUS_ON) + ";"
	rows, err := DB.Query(str_sql)
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	data := make([]*DoubleContinuity, 0)

	for rows.Next() {
		rows.Columns()
		double := new(DoubleContinuity)
		err := rows.Scan(
			&double.Id,
			&double.Alias,
			&double.Package_a,
			&double.Package_b,
			&double.Status,
			&double.Start,
			&double.End,
			&double.Number,
		)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, double)
	}

	return data
}