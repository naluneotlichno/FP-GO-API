package nextdate

import (
    "fmt"     // 📜 Работаем со строками, форматируем данные
    "strconv" // 🔢 Конвертация строк в числа
    "strings" // ✂️ Работа со строками (разделение, удаление пробелов и т. д.)
    "time"    // ⏳ Работа с датами и временем
    "log"     // 📝 Работаем с логами
)

// NextDate вычисляет следующую дату задачи на основе правила повторения.
// Возвращает дату в формате `20060102` (YYYYMMDD) или ошибку, если правило некорректно.
func NextDate(now time.Time, date string, repeat string) (string, error) {
    log.Printf("🔍 [DEBUG] Входные данные: now=%s, date=%s, repeat=%s",
        now.Format("20060102"), date, repeat)

    // 1. Парсим входную дату (date)
    parsedDate, err := time.Parse("20060102", date)
    if err != nil {
        return "", fmt.Errorf("❌ Ошибка: Некорректная дата '%s'", date)
    }
    log.Printf("✅ [DEBUG] Парсинг даты успешен: %s", parsedDate.Format("20060102"))

    // 2. Проверяем, указано ли правило повторения
    if repeat == "" {
        return "", fmt.Errorf("❌ Ошибка: Задача не повторяется, можно удалить")
    }

    // Определяем, прошла ли дата или она сегодня (parsedDate <= now)
    // Если parsedDateAfterNow == true, значит parsedDate > now
    parsedDateAfterNow := parsedDate.After(now)

    // --- Обработка различных типов repeat ---
    // 1) Ежегодное повторение: repeat = "y"
    if repeat == "y" {
        nextDate := parsedDate

        if !parsedDateAfterNow {
            // Если дата <= now, крутим +1 год, пока не станет > now
            for !nextDate.After(now) {
                nextDate = nextDate.AddDate(1, 0, 0)
            }
        } else {
            // Если исходная дата > now, просто делаем +1 год один раз
            nextDate = nextDate.AddDate(1, 0, 0)
        }

        log.Printf("✅ [DEBUG] Ежегодное повторение! Следующая дата: %s", nextDate.Format("20060102"))
        return nextDate.Format("20060102"), nil
    }

    // 2) Повтор через N дней: repeat = "d N"
    if strings.HasPrefix(repeat, "d ") {
        parts := strings.Split(repeat, " ")
        if len(parts) != 2 {
            return "", fmt.Errorf("❌ Ошибка: Неверный формат правила '%s'", repeat)
        }

        days, err := strconv.Atoi(parts[1])
        if err != nil || days < 1 || days > 400 {
            return "", fmt.Errorf("❌ Ошибка: Некорректное количество дней '%s'", parts[1])
        }

        nextDate := parsedDate

        if !parsedDateAfterNow {
            // Если исходная дата <= now, крутим по days, пока не станет > now
            for !nextDate.After(now) {
                nextDate = nextDate.AddDate(0, 0, days)
            }
        } else {
            // Если исходная дата > now, делаем одну прибавку на days
            nextDate = nextDate.AddDate(0, 0, days)
        }

        log.Printf("✅ [DEBUG] Повтор каждые %d дней. Следующая дата: %s", days, nextDate.Format("20060102"))
        return nextDate.Format("20060102"), nil
    }

    // 3) Если правило не поддерживается
    return "", fmt.Errorf("❌ Ошибка: Неподдерживаемый формат повторения '%s'", repeat)
}
