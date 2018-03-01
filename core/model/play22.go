package model

type Play22 struct {
	Id     int
	Number int
	Start  int
	End    int
	Status int
	Time   int
}

var Play22Consecutive string = "1"

func (model *Play22) Query(Type string) []*Play22 {
	//只查询开启了报警的数据包
	sql_str := `SELECT * FROM play22 WHERE status = 1;`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	data := make([]*Play22, 0)

	for rows.Next() {
		rows.Columns()
		play22 := new(Play22)
		err := rows.Scan(&play22.Id, &play22.Number, &play22.Start, &play22.End, &play22.Status, &play22.Time)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, play22)
	}

	return data
}

