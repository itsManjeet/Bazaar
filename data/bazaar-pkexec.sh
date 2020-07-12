#!/bin/sh

#if [[ $(ps -e | grep 'bazaar') == 0 ]] ; then
#    echo "bazaar is already running"
#    echo "only one process of bazaar is permitted"
#    exit 1
#fi

UI_CMD="/bin/bazaar"

pkexec $UI_CMD
