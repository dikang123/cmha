package check


func ShowSlave(user, password, host, port, defaultDb, timeout string) (string, error) {
	db, err := Conn(user, password, host, port, defaultDb, timeout)
	if err != nil {
		//log.Println("conn mysql error:", err)
		return "", err
	}
	row, err := db.Query("show slave status")
	if err != nil {
		//log.Println("exec show slave status error:", err)
		return "", err
	}
	cols, _ := row.Columns()
	buffer := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range buffer {
		buffer[i] = &data[i]
	}
	for row.Next() {
		err = row.Scan(buffer...)
		if err != nil {
	//		fmt.Println(err)
			return "", err
		}
	}
	mapField2Data := make(map[string]interface{}, len(cols))
	for k, col := range data {
		mapField2Data[cols[k]] = col
	}

	Slave_IO_Running := mapField2Data["Slave_IO_Running"]
	Slave_IO_Running_String := string(Slave_IO_Running.([]uint8))
//	log.Println("Slave_IO_Running:%s", Slave_IO_Running_String)
	return Slave_IO_Running_String, nil
}
