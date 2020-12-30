package slackbi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (b *SBI) serve() error {
	http.HandleFunc(b.config.Service.Path, func(w http.ResponseWriter, r *http.Request) {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			b.logger.Tracef(MF_MSG_RECEIVED_SData.Format(string(body)))

			// ====================================================================
			// Regular message
			// ====================================================================
			var msg_received MsgIncoming
			if err := json.Unmarshal(body, &msg_received); err != nil {
				b.logger.Debugf(MF_MSG_UNMARSHAL_FAILED_SErr_SData.Format(err.Error(), string(body)))
				b.logger.Errorf(MF_MSG_UNMARSHAL_FAILED_SErr.Format(err.Error()))
			} else if msg_received.Type == "url_verification" {
				// {"token":"XLVmO4XWXDqbcxfW0h","challenge":"1GuIC3C5d8DmQ9Jza3VYtJgTdxHCItVzq7v","type":"url_verification"}
				b.logger.Infof(MF_SLACK_BOT_CHALLENGE_SChallenge.Format(msg_received.Challenge))
				w.Write([]byte(msg_received.Challenge))
			} else if msg_received.Event.BotID == "" &&
				msg_received.Event.SubType != "file_share" &&
				msg_received.Event.SubType != "bot_message" {

				// Consider a legit message here
				b.logger.Debugf(MF_MSG_RECEIVED_SFrom_SMsg_SThread.Format(msg_received.Event.User, msg_received.Event.Text, msg_received.Event.TS))

				// If the command starts with "sys" (SYS_COMMAND), intercept that.
				// To avoid modules starting with sys such as `@shorty sysabc` type,
				// but allows `@shorty sys`, this will lowercase all chars, then split it first.
				if strings.Split(strings.ToLower(msg_received.Text()), " ")[0] == SYS_COMMAND {
					// THIS CALLS FOR SYS
					// TODO: if func not nil, run it, otherwise, say sys command not set.

				} else {
					// THIS CALLS FOR MODULES
					// `command` will handle parse and sending response back.
					// Bot will kick off command, but will not wait for the result to come
					// out. This is due to some process take unexpected time.
					go b.command(msg_received)
				}

			} else {
				// Consider a reply of bot itself.... will ignore for now..
				b.logger.Tracef(MF_MSG_RECEIVED_IGNORE_SThread_SBotID_SSubType.Format(msg_received.Event.TS, msg_received.Event.BotID, msg_received.Event.SubType))
			}
			w.Write([]byte(M_HTTP_RESP_OK.String()))
		} else {
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf(MF_HTTP_RESP_UNEXPECTED_SErr.Format(err.Error()))))
			b.logger.Errorf(MF_MSG_READ_FAILED_SErr.Format(err.Error()))
		}
	})
	return http.ListenAndServe(fmt.Sprintf("%s%s", b.config.Service.Address, b.config.Service.Port), nil)
}
