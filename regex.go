package slackbi

import "regexp"

var reMention = regexp.MustCompile("\\s?<@[A-Z0-9]+>\\s?")

func RemoveMention(s string) string {
	return reMention.ReplaceAllString(s, "")
}
