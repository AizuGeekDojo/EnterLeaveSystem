#!/usr/bin/env bash

export DISPLAY=:0.0
chromium --noerrdialogs --kiosk --incognito http://localhost:3000/ >/dev/null 2>&1 &
