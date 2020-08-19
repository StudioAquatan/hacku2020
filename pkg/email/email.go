package email

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"

	idle "github.com/emersion/go-imap-idle"

	"github.com/emersion/go-imap/client"
)

func WatchEmail(body chan string, addr, box, userName, pass string) {
	c, err := login(addr, userName, pass)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	defer logout(c)

	// select mailbox
	_, err = c.Select(box, false)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	idleClient := idle.NewClient(c)

	// Create a channel to receive mailbox updates
	updates := make(chan client.Update)
	c.Updates = updates

	// Start idling
	done := make(chan error, 1)
	go func() {
		done <- idleClient.IdleWithFallback(nil, 0)
	}()

	// Listen for updates
	for {
		select {
		case update := <-updates:
			if _, ok := update.(*client.MailboxUpdate); ok {
				log.Println("[INFO] Mailbox has updated")
				fetchBody(body, addr, box, userName, pass)
			}
		case err := <-done:
			if err != nil {
				log.Fatalf("[ERROR] %s", err)
			}
			log.Println("Not idling anymore")
			return
		}
	}
}

func fetchBody(body chan string, addr, box, userName, pass string) {
	c, err := login(addr, userName, pass)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	defer logout(c)

	mbox, err := c.Select(box, false)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	// Get the last message
	if mbox.Messages == 0 {
		return
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatalf("[ERROR] %s", err)
		}
	}()
	msg := <-messages
	if msg == nil {
		log.Fatal("[ERROR] Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("[ERROR] Server didn't returned message body")
	}
	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)

	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("[ERROR] %s", err)
		}

		switch p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, err := ioutil.ReadAll(p.Body)
			if err != nil {
				log.Fatalf("[ERROR] %s", err)
			}
			body <- string(b)
		}
	}
}

func login(addr, userName, pass string) (*client.Client, error) {
	log.Println("[INFO] Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(addr, nil)
	if err != nil {
		return nil, err
	}
	log.Println("[INFO] Connected")

	// Login
	if err := c.Login(userName, pass); err != nil {
		return nil, err
	}
	log.Println("[INFO] Logged in")
	return c, nil
}

func logout(c *client.Client) {
	err := c.Logout()
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	log.Println("[INFO] Logout")
}
