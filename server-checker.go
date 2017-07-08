/*
 * PSO Server-checker
 * Copyright (C) 2017 Christen Gottschlich
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"fmt"
	crypto "github.com/dcrodman/bb_reverse_proxy/encryption"
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
	cCrypt *crypto.PSOCrypt
	sCrypt *crypto.PSOCrypt
)

func main() {
	// Open configuration file
	file, err := os.Open("/usr/local/etc/psobb-server-checker/server-checker.conf")
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Open("server-checker.conf")
			if err != nil {
				fmt.Println("Could not open conf file: " + err.Error())
				os.Exit(1)
			}
		}
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
	scon := connect("ship")

	// Monitor connections
	ch := make(chan int)
	go read(ch, pcon, "patch", 1)
	go read(ch, lcon, "login", 2)
	go read(ch, scon, "ship", 3)
	for {
		sig := <-ch
		switch sig {
			case 1: alert("patch")
			case 2: alert("login")
			case 3: alert("ship")
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

func read(ch chan int, conn net.Conn, name string, code int) {
	buff := make([]byte, 400)
	for {
		nbytes, err := conn.Read(buff)
		if err != nil { // If disconnected
			log.Printf("Server (%v) closed the connection.", name)
			//log.Println(err)
			ch <- code
			break
		}

		// Let's do stuff with it
		if nbytes == 200 && name == "ship" { // 200 means its an auth packet
			var sKey [48]byte
			var cKey [48]byte
			copy(sKey[:], buff[104:152])
			copy(cKey[:], buff[152:200])
			sCrypt = crypto.NewBBCrypt(sKey)
			cCrypt = crypto.NewBBCrypt(cKey)
		}
		if nbytes == 8 { // This is a ping, return it
			sCrypt.Decrypt(buff, 8)
			cCrypt.Encrypt(buff, 8)
			send(buff[:8], conn)
		}
	}
	return
}

func send(packet []byte, conn net.Conn) {
	_, err := conn.Write(packet)
	if err != nil {
		log.Println("Tried to send, but failed.")
	}
}
