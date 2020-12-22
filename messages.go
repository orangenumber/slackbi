package slackbi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/orangenumber/slackbi/lib/postfile"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// Written by Gon Yi
// As of 12/22/2020

// =====================================================================================================================
// COLOR CODE
// =====================================================================================================================
const (
	COLOR_RED        = "#ff0000"
	COLOR_BLUE       = "#0000ff"
	COLOR_ORANGE     = "#ff7b00"
	COLOR_GREEN      = "#00d41c"
	COLOR_YELLOW     = "#fffb00"
	COLOR_GRAY       = "#4a4a4a"
	COLOR_GRAY_LIGHT = "#a8a8a8"
)

// =====================================================================================================================
// INCOMING
// =====================================================================================================================
type MsgIncoming struct {
	Token    string `json:"token,omitempty"`      // "etlXATMDxTR3iWNXir47ksC7"
	TeamID   string `json:"team_id,omitempty"`    // "TCVUCDDDY",
	ApiAppID string `json:"api_app_id,omitempty"` // "AGTCUQE0J",
	Event    struct {
		ClientMsgID string     `json:"client_msg_id,omitempty"` // "f5f5b9e5-69b8-4fb9-af18-5bf8940053f9",
		Type        string     `json:"type,omitempty"`          // "message",
		SubType     string     `json:"subtype,omitempty"`       // "",
		Text        string     `json:"text,omitempty"`          // "hello",
		User        string     `json:"user,omitempty"`          // "UFLEM86PP",
		TS          string     `json:"ts,omitempty"`            // "1601562074.000300",
		Team        string     `json:"team,omitempty"`          // "TCVUCDDDY",
		Blocks      []MsgBlock `json:"blocks,omitempty"`
		Channel     string     `json:"channel,omitempty"`      // "D015QGYPV0F",
		EventTS     string     `json:"event_ts,omitempty"`     // "1601562074.000300",
		ChannelType string     `json:"channel_type,omitempty"` // "im"
		BotID       string     `json:"bot_id"`
	} `json:"event,omitempty"`
	Type         string   `json:"type,omitempty"`          // "event_callback",
	EventID      string   `json:"event_id,omitempty"`      // "Ev01BBQ6GCMD",
	EventTime    int64    `json:"event_time,omitempty"`    // 1601562074,
	AuthedUsers  []string `json:"authed_users,omitempty"`  // ["U0151J9KL3U"],
	EventContext string   `json:"event_context,omitempty"` // "1-message-TCVUCDDDY-D015QGYPV0F"
}

func (in *MsgIncoming) Text() string {
	return strings.TrimSpace(RemoveMention(in.Event.Text))
}
func (in *MsgIncoming) TextRaw() string {
	return in.Event.Text
}

func (in *MsgIncoming) ResponseText(sbi *SBI, useThread bool, text string) error {
	var out MsgOutgoing
	out.Token = sbi.config.SlackToken
	out.Channel = in.Event.Channel
	if useThread {
		out.ThreadTs = in.Event.TS
	}
	out.Text = text
	return out.Send(sbi)
}

func (in *MsgIncoming) Response(sbi *SBI, module string, response *MsgOutgoing) error {
	// Override values and create MsgOutgoing
	var out MsgOutgoing
	out.Token = sbi.config.SlackToken
	out.Channel = in.Event.Channel
	if response.Custom.ReplyInThread {
		out.ThreadTs = in.Event.TS
	}
	out.Text = response.Text
	out.Attachments = response.Attachments
	out.Blocks = response.Blocks

	if (out.Blocks == nil || len(out.Blocks) == 0) && out.Attachments == nil && out.Text == "" {
		if len(response.Custom.Files) == 0 {
			sbi.logger.Warnf("response doesn't have any message, incomingMsg.TS=%s", in.Event.TS)
			return fmt.Errorf("response doesn't have any message, incomingMsg.TS=%s", in.Event.TS)
		}

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		for _, v := range response.Custom.Files {
			file := path.Join(pwd, sbi.config.Module.Dir, module, v)
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			if err := out.SendFile(sbi, v, f); err != nil {
				return err
			}
		}
	} else {
		out.Send(sbi)
	}

	return nil
}

// =====================================================================================================================
// OUTGOING
// =====================================================================================================================
type MsgOutgoing struct {
	Token       string          `json:"token"`               // "xoxb-437964455474-1171621666130-136vHKVZuCqFYexehzNvNnBm",
	Channel     string          `json:"channel"`             // "D015QGYPV0F",
	ThreadTs    string          `json:"thread_ts,omitempty"` // "1601562074.000300"
	Text        string          `json:"text,omitempty"`      // "did you just say `hello`",
	Attachments []MsgAttachment `json:"attachments,omitempty"`
	Blocks      []MsgBlock      `json:"blocks,omitempty"`
	Custom      MsgCustom       `json:"custom,omitmepty"`
}

func (out *MsgOutgoing) Send(sbi *SBI) error {
	payload, err := json.Marshal(out)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", SLACK_ENDPOINT_MSG, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sbi.config.SlackToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		sbi.logger.Errorf("sent request, error=%s", err.Error())
		return err
	}
	sbi.logger.Tracef("sent request, received data=%s", buf.String())

	return nil
}

func (out *MsgOutgoing) SendFile(sbi *SBI, filename string, ior io.Reader) error {
	p := postfile.New()
	if _, err := p.AddFileReader(filename, ior); err != nil {
		return err
	}
	p.AddField("channels", out.Channel)
	p.AddHTTPHeader("Authorization", "Bearer "+sbi.config.SlackToken)

	var buf bytes.Buffer
	if _, err := p.Send("POST", SLACK_ENDPOINT_FILE_UPLOAD, &buf); err != nil {
		return err
	}
	sbi.logger.Tracef("sent request, received data=%s", buf.String())
	return nil
}

// ====================================================================
// CUSTOM
// ====================================================================
type MsgCustom struct {
	ReplyInThread bool     `json:"reply_in_thread"`
	Files         []string `json:"files"`
}

// ====================================================================
// ATTACHMENT
// ====================================================================
type MsgAttachment struct {
	Fallback   string     `json:"fallback,omitempty"`    // Plain-text summary of the attachment.",
	Color      string     `json:"color,omitempty"`       // #2eb886",
	Pretext    string     `json:"pretext,omitempty"`     // Optional text that appears above the attachment block",
	AuthorName string     `json:"author_name,omitempty"` // Bobby Tables",
	AuthorLink string     `json:"author_link,omitempty"` // http://flickr.com/bobby/",
	AuthorIcon string     `json:"author_icon,omitempty"` // http://flickr.com/icons/bobby.jpg",
	Title      string     `json:"title,omitempty"`       // Slack API Documentation",
	TitleLink  string     `json:"title_link,omitempty"`  // https://api.slack.com/",
	Text       string     `json:"text,omitempty"`        // Optional text that appears within the attachment",
	Fields     []MsgField `json:"fields,omitempty"`
	ImageURL   string     `json:"image_url,omitempty"`   // http://my-website.com/path/to/image.jpg",
	ThumbURL   string     `json:"thumb_url,omitempty"`   // http://example.com/path/to/thumb.png",
	Footer     string     `json:"footer,omitempty"`      // Slack API",
	FooterIcon string     `json:"footer_icon,omitempty"` // https://platform.slack-edge.com/img/default_application_icon.png",
	Ts         int64      `json:"ts,omitempty"`          // 123456789
	Blocks     []MsgBlock `json:"blocks,omitempty"`
}

// ====================================================================
// MSG BLOCK
// ====================================================================
type MsgBlock struct {
	Type      string        `json:"type,omitempty"`
	BlockID   string        `json:"block_id,omitempty"` // "D7US",
	Text      *MsgText      `json:"text,omitempty"`
	Fields    []MsgText     `json:"fields,omitempty"`
	Accessory *MsgAccessory `json:"accessory,omitempty"`
	Elements  []MsgElement  `json:"elements,omitempty"`
	Title     *MsgText      `json:"title,omitempty"`
	ImageURL  string        `json:"image_url,omitempty"`
	AltText   string        `json:"alt_text,omitempty"`
}

type MsgAccessory struct {
	Type        string `json:"type,omitempty"`
	Placeholder struct {
		Type  string `json:"type,omitempty"`
		Text  string `json:"text,omitempty"`
		Emoji bool   `json:"emoji,omitempty"`
	} `json:"placeholder,omitempty"`
	Options []struct {
		Text        MsgText `json:"text,omitempty"`
		Description struct {
			Type string `json:"type,omitempty"`
			Text string `json:"text,omitempty"`
		} `json:"description,omitempty"`
		Value string `json:"value,omitempty"`
	} `json:"options,omitempty"`
	Text        MsgText `json:"text,omitempty"`
	Value       string  `json:"value,omitempty"`
	URL         string  `json:"url,omitempty"`
	ActionID    string  `json:"action_id,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	AltText     string  `json:"alt_text,omitempty"`
	InitialTime string  `json:"initial_time,omitempty"`
}

type MsgElement struct {
	Type        string  `json:"type,omitempty"`
	Placeholder MsgText `json:"placeholder,omitempty"`
	ActionID    string  `json:"action_id,omitempty"`
	Options     []struct {
		Text  MsgText `json:"text,omitempty"`
		Value string  `json:"value,omitempty"`
	} `json:"options,omitempty"`
}

// ====================================================================
// FIELD
// ====================================================================
type MsgField struct {
	Title string `json:"title,omitempty"` // Priority",
	Value string `json:"value,omitempty"` // High",
	Short bool   `json:"short,omitempty"` // false
}

// ====================================================================
// TEXT
// ====================================================================
type MsgText struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}
