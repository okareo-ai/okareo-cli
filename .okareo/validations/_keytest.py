#!/usr/bin/env python3

import os

OKAREO_API_KEY = os.getenv("OKAREO_API_KEY")
OKAREO_RUN_ID = os.getenv("OKAREO_RUN_ID")

if not OKAREO_API_KEY:
    print("OKAREO_API_KEY is not set")

if not OKAREO_RUN_ID:
    print("OKAREO_RUN_ID is not set")

if OKAREO_API_KEY and OKAREO_RUN_ID:
    print('All keys are available.')
