# GopherVault

"GopherVault: Go Secret Sentinel" — клиент-серверная система, позволяющая пользователю надежно и безопасно хранить логины, пароли, двоичные данные и другую конфиденциальную информацию. Реализована как часть курса от Yandex Practicum "Advanced Go Developer".

## Общее устройство

- Клиент распространяется в виде CLI-приложения;
- В качестве хранилища данных используется PostgreSQL;
- Клиент и сервер обмениваются данными по HTTP-протоколу;
- Чувствительные данные хранятся в зашифрованном виде;
- Механизм конфигурируется через следующие переменные окружения:
  - `POSTGRES_HOST` - хост хранилища
  - `POSTGRES_PORT` - порт хранилища
  - `POSTGRES_USER` - пользователь ` GopherVault `
  - `POSTGRES_PASSWORD` - пароль пользователя ` GopherVault `
  - `POSTGRES_DB` - имя базы данных, в которой хранится вся пользовательская информация;
  - `APPLICATION_PORT` - порт приложения ` GopherVault `
  - `APPLICATION_HOST` - хост приложения ` GopherVault `
  - `KEEPER_ENCRYPTION_KEY` - ключ для шифрования чувствительной информации
- В хранилище ` GopherVault ` существуют следующие системные таблицы:
  - `registered_users` - таблица пользователей, зарегистрированных в ` GopherVault `
  - `credentials` - таблица с сохраненными логинами/паролями пользователей. Каждый пользователь
    через приложение может получить только свои логины/пароли. Пароли хранятся в зашифрованном виде
  - `notes` - таблица, в которой хранится произвольная пользовательская информация - различные
    заметки, бинарные данные etc. Все содержимое хранится в зашифрованном виде. Каждый пользователь
    через приложение может получить только свои данные
  - `cards` - данные банковских карт: имя банка, номер карты, cv-код, пароль от банковского приложения.
    CV и пароли хранятся в зашифрованном виде. Каждый пользователь через приложение может получить данные
    только своих карт

## Cхема взаимодействия с системой

**Для нового пользователя**:

- Пользователь получает клиент под необходимую ему платформу
- Пользователь проходит процедуру первичной регистрации
- Пользователь добавляет в клиент новые данные
- Клиент синхронизирует данные с сервером

**Для существующего пользователя**:

- Пользователь получает клиент под необходимую ему платформу
- Пользователь проходит процедуру аутентификации
- Клиент синхронизирует данные с сервером
- Пользователь запрашивает данные
- Клиент отображает данные для пользователя

## Установка приложения для своей платформы

- Склонировать репозиторий
- Собрать приложение

    ```shell
    make install
    ```

- Остановить приложение

    ```shell
    make stop
    ```

Команда `make install` поднимает PostgreSQL в докере, применяет необходимые для работы сервиса миграции,
а затем запускает HTTP сервер приложения, который начинает принимать запросы.

## Установка дополнительного ПО и сборка в Windows

Если вы хотите запустить команды из Makefile в среде Windows 11, у вас есть несколько вариантов. Один из самых популярных и удобных способов — использование утилиты Windows Subsystem for Linux (WSL), которая позволяет запускать Linux-команды непосредственно на Windows. Рассмотрим подробнее шаги этого подхода.

### Использование WSL

#### Установка WSL и Linux-дистрибутива:

1. Откройте PowerShell от имени администратора и выполните команду для установки WSL и дистрибутива Ubuntu:
```powershell
   wsl --install
```


2. Подождите, пока установка завершится, и следуйте инструкциям на экране для настройки Ubuntu.
  Установка Make:
  После установки и настройки Ubuntu откройте WSL-терминал и выполните следующие команды для обновления пакетного менеджера и установки Make:
  ```bash
  sudo apt update
  sudo apt install build-essential
  ```
  
3. Работа с Makefile:
Перейдите в директорию с вашим проектом (где находится Makefile), используя команду cd. Например:
 ```bash
  cd /mnt/c/путь/к/проекту
```
- Запустите Makefile, используя `make`. Например:
 ```bash
  make
```
### Docker Desktop 
Убедитесь, что Docker Desktop интегрирован с WSL и docker-compose правильно установлен. Ниже приведены шаги для выполнения этой задачи:

#### Шаг 1: Убедитесь, что Docker Desktop установлен и работает
Скачайте и установите Docker Desktop: Если Docker Desktop еще не установлен, его можно скачать и установить с https://www.docker.com/products/docker-desktop.
Запустите Docker Desktop: Убедитесь, что Docker Desktop запущен. Проверьте, что иконка Docker отображается в системном трее и статус Docker — "Running".
#### Шаг 2: Включите интеграцию WSL 2 с Docker Desktop
Откройте Docker Desktop: Нажмите правой кнопкой мыши по иконке Docker в системном трее и выберите "Settings".
Перейдите в раздел Resources: В меню слева выберите "Resources" -> "WSL Integration".
Включите WSL 2 интеграцию: Поставьте галочку "Enable integration with my default WSL distro" и выберите нужную WSL 2 дистрибуцию в списке.
Сохраните изменения: Нажмите "Apply & Restart".
#### Шаг 3: Убедитесь, что docker-compose доступен в WSL 2
Проверьте наличие docker-compose: В WSL терминале выполните команду:
```
docker-compose --version
```
Если команда выполнится успешно и выведет версию docker-compose, значит все в порядке.

Установка docker-compose (опционально): Если команда `docker-compose` не найдена, выполните установку:
```
sudo apt update
sudo apt install docker-compose
```
#### Шаг 4: Запуск docker-compose up --detach
Переходите в директорию вашего проекта и запустите контейнеры в фоновом режиме с docker-compose:

```
cd /mnt/c/путь/к/проекту
docker-compose up --detach
```
Дополнительные советы
Обновление Docker Desktop: Иногда ошибки могут быть вызваны устаревшей версией Docker Desktop. Убедитесь, что у вас установлена последняя версия.
Проверка интеграции: Вы можете также проверить, что Docker правильно интегрирован с WSL, выполнив команды docker version и docker info в WSL терминале.
Теперь docker-compose должен работать корректно в вашей WSL 2 среде. Если проблема остается, возможно стоит пересмотреть установки или обратитесь к документации Docker для более детального исследования проблемы.



## Возможности приложения

TLDR: для каждой команды доступна справка с примерами использования.

**Посмотреть дату сборки приложения**

```shell
GopherVault build-date
```

**Посмотреть версию приложения**

```shell
GopherVault --version
```

**Регистрация в приложении**

```shell
GopherVault register --login <user-system-login> --password <user-system-password>
```

**Вход в приложение**

```shell
GopherVault login --login <user-system-login> --password <user-system-password>
```

**Добавить данные о банковской карте**

```shell
GopherVault  add-card --user <user-system-login> --bank <bank-name> --number <card-number> --cv <card-cv> --password <password>
```

Номер карты должен содержать 16 знаков, cv - 3 знака. Можно добавить метаинформацию о карте:

```shell
GopherVault  add-card --user <user-system-login> --bank <bank-name> --number <card-number> --cv <card-cv> --password <password> --metadata <some metadata>
```

**Добавить логин/пароль**

```shell
GopherVault add-credentials --user <user-name> --login <user-login> --password <password to store> --metadata <some description>
```

**Добавить произвольную текстовую информацию**

```shell
GopherVault add-note --user <user-name> --title <note title> --content <note content> --metadata <note metadata>
```

**Удалить логин/пароль**

```shell
GopherVault delete-credentials --user <user-name> --login <user-login>
```

Можно удалить все сохраненные пары логин/пароль для пользователя, не указывая конкретный логин:

```shell
GopherVault delete-credentials --user <user-name>
```

**Удалить произвольную информацию**

```shell
GopherVault delete-note --user <user-name> --title <note title>
```

Можно удалить все данные для пользователя, если не указывать идентификатор данных:

```text
GopherVault delete-note --user <user-name>
```

**Удалить данные банковских карт**

Эта команда удалит все данные карт банка `<bank-name>` пользователя `<user-name>`

```shell
GopherVault  delete-card --user <user-name> --bank <bank-name>
```

Эта команда удалит данные карты с номером `<card-number>` пользователя `<user-name>`

```shell
GopherVault  delete-card --user <user-name> --number `<card-number>`
```

**Получить сохраненные пары логин/пароль**

```shell
GopherVault get-credentials --user <user-name>
```

Можно получить информацию для конкретного логина:

```text
GopherVault get-credentials --user <user-name> --login <login>
```

**Получить сохраненные произвольные данные**

```shell
GopherVault get-note --user <user-name>
```

Можно получить информацию для конкретного идентификатора:

```text
GopherVault get-credentials --user <user-name> --title <note title>
```

**Получить данные сохраненных банковских карт**

```shell
GopherVault get-card --user <user-name>
```

Можно получить информацию по картам конкретного банка:

```text
GopherVault get-credentials --user <user-name> --bank <bank-name>
```

Можно получить информацию по карте с конкретным номером:

```text
GopherVault get-credentials --user <user-name> --number <card-number>
```

**Изменить пароль для сохраненного логина**

```text
GopherVault update-credentials --user <user-name> --login <saved-login> --password <new-password>
```

**Отредактировать сохраненные произвольные данные**

```text
GopherVault update-notes --user <user-name> --title <note-title> --content <new-content>
```
