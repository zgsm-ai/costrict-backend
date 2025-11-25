#!/usr/bin/env python

import sys

this_python = sys.version_info[:2]
min_version = (3, 9)
if this_python < min_version:
    message_parts = [
        "This script does not work on Python {}.{}.".format(*this_python),
        "The minimum supported Python version is {}.{}.".format(*min_version),
        "Please use https://bootstrap.pypa.io/pip/{}.{}/get-pip.py instead.".format(*this_python),
    ]
    print("ERROR: " + " ".join(message_parts))
    sys.exit(1)

import os.path
import pkgutil
import shutil
import tempfile

def main():
    print("Python version check passed")
    print("Minimum required version: 3.9")
    print("Current version: {}.{}".format(*this_python))

