package slackbi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (b *SBI) serve() error {
	http.HandleFunc(b.config.Service.Path, func(w http.ResponseWriter, r *http.Request) {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			b.logger.Tracef("Received a call, data=%s", string(body))

			// ====================================================================
			// Regular message
			// ====================================================================
			var msg_received MsgIncoming
			if err := json.Unmarshal(body, &msg_received); err != nil {
				b.logger.Debugf("MSG received but can't unmarshal, data=%s, err=%s", string(body), err.Error())
				b.logger.Errorf("MSG received but can't unmarshal, err=%s", err.Error())
			} else if msg_received.Type == "url_verification" {
				// {"token":"XLVmO4XWXDqbcxfW0h","challenge":"1GuIC3C5d8DmQ9Jza3VYtJgTdxHCItVzq7v","type":"url_verification"}
				b.logger.Infof("Received a slack bot challenge request; challenge=%s", msg_received.Challenge)
				w.Write([]byte(msg_received.Challenge))
			} else if msg_received.Event.BotID == "" &&
				msg_received.Event.SubType != "file_share" &&
				msg_received.Event.SubType != "bot_message" {
				// ====================================================================
				// Consider a legit message here
				// ====================================================================
				b.logger.Debugf("MSG msg_received, From=%s, Msg=%s, Thread=%s", msg_received.Event.User, msg_received.Event.Text, msg_received.Event.TS)

				// ====================================================================
				// `command` will handle parse and sending response back.
				// Bot will kick off command, but will not wait for the result to come
				// out. This is due to some process take unexpected time.
				// ====================================================================
				go b.command(msg_received)

			} else {
				// Consider a reply of bot itself.... will ignore for now..
				b.logger.Tracef("MSG msg_received but ignored, TS=%s, BotID=%s, SubType=%s", msg_received.Event.TS, msg_received.Event.BotID, msg_received.Event.SubType)
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
