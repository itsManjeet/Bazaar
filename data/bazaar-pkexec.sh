#!/bin/sh

if test "z`ps -e | grep bazaar`" != "z" ; then
    echo "bazaar is already running"
    echo "only one process of bazaar is permitted"
    exit 1
fi

UI_CMD="/usr/bin/bazaar"

pkexec $UI_CMD
