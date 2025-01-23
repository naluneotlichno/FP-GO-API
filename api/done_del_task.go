package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/naluneotlichno/FP-GO-API/database"
	"github.com/naluneotlichno/FP-GO-API/nextdate"
)

// DoneTaskHandler обрабатывает POST /api/task/done?id=...
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 [DoneTaskHandler] Запрос на /api/task/done получен...")

	idStr := r.URL.Query().Get("id")
	log.Printf("🔍 [DoneTaskHandler] ID из запроса: %s\n", idStr)
	if idStr == "" {
		log.Println("🚨 [DoneTaskHandler] ID не указан")
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("🚨 [DoneTaskHandler] Ошибка парсинга ID=%s: %v\n", id, err)
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Некорректный идентификатор"})
		return
	}

	task, err := database.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, database.ErrTask) {
			jsonResponse(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
			return
		}
		log.Printf("🚨 [DoneTaskHandler] Ошибка получения задачи ID=%s: %v\n", id, err)
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении задачи"})
		return
	}

	log.Printf("✅ [DoneTaskHandler] Найдена задача: %#v\n", task)

	if task.Repeat == "" {
		log.Printf("🔍 [DoneTaskHandler] repeat пустой. Удаляем задачу ID=%s\n", id)
		if err := database.DeleteTask(id); err != nil {
			log.Printf("🚨 [DoneTaskHandler] Ошибка при удалении задачи ID=%s: %v\n", id, err)
			jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении задачи"})
			return
		}
	} else {
		now := time.Now()
		nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat, "done")
		if err != nil {
			log.Printf("🚨 [DoneTaskHandler] Ошибка вычисления следующей даты: %v\n", err)
			jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при вычислении следующей даты"})
			return
		}

		task.Date = nextDate
		err = database.UpdateTask(task)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении задачи"})
			return
		}
	}

	jsonResponse(w, http.StatusOK, map[string]any{})
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	log.Printf("📤 [jsonResponse] Отправляем ответ: статус=%d, payload=%#v\n", status, payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// DeleteTaskHandler обрабатывает DELETE /api/task?id=...
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 [DeleteTaskHandler] Запрос на DELETE /api/task получен...")

	idStr := r.URL.Query().Get("id")
	log.Printf("🔍 [DeleteTaskHandler] ID из запроса: %s\n", idStr)
	if idStr == "" {
		log.Println("🚨 [DeleteTaskHandler] ID не указан")
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("🚨 [DeleteTaskHandler] Ошибка парсинга ID=%s: %v\n", id, err)
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Некорректный идентификатор"})
		return
	}

	log.Printf("🔍 [DeleteTaskHandler] Пытаемся удалить задачу с ID=%s\n", id)
	if err := database.DeleteTask(id); err != nil {
		if errors.Is(err, fmt.Errorf("задача не найдена")) {
			jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		log.Printf("🚨 [DeleteTaskHandler] Ошибка удаления задачи ID=%s: %v\n", id, err)
		jsonResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("✅ [DeleteTaskHandler] Задача ID=%s успешно удалена\n", id)
	jsonResponse(w, http.StatusOK, map[string]interface{}{})
}

// NextDateAdapter — переходник между твоей NextDate(...) и тем, что ожидают тесты.
func NextDateAdapter(oldDate time.Time, repeat string) (time.Time, error) {
	log.Printf("🔍 [NextDateAdapter] Параметры: oldDate=%s, repeat=%s\n", oldDate.Format("20060102"), repeat)

	repeatParts := strings.Split(repeat, " ")
	if len(repeatParts) != 2 || repeatParts[0] != "d" {
		err := fmt.Errorf("некорректный формат repeat: %s", repeat)
		log.Printf("🚨 [NextDateAdapter] %v\n", err)
		return time.Time{}, err
	}

	days, err := strconv.Atoi(repeatParts[1])
	if err != nil {
		log.Printf("🚨 [NextDateAdapter] Ошибка парсинга дней: %v\n", err)
		return time.Time{}, err
	}

	newDate := oldDate.AddDate(0, 0, days)
	log.Printf("✅ [NextDateAdapter] Новая дата: %s\n", newDate.Format("20060102"))
	return newDate, nil
}
