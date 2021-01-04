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
			b.logger.Tracef(mf_msg_received_sData.Format(string(body)))

			// ====================================================================
			// Regular message
			// ====================================================================
			var msgin MsgIncoming
			if err := json.Unmarshal(body, &msgin); err != nil {
				b.logger.Debugf(mf_msg_unmarshal_failed_sErr_sData.Format(err.Error(), string(body)))
				b.logger.Errorf(mf_msg_unmarshal_failed_sErr.Format(err.Error()))
			} else if msgin.Type == "url_verification" {
				// {"token":"XLVmO4XWXDqbcxfW0h","challenge":"1GuIC3C5d8DmQ9Jza3VYtJgTdxHCItVzq7v","type":"url_verification"}
				b.logger.Infof(mf_slack_bot_challenge_sChallenge.Format(msgin.Challenge))
				w.Write([]byte(msgin.Challenge))
			} else if msgin.Event.BotID == "" &&
				msgin.Event.SubType != "file_share" &&
				msgin.Event.SubType != "bot_message" {

				// Consider a legit message here
				b.logger.Debugf(mf_msg_received_sFrom_sMsg_sThread.Format(msgin.Event.User, msgin.Event.Text, msgin.Event.TS))

				// If the command starts with "sys" (SYS_COMMAND), intercept that.
				// To avoid modules starting with sys such as `@shorty sysabc` type,
				// but allows `@shorty sys`, this will lowercase all chars, then split it first.
				if strings.Split(strings.ToLower(msgin.Text()), " ")[0] == b.config.Service.SysCommand {
					// THIS CALLS FOR SYS
					go b.sysCommand(msgin)

				} else {
					// THIS CALLS FOR MODULES
					// `command` will handle parse and sending response back.
					// Bot will kick off command, but will not wait for the result to come
					// out. This is due to some process take unexpected time.
					go b.command(msgin)
				}

			} else {
				// Consider a reply of bot itself.... will ignore for now..
				b.logger.Tracef(mf_msg_received_ignore_sThread_sBotID_sSubtype.Format(msgin.Event.TS, msgin.Event.BotID, msgin.Event.SubType))
			}
			w.Write([]byte(m_http_resp_ok.String()))
		} else {
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf(MF_HTTP_RESP_UNEXPECTED_SErr.Format(err.Error()))))
			b.logger.Errorf(mf_msg_read_failed_sErr.Format(err.Error()))
		}
	})
	return http.ListenAndServe(fmt.Sprintf("%s%s", b.config.Service.Address, b.config.Service.Port), nil)
}
