package storage

import (
	"fmt"
	"github.com/msgbox/relay/structs"
)

func InsertMessage(user string, box string, msg *structs.Message) error {
	db := Connect()
	defer db.Close()

	_, err := db.Exec("INSERT INTO messages (id, user_id, box_id, creator, created_at, payload) VALUES ($1, $2, $3, $4, $5, $6)", msg.GetId(), user, box, msg.GetCreator(), msg.GetCreatedAt(), msg.GetPayload())
	if err != nil {
		return fmt.Errorf("Error writing message to database: %s", err)
	}

	return nil
}
