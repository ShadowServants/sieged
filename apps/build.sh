#!/bin/bash

cd service_contoller && go build && cd ..

cd flag_handler && go build && cd ..

cd flag_adder && go build && cd ..