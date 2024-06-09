# GopherVault
"GopherVault: Go Secret Sentinel"  client-server system that allows the user to reliably and securely store logins, passwords, binary data and other private information. Implemented as part of the course from Yandex Practicum "Advanced Go Developer"


Установка дополнительного ПО и сборка в Windows

Если вы хотите запустить команды из Makefile в среде Windows 11, у вас есть несколько вариантов. Один из самых популярных и удобных способов — использование утилиты Windows Subsystem for Linux (WSL), которая позволяет запускать Linux-команды непосредственно на Windows. Рассмотрим подробнее шаги этого подхода, а также альтернативные методы.

 Использование WSL
Установка WSL и Linux-дистрибутива:
Откройте PowerShell от имени администратора и выполните команду для установки WSL и дистрибутива Ubuntu:
powershell
Копировать
wsl --install
Копировать
- Подождите, пока установка завершится, и следуйте инструкциям на экране для настройки Ubuntu.
  Установка Make:
  После установки и настройки Ubuntu откройте WSL-терминал и выполните следующие команды для обновления пакетного менеджера и установки Make:
  bash
  Копировать
  sudo apt update
  sudo apt install build-essential
  Работа с Makefile:
  Перейдите в директорию с вашим проектом (где находится Makefile), используя команду cd. Например:
  bash
  Копировать
  cd /mnt/c/путь/к/проекту
  Копировать
- Запустите Makefile, используя `make`. Например:
  bash
  Копировать
  make

Способ 2: Использование MSYS2
Установка MSYS2:
Скачайте и установите MSYS2 с официального сайта: https://www.msys2.org/.
Запустите MSYS2 и обновите пакетный менеджер и установите необходимые утилиты:
bash
Копировать
pacman -Syu
pacman -S base-devel
Пользование MSYS2:
Откройте MSYS2-терминал и перейдите в директорию с Makefile, используя команду cd. Например:
bash
Копировать
cd /c/путь/к/проекту
Копировать
- Запустите Makefile, используя `make`. Например:
  bash
  Копировать
  make
  Способ 3: Использование MinGW-w64 или Cygwin
  Установка MinGW-w64:
  Скачайте и установите MinGW-w64 с официального сайта или через https://www.msys2.org/.
  Установите необходимые пакеты для поддержки make, gcc и других утилит:
  bash
  Копировать
  pacman -S mingw-w64-x86_64-gcc
  pacman -S mingw-w64-x86_64-make
  Или установка Cygwin:
  Скачайте и установите Cygwin с сайта https://www.cygwin.com/.
  Во время установки выберите пакет make для установки.
  Запуск Makefile с помощью MinGW-w64 или Cygwin:
  Откройте соответствующий терминал и перейдите в директорию с вашим Makefile.
  Используйте команду make, как и в предыдущих способах:
  bash
  Копировать
  make
  Заключение
  Использование WSL, MSYS2, MinGW-w64 или Cygwin позволяет запускать команды из Makefile в среде Windows 11 так же, как это делается в Unix-подобных системах. Каждый из этих способов имеет свои преимущества, и выбор зависит от ваших предпочтений и окружения разработки.


  Убедитесь, что Docker Desktop интегрирован с WSL и docker-compose правильно установлен. Ниже приведены шаги для выполнения этой задачи:

Шаг 1: Убедитесь, что Docker Desktop установлен и работает
Скачайте и установите Docker Desktop: Если Docker Desktop еще не установлен, его можно скачать и установить с https://www.docker.com/products/docker-desktop.
Запустите Docker Desktop: Убедитесь, что Docker Desktop запущен. Проверьте, что иконка Docker отображается в системном трее и статус Docker — "Running".
Шаг 2: Включите интеграцию WSL 2 с Docker Desktop
Откройте Docker Desktop: Нажмите правой кнопкой мыши по иконке Docker в системном трее и выберите "Settings".
Перейдите в раздел Resources: В меню слева выберите "Resources" -> "WSL Integration".
Включите WSL 2 интеграцию: Поставьте галочку "Enable integration with my default WSL distro" и выберите нужную WSL 2 дистрибуцию в списке.
Сохраните изменения: Нажмите "Apply & Restart".
Шаг 3: Убедитесь, что docker-compose доступен в WSL 2
Проверьте наличие docker-compose: В WSL терминале выполните команду:
bash
Копировать
docker-compose --version
Если команда выполнится успешно и выведет версию docker-compose, значит все в порядке.

Установка docker-compose (опционально): Если команда docker-compose не найдена, выполните установку:
bash
Копировать
sudo apt update
sudo apt install docker-compose
Шаг 4: Запуск docker-compose up --detach
Переходите в директорию вашего проекта и запустите контейнеры в фоновом режиме с docker-compose:

bash
Копировать
cd /mnt/c/путь/к/проекту
docker-compose up --detach
Дополнительные советы
Обновление Docker Desktop: Иногда ошибки могут быть вызваны устаревшей версией Docker Desktop. Убедитесь, что у вас установлена последняя версия.
Проверка интеграции: Вы можете также проверить, что Docker правильно интегрирован с WSL, выполнив команды docker version и docker info в WSL терминале.
Теперь docker-compose должен работать корректно в вашей WSL 2 среде. Если проблема остается, возможно стоит пересмотреть установки или обратитесь к документации Docker для более детального исследования проблемы.