package main

import (
	"strconv"
	"bufio"
	"os"
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"net"
)

func main() {
	// Read settings for bot
	file, err := os.Open("servercheck.conf")
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 100)
	scan := bufio.NewScanner(file)
	var botkey string
	var chatid int64
	for i := 0; scan.Scan(); {
		data = scan.Bytes()
		if data == nil {
			// Don't do anything
		} else if err = scan.Err(); err != nil {
			// Again, don't do anything
		} else if data[0] != '#' {
			switch i {
				case 0: botkey = scan.Text()
						i++
				case 1: str := scan.Text()
						chatid, err = strconv.ParseInt(str, 10, 64)
						if err != nil { log.Fatal(err) }
			}
		}
	}

	// Setup bot
	bot, err := tgbotapi.NewBotAPI(botkey)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Connect to patch server
	conn, err := net.Dial("tcp", "127.0.0.1:11000")
	if err != nil { 
		print("Can't connect to patch.")
		log.Fatal(err)
	} else {
		print("Connected to patch.\n")
		buf := make([]byte, 4096)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				log.Printf("Server closed connection.\n")
				break
			}
		}
	}

	// Send message
	msg := tgbotapi.NewMessage(chatid, "patch server down")
	bot.Send(msg)
}

