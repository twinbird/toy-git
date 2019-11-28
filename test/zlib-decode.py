#!/usr/bin/python
# -*- coding: utf-8 -*-

import zlib
import sys

with open(sys.argv[1], 'rb') as f:
    compressed = f.read()
    decompressed = zlib.decompress(compressed)
    print(decompressed)
