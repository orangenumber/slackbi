package slackbi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/orangenumber/areq"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type msgReq struct {
	Event struct {
		BotID        string `json:"bot_id"`
		Type         string `json:"type"`
		Subtype      string `json:"subtype"`
		User         string `json:"user"`
		Text         string `json:"text"`
		TS           string `json:"ts"`
		Channel      string `json:"channel"`
		ParentUserID string `json:"parent_user_id"`
		EventTS      string `json:"event_ts"`
		ChannelType  string `json:"channel_type"`
	} `json:"event"`
}

// Text() will return without mention
func (r *msgReq) Text() string {
	return RemoveMention(r.Event.Text)
}

func (r *msgReq) TextRaw() string {
	return r.Event.Text
}

func (r *msgReq) RespFile(slackToken string, filename string) (err error) {
	//curl -F file=@test.txt \
	//	-F "initial_comment=Hello, Leadville" \
	//	-F channels=D01HSN5P7JM \
	//	-H "Authorization: Bearer xoxb-132329662322-1592329468787-GwGCuTncBxv6ioP8nLpOiLku" \
	//	https://slack.com/api/files.upload

	cmd := exec.Command("curl",
		`-F file=@`+filename,
		`-F "initial_comment=result"`,
		`-F "channels=`+r.Event.Channel+`"`,
		`-H "Authorization: Bearer`+slackToken+`"`,
		`https://slack.com/api/files.upload`)
	log.Printf("Running command and waiting for it to finish...")

	for _, v := range cmd.Args {
		println(v)
	}
	return cmd.Run()
}

func (r *msgReq) RespText(slackToken string, replyToThread bool, txt string) error {
	msg_resp := msgResp{
		Token:   slackToken,
		Channel: r.Event.Channel,
		Text:    txt, // "you said " + RemoveMention(msg_req.Event.Text),
		//ThreadTS: msg_req.Event.TS,
	}
	if replyToThread {
		msg_resp.ThreadTS = r.Event.TS
	}
	return msg_resp.Send()
}

type msgResp struct {
	Token    string `json:"token"`
	Channel  string `json:"channel"`
	Text     string `json:"text"`
	ThreadTS string `json:"thread_ts"`
}

func (r *msgResp) Send() error {
	reqData, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = areq.Request("POST", "https://slack.com/api/chat.postMessage",
		areq.PluginFn{
			Name: "Authorization Header",
			FnPre: func(d *areq.AReq) {
				d.Req.Header.Set("Authorization", "Bearer "+r.Token)
			},
		},
		areq.Plugin.SetBody(bytes.NewReader(reqData), "json"),
	)
	// Response
	if err != nil {
		return err
	}
	return nil
}

type msgChallenge struct {
	Type      string `json:"type,omitEmpty"`
	Token     string `json:"token,omitEmpty"`
	Challenge string `json:"challenge,omitEmpty"`
}

func (b *SlackBot) serve() error {
	http.HandleFunc(b.config.Service.Path, func(w http.ResponseWriter, r *http.Request) {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			// Handling challenge
			if b.config.Service.AcceptChallege {
				if bytes.ContainsAny(body, "challenge") {
					var chall msgChallenge
					if err := json.Unmarshal(body, &chall); err == nil {
						b.logger.Infof("Received a slack bot challenge request; challenge=%s", chall.Challenge)
						w.Write([]byte(chall.Challenge))
						return
					}
				}
			}

			// {"token":"XLVmORnn6ha4XWXDqbcxfW0h","challenge":"1Gu5OMmDSab2rbxRKk6gIC3C5d8DmQ9Jza3VYtJgTdxHCItVzq7v","type":"url_verification"}
			var msg_req msgReq

			//if json.Unmarshal(body, &msg_req) == nil && (msg_req.Event.User != "" && msg_req.Event.Subtype != "bot_message") {
			if json.Unmarshal(body, &msg_req) == nil && (msg_req.Event.BotID == "") {
				b.logger.Debugf("MSG received, From=%s, Msg=%s, Thread=%s", msg_req.Event.User, msg_req.Event.Text, msg_req.Event.TS)
				b.logger.Tracef("MSG raw=%s", string(body))

				//msg_req.RespFile(b.config.SlackToken, []byte("hello"), msg_req.Text(), "txt")

				//msg_resp := msgResp{
				//	Token:   b.config.SlackToken,
				//	Channel: msg_req.Event.Channel,
				//	Text:    "you said " + RemoveMention(msg_req.Event.Text),
				//	//ThreadTS: msg_req.Event.TS,
				//}

				if err := msg_req.RespFile(b.config.SlackToken, "test.go"); err != nil {
					b.logger.Errorf("MSG respond failed, Thread=%s, Channel=%s, Msg=%s, Err=%s", msg_req.Event.TS, msg_req.Event.Channel, msg_req.Event.Text, err.Error())
				} else {
					b.logger.Debugf("MSG respond, Thread=%s, Channel=%s", msg_req.Event.TS, msg_req.Event.Channel)
				}

				w.Write([]byte("ok"))
			} else {
				b.logger.Tracef("MSG received but ignored, TS=%s", msg_req.Event.TS)
			}
		} else {
			w.WriteHeader(200)
			w.Write([]byte("Unexpected: " + err.Error()))
		}
	})
	return http.ListenAndServe(fmt.Sprintf("%s:%d", b.config.Service.Address, b.config.Service.Port), nil)
}
