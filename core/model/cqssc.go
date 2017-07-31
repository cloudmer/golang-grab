package model

type Cqssc struct {
	Id    int
	Qishu string
	One   string
	Two   string
	Three string
	Four  string
	Five  string
	Time  int
}

func (model *Cqssc) Query(limit string) []*Cqssc {
	sql_str := `SELECT * FROM (
					SELECT id,qishu,one,two,three,four,five,time FROM cqssc  ORDER BY time DESC LIMIT ` + limit + `
				) AS ssc ORDER BY time ASC`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	data := make([]*Cqssc, 0)

	for rows.Next() {
		rows.Columns()
		cq := new(Cqssc)
		var err error
		err = rows.Scan(&cq.Id, &cq.Qishu, &cq.One, &cq.Two, &cq.Three, &cq.Four, &cq.Five, &cq.Time)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, cq)
	}
	return  data
}

//获取最新一期的开奖号码
func (model *Cqssc) GetNesCode() (string, error) {
	str_sql := `SELECT one,two,three,four,five FROM cqssc ORDER BY time DESC LIMIT 1;`
	cq := new(Cqssc)
	err := DB.QueryRow(str_sql).Scan(&cq.One, &cq.Two, &cq.Three, &cq.Four, &cq.Five)
	if err != nil {
		return "", err
	}
	return cq.One + cq.Two + cq.Three + cq.Four + cq.Five, err
}
