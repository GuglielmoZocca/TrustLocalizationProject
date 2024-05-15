#!/usr/bin/env python3

#Code to calculate the average period of time to do the ranging, from data in log file

import os
import sys
import re

#Calculate the avg interval to range a distance between a target and a device
def parse_file(log_file):

    # regular expressions to match
    # Ex.: RANGING OK [11:0c->19:15] 169 mm [1.000][0.0033]
    regex_rng = re.compile(r".*RANGING OK \[(?P<init>\w\w:\w\w)->(?P<resp>\w\w:\w\w)\] (?P<dist>\d+) mm \[(?P<conf>[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+))\]\[(?P<temp>[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+))\].*\n")

    # open log and read line by line
    with open(log_file, 'r') as f:

        sumTemp = float(0)
        num = 0
        for line in f:

            # match transmissions strings
            m = regex_rng.match(line)
            if m:
                num = num + 1
                d = m.groupdict()
                ini = d['init'] #id of device
                sumTemp = sumTemp + float(d['temp'])*1000 #Sum mil seconds interval



        f = open(f"/ParsingCode/AVGtemp[{ini}].txt", "w")
        avgTemp = sumTemp/num #avg time
        f.write(f"avg interval in milliseconds:{avgTemp}")
        f.close()


if __name__ == '__main__':

    if len(sys.argv) < 2:
        print("Error: Missing log file.")
        sys.exit(1)

    # get the log path and check that it exists
    log_file = sys.argv[1]
    if not os.path.isfile(log_file) or not os.path.exists(log_file):
        print("Error: Log file not found.")
        sys.exit(1)

    # parse file
    parse_file(log_file)
