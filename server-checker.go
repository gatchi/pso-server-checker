package main

import (
	"strconv"
	"bufio"
	"os"
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"net"
)

var botkey string
var chatid int64
var patch string  // Addresses

func main() {
	// Open configuration file
	file, err := os.Open("server-checker.conf")
	if err != nil {
		log.Fatal(err)
	}

	// Read settings
	data := make([]byte, 100)
	scan := bufio.NewScanner(file)
	for i := 0; scan.Scan(); {
		data = scan.Bytes()
		if data[0] != '#' {
			switch i {
				case 0: botkey = scan.Text()
						i++
				case 1: str := scan.Text()
						chatid, err = strconv.ParseInt(str, 10, 64)
						if err != nil { log.Fatal(err) }
						i++
				case 2: patch = scan.Text()
						i++
			}
		}
	}

	// Setup bot
	bot, err := tgbotapi.NewBotAPI(botkey)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Uncomment for program to display bot responses
	//bot.Debug = true

	// Connect to patch server
	conn, err := net.Dial("tcp", patch)
	if err != nil {
		println("Can't connect to patch.")
		log.Fatal(err)
	} else {
		println("Connected to patch.")
		buf := make([]byte, 4096)
		for {
			nbytes, err := conn.Read(buf)
			if err != nil {
				log.Println("Server closed connection.")
				log.Println(err)
				break
			}
			log.Printf("%v bytes read.\n", nbytes)
		}
	}

	// Send message
	msg := tgbotapi.NewMessage(chatid, "patch server down")
	bot.Send(msg)
}

