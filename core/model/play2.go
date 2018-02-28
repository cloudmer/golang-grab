package model

type Play2 struct {
	Id     int
	Number int
	Start  int
	End    int
	Status int
	Type   int
	Time   int
	Cycle  int
}

var Play2Consecutive string = "1"

func (model *Play2) Query(Type string) []*Play2 {
	//只查询开启了报警的数据包
	sql_str := `SELECT * FROM play2 WHERE status = 1 AND type=` + Type + `;`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	data := make([]*Play2, 0)

	for rows.Next() {
		rows.Columns()
		play2 := new(Play2)
		err := rows.Scan(&play2.Id, &play2.Number, &play2.Start, &play2.End, &play2.Status, &play2.Type, &play2.Time, &play2.Cycle)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, play2)
	}

	return data
}

