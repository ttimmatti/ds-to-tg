package tg_msgs

import (
	"errors"
	"log"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

func handleError(msg Msg, err error) {
	var ierr *errror.Error
	if !errors.As(err, &ierr) {
		log.Println("error wosn't of type Error: ", err)
		return
	}

	log.Printf("err: %s, Msg: %s --> %s", err, msg.From.Username, msg.Text)

	switch ierr.Code() {
	case errror.ErrorCodeFailure,
		errror.ErrorCodeUnknown:
		msg := SendMsg{
			Text:       "errFailure: " + ierr.Error(),
			Chat_id:    ADMIN_ID,
			Parse_mode: "",
		}
		if err := sendMsg(msg); err != nil {
			log.Println("handleError: Couldn't send note about fatal error to admin!!!")
		}
	case errror.ErrorCodeInvalidArgument:
		msg := SendMsg{
			Text:       "Error. Code: Invalid Argument",
			Chat_id:    msg.Chat.Id,
			Parse_mode: "",
		}
		if err := sendMsg(msg); err != nil {
			log.Println("handleError: Couldn't send error msg to user!!!")
		}
	case errror.ErrorCodeNotFound:
		msg := SendMsg{
			Text:       "Error. Code: Not Found",
			Chat_id:    msg.Chat.Id,
			Parse_mode: "",
		}
		if err := sendMsg(msg); err != nil {
			log.Println("handleError: Couldn't send error msg to user!!!")
		}
	case errror.ErrorCodeWrongCmd:
		msg := SendMsg{
			Text:       "Wrong cmd\n---------\nНеправильная команда",
			Chat_id:    msg.Chat.Id,
			Parse_mode: "",
		}
		if err := sendMsg(msg); err != nil {
			log.Println("handleError: Couldn't send error msg to user!!!")
		}
	}
}
