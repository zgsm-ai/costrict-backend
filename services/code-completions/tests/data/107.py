#!/usr/bin/env python
# -*- coding: utf-8 -*-

import csv
import os
import sys

__version__ = '2.1.4b1'

class FakeShutdownEvent(object):
    """Class to fake a threading.Event.isSet so that users of this module
    are not required to register their own threading.Event()
    """

    @staticmethod
    def isSet():
        "Dummy method to always return false"""
        return False

    
