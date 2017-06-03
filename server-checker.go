package main

import (
	"fmt"
	"strconv"
	"bufio"
	"os"
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"net"
)

var botkey string
var chatid int64
var patch,login,ship string  // Addresses

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
				case 3: login = scan.Text()
						i++
				case 4: ship = scan.Text()
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

	// Connect to the servers
	pcon := dock(patch, "patch")
	lcon := dock(login, "login")
	scon := dock(ship, "ship")

	// Take turns reading from each connection
	pch := make(chan int)
	lch := make(chan int)
	sch := make(chan int)
	go check(pch, pcon, "Patch")
	go check(lch, lcon, "Login")
	go check(sch, scon, "Ship")
	//go check(sch, scon, "Ship")
	sc := 3  // Server counter
	for {
		select {
			case <-pch: msg := tgbotapi.NewMessage(chatid, "patch server down")
						bot.Send(msg)
						sc--
			case <-lch: msg := tgbotapi.NewMessage(chatid, "login server down")
						bot.Send(msg)
						sc--
			case <-sch: msg := tgbotapi.NewMessage(chatid, "ship disconnect")
						bot.Send(msg)
						sc--
		}
		if sc == 0 {
			println("No more active servers.")
			break
		}
	}
}

// Checks to see if still connected by trying to read
func check(ch chan int, conn net.Conn, name string) {
	buff := make([]byte, 400)
	for {
		nbytes, err := conn.Read(buff)
		if err != nil {
			log.Printf("%v server closed the connection.", name)
			//log.Println(err)
			ch <- 1
			break
		}
		log.Printf("%v bytes read from %v server.\n", nbytes, name)
		/*for i := 0; i < nbytes; i++ {
			fmt.Printf("%X ", buff[i])
		}*/
		println()
	}
	return
}

func dock(addr, name string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Can't connect to %v.\n", name)
		os.Exit(1)
	} else {
		log.Printf("Connected to %v.\n", name)
	}
	return conn
}
