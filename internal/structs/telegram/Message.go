package telegram

import "table10/internal/models"

// Message from bot
type Message struct {
	User  *models.User //Send to
	Text  string       //text message
	Photo Photo        //photo struct
}
