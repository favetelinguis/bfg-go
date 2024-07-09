package betfair

import (
	"bufio"
	"fmt"
	"net"
)

type Stream struct {
	conn net.Conn
}

func NewStream(appKey string, token string) *Stream {
	// Socket connection options
	// host := "stream-api.betfair.com"
	host := "stream-api-integration.betfair.com"
	port := "443"

	// Establish connection to the socket
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}
	// TODO how do i handle this when i open a new socket?
	defer conn.Close()
	fmt.Println("Connected")

	// Send authentication message
	authMsg := `{"op": "authentication", "appKey": "<your-appkey>", "session": "<your-session>"}` + "\r\n"
	_, err = conn.Write([]byte(authMsg))
	if err != nil {
		panic(err)
	}

	// Maybe someting like this?
	// go handleClient(conn)
	// TODO up until here
	return &Stream{
		conn,
	}
	// TODO send and read end with CRLF JSON make sure we havve \r\n
	// Subscribe to order/market stream
	// orderMsg := `{"op": "orderSubscription", "orderFilter": {"includeOverallPosition": false, "customerStrategyRefs": ["betstrategy1"], "partitionMatchedByStrategyRef": true}, "segmentationEnabled": true}` + "\r\n"
	// _, err = conn.Write([]byte(orderMsg))
	// if err != nil {
	// 	log.Fatalf("Failed to write data: %v", err)
	// }

	// marketMsg := `{"op":"marketSubscription","id":2,"marketFilter":{"marketIds":["1.120684740"],"bspMarket":true,"bettingTypes":["ODDS"],"eventTypeIds":["1"],"eventIds":["27540841"],"turnInPlayEnabled":true,"marketTypes":["MATCH_ODDS"],"countryCodes":["ES"]},"marketDataFilter":{}}` + "\r\n"
	// _, err = conn.Write([]byte(marketMsg))
	// if err != nil {
	// 	log.Fatalf("Failed to write data: %v", err)
	// }

	// // Handle incoming data
	// go func() {
	// 	buffer := make([]byte, 4096)
	// 	for {
	// 		n, err := conn.Read(buffer)
	// 		if err != nil {
	// 			log.Printf("Error reading data: %v", err)
	// 			return
	// 		}
	// 		fmt.Println("Received: " + string(buffer[:n]))
	// 	}
	// }()

	// // Handle OS signals to keep the program running
	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// <-sig

	// fmt.Println("Connection closed")
}

func (p *Stream) Send(msg string) {
	fmt.Fprintf(p.conn, msg+"\r\n")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// TODO is \n good enogh or do i really need \r\n?
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Process and use data
		fmt.Printf("Received: %s\n", message)
	}
}

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"strings"
// )

// func main() {
// 	conn, err := net.Dial("tcp", "example.com:80")
// 	if err != nil {
// 		fmt.Println("Error connecting:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	bufferedReader := bufio.NewReader(conn)
// 	data, err := readUntilCRLF(bufferedReader)
// 	if err != nil {
// 		fmt.Println("Error reading:", err)
// 		return
// 	}

// 	fmt.Println("Received data:", string(data))
// }

// func readUntilCRLF(reader *bufio.Reader) ([]byte, error) {
// 	var data []byte
// 	for {
// 		line, isPrefix, err := reader.ReadLine()
// 		if err != nil {
// 			return nil, err
// 		}
// 		data = append(data, line...)
// 		if !isPrefix {
// 			data = append(data, '\n')
// 			if strings.HasSuffix(string(data), "\r\n") {
// 				break
// 			}
// 		}
// 	}
// 	return data, nil
// }
