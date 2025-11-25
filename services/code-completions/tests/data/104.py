#!/usr/bin/env python
# -*- coding: utf-8 -*-

import csv
import datetime
import os
import sys

__version__ = '2.1.4b1'

def version():
    """Print version"""
    print('speedtest-cli %s' % __version__)
    print('Python %s' % sys.version.replace('\n', ''))

if __name__ == "__main__":
    
