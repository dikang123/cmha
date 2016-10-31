package check



func ShowSlave(user, password, host, port, defaultDb, timeout string) (string, error) {
	db, err := Conn(user, password, host, port, defaultDb, timeout)
	if err != nil {
		return "", err
	}
	row, err := db.Query("show slave status")
	if err != nil {
		return "", err
	}
	defer db.Close()
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
			return "", err
		}
	}
	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}

	Slave_IO_Running := mapField2Data["Slave_IO_Running"]
	Slave_IO_Running_String := string(Slave_IO_Running.([]uint8))
	return Slave_IO_Running_String, nil
}
