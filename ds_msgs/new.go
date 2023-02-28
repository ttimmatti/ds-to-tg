package ds_msgs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

var PROXY string

func getMsgs(channel_id string, tryN int) ([]Msg, error) {
	if len(TOKEN) < 2 {
		return nil, errror.NewErrorf(errror.ErrorCodeFailure,
			"getMsgs: discord token empty", TOKEN)
	}
	if len(PROXY) < 2 {
		return nil, errror.NewErrorf(errror.ErrorCodeFailure,
			"getMsgs: proxy empty", TOKEN)
	}

	r, err := http.NewRequest(http.MethodGet, DIS_CHANNELS_API+channel_id+"/messages?limit=50", nil)
	if err != nil {
		return nil, errror.WrapErrorF(
			err,
			errror.ErrorCodeFailure,
			fmt.Sprintf("1getMsgs for %s", channel_id),
		)
	}

	proxy, err := url.Parse(PROXY)
	if err != nil {
		return nil, errror.WrapErrorF(
			err,
			errror.ErrorCodeFailure,
			fmt.Sprintf("2getMsgs for %s", channel_id),
		)
	}

	r.Header["Authorization"] = []string{TOKEN}

	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errror.WrapErrorF(
			err,
			errror.ErrorCodeFailure,
			fmt.Sprintf("3getMsgs for %s", channel_id),
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 500 && tryN < 1 {
		time.Sleep(2 * time.Second)
		msgs, err := getMsgs(channel_id, tryN+1)
		if err == nil {
			return msgs, nil
		}
		return nil, err
	}

	b, _ := io.ReadAll(resp.Body)

	msgs := []Msg{}
	if err := json.Unmarshal(b, &msgs); err != nil {
		return nil, errror.WrapErrorF(
			err,
			errror.ErrorCodeFailure,
			fmt.Sprintf("4getMsgs for %s; Body: %s", channel_id, string(b)),
		)
	}

	return msgs, nil
}
