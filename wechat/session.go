package wechat

import (
	"fmt"
	"log"
	"time"
)

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

	go session.LoginAndServe()
	return session, nil
}

func waitForScan(session *Session) error {
	for {
		select {
		case <-time.After(1 * time.Second):
			log.Println("Going to polling...")
			code, err := PollLogin(session.UUID)
			if err != nil {
				log.Panicf("Polling login state encounter error %v", err)
				return err
			}
			switch code {
			case "200":
				return nil
			case "400":
				return fmt.Errorf("Polling stopped")
			default:
				log.Printf("Getting login state code: %s", code)
			}
		}
	}
}

// LoginAndServe wait for user to login and take over WeChat message receiving
// and responding.
func (s *Session) LoginAndServe() error {
	if err := waitForScan(s); err != nil {
		return err
	}
	return nil
}
