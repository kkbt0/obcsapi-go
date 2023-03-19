#!/bin/bash

if [ -z $1 ]; then
        echo "you must input a port"
        exit 0
fi

PID=$(netstat -nlp | grep ":$1" | awk '{print $7}' | awk -F '[ / ]' '{print $1}')

if [ $? == 0 ]; then
        echo "process id is:${PID}"
else
        echo "process $1 no exit"
        exit 0
fi

kill -9 ${PID}

if [ $? == 0 ]; then
        echo "kill $1 success"
else
        echo "kill $1 fail"
fi
