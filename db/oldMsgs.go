package db

import (
	"context"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

func ReadChannelMsgs(channel_id string) ([]Msg, error) {
	rows, err := DB.QueryContext(context.Background(),
		"select channel_id,msg_id,timestamp from "+MSGS_DB)
	if err != nil {
		return nil, errror.WrapErrorF(err,
			errror.ErrorCodeFailure,
			"ReadChannelMsgs_query_err", channel_id)
	}

	var msgs []Msg

	for i := 0; rows.Next(); i++ {
		var (
			channel_id, msg_id, timestamp string
		)

		_ = rows.Scan(&channel_id, &msg_id, &timestamp)

		msgs = append(msgs, Msg{
			Channel_id: channel_id,
			Id:         msg_id,
			Timestamp:  timestamp,
		})
	}

	return msgs, nil
}

type Msg struct {
	Id         string
	Content    string
	Channel_id string
	Author     struct {
		Username string
	}
	Timestamp string
}
