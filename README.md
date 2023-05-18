# Wallet and Stats web-servers

## Описание проекта:
> Полное описание проекта можно найти в  [servers-task.md](https://github.com/PolinaShirinskaya/wallet-and-stats-servers/blob/main/servers-task.md)
>
В данном проекте реализованы два веб сервера на Golang - `wallet`, `stats` c `REST API` интерфейсом взаимодействия. 

- `wallet`: сервис отвечающий за взаимодействие с кошельками - создание новых, операции со средствами на кошельке
    - Публикует события в **Kafka** после ряда операций.
- `stats`: сервис сбора статистики по кошелькам. Собирает статистику по действиям с самими кошельками (таким как создание), так и по самим операциям (например пополнение счета).
    - Потребляет опубликованные события сервисом `wallet`.

- В качества брокера очередей используется **Kafka**
- Сервер написан с помощью стандартной библиотеки `"net/http"`
- Для Kafka используется библиотека `"github.com/Shopify/sarama"`

## Запуск проекта:
**Для начала необходимо запустить **Kafka** с помощью Docker**
```
docker-compose up -d
```
**Для запуска проекта нам потребуется 4 окна терминала:**

**1. Для запуска сервера `wallet` (producer Kafka)**
```
//в каталоге wallet_producer
$>go run main.go
```

**2. Для запуска сервера `stats` (consumer Kafka)**
```
//в каталоге stats_consumer
go run main.go
```

**3. Для отправки HTTP-запросов на сервер `wallet` с помощью скрипта [request.sh](https://github.com/PolinaShirinskaya/wallet-and-stats-servers/blob/main/wallet_producer/requests.sh)**

При запуске скрипта Вы можете получить список возможных команд:
```
bash request.sh

//вывод команды
Please enter the command:
'POST' - create wallets
'GET ALL' - get list of all wallets
'GET' + 'wallet ID' - get wallet by ID
'DELETE' + 'wallet ID' - delete wallet by ID
'PUT' + 'wallet ID' + 'new name' - update wallet name by ID
'POST DEPOSIT' + 'wallet ID' - deposit wallet by ID
'POST WITHDRAW' + 'wallet ID' - withdraw wallet by ID
'POST DEPOSIT' + 'wallet ID' + 'transfer to ID' - transfer by IDs

```

Например, для POST-запроса на создание кошелька нам понадобится:
```
bash request.sh "POST"
```

**4. Для отправки GET запроса на сервер `stats` с помощью скрипта [request.sh](https://github.com/PolinaShirinskaya/wallet-and-stats-servers/blob/main/stats_consumer/requests.sh)**
```
bash request.sh "GET"
```
---
---
## Как это будет выглядеть
![This is a alt text.](/images/working_servers.png "This is workin servers.")
