# СЕРВЕР ЁП 😎🔥🚀

   Проект создан, чтобы набить руку (или своё многострадальное лицо) в процессе написания кода. Выполняет примитивную функцию заметок и ежедневника. В процессе разработки, возможно, часть моего мозга испарилась, и я стал менее эмоциональным. Думаю, это плюс. 😅💻📚

# ЗАДАНИЕ СО ЗВЕЗДОЧКОЙ 🌟🌌💥

   Пока еще не сделано, но в процессе. Звёзды упадут и всё полыхнёт, но пока тут зияющая дыра, наполненная хаосом и надеждой на лучшее. 🎭🌀✨

# ЗАПУСК ШОЛОПАЯ 🚀🛠️🔧

1. **Установите зависимости**

   ```bash
   go mod tidy
   ```

   Убедитесь, что у вас установлен Go версии 1.18 или выше. Иначе сервер обидится и не запустится. 😤⚙️📂

2. **Запустите сервер**

   ```bash
   go run main.go
   ```

   По умолчанию сервер стартанёт на порту **7540**, но вы можете настроить свой кастомный порт через переменную окружения `TODO_PORT`. Пример:

   ```bash
   TODO_PORT=8080 go run main.go
   ```

3. **Откройте браузер**
   Шлёпните в адресную строку:

   ```
   http://localhost:7540
   ```

   И наслаждайтесь своим новым ежедневником. 🌐🎉📖

# ТЕСТЫ В БОЮ 🧪✅🎯

Чтобы убедиться, что всё работает как часы:

1. Вбейте следующую команду для запуска тестов:
   ```bash
   go test ./tests
   ```
2. В настройках `tests/settings.go` есть несколько интересных переменных:
   - `Port`: порт сервера для тестов (по умолчанию 7540).
   - `DBFile`: путь до базы данных. Вы можете использовать свой SQLite файл.

Убедитесь, что тесты зелёные, иначе придётся рыдать вместе с сервером. 😭🔍💔

# СБОРКА ЧЕРЕЗ ДОКЕР (ЕСЛИ ЕСТЬ) 🐳🛳️✨

Если вы настроили Docker (а если нет, пора начать):

1. Соберите образ:
   ```bash
   docker build -t todo-server .
   ```
2. Запустите контейнер:
   ```bash
   docker run -d -p 7540:7540 todo-server
   ```
3. Откройте браузер и введите:
   ```
   http://localhost:7540
   ```

Поздравляю, вы официально задокерили этот шедевр! 🎊🤖⚡

---

С этим README ваш проект становится эпическим. Ревьюверы, готовьтесь падать в обморок от его величия! 🌟👑💥

