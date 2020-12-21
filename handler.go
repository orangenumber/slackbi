package slackbi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/orangenumber/areq"
	"io/ioutil"
	"net/http"
)

type msgReq struct {
	Event struct {
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		User    string `json:"user"`
		Text    string `json:"text"`
		Ts      string `json:"ts"`
		Channel string `json:"channel"`
	} `json:"event"`
}

type msgResp struct {
	Token    string `json:"token"`
	Channel  string `json:"channel"`
	Text     string `json:"text"`
	ThreadTS string `json:"thread_ts"`
}

func (b *SlackBot) serve() error {
	http.HandleFunc(b.config.Service.Path, func(w http.ResponseWriter, r *http.Request) {
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			var msg_req msgReq

			if json.Unmarshal(body, &msg_req) == nil && (msg_req.Event.User != "" && msg_req.Event.Subtype != "bot_message") {
				msg_resp := msgResp{
					Token: b.config.SlackToken,
					Channel: msg_req.Event.Channel,
					Text: msg_req.Event.Text,
					ThreadTS: msg_req.Event.Ts,
				}

				if reqData, err := json.Marshal(msg_resp); err == nil {
					sc, err := areq.Request("POST", "https://slack.com/api/chat.postMessage",
						areq.PluginFn{
							Name: "Authorization Header",
							FnPre: func(d *areq.AReq) {
								d.Req.Header.Set("Authorization", "Bearer "+msg_resp.Token)
								// d.Logf("[%s] applied (id: %s)", Name, id)
							},
						},
						areq.Plugin.SetBody(bytes.NewReader(reqData), "json"),
					)
					fmt.Printf("responded for %s [%d]: %s\n", msg_req.Event.User, sc, string(reqData))
					fmt.Println(string(body))
					if err != nil {
						fmt.Printf("err: %s\n\n")
					}
				}
				w.Write([]byte("ok"))
			}
		} else {
			w.WriteHeader(200)
			w.Write([]byte("Unexpected: " + err.Error()))
		}
	})
	return http.ListenAndServe(fmt.Sprintf("%s:%d", b.config.Service.Address, b.config.Service.Port), nil)
}
