#!/bin/sh
sudo launchctl limit maxfiles 6553500
ulimit -n 6553500
tower