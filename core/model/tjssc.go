package model

type Tjssc struct {
	Id    int
	Qishu string
	One   string
	Two   string
	Three string
	Four  string
	Five  string
	Time  int
}

func (model *Tjssc) Query(limit string) []*Tjssc {
	sql_str := `SELECT * FROM (
					SELECT id,qishu,one,two,three,four,five,time FROM tjssc  ORDER BY time DESC LIMIT ` + limit + `
				) AS ssc ORDER BY time ASC`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	data := make([]*Tjssc, 0)

	for rows.Next() {
		rows.Columns()
		ssc := new(Tjssc)
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
func (model *Tjssc) GetNesCode() (string, error) {
	str_sql := `SELECT one,two,three,four,five FROM tjssc ORDER BY time DESC LIMIT 1;`
	tj := new(Tjssc)
	err := DB.QueryRow(str_sql).Scan(&tj.One, &tj.Two, &tj.Three, &tj.Four, &tj.Five)
	if err != nil {
		return "", err
	}
	return tj.One + tj.Two + tj.Three + tj.Four + tj.Five, err
}
