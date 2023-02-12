package db

import (
	"context"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

const CHANNELS_DB = "channels"

func ReadChannels() ([]Channel, error) {
	rows, err := DB.QueryContext(context.Background(),
		"select channel_id,name,tg_channel_id from "+CHANNELS_DB)
	if err != nil {
		return nil, errror.WrapErrorF(err,
			errror.ErrorCodeFailure,
			"ReadChannels_query_err")
	}

	var channels []Channel

	for i := 0; rows.Next(); i++ {
		var (
			channel_id, name, tg_channel_id string
		)

		_ = rows.Scan(&channel_id, &name, &tg_channel_id)

		channels = append(channels, Channel{
			Channel_id:    channel_id,
			Name:          name,
			Tg_channel_id: tg_channel_id,
		})
	}

	return channels, nil
}

func AddChannel(channel_id, name, tg_channel_id string) error {
	sqlresult, err := DB.ExecContext(context.Background(),
		"insert into "+CHANNELS_DB+"(channel_id,name,tg_channel_id) values($1,$2,$3)",
		channel_id, name, tg_channel_id)
	if err != nil {
		return errror.WrapErrorF(err, errror.ErrorCodeFailure,
			"AddChannel_corrupt_field (channel_id,name):",
			channel_id, name, tg_channel_id)
	}

	rows, _ := sqlresult.RowsAffected()
	if rows == 0 {
		return errror.NewErrorf(errror.ErrorCodeFailure,
			"AddChannel_rows_affected_0 (channel_id,name):",
			channel_id, name, tg_channel_id)
	}

	return nil
}

func DeleteChannel(channel_id string) error {
	sqlresult, err := DB.ExecContext(context.Background(),
		"delete from "+CHANNELS_DB+" where channel_id=$1",
		channel_id)
	if err != nil {
		return errror.WrapErrorF(err, errror.ErrorCodeFailure,
			"DeleteChannel_corrupt_field (channel_id):",
			channel_id)
	}

	rows, _ := sqlresult.RowsAffected()
	if rows == 0 {
		return errror.NewErrorf(errror.ErrorCodeFailure,
			"DeleteChannel_rows_affected_0 (channel_id):",
			channel_id)
	}

	return nil
}

func GetChannelsForChat(tg_channel_id string) ([]Channel, error) {
	rows, err := DB.QueryContext(context.Background(),
		"select channel_id,name from "+CHANNELS_DB+" where tg_channel_id=$1",
		tg_channel_id)
	if err != nil {
		return nil, errror.WrapErrorF(err,
			errror.ErrorCodeFailure,
			"GetChannelsForChat_query_err")
	}

	var channels []Channel

	for i := 0; rows.Next(); i++ {
		var (
			channel_id, name string
		)

		_ = rows.Scan(&channel_id, &name)

		channels = append(channels, Channel{
			Channel_id: channel_id,
			Name:       name,
		})
	}

	return channels, nil
}

type Channel struct {
	Channel_id    string
	Name          string
	Tg_channel_id string
	Msgs          []Msg
}
