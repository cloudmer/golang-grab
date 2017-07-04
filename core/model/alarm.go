package model

type Alarm struct {
	Id int
	Number int
	Start  int
	End    int
	Status int
	Type   int
	Time   int
}

var AlarmConsecutive string = "1"

func (model *Alarm) Query(Type string) []*Alarm {
	//只查询开启了报警的数据包
	sql_str := `SELECT * FROM alarm WHERE status = 1 AND type=` + Type + `;`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	data := make([]*Alarm, 0)

	for rows.Next() {
		rows.Columns()
		alarm := new(Alarm)
		err := rows.Scan(&alarm.Id, &alarm.Number, &alarm.Start, &alarm.End, &alarm.Status, &alarm.Type, &alarm.Time)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, alarm)
	}

	return data
}

