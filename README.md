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
- [ ] добавить нормальное описание (что делает бот, какие есть команды, примеры)
  - [ ] синтаксис команд для лички и для группы 
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
This bot collects information about accomodation in cities, regions and countries of Europe. To get the latest send a city, country or region to the bot.
```

### Описание команд

#### В группе
```
/start, /help


```

#### В личке
```
/start, /help

Send the name of a city, country or region to the bot(for example: Berlin), and it will return the latest information about
accommodation and conditions 

Use Older/Newer buttons to access older updates.  
 
For example:
`Berlin`
```

```
ответ на локацию:


```

```
ответ на локацию если ничего не найдено:

**Berlin**
Нічого не знайдено, спробуйте інше місто, регіон або країну.
Ничего не найдено, попробуйте другой город, регион или страну.
Nothing found, try another city, region or country.
```

### Другие настройки

* должен иметь доступ к сообщениям в группах
* возможно должен иметь inline режим
