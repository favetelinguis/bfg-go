package betfair

import (
	"crypto/tls"
	"os"
	"path/filepath"
	"time"

	"github.com/favetelinguis/bfg-go/betfair/cache"
)

const (
	login_url       = "https://identitysso-cert.betfair.se/api/"
	identity_url    = "https://identitysso.betfair.se/api/"
	api_betting_url = "https://api.betfair.com/exchange/betting/json-rpc/v1"
	api_account_url = "https://api.betfair.com/exchange/account/json-rpc/v1"
	stream_url      = "stream-api.betfair.com:443"
)

// holds login data
type LoginConfig struct {
	Username string
	Password string
	AppKey   string
	CertFile string
	KeyFile  string
	Locale   string
}

// holds data from successful login
type Token struct {
	SessionToken string
	LoginTime    time.Time
}

// main client
type Session struct {
	loginConfig  *LoginConfig
	token        *Token
	certificates *tls.Certificate
	Betting      *Betting
	Account      *Account
	Streaming    *Streaming
}

type Betting struct {
	Session *Session
}

type Account struct {
	Session *Session
}

type Streaming struct {
	Session     *Session
	conn        *tls.Conn
	closeCh     chan struct{}
	msgCount    int
	StatusCache *cache.StatusCache
}

// creates a new client, this is the central object for interactions with betfair
func NewSession(loginConfig *LoginConfig) (*Session, error) {
	c := &Session{}

	c.token = &Token{}
	var cert tls.Certificate
	var err error

	// load client cert
	homeDir, _ := os.UserHomeDir()
	certPath := filepath.Join(homeDir, ".config", "bfg", loginConfig.CertFile)
	keyPath := filepath.Join(homeDir, ".config", "bfg", loginConfig.KeyFile)
	cert, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	c.certificates = &cert

	// set login config
	c.loginConfig = loginConfig

	// create betting
	c.Betting = &Betting{}
	c.Betting.Session = c

	// create account
	c.Account = &Account{}
	c.Account.Session = c

	// create streaming
	c.Streaming = &Streaming{}
	c.Streaming.Session = c

	return c, nil
}

// check if more then 12h has passed since login/keepAlive
func (s *Session) IsSessionExpired() bool {
	if s.token.SessionToken == "" {
		return true
	}
	duration := time.Since(s.token.LoginTime)

	// TODO not sure its 12 hours for swe but i can test
	return duration.Hours() > 12
}
