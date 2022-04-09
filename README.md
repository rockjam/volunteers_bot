## TODO:

- [x] починить VPC в serverless
- [x] CI в github (+ добавить секреты)
- [ ] В личке добавить inline кнопки Newer/Older для того чтобы получать более новые/старые сообщения по выбранной локации  
  - [x] сделать чтобы кнопка Newer отображалась только когда есть более новые сообщения
  - [x] сделать чтобы кнопка Older отображалась только когда есть более старые сообщения
  - [ ] не отправлять каждый раз новое сообщение, а обновлять существующее
- [x] оставить только команды /start и /help
- [x] отвечать на /start и /help. В остальных случаях воспринимать текст как локацию
- [x] отвечать ошибкой когда неизвестная / команда
- [x] добавить нормальное описание (что делает бот, какие есть команды, примеры)
  - [x] синтаксис команд для лички и для группы 
- [ ] сделать inline режим для бота, чтобы `@bot Ber` подсказывал результаты
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

## Конфигурация бота

### Название бота

Digital Volunteers Arrivals
`@dv_arrivals_bot`

Test Digital Volunteers Arrivals
`@dv_arrivals_test_bot`

### Команды
```
start - як користуватися/как пользоваться/how to use
help - як користуватися/как пользоваться/how to use
```

### About/О боте
```
Надає інформацію про розміщення біженців з України в ЄС
Предоставляет информацию о размещении беженцев из Украины в ЕС 
```

### Описание

```
Цей бот збирає інформацію про можливості розміщення в містах, регіонах і країнах Європи. Щоб дізнатися актуальну інформацію, відправте боту місто, регіон або країну.
Этот бот собирает информацию о возможностях размещения в городах, регионах и странах Европы. Чтобы узнать актуальную информацию, отправьте боту город, регион или страну.
This bot collects information about accommodation in cities, regions and countries of Europe. To get the latest send a city, country or region to the bot.
```

### Описание команд

#### В группе
```
/start, /help

Hello from <b>Digital Volunteers Arrivals Bot</b>.
It collects information about accommodation in cities, regions and countries of Europe shared in this group.
When you share an update, don't forget to include hashtags for location, e.g: #Germany #Berlin

Write a DM to @dv_arrivals_bot to browse all updates.
```

#### В личке
```
/start, /help

<b>UKR:</b> Надішліть місто, регіон або країну боту (наприклад: Берлін), щоб отримати актуальну інформацію про можливості та умови розміщення там. Використовуйте кнопки Older/Newer, щоб отримати старі/нові оновлення.
<b>RU:</b> Отправьте город, регион или страну боту(например: Берлин), чтобы получить актуальную информацию о возможностях и условиях размещения там. Используйте кнопки Older/Newer, чтобы получить старые/новые обновления.    
<b>EN:</b> Send a city, country or region to the bot(e.g.: Berlin) to get the latest updates on accomodation for it. To browse though updates use Older/Newer buttons next to the message.

```

#### Ответы

```
ответ на локацию если ничего не найдено:

<b>Berlin</b>:
<b>UKR:</b> Нічого не знайдено, спробуйте інше місто, регіон або країну.
<b>RU:</b> Ничего не найдено, попробуйте другой город, регион или страну.
<b>EN:</b> Nothing found, try another city, region or country.
```

```
ответ на неизвестную команду:

<b>UKR:</b> Невідома команда. Надішліть місто, регіон або країну, щоб отримати актуальну інформацію про можливості та умови розміщення, або /start щоб дізнатися як користуватися ботом.
<b>RU:</b> Неизвестная команда. Отправьте город, регион или страну, чтобы получить актуальную информацию о возможностях и условиях размещения, или /start чтобы узнать как пользоваться ботом.
<b>EN:</b> Unknown command. Send a city, country or region to get a latest update, or /start for help.
```

### Другие настройки

* должен иметь доступ к сообщениям в группах
* возможно должен иметь inline режим
