#!/usr/bin/env python3

# Parsing code of log data of testbed environment. It create a file usable by encryption code

import os
import sys
import re
import pandas as pd

def parse_file(log_file):

    # regular expressions to match
    # Ex.: RANGING OK [11:0c->19:15] 169 mm [1.0000]
    regex_rng = re.compile(r".*RANGING OK \[(?P<init>\w\w:\w\w)->(?P<resp>\w\w:\w\w)\] (?P<dist>\d+) mm \[(?P<conf>[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+))\].*\n")


    # open log and read line by line
    with open(log_file, 'r') as f:
        for line in f:

            # match transmissions strings
            m = regex_rng.match(line)
            if m:

                # get dictionary of matched groups
                d = m.groupdict()


                init = d['init'] #id of device
                resp = d['resp'] #id of target
                dist = int(d['dist']) #Distance calculate by the device from the target
                conf = float(d['conf']) #Confidence of device


                f = open(f"../DeviceFiles/Device[{init}].txt", "a")
 
                f.write(f"{init},{resp},{dist},{conf}\n")
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
