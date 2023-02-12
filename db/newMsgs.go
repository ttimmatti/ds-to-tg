package db

import (
	"context"
	"fmt"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

func AddNewMsgs(msgs []Msg) (int, error) {
	scss := 0
	var gerr error
	for _, m := range msgs {
		if err := AddMsg(m); err == nil {
			gerr = err
			scss++
		}
	}
	if gerr != nil {
		return scss, fmt.Errorf("AddNewMsgs: %w", gerr)
	}
	return scss, nil
}

func AddMsg(msg Msg) error {
	sqlresult, err := DB.ExecContext(context.Background(),
		"insert into "+MSGS_DB+"(channel_id,msg_id,timestamp) values($1,$2,$3)",
		msg.Channel_id, msg.Id, msg.Timestamp)
	if err != nil {
		return errror.WrapErrorF(err, errror.ErrorCodeFailure,
			"AddMsg_corrupt_field (channel_id,msg_id,timestamp):",
			msg.Channel_id, msg.Id, msg.Timestamp)
	}

	rows, _ := sqlresult.RowsAffected()
	if rows == 0 {
		return errror.NewErrorf(errror.ErrorCodeFailure,
			"AddMsg_affected_0 (channel_id,msg_id,timestamp):",
			msg.Channel_id, msg.Id, msg.Timestamp)
	}

	return nil
}
