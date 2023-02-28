package tg_msgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ttimmatti/discord-tg_parser/db"
	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

var CHARACTERS_NEED_ESCAPING []string = []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

const MarkdownV2 = "MarkdownV2"

func SendNewMsgs(chs []db.Channel) []error {
	var errs []error
	for _, ch := range chs {
		for _, m := range ch.Msgs {
			msg := constructDiscordMsg(m, ch)
			if err := sendMsg(msg); err != nil {
				errs = append(errs, fmt.Errorf("SendNewMsgs: %w", err))
			}
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func constructDiscordMsg(msg db.Msg, ch db.Channel) SendMsg {
	resText := "\\#*" + escapeMarkdown(ch.Name) + "*\n\n"
	resText += "@_" + escapeMarkdown(msg.Author.Username) + "_\n\n"

	dstext := msg.Content
	resText += escapeMarkdown("=======================================\n" + dstext + "\n=======================================")

	chat_id, _ := strconv.ParseInt(ch.Tg_channel_id, 10, 64)
	reply := DefaultReply(chat_id, resText, MarkdownV2)
	return *reply
}

func escapeMarkdown(str string) string {
	for _, esc_char := range CHARACTERS_NEED_ESCAPING {
		str = strings.Join(strings.Split(str, esc_char), "\\"+esc_char)
	}
	return str
}

func sendMsg(msg SendMsg) error {
	if len(TG_API) < 2 {
		return errror.NewErrorf(errror.ErrorCodeFailure,
			"sendMsg: tg_api empty", TG_API)
	}

	respByte, err := json.Marshal(msg)
	if err != nil {
		return errror.WrapErrorF(err,
			errror.ErrorCodeFailure,
			"sendMsg_json_marshal_err")
	}

	resp, err := http.Post(TG_API+"/sendMessage", "Content-Type: application/json", bytes.NewBuffer(respByte))
	if err != nil {
		time.Sleep(2 * time.Second)
		resp, err = http.Post(TG_API+"/sendMessage", "Content-Type: application/json", bytes.NewBuffer(respByte))
		if err != nil {
			return errror.WrapErrorF(err,
				errror.ErrorCodeFailure,
				"sendMsg_post_msg")
		}
	}

	defer resp.Body.Close()

	// remove for prod
	response, _ := io.ReadAll(resp.Body)
	var result map[string]bool
	json.Unmarshal(response, &result)
	ok := result["ok"]
	log.Println("sendMsg: ok:", ok)

	if !ok {
		return errror.NewErrorf(errror.ErrorCodeFailure,
			"sendMsg: ok return false. TG response: ", string(response), "Post body: "+string(msg.Text))
	}

	return nil
}
func DefaultReply(chat_id int64, text, parse_mode string) *SendMsg {
	msg := &SendMsg{
		Chat_id:                  chat_id,
		Text:                     text,
		Parse_mode:               parse_mode,
		Disable_web_page_preview: true,
	}

	return msg
}

type SendMsg struct {
	Chat_id                  int64  `json:"chat_id"`
	Text                     string `json:"text"`
	Parse_mode               string `json:"parse_mode"`
	Disable_web_page_preview bool   `json:"disable_web_page_preview"`
}
