package api
// 1️⃣ ШО, ОПЯТЬ ГОЛЫЙ package?!!? 😱
// Да, бро, это святая обязанность Go-программиста — сказать «package api»,
// чтобы гонять наших прекрасных функций в рамках пространства имён "api".
// Иначе Go взорвётся рандомным выстрелом из лазерного кота! 🐱🔫

import (
    "encoding/json"
    // 2️⃣ ОГО, КУДА МЫ ПОПАЛИ? В ДЕБРИ JSON! 🤯
    // Этот пакет помогает нам общаться с клиентами на языке JSON,
    // превращая нашу Go-структуру в унылый текст и обратно. 

    "fmt"
    // 3️⃣ Разрешите представиться: пакет "fmt". Звезда форматирования и лихих прибауток. 🤡
    // Если вам нужна строка по типу "Заголовок: %s" — сюда, вот он вас оформит по красоте.

    "log"
    // 4️⃣ Пакет "log" — ваш лучший друг, когда всё идет **не** по плану. 😈
    // Впихиваешь туда ошибку — он тебе любезно выворачивает её наизнанку
    // вместе с конфетти и фейерверками. 🎉

    "net/http"
    // 5️⃣ Мать его, да это же "net/http"! 🚀
    // Короче, без этого жить было бы скучно, ведь именно он поднимает наши HTTP-сервера.
    // Как говорится: "Hello, server! Hello, client! Let's dance! 💃".

    "time"
    // 6️⃣ "time" — лучший чувак для путешествий во времени, чтобы улететь на 3 дня назад и вернуться к дедлайну. 🔮
    // Будем использовать, чтобы не забыть, что сегодня уже завтра.

    "database/sql"
    // 7️⃣ ОПА! "database/sql"? 🔥
    // Тут, как говорится, всё серьёзно. Этим API мы завязываем нашу жизнь с SQL-базами.
    // Запросы, коннекты, стаканы чая... всё в одном флаконе.

    "github.com/naluneotlichno/FP-GO-API/database"
    // 8️⃣ Щас будет магия локальных пакетов. 🎩
    // "database" — там где-то наша волшебная функция GetDB()
    // (наверное, она вызовет единорога, чтобы нас связать с MySQL или ещё чем).

    "github.com/naluneotlichno/FP-GO-API/api"
    // 9️⃣ Ну всё, приготовься, — "nextdate" ворвётся и перескочит тебя во времени.
    // Если дедлайн прошел вчера, он телепортирует задачу на послезавтра. 🤷
)

// TaskRequest — структура входного JSON-запроса
type TaskRequest struct {
    Date    string `json:"date"`
    // 😺 Даже не знаю, что тут может пойти не так. 😼
    // "date" — строка. Если пустая — будем плясать вокруг костра и ставить сегодняшнюю дату.

    Title   string `json:"title"`
    // 😎 "title" — наше всё. Без него жизнь тлен. Поле обязательно! Пустое — не взлетим.

    Comment string `json:"comment,omitempty"`
    // 🤫 comment: если нет комментариев — и не надо... 😉
    // Поле не попадёт в JSON, если пустое (спасибо "omitempty").

    Repeat  string `json:"repeat,omitempty"`
    // 🔁 Если нужно повторять задачу снова и снова, запихай сюда "DAILY", "WEEKLY" или что угодно.
    // И тогда процесс никогда не кончится! 😈
}

// TaskResponse — структура ответа (id или ошибка)
type TaskResponse struct {
    ID    int64  `json:"id,omitempty"`
    // 😍 Ах, наш дорогой ID! Если запись в БД хорошо зашла, будет тут.

    Error string `json:"error,omitempty"`
    // 💀 Если мир рушится и SQL нас предал, тут пишем "Error: всё плохо, уходите". 
}

// AddTaskHandler обрабатывает POST-запросы на /api/task
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
    // 1️⃣ Здесь мы устраиваем гениальную засаду:
    // принимаем клиента, проверяем, постучался ли он POST'ом,
    // и творим магию вставки задачи.

    if r.Method != http.MethodPost {
        // 🔨 Если мужик не "POST", то мы его кувалдой по 405 Method Not Allowed! 🔥
        http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
        // 💅 При этом говорим "сорян, брат, мы тут только POST-методы любим". 
        return
    }

    // ✅ Декодируем JSON-запрос
    var req TaskRequest
    // Это наши почтовые ящики для данных от клиента.

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // 🏆 ОЙ-ОЙ, НЕ СМОГЛИ РАСПАРСИТЬ JSON?
        log.Printf("❌ Ошибка десериализации JSON: %v", err)
        // Логируем позор на всю консоль, чтобы потом было, о чём детям рассказать. 😭
        http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
        // Выдаем ошибку 400 и говорим: "Не умеешь JSON писать — не говори с нами".
        return
    }

    // ✅ Проверяем обязательные поля
    if req.Title == "" {
        // 😂 "title" пустой? Это как выйти на улицу голым! 
        // Не ну правда, как так-то без заголовка?
        http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
        // Жалуемся клиенту, что "Title" — наше святое. 👼
        return
    }

    // ✅ Если дата пустая — подставляем текущую
    if req.Date == "" {
        // ⏳ Если "date" не задали, мы юзаем сегодняшнюю.
        // Это как прийти в клуб без даты — берёшь случайную прямо на месте! 🤪
        req.Date = time.Now().Format("20060102")
    }

    // ✅ Парсим дату, если формат кривой — шлём ошибку
    taskDate, err := time.Parse("20060102", req.Date)
    // 😏 time.Parse — чувак суровый, если не видит 8 цифр (YYYYMMDD), он пошлёт нас к чёрту.

    if err != nil {
        // 🥴 Видимо формат кривой: "2025-01-16" или "ололо123"? Мимо кассы!
        http.Error(w, `{"error": "Дата указана некорректно"}`, http.StatusBadRequest)
        // Гоним взашей клиента, ибо не умеет в нашу дату!
        return
    }

    // ✅ Если дата в прошлом — применяем правило повторения
    if taskDate.Before(time.Now()) {
        // ⏲ .Before() проверяет, стоит ли дата до текущего момента.
        // Если да, значит «ГГ», мы просрочили дедлайн. 😱

        if req.Repeat != "" {
            // 🌀 Если есть правило повторения, крутим наш волчок дальше!
            nextDate, err := nextdate.NextDate(time.Now(), req.Date, req.Repeat)
            // 🔮 "NextDate" из пакета "nextdate" магическим образом вычислит, 
            // куда переносить просроченную задачу (например, на завтра).

            if err != nil {
                // 💔 Если формат повтора нам не понравился, бахаем 400.
                http.Error(w, fmt.Sprintf(`{"error": "Неверный формат правила повторения: %s"}`, err.Error()), http.StatusBadRequest)
                // Всё, до свидульки.
                return
            }
            req.Date = nextDate
            // 🤙 Вот и всё, у нас новенькая дата в будущем. Начинаем тусоваться заново!

        } else {
            // 🤷 Если повтора нет, а дата ушла в прошлое — значит ставим «сегодня» и не морочим голову.
            req.Date = time.Now().Format("20060102")
        }
    }

    // ✅ Подключаемся к базе данных
    db, err := database.GetDB()
    // 🔑 Вызываем нашу волшебную функцию GetDB() из пакета "database". 
    // Надеемся, что откроется портал в мир SQL.

    if err != nil {
        // 🔒 Если портал закрыт, значит всё, взрываемся 500 (Internal Server Error).
        http.Error(w, `{"error": "Ошибка подключения к БД"}`, http.StatusInternalServerError)
        return
    }

    // ✅ Вставляем новую задачу в базу
    query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
    // SQL-запрос: "Сэр, положите дату, заголовок, комментарий и повтор в таблицу scheduler!"
    // Иначе все взорвётся! 💣

    res, err := db.Exec(query, req.Date, req.Title, req.Comment, req.Repeat)
    // Вот она: шик, блеск, красота. Отправляем запрос в базу:
    // "Привет, база, дай нам новую задачу!"

    if err != nil {
        // 🤬 Если база заявила: "Фиг тебе, а не INSERT", 
        // тогда возвращаем код 500 и валим в закат.
        http.Error(w, `{"error": "Ошибка записи в БД"}`, http.StatusInternalServerError)
        return
    }

    // ✅ Получаем ID новой задачи
    taskID, err := res.LastInsertId()
    // lastInsertId() — это как спрашивать у базы: "Ну и какой у меня теперь ID, шеф?" 
    // Она такая: "Ну держи, 108". 

    if err != nil {
        // 🙈 Если база сказала "не знаю я никаких ID",
        // выдаём 500, ибо беда.
        http.Error(w, `{"error": "Ошибка получения ID записи"}`, http.StatusInternalServerError)
        return
    }

    // ✅ Возвращаем JSON-ответ в формате, который ожидает тест
    resp := TaskResponse{ID: taskID}
    // Создаём ответ, где ID = волшебное число, которое нам вернула база. 🔮

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    // Игриво шепчем клиенту в заголовках, что мы будем говорить на языке JSON в UTF-8. 🦜

    json.NewEncoder(w).Encode(resp)
    // И, наконец, кодируем нашу структуру в JSON и впихиваем в ответ клиенту.
    // БАМ! Всё, задачка записана, ID отправлен. Профит! 🍾
}