package tg_msgs

import (
	"fmt"

	"github.com/ttimmatti/discord-tg_parser/db"
	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

const CMDS = `1. Add (enter discord first)
/add <channel_id> <name>

2. Delete
/delete <channel_id>

3. View channels assigned to this group
/channels

4. View this message
/help`

func handleCmd(cmd string, msg Msg) error {
	switch cmd {
	case "/add":
		if err := handleAdd(msg); err != nil {
			return fmt.Errorf("handleCmd: %w", err)
		}
	case "/delete":
		if err := handleDelete(msg); err != nil {
			return fmt.Errorf("handleCmd: %w", err)
		}
	case "/channels":
		if err := handleChannels(msg); err != nil {
			return fmt.Errorf("handleCmd: %w", err)
		}
	case "/help":
		if err := handleHelp(msg); err != nil {
			return fmt.Errorf("handleCmd: %w", err)
		}
	default:
		return errror.NewErrorf(errror.ErrorCodeWrongCmd, "default")
	}

	reply := DefaultReply(msg.Chat.Id, "ok", "")
	if err := sendMsg(*reply); err != nil {
		return fmt.Errorf("handleCmd: %w", err)
	}

	return nil
}

func handleHelp(msg Msg) error {
	reply := DefaultReply(msg.Chat.Id, CMDS, "")
	if err := sendMsg(*reply); err != nil {
		return fmt.Errorf("handleChannels: %w", err)
	}
	return nil
}

func handleChannels(msg Msg) error {
	chs, err := db.GetChannelsForChat(fmt.Sprintf("%d", msg.Chat.Id))
	if err != nil {
		return fmt.Errorf("handleChannels: %w", err)
	}

	var text string
	if len(chs) == 0 {
		text = "No channels here"
	}

	for i, ch := range chs {
		if i != 0 {
			text += "\n=========================\n"
		}

		text += fmt.Sprintf("Channel_id: %s; Name: %s", ch.Channel_id, ch.Name)
	}

	text = escapeMarkdown(text)

	reply := DefaultReply(msg.Chat.Id, text, MarkdownV2)
	if err := sendMsg(*reply); err != nil {
		return fmt.Errorf("handleChannels: %w", err)
	}

	return nil
}

func handleAdd(msg Msg) error {
	_, ch_id, name, err := parseCmd(msg.Text)
	if err != nil {
		return fmt.Errorf("handleAdd: %w", err)
	}
	if err := db.AddChannel(ch_id, name, fmt.Sprintf("%d", msg.Chat.Id)); err != nil {
		return fmt.Errorf("handleAdd: %w", err)
	}
	return nil
}

func handleDelete(msg Msg) error {
	_, ch_id, _, err := parseCmd(msg.Text)
	if err != nil {
		return fmt.Errorf("handleAdd: %w", err)
	}
	if err := db.DeleteChannel(ch_id); err != nil {
		return fmt.Errorf("handleAdd: %w", err)
	}
	return nil
}
