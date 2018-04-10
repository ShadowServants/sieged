#!/usr/bin/env bash
cd flag_adder && gox -osarch="darwin/amd64"  && mv flag_adder_darwin_amd64 ../../ansible/supervisor/binaries/darwin64/flag_adder && cd ..
cd flag_handler && gox -osarch="darwin/amd64" && mv flag_handler_darwin_amd64 ../../ansible/supervisor/binaries/darwin64/flag_handler && cd ..
cd round_handler && gox -osarch="darwin/amd64" && mv round_handler_darwin_amd64 ../../ansible/supervisor/binaries/darwin64/round_handler && cd ..
cd router && gox -osarch="darwin/amd64" && mv router_darwin_amd64 ../../ansible/supervisor/binaries/darwin64/router  && cd ..
