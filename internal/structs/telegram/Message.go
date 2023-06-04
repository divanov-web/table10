package telegram

import "table10/internal/models"

// Message from bot
type Message struct {
	User   *models.User //Send to
	ChatId *int         //ChatId uses against the User
	Text   string       //text message
	Photo  Photo        //photo struct
}
