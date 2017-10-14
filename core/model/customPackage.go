package model

type CustomPackage struct {
	Id         int
	Alias      string
	Package   string
	Status     int
	Start      int
	End        int
	Continuity int
	Number     int
}

func (model *CustomPackage) Query() []*CustomPackage  {
	//只查询开启了报警的数据
	str_sql := "SELECT * FROM `custom_package` WHERE status=1;"
	rows, err := DB.Query(str_sql)
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	data := make([]*CustomPackage, 0)

	for rows.Next() {
		rows.Columns()
		obj := new(CustomPackage)
		err := rows.Scan(
			&obj.Id,
			&obj.Alias,
			&obj.Package,
			&obj.Status,
			&obj.Start,
			&obj.End,
			&obj.Continuity,
			&obj.Number,
		)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, obj)
	}

	return data
}