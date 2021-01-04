package slackbi

import (
	"fmt"
	"github.com/gonyyi/afmt"
	"runtime"
	"sort"
	"strings"
	"time"
)

// sysCommand does not use NLP.
// 1. This is for admin to set. when new command is given, it will search if that is the one
//    who has access by asking questions.
// 2. admin should able to
//    1. check memory AND force GC
//    2. pull recent logs AND export few lines
//    3. set different log level
//    4. refresh/reindex modules
//    5. enable and disable modules
//    6. show current version information
func (b *SBI) sysCommand(in MsgIncoming) {
	switch in.TextNorm() {
	case "sys":
		in.ResponseMarkdown(b, false, syscmdHelp(b))
	case "sys mod":
		in.ResponseText(b, false, syscmdMod(b))
	case "sys mod refresh", "sys mod reload", "sys mod load":
		if err := modRefresh(b); err != nil {
			in.ResponseTextf(b, false, "Mod reload failed, err=%s", err.Error())
		} else {
			in.ResponseTextf(b, false, "Mod reloaded: %d modules loaded", len(b.modules.Items))
		}
	case "sys time", "sys date", "sys what time is it", "sys what time is it?":
		in.ResponseText(b, false, syscmdTime())
	case "sys mem", "sys memory", "sys ram":
		in.ResponseMarkdown(b, false, syscmdMemory())
	case "sys whoisdaddy?", "sys who is the daddy?", "sys who's your daddy?", "sys who is your daddy?",
		"sys who is the daddy", "sys who's daddy", "sys who's your daddy", "sys who is your daddy":
		in.ResponseText(b, false, "Awesome Gon is my daddy!")

		// TODO: need some way to get/set basic log status like logger level, and log file location..
	// case "sys set log level to trace":
	// case "sys set log level to debug":
	// case "sys set log level to info":
	// case "sys set log level to warn":
	// case "sys set log level to error":

	default:
		in.ResponseText(b, false,
			"unrecognized command")
	}
}

func modRefresh(b *SBI) error {
	return b.modules.Load()
}

func modListing(b *SBI) []string {
	var out []string
	for name, v := range b.modules.Items {
		out = append(out, fmt.Sprintf("%s: %s", name, v.Info.Name))
	}
	return out
}

func syscmdMod(b *SBI) string {
	listing := modListing(b)
	sort.Strings(listing)
	out := "```"
	out += strings.Join(listing, "\n")
	return out + b.modules.lastRefresh.Format("\n(Last module indexing: 01/02/2006 15:04:05)```")
}

func syscmdHelp(b *SBI) string {
	lastRefreshAgo := time.Now().Sub(b.modules.lastRefresh)
	return "```" + mf_sbi_sVersion.Format(SBI_VERSION) + "\n" +
		mf_app_sName_sVersion.Format(b.config.BotName, b.config.BotVersion) + "\n" +
		"Last module indexing: " + b.modules.lastRefresh.Format("01/02/2006 15:04:05") + " (" + lastRefreshAgo.String() + " ago)\n" +
		"```"
}

func syscmdTime() string {
	now := time.Now()
	var prefix string

	time10pm, _ := time.Parse("15:04:05", "22:00:00")
	time06am, _ := time.Parse("15:04:05", "06:00:00")

	if now.After(time10pm) || now.Before(time06am) {
		prefix += "It's basically time to go bed!\n"
	}
	timeStr := fmt.Sprintf("It's %s", now.Format("03:04:05 PM Mon, 01/02/2006"))
	return prefix + timeStr
}

func syscmdMemory() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m) // See: https://golang.org/pkg/runtime/#MemStats
	return fmt.Sprintf("```Alloc:      %s\nTotalAlloc: %s\nSys:        %s\nNumGC:      %s\n```",
		afmt.HumanBytes(int64(m.Alloc), 1),
		afmt.HumanBytes(int64(m.TotalAlloc), 1),
		afmt.HumanBytes(int64(m.Sys), 1),
		afmt.HumanNumber(int64(m.NumGC), 1))
}
