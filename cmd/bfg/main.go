package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/favetelinguis/bfg-go/betfair"
)

func main() {
	// listMarketsCmd := flag.NewFlagSet("list-markets", flag.ExitOnError)

	// createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	// createMarket := createCmd.String("market", "", "the market to use")

	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)

	logoutCmd := flag.NewFlagSet("logout", flag.ExitOnError)

	// todoCmd := flag.NewFlagSet("todo", flag.ExitOnError)
	// todoId := todoCmd.Int("id", 0, "the id")
	// todoBody := todotodoCmd.String("body", "", "the body")
	// todoCompleted := todoCmd.Bool("completed", false, "mark as completed")

	if len(os.Args) < 2 {
		fmt.Println("missing subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "login":
		if err := loginCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("failed parsing args")
			os.Exit(1)
		}
		var err error
		fmt.Println("running login")
		conf := newLoginConfig()
		client, err := betfair.NewClient(&conf)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := client.Logout()
			if err != nil {
				panic(err)
			}
		}()

		err = client.Login()
		if err != nil {
			panic(err)
		}

		if client.IsSessionExpired() {
			panic("Login failed")
		}
		val, err := client.Betting.ListMarketCatalogue()
		if err != nil {
			panic(err)
		}
		// fmt.Printf("client %+v\n", *client.session)
		fmt.Printf("markets %+v\n", val)

		err = client.Streaming.Connect()
		if err != nil {
			panic(err)
		}
		defer client.Streaming.Close()

		err = client.Streaming.Authenticate()
		if err != nil {
			panic(err)
		}

		err = client.Streaming.SubscribeToMarkets([]string{"1.230518311"})
		if err != nil {
			panic(err)
		}

		// Allow for some time to see response
		time.Sleep(30 * time.Second)

		// client.Account.GetAccountFunds()

		// val, err := session.CallListMarketCatalogue()
		// _ = betfair.NewStream(session)
		// if err != nil {
		// 	panic(err)
		// }

	case "logout":
		if err := logoutCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("failed parsing args")
			os.Exit(1)
		}
		fmt.Println("running logout")

	default:
		fmt.Println("unimplemented subcommand")
		os.Exit(1)
	}
}
