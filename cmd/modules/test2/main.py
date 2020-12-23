#!/usr/bin/python

import sys,json

if len(sys.argv):
    rawJ = sys.stdin.readlines()
    j = json.loads("".join(rawJ))
    text = j.get("event", {}).get("text", "")

    out = {"blocks": [
        {"type": "section", "text": {"type": "mrkdwn", "text": "```"+text+"```"}},
        {"type": "divider"},
        {"type": "section", "text": {"type": "mrkdwn", "text": "```whatever "+text+"```"}}
    ]}
    sys.stdout.write(json.dumps(out))
