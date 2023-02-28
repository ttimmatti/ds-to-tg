package ds_msgs

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ttimmatti/discord-tg_parser/db"
)

const DIS_CHANNELS_API = "https://discord.com/api/v9/channels/"

var TOKEN string

func GetAllNew() ([]db.Channel, []error) {
	errs := []error{}

	chs, err := db.ReadChannels()
	if err != nil {
		return nil, []error{fmt.Errorf("GetAllNew: %w", err)}
	}

	for i, ch := range chs {
		if i != 0 {
			time.Sleep(time.Duration(rand.Int()%18+3) * time.Second)
		}

		msgs, err := getNew(ch.Channel_id)
		if err != nil {
			errs = append(errs, fmt.Errorf("[ERROR] GetAllNew (Probably missing Access) :: %w", err))
		}
		chs[i].Msgs = msgs.toDb()

		log.Printf("[INFO] GetAllNew: Checked %s", ch.Name)
	}

	if len(errs) != 0 {
		return chs, errs
	}

	return chs, nil
}

func getNew(channel_id string) (Msgs, error) {
	msgsNew := Msgs{}

	msgsOld, err := getOldMsgs(channel_id)
	if err != nil {
		return msgsNew, fmt.Errorf("getNew: %w", err)
	}

	msgs, err := getMsgs(channel_id, 0)
	if err != nil {
		return msgsNew, fmt.Errorf("getNew: %w", err)
	}

	// msgsNew = msgs - msgs(not from today by timestamp) - msgsOld
	msgsNew = filterOld(msgs, msgsOld)

	scss, err := db.AddNewMsgs(msgsNew.toDb())
	if err != nil {
		log.Printf("ds_msgs: GetNew: added %d out of %d msgs to db. err: %s",
			scss, len(msgsNew), err)
	}

	return msgsNew, nil
}

//

//

//

//

//

//

//

//

//

//

//

//

func printMsgs(msgs []Msg) error {
	for _, msg := range msgs {
		fmt.Printf(
			"\033[34mNEWMESSAGE!\033[0m\nId: %s\nAuthor: %s -->\nMsg: %s\nChannel_id: %s\nTime: %s\n\n\n", msg.Id, msg.Author.Username, msg.Content, msg.Channel_id, msg.Timestamp,
		)
	}

	return nil
}
