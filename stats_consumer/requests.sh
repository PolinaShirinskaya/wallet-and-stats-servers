#!/bin/bash

set -e
#set -u flag for checking variables
#set -x flag for debugging
set -o pipefail

SERVERPORT=5200
SERVERADDR=localhost:${SERVERPORT}

if [ -z "$1" ]
then
echo -e "\e[1;35;4mPlease enter the command:\e[0m"

echo -e "\e[1;35;1m'GET' -\e[0m get wallets stats\e[0m"


# GET /wallets/stats/ - получение статистика по кошелькам (в т.ч. деактивированным)
elif [ -n "$1" ] && [ "$1" = "GET" ]
then
echo -e "\e[1;36;1mGET /wallets/stats/- get wallets statistic \e[0m"
curl -iL -w "\n" -X GET -H "Content-Type: application/json" ${SERVERADDR}/wallets/stats/


fi
