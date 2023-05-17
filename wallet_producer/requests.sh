#!/bin/bash

set -e
#set -u flag for checking variables
#set -x flag for debugging
set -o pipefail

SERVERPORT=5100
SERVERADDR=localhost:${SERVERPORT}

if [ -z "$1" ]
then
echo -e "\e[1;35;4mPlease enter the command:\e[0m"
echo -e "\e[1;35;1m'POST' -\e[0m create wallets\e[0m"
echo -e "\e[1;35;1m'GET ALL' -\e[0m get list of all wallets\e[0m"
echo -e "\e[1;35;1m'GET' + 'wallet ID' -\e[0m get wallet by ID\e[0m"
echo -e "\e[1;35;1m'DELETE' + 'wallet ID' -\e[0m delete wallet by ID\e[0m"
echo -e "\e[1;35;1m'PUT' + 'wallet ID' + 'new name' -\e[0m update wallet name by ID\e[0m"
echo -e "\e[1;35;1m'POST DEPOSIT' + 'wallet ID' -\e[0m deposit wallet by ID\e[0m"
echo -e "\e[1;35;1m'POST WITHDRAW' + 'wallet ID' -\e[0m withdraw wallet by ID\e[0m"
echo -e "\e[1;35;1m'POST DEPOSIT' + 'wallet ID' + 'transfer to ID' -\e[0m transfer by IDs\e[0m"



# POST/wallets создание кошелька
elif [ -n "$1" ] && [ "$1" = "POST" ]
then
echo -e "\e[1;36;1mPOST /wallets/ - creating wallets\e[0m"
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"name":"John"}' ${SERVERADDR}/wallets/
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"name":"Anna"}' ${SERVERADDR}/wallets/
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"name":"Paul"}' ${SERVERADDR}/wallets/



# GET /wallets - получение списка всех доступных кошельков (в т.ч. деактивированных)
elif [ -n "$1" ] && [ "$1" = "GET ALL" ]
then
echo -e "\e[1;36;1mGET /wallets/- get all wallets\e[0m"
curl -iL -w "\n" -X GET -H "Content-Type: application/json" ${SERVERADDR}/wallets/



# GET /wallets/{id} - получение кошелька по его идентификатору
elif [ -n "$1" ] && [ "$1" = "GET" ] && [ -n "$2" ]
then
echo -e "\e[1;36;1mGET /wallets/{id} - find wallet by ID\e[0m"
curl -iL -w "\n" -X GET -H "Content-Type: application/json" ${SERVERADDR}/wallets/"$2"


# DELETE /wallets/{id} - деактивация кошелька по идентификатору в path params
elif [ -n "$1" ] && [ "$1" = "DELETE" ] && [ -n "$2" ]
then
echo -e "\e[1;36;1mDELETE /wallets/{id} - delete wallet by ID\e[0m"
curl -iL -w "\n" -X DELETE -H "Content-Type: application/json" ${SERVERADDR}/wallets/"$2"


# PUT /wallets/{id} - обновление кошелька по его идентификатору
elif [ -n "$1" ] && [ "$1" = "PUT" ] && [ -n "$2" ] && [ -n "$3" ] 
then
echo -e "\e[1;36;1mPUT /wallets/{id} - update wallet by ID\e[0m"
curl -iL -w "\n" -X PUT -H "Content-Type: application/json" --data '{"name":"'$3'"}' ${SERVERADDR}/wallets/"$2"


# POST /wallets/{id}/deposit - метод пополнения кошелька
elif [ -n "$1" ] && [ "$1" = "POST DEPOSIT" ] && [ -n "$2" ]
then
echo -e "\e[1;36;1mPUT /wallets/{id}/deposit - deposit wallet by ID\e[0m"
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"amount": 1000.42}' ${SERVERADDR}/wallets/"$2"/deposit


# POST /wallets/{id}/withdraw - метод снятия средств с кошелька
elif [ -n "$1" ] && [ "$1" = "POST WITHDRAW" ] && [ -n "$2" ]
then
echo -e "\e[1;36;1mPUT /wallets/{id}/withdraw - withdraw wallet by ID\e[0m"
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"amount": 500.42}' ${SERVERADDR}/wallets/"$2"/withdraw


# POST /wallets/{id}/transfer - перевод между двумя кошельками.
elif [ -n "$1" ] && [ "$1" = "POST TRANSFER" ] && [ -n "$2" ] && [ -n "$3" ]
then
echo -e "\e[1;36;1mPUT /wallets/{id}/transfer - transfer by IDs\e[0m"
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"amount": 500.42, "transfer_to": "'$3'"}' ${SERVERADDR}/wallets/"$2"/transfer


fi
