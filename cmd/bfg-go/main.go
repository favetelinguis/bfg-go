package main

import (
	"flag"
	"fmt"
	"os"

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
		fmt.Println("running login")
		session := betfair.NewSession()
		val, err := session.CallListMarketCatalogue()
		if err != nil {
			panic(err)
		}
		fmt.Printf("config %+v\n", val)

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
