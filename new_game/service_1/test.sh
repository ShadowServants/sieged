#!/bin/bash

bash build.sh
nohup redis-server --port 6378 &
supervisord -c supervisord.conf

