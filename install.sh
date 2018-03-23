#! /bin/sh

YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NOCOLOR='\033[0m'

# install the golang binary (for the moment you need golang to install)
printf "${YELLOW}Building Dockerlang Compterpreter..."
cd docdocdoc && go install && cd ..
printf "${GREEN}ok! ${NOCOLOR} \n"

# build the dockerlang memorycell image
printf "${YELLOW}Building Dockerlang memory cells..."
docker build . -t "dockerlang" >/dev/null
printf "${GREEN}ok! ${NOCOLOR} \n"
