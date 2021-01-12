package plugin

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
)

func GetMessage() (*MsgIncoming, error) {
	data := readStdin()
	if len(data) == 0 {
		return nil, errors.New("no message received")
	}

	var msg MsgIncoming
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func readStdin() []byte {
	// check if there is somethinig to read on STDIN
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var stdin []byte
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			os.Stderr.WriteString("oopsie.. error= " + err.Error())
		}
		return stdin
	} else {
		os.Stderr.WriteString("unexpected: expecting stdin but did not received it")
	}
	return nil
}

type MsgIncoming struct {
	Token    string `json:"token,omitempty"`      // "etlXATMDxTR3i"
	TeamID   string `json:"team_id,omitempty"`    // "TCVUCDDDY",
	ApiAppID string `json:"api_app_id,omitempty"` // "AGTCUQE0J",
	Event    struct {
		ClientMsgID string `json:"client_msg_id,omitempty"` // "f5f5b9e5-69b8-4fb9-af18-5bf8940053f9",
		Type        string `json:"type,omitempty"`          // "message",
		SubType     string `json:"subtype,omitempty"`       // "",
		Text        string `json:"text,omitempty"`          // "hello",
		User        string `json:"user,omitempty"`          // "UFLEM86PP",
		TS          string `json:"ts,omitempty"`            // "1601562074.000300",
		Team        string `json:"team,omitempty"`          // "TCVUCDDDY",
		// Blocks      []MsgBlock `json:"blocks,omitempty"`
		Channel     string `json:"channel,omitempty"`      // "D015QGYPV0F",
		EventTS     string `json:"event_ts,omitempty"`     // "1601562074.000300",
		ChannelType string `json:"channel_type,omitempty"` // "im"
		BotID       string `json:"bot_id"`
	} `json:"event,omitempty"`
	Type         string   `json:"type,omitempty"`          // "event_callback",
	Challenge    string   `json:"challenge,omitempty"`     // this is only for slack API's challenge
	EventID      string   `json:"event_id,omitempty"`      // "Ev01BBQ6GCMD",
	EventTime    int64    `json:"event_time,omitempty"`    // 1601562074,
	AuthedUsers  []string `json:"authed_users,omitempty"`  // ["U0151J9KL3U"],
	EventContext string   `json:"event_context,omitempty"` // "1-message-TCVUCDDDY-D015QGYPV0F"
}

func (in *MsgIncoming) Text() string {
	return strings.TrimSpace(removeMention(in.Event.Text))
}

var reMention = regexp.MustCompile("\\s?<@[A-Z0-9]+>\\s?")

func removeMention(s string) string {
	return reMention.ReplaceAllString(s, "")
}

func NewResponse() *MsgOutgoing {
	var out MsgOutgoing
	return &out
}

type MsgOutgoing struct {
	Text string `json:"text,omitempty"` // "did you just say `hello`",
	// Attachments []MsgAttachment `json:"attachments,omitempty"`
	// Blocks      MsgBlocks       `json:"blocks,omitempty"`
	Custom MsgCustom `json:"custom,omitmepty"`
}

type MsgCustom struct {
	ReplyInThread bool     `json:"reply_in_thread"`
	Files         []string `json:"files,omitempty"`
}

func (out *MsgOutgoing) AsText(s string) *MsgOutgoing {
	out.Text = s
	return out
}

func (out *MsgOutgoing) ReplyToThread(b bool) *MsgOutgoing {
	out.Custom.ReplyInThread = b
	return out
}

func (out *MsgOutgoing) AddFiles(filenames ...string) *MsgOutgoing {
	out.Custom.Files = append(out.Custom.Files, filenames...)
	return out
}

func (out *MsgOutgoing) AsError(s string) *MsgOutgoing {
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	os.Stderr.WriteString(s)
	return out
}

func (out *MsgOutgoing) Send() error {
	b, err := json.Marshal(out)
	if err != nil {
		return err
	}
	os.Stdout.Write(b)
	return nil
}
