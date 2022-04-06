## TODO:

- [x] починить VPC в serverless
- [x] CI в github (+ добавить секреты)
- [ ] В личке добавить inline кнопки Newer/Older для того чтобы получать более новые/старые сообщения по выбранной локации  
  - [x] сделать чтобы кнопка Newer отображалась только когда есть более новые сообщения
  - [x] сделать чтобы кнопка Older отображалась только когда есть более старые сообщения
  - [ ] не отправлять каждый раз новое сообщение, а обновлять существующее
- [ ] сделать inline режим для бота, чтобы `@bot Ber` подсказывал результаты
- [ ] добавить нормальное описание (что делает бот, какие есть команды, примеры)
  - синтаксис команд для лички и для группы 
- [ ] бонус: если `/start` команада вызвана из канала - выводить подробное описание 1 раз

## How to run locally

```shell
scripts/run_with_env.sh .env-local go run main.go
```

where `.env-local` is the file with the following env variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=<db user name>
DB_PASSWORD=<db password>
DB_NAME=volunteers_bot
GROUP_CHAT_ID=<main group chat ID in telegram>
BOT_TOKEN=<bot token>
BOT_NAME=<bot name>
```
