#!/bin/bash


cd service_1 && bash test.sh && cd ..
cd service_2 && bash test.sh && cd ..
cd router && go build && cd ..