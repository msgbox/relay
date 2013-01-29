package storage

import ()

type Box struct {
	Id   string
	Name string
}

func FindBox(user string, name string) *Box {
	db := Connect()
	defer db.Close()
	row := db.QueryRow("SELECT id, name FROM boxes WHERE account_id = ($1) AND name = ($2)", user, name)
	box := new(Box)
	row.Scan(&box.Id, &box.Name)
	return box
}
