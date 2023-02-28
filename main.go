package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/ttimmatti/discord-tg_parser/db"
	"github.com/ttimmatti/discord-tg_parser/ds_msgs"
	"github.com/ttimmatti/discord-tg_parser/env"
	"github.com/ttimmatti/discord-tg_parser/tg_msgs"
)

const WD = "/home/ttimmatti/my_scripts/go/discord_bot/"

const DS_REPEAT_MINUTES = 20

//https://discord.com/api/v9/channels/1017155496026841118/messages?limit=50

func main() {
	err := godotenv.Load(WD + ".env")
	if err != nil {
		log.Fatalln("Couldnt get Environment")
	}

	db.OnExit()
	db.DB = db.SetConn(env.GetDbEnv())

	tg_msgs.TG_API = env.GetTGApiEnv()
	tg_msgs.ADMIN_ID = env.GetAdminIdEnv()

	go tg_msgs.StartReceiving(
		env.GetTGApiEnv(),
		2,
	)

	ds_msgs.PROXY = env.GetProxyEnv()
	ds_msgs.TOKEN = env.GetDisToken()

	if err := startReceivingDs(); err != nil {
		log.Printf("main: %s", err)
	}
}

func startReceivingDs() error {
	log.Println("startReceiving")
	for {
		handleDsMsgs()
		time.Sleep(DS_REPEAT_MINUTES * time.Minute)
	}
}

func handleDsMsgs() {
	chs, errs := ds_msgs.GetAllNew()
	if errs != nil {
		if err := tg_msgs.HandleErrors(errs); err != nil {
			log.Fatal(err)
		}
	}

	i := 0
	for _, ch := range chs {
		i += len(ch.Msgs)
	}
	log.Printf("-- handleDsMsgs: %d new msgs", i)

	if errs := tg_msgs.SendNewMsgs(chs); errs != nil {
		if err := tg_msgs.HandleErrors(errs); err != nil {
			log.Fatal(err)
		}
	}
}
