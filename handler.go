package slackbi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type msgChallenge struct {
	Type      string `json:"type,omitEmpty"`
	Token     string `json:"token,omitEmpty"`
	Challenge string `json:"challenge,omitEmpty"`
}

func (b *SBI) serve() error {
	http.HandleFunc(b.config.Service.Path, func(w http.ResponseWriter, r *http.Request) {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			// ====================================================================
			// Handling challenge
			// ====================================================================
			if b.config.Service.AcceptChallege {
				// {"token":"XLVmORnn6ha4XWXDqbcxfW0h","challenge":"1Gu5OMmDSab2rbxRKk6gIC3C5d8DmQ9Jza3VYtJgTdxHCItVzq7v","type":"url_verification"}
				if bytes.ContainsAny(body, "challenge") {
					var chall msgChallenge
					if err := json.Unmarshal(body, &chall); err == nil {
						b.logger.Infof("Received a slack bot challenge request; challenge=%s", chall.Challenge)
						w.Write([]byte(chall.Challenge))
						return
					}
				}
			}

			// ====================================================================
			// Regular message
			// ====================================================================
			var msg_received MsgIncoming
			if err := json.Unmarshal(body, &msg_received); err != nil {
				b.logger.Debugf("MSG received but can't unmarshal, data=%s, err=%s", string(body), err.Error())
				b.logger.Errorf("MSG received but can't unmarshal, err=%s", err.Error())
			} else if msg_received.Event.BotID == "" || (msg_received.Event.User != "" && msg_received.Event.SubType != "bot_message") {
				// ====================================================================
				// Consider a legit message here
				// ====================================================================
				b.logger.Debugf("MSG msg_received, From=%s, Msg=%s, Thread=%s", msg_received.Event.User, msg_received.Event.Text, msg_received.Event.TS)
				b.logger.Tracef("MSG raw=%s", string(body))

				// ====================================================================
				// `command` will handle parse and sending response back.
				// Bot will kick off command, but will not wait for the result to come
				// out. This is due to some process take unexpected time.
				// ====================================================================
				go b.command(msg_received)

			} else {
				// Consider a reply of bot itself.... will ignore for now..
				b.logger.Tracef("MSG msg_received but ignored, TS=%s", msg_received.Event.TS)
			}
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(200)
			w.Write([]byte("Unexpected: " + err.Error()))
			b.logger.Errorf("MSG received but failed to read, err=%s", err.Error())
		}
	})
	return http.ListenAndServe(fmt.Sprintf("%s%s", b.config.Service.Address, b.config.Service.Port), nil)
}
