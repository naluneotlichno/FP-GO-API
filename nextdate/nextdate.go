package nextdate

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 🔥 HandleNextDate обрабатывает запросы на /api/nextdate
func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	log.Println("✅ Запрос на расчет даты получен!")

	// ✅ Извлекаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// ✅ Парсим параметр `now` в формате time.Time
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Некорректная дата now", http.StatusBadRequest)
		return
	}

	// ✅ Вызываем функцию NextDate
	nextDate, err := NextDate(now, dateStr, repeat, "nextdate")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка расчета следующей даты: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// ✅ Возвращаем результат клиенту
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

// NextDate вычисляет следующую дату задачи на основе правила повторения.
// Возвращает дату в формате `20060102` (YYYYMMDD) или ошибку, если правило некорректно.
func NextDate(now time.Time, dateStr string, repeat string, status string) (string, error) {
    if dateStr == "" {
        return "", errors.New("дата не может быть пустой")
    }

    parsedDate, err := time.Parse("20060102", dateStr)
    if err != nil {
        return "", fmt.Errorf("ошибка разбора даты: %w", err)
    }

    if repeat == "" {
        if parsedDate.After(now) {
            return parsedDate.Format("20060102"), nil
        }
        return "", errors.New("дата в прошлом и правило повторения не задано")
    }

    if strings.HasPrefix(repeat, "d ") {
        daysStr := strings.TrimPrefix(repeat, "d ")
        days, err := strconv.Atoi(daysStr)
        if err != nil || days < 1 || days > 400 {
            return "", errors.New("некорректное значение дней")
        }

        nextDate := parsedDate.AddDate(0, 0, days)
        for !nextDate.After(now) {
            nextDate = nextDate.AddDate(0, 0, days)
        }

        return nextDate.Format("20060102"), nil
    }

    if strings.HasPrefix(repeat, "w ") {
        weeksStr := strings.TrimPrefix(repeat, "w ")
        weeks, err := strconv.Atoi(weeksStr)
        if err != nil || weeks < 1 || weeks > 52 {
            return "", errors.New("некорректное значение недель")
        }

        nextDate := parsedDate.AddDate(0, 0, weeks*7)
        for !nextDate.After(now) {
            nextDate = nextDate.AddDate(0, 0, weeks*7)
        }

        return nextDate.Format("20060102"), nil
    }

    if strings.HasPrefix(repeat, "m ") {
        monthsStr := strings.TrimPrefix(repeat, "m ")
        months, err := strconv.Atoi(monthsStr)
        if err != nil || months < 1 || months > 120 {
            return "", errors.New("некорректное значение месяцев")
        }

        nextDate := parsedDate.AddDate(0, months, 0)
        for !nextDate.After(now) {
            nextDate = nextDate.AddDate(0, months, 0)
        }

        return nextDate.Format("20060102"), nil
    }

    if repeat == "y" {
        nextDate := parsedDate.AddDate(1, 0, 0)
        if parsedDate.Month() == time.February && parsedDate.Day() == 29 {
            if nextDate.Month() != time.February || nextDate.Day() != 29 {
                nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, nextDate.Location())
            }
        }

        for !nextDate.After(now) {
            nextDate = nextDate.AddDate(1, 0, 0)
        }

        return nextDate.Format("20060102"), nil
    }

    log.Printf("❌ Некорректный формат repeat: %s", repeat)
    return "", errors.New("неподдерживаемый формат повторения")
}

