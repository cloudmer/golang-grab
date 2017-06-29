package model

type Packet struct {
	Id           int
	Alias        string
	DataTxt      string
	Start        int
	End          int
	RegretNumber int
	Forever      int
	State        int
	Type         int
	Time         int
}

func (model *Packet) Query() []*Packet {
	//只查询开启了报警的数据包
	sql_str := `SELECT * FROM packet WHERE state = 1;`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	data := make([]*Packet, 0)

	for rows.Next() {
		rows.Columns()
		packet := new(Packet)
		err := rows.Scan(&packet.Id, &packet.Alias, &packet.DataTxt, &packet.Start, &packet.End, &packet.RegretNumber, &packet.Forever, &packet.Start, &packet.Type, &packet.Time)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, packet)
	}

	return data
}