#!/usr/bin/python

import sys,json

if len(sys.argv):
    rawJ = sys.stdin.readlines()
    j = json.loads("".join(rawJ))
    text = j.get("event", {}).get("text", "")

    f = open("test1.txt", "w")
    f.write("file 111!")
    f.close()

    f = open("test2.txt", "w")
    f.write("file 222")
    f.close()

    f = open("test3.txt", "w")
    f.write("file 2333")
    f.close()

    out = {"custom": {"files": [
        "test1.txt", "test2.txt", "test3.txt"
    ]}}
    sys.stdout.write(json.dumps(out))
