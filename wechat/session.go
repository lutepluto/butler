package wechat

import "log"

// Session is responsible for WeChat login session.
type Session struct {
	UUID        string
	RedirectURI string
	PassTicket  string
}

// DefaultSession is the default WeChat login session.
var DefaultSession *Session

func init() {
	session, err := CreateSession()
	if err != nil {
		log.Panicf("Creating session encounters error %v", err)
	}
	DefaultSession = session
}

// CreateSession creates the WeChat login session.
func CreateSession() (*Session, error) {

	uuid, err := UUID()
	if err != nil {
		return nil, err
	}

	session := &Session{
		UUID: uuid,
	}

	return session, nil
}
