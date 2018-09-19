#!/usr/bin/env bash
cd flag_adder && gox -osarch="linux/amd64"  && mv flag_adder_linux_amd64 ../../ansible/supervisor/binaries/linux64/flag_adder && cd ..
cd flag_handler && gox -osarch="linux/amd64" && mv flag_handler_linux_amd64 ../../ansible/supervisor/binaries/linux64/flag_handler && cd ..
cd round_handler && gox -osarch="linux/amd64" && mv round_handler_linux_amd64 ../../ansible/supervisor/binaries/linux64/round_handler && cd ..
cd router && gox -osarch="linux/amd64" && mv router_linux_amd64 ../../ansible/supervisor/binaries/linux64/router  && cd ..
cd http_router && gox -osarch="linux/amd64" && mv http_router_linux_amd64 ../../ansible/supervisor/binaries/linux64/http_router  && cd ..
cd tokens && gox -osarch="linux/amd64" && mv tokens_linux_amd64 ../../ansible/supervisor/binaries/linux64/tokens && cd ..