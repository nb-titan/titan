// Package ccs provides GCM Cloud Connection Server (XMPP) client implementation.
// https://developer.android.com/google/gcm/ccs.html
package ccs

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/soygul/nbusy-server/xmpp"
)

const (
	gcmXML    = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`
	gcmACK    = `{"to": "%v", "message_id": "%v", "message_type": "ack"}`
	gcmDomain = "gcm.googleapis.com"
)

// Conn is a GCM CCS connection.
type Conn struct {
	Host, SenderID, APIKey string
	Debug                  bool
	xmppConn               *xmpp.Client
}

// Connect connects to GCM CCS server denoted by host (production or staging CCS endpoint URI) along with relevant credentials.
// Debug mode dumps all CSS communications to stdout.
func Connect(host, senderID, apiKey string, debug bool) (*Conn, error) {
	if !strings.Contains(senderID, gcmDomain) {
		senderID += "@" + gcmDomain
	}

	c, err := xmpp.NewClient(host, senderID, apiKey, debug)
	if debug {
		if err == nil {
			log.Printf("New CCS connection established with XMPP parameters: %+v\n", c)
		} else {
			log.Printf("New CCS connection failed to establish with XMPP parameters: %+v and with error: %v\n", c, err)
		}
	}

	if err != nil {
		return nil, err
	}

	return &Conn{
		Host:     host,
		SenderID: senderID,
		APIKey:   apiKey,
		Debug:    debug,
		xmppConn: c,
	}, nil
}

// Receive retrieves the next incoming messages from the CCS connection.
func (c *Conn) Receive() (*InMsg, error) {
	event, err := c.xmppConn.Recv()
	if err != nil {
		c.Close()
		return nil, err
	}

	switch v := event.(type) {
	case xmpp.Chat:
		isGcmMsg, message, err := c.handleMessage(v.Other[0])
		if err != nil {
			return nil, err
		}
		if isGcmMsg {
			return nil, nil
		}
		return message, nil
	}

	return nil, nil
}

func (c *Conn) handleMessage(msg string) (isGcmMsg bool, message *InMsg, err error) {
	log.Printf("Incoming raw CCS message: %+v\n", msg)
	var m InMsg
	err = json.Unmarshal([]byte(msg), &m)
	if err != nil {
		return false, nil, errors.New("unknow message")
	}

	if m.MessageType != "" {
		switch m.MessageType {
		case "ack":
			return true, nil, nil
		case "nack":
			errFormat := "From: %v, Message ID: %v, Error: %v, Error Description: %v"
			result := fmt.Sprintf(errFormat, m.From, m.ID, m.Err, m.ErrDesc)
			return true, nil, errors.New(result)
		}
	} else {
		ack := fmt.Sprintf(gcmACK, m.From, m.ID)
		c.xmppConn.SendOrg(fmt.Sprintf(gcmXML, ack))
	}

	if m.From != "" {
		return false, &m, nil
	}

	return false, nil, errors.New("unknow message")
}

// Send sends a message to GCM CCS server and returns the number
// of bytes written and any net.Conn write error encountered.
func (c *Conn) Send(message *OutMsg) (n int, err error) {
	res := fmt.Sprintf(gcmXML, message)
	return c.xmppConn.SendOrg(res)
}

// Close a CSS connection.
func (c *Conn) Close() error {
	return c.xmppConn.Close()
}
