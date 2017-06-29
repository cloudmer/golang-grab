package model

type Mailbox struct {
	Address string
}

func (model *Mailbox) Query() []*Mailbox {
	//只查询收件人
	sql_str := `SELECT email_address FROM mailbox WHERE type = 1;`
	rows, err := DB.Query(sql_str)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	data := make([]*Mailbox, 0)

	for rows.Next() {
		rows.Columns()
		mail := new(Mailbox)
		var err error
		err = rows.Scan(&mail.Address)
		if err != nil {
			panic(err.Error())
		}

		data = append(data, mail)
	}
	return  data
}