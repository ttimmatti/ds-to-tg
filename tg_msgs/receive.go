package tg_msgs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

var ADMIN_ID int64
var LAST_MSG_INDEX int64

func StartReceiving(tg_api string, update_sec int64) {
	TG_API = tg_api
	for {
		// update every $n secs
		time.Sleep(time.Duration(update_sec) * time.Second)

		// get messages
		updates, err := checkMsgs(LAST_MSG_INDEX) //TODO: what if error sending response and etc. bulletproof
		if len(updates) < 1 || err != nil {
			if err != nil {
				log.Println(err)
			}
			// if empty then sleep
			continue
		}

		for _, update := range updates {
			go HandleMsg(update.Message)
		}
	}
}

func HandleMsg(msg Msg) {
	//check that the message is from a private chat, if not skip it
	if msg.From.Id != ADMIN_ID {
		return
	}

	if msg.Text == "" {
		return
	}

	text := msg.Text

	log.Printf("HandleMsg: Received: %s --> %s", msg.From.Username, msg.Text)

	cmd, _, _, err := parseCmd(text)
	if err != nil {
		handleError(msg, err)
		return
	}

	if err := handleCmd(cmd, msg); err != nil {
		handleError(msg, fmt.Errorf("HandleMsg: %w", err))
	}
}

func checkMsgs(index int64) ([]MsgUpd, error) {
	response, err := http.Get(TG_API + "/getUpdates" + "?offset=" + fmt.Sprintf("%d", index+1))
	if err != nil {
		return nil, errror.WrapErrorF(err,
			errror.ErrorCodeFailure,
			"checkMsgs:")
	}
	defer response.Body.Close()

	rBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errror.WrapErrorF(err,
			errror.ErrorCodeFailure,
			"checkMsgs: ")
	}

	update := &UpdatesResponse{}
	json.Unmarshal(rBody, update)

	updates := update.Result

	msgs_ln := len(updates)
	if msgs_ln > 0 {
		last_message := updates[msgs_ln-1]
		LAST_MSG_INDEX = last_message.Update_id
	}

	return updates, nil
}

func getLastMsgIndex() int64 {
	responseU, err := http.Get(TG_API + "/getUpdates")
	if err != nil {
		log.Println(err)
	}
	defer responseU.Body.Close()

	rBody, err := io.ReadAll(responseU.Body)
	if err != nil {
		log.Println(err)
	}

	update := &UpdatesResponse{}
	json.Unmarshal(rBody, update)

	messages := update.Result
	msgs_ln := len(messages)
	if msgs_ln > 0 {
		last_message := messages[msgs_ln-1]
		return last_message.Update_id
	}

	return 0
}

func parseCmd(text string) (string, string, string, error) {
	//parses the command

	//if it's not a cmd
	if !strings.HasPrefix(text, "/") {
		return "", "", "", fmt.Errorf("not a cmd")
	}

	textS := strings.Split(text, " ")
	textN := len(textS)
	if textN < 1 {
		//TODO: let me know about this error
		return "", "", "", fmt.Errorf("!FOR SOME REASON TEXT WAS EMPTY!!! returning")
	}

	// for /start and /delete we need only cmd and value
	// for /update we need cmd and two values
	// for /read we need no values
	switch textN {
	case 1:
		// only cmd -- read
		return textS[0], "", "", nil
	case 2:
		// cmd and val -- add/delete
		return textS[0], textS[1], "", nil
	case 3:
		// cmd and 2 vals -- update
		if strings.Contains(textS[2], "\\\\") {
			textS[2] = strings.Join(strings.Split(textS[2], "\\\\"), "")
		} else if strings.Contains(textS[2], "\\") {
			textS[2] = strings.Join(strings.Split(textS[2], "\\"), "")
		}
		return textS[0], textS[1], textS[2], nil
	default:
		return "", "", "", errror.NewErrorf(
			errror.ErrorCodeWrongCmd,
			"wrongCmd")
		//TODO: RETURN ERROR TO USER
	}
}

// //////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////
// required types
type UpdatesResponse struct {
	Ok     bool
	Result []MsgUpd
}
type MsgUpd struct {
	Update_id int64
	Message   Msg
}
type Msg struct {
	Message_id int64
	From       struct {
		Id         int64
		First_name string
		Username   string
	}
	Chat struct {
		Id int64
	}
	Date     int64
	Text     string
	Entities struct {
		Offset int64
		Length int64
		Type   string
	}
}
