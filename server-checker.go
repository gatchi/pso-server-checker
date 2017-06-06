package main

import (
	"bufio"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	botkey string
	chatid int64
	addrs = map[string]string {
		"patch": "127.0.0.1",
		"login": "127.0.0.1",
		"ship":  "127.0.0.1",
	}
	ports = map[string]string {
		"patch": "11000",
		"login": "12000",
		"ship":  "5278",
	}
	bot *tgbotapi.BotAPI
)

func main() {
	// Open configuration file
	file, err := os.Open("server-checker.conf")
	if err != nil {
		log.Fatal(err)
	}

	// Read settings
	data := make([]byte, 100)
	scan := bufio.NewScanner(file)
	var i int
	for i = 0; scan.Scan(); {
		data = scan.Bytes()
		if data[0] != '#' {
			switch i {
				case 0: botkey = scan.Text()
						i++
				case 1: str := scan.Text()
						chatid, err = strconv.ParseInt(str, 10, 64)
						if err != nil { log.Fatal(err) }
						i++
				case 2: addrs["patch"] = scan.Text()
						i++
				case 3: ports["patch"] = scan.Text()
						i++
				case 4: addrs["login"] = scan.Text()
						i++
				case 5: ports["login"] = scan.Text()
						i++
				case 6: addrs["ship"] = scan.Text()
						i++
				case 7: ports["ship"] = scan.Text()
						i++
			}
		}
	}
	if i < 7 {
		println("Config file missing fields.")
		println("Filling with defaults.")
	}

	// Setup bot
	bot, err = tgbotapi.NewBotAPI(botkey)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Uncomment to display bot responses
	//bot.Debug = true

	// Connect to the servers
	pcon := connect("patch")
	lcon := connect("login")

	// Monitor connections
	ch := make(chan int)
	go read(ch, pcon, "patch", 1)
	go read(ch, lcon, "login", 2)
	for {
		sig := <-ch
		switch sig {
			case 1: alert("patch")
			case 2: alert("login")
		}
	}
}

func alert(name string) {
	msg := tgbotapi.NewMessage(chatid, name + " server down")
	bot.Send(msg)
}

func connect(name string) net.Conn {
	conn, err := net.Dial("tcp", addrs[name] + ":" + ports[name])
	if err != nil {
		fmt.Printf("Can't connect to %v.\n", name)
		os.Exit(1)
	} else {
		log.Printf("Connected to %v.\n", name)
	}
	return conn
}

// Checks to see if still connected by trying to read
func read(ch chan int, conn net.Conn, name string, code int) {
	buff := make([]byte, 400)
	for {
		nbytes, err := conn.Read(buff)
		if err != nil {
			log.Printf("%v server closed the connection.", name)
			//log.Println(err)
			ch <- code
			break
		}
		log.Printf("%v bytes read from %v server.\n", nbytes, name)
	}
	return
}
