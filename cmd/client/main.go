package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

var (
	server   string
	secret   string
	username string
	password string

	count int
)

func init() {
	args := os.Args
	if len(args) != 6 {
		fmt.Println(`

Usage:
      ./radclient radius_server server_secret username password message_counts 

	`)
		os.Exit(0)
	} else {
		var err error
		server = args[1]
		secret = args[2]
		username = args[3]
		password = args[4]
		count, err = strconv.Atoi(args[5])
		if err != nil {
			fmt.Println("please enter a number for message_count")
			os.Exit(0)
		}
		fmt.Println(server, secret, username, password, count)
	}
}

func main() {

	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			packet := radius.New(radius.CodeAccessRequest, []byte(secret))
			rfc2865.UserName_SetString(packet, username)
			rfc2865.UserPassword_SetString(packet, password)
			response, err := radius.Exchange(context.Background(), packet, server)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Code:", response.Code)
			if response.Code == radius.CodeAccessChallenge {
				state := rfc2865.State_GetString(response)
				var otpCode string
				fmt.Println(rfc2865.ReplyMessage_GetString(response))
				if _, err := fmt.Scanln(&otpCode); err != nil {
					log.Println(err)
				}

				rfc2865.UserPassword_SetString(packet, otpCode)
				rfc2865.State_SetString(packet, state)
				challengeResponse, err := radius.Exchange(context.Background(), packet, server)
				if err != nil {
					log.Fatal(err)
				}
				log.Print("Code:", challengeResponse.Code)

			}

			wg.Done()

		}()
	}
	wg.Wait()

}
