package model

type Xjssc struct {
	Id    int
	Qishu string
	One   string
	Two   string
	Three string
	Four  string
	Five  string
	Time  int
}

func (model *Xjssc) Query(limit string) []*Xjssc {
	sql_str := `SELECT * FROM (
					SELECT id,qishu,one,two,three,four,five,time FROM xjssc  ORDER BY time DESC LIMIT ` + limit + `
				) AS ssc ORDER BY time ASC`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	data := make([]*Xjssc, 0)

	for rows.Next() {
		rows.Columns()
		ssc := new(Xjssc)
		var err error
		err = rows.Scan(&ssc.Id, &ssc.Qishu, &ssc.One, &ssc.Two, &ssc.Three, &ssc.Four, &ssc.Five, &ssc.Time)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, ssc)
	}
	return  data
}

//获取最新一期的开奖号码
func (model *Xjssc) GetNesCode() (string, error) {
	str_sql := `SELECT one,two,three,four,five FROM xjssc ORDER BY time DESC LIMIT 1;`
	xj := new(Xjssc)
	err := DB.QueryRow(str_sql).Scan(&xj.One, &xj.Two, &xj.Three, &xj.Four, &xj.Five)
	if err != nil {
		return "", err
	}
	return xj.One + xj.Two + xj.Three + xj.Four + xj.Five, err
}