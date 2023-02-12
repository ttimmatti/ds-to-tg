package ds_msgs

import (
	"fmt"
	"strings"
	"time"

	"github.com/ttimmatti/discord-tg_parser/db"
)

func (msgToFind *Msg) isIn(msgs []Msg) bool {
	for _, m := range msgs {
		if m.Id == msgToFind.Id {
			return true
		}
	}
	return false
}

func (msgToCheck *Msg) timestampIs(t string) bool {
	return strings.Contains(msgToCheck.Timestamp, t)
}

func getOldMsgs(channel_id string) ([]Msg, error) {
	msgs, err := db.ReadChannelMsgs(channel_id)
	if err != nil {
		return nil, fmt.Errorf("getOldMsgs: %w", err)
	}

	var resMsgs []Msg
	for _, m := range msgs {
		resMsgs = append(resMsgs, Msg{
			m,
		})
	}

	return resMsgs, nil
}

func filterOld(msgs []Msg, msgsOld []Msg) []Msg {
	resultMsgs := []Msg{}

	t := time.Now().UTC().Format("2006-01-02")

	for _, msg := range msgs {
		if !msg.isIn(msgsOld) && msg.timestampIs(t) {
			resultMsgs = append(resultMsgs, msg)
		}
	}

	return resultMsgs
}

type Msg struct {
	db.Msg
}

type Msgs []Msg

func (ms Msgs) toDb() []db.Msg {
	resMsgs := []db.Msg{}
	for _, m := range ms {
		resMsgs = append(resMsgs, db.Msg{
			Id:         m.Id,
			Content:    m.Content,
			Author:     m.Author,
			Channel_id: m.Channel_id,
			Timestamp:  m.Timestamp,
		})
	}
	return resMsgs
}
