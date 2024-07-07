package betfair

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type loginResponse struct {
	SessionToken string
	LoginStatus  string
}

func login(appKey string, username string, password string) string {
	// Load client cert
	homeDir, _ := os.UserHomeDir()
	certPath := filepath.Join(homeDir, ".config", "bfg", "betfair-2048.crt")
	keyPath := filepath.Join(homeDir, ".config", "bfg", "betfair-2048.key")
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		panic(err)
	}

	// Create tls.Config with the client certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Create http.Transport with the tls.Config
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Create http.Client with the transport
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	// Create body
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	request, err := http.NewRequest(http.MethodPost,
		"https://identitysso-cert.betfair.se/api/certlogin",
		strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Application", appKey)

	// Make a request
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Handle response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var response loginResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}
	if response.LoginStatus != "SUCCESS" {
		panic(fmt.Sprintf("login status not SUCCESS is  %s", response.LoginStatus))
	}
	return response.SessionToken
}

// HOW TOSENDJSON
// type Post struct {
//     Title  string `json:"title"`
//     Body   string `json:"body"`
//     UserID int    `json:"userId"`
// }

// bytes, err := json.Marshal(Post{Title: "foo", Body: "bar", UserID: 1})
// if err != nil {
//     log.Fatal(err)
// }

// request, err := http.NewRequest(http.MethodPost, "https://jsonplaceholder.typicode.com/posts", bytes.NewBuffer(bytes))
