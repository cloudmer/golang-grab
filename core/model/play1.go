package model

type Play1 struct {
	Id                   int
	Alias                string
	Package_a            string
	Package_b            string
	Status               int
	Start                int
	End                  int
	ContinuityNumber     int
	Number               int
}

func (model *Play1) Query() []*Play1  {
	//只查询开启了报警的数据
	str_sql := "SELECT * FROM `play1` WHERE status=1;"
	rows, err := DB.Query(str_sql)
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	data := make([]*Play1, 0)

	for rows.Next() {
		rows.Columns()
		obj := new(Play1)
		err := rows.Scan(
			&obj.Id,
			&obj.Alias,
			&obj.Package_a,
			&obj.Package_b,
			&obj.Status,
			&obj.Start,
			&obj.End,
			&obj.ContinuityNumber,
			&obj.Number,
		)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, obj)
	}

	return data
}