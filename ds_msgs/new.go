package ds_msgs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

var PROXY string

func getMsgs(channel_id string) ([]Msg, error) {
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
		return nil, fmt.Errorf("getMsgs for %s: %w", channel_id, err)
	}

	proxy, err := url.Parse(PROXY)
	if err != nil {
		return nil, fmt.Errorf("getMsgs for %s: %w", channel_id, err)
	}

	r.Header["Authorization"] = []string{TOKEN}

	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("getMsgs for %s: %w", channel_id, err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	msgs := []Msg{}
	if err := json.Unmarshal(b, &msgs); err != nil {
		return nil, fmt.Errorf("getMsgs for %s: %w\nBody: %s", channel_id, err, string(b))
	}

	return msgs, nil
}
