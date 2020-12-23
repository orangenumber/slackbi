#!/usr/bin/python

import sys

if len(sys.argv):  # this should be always 3: self, arg1(userid), arg2(rest)
    print "stdIN: <%s>" % sys.stdin.readlines()

    sys.stderr.write("py-stderr\n")
    sys.stdout.write("py-stdout\n")
