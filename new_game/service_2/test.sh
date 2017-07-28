#!/bin/bash

bash build.sh
nohup redis-server --port 6377 &
supervisord -c supervisord.conf


