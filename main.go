package main

import (
	"log"
	"net/http"
	"os"

	"github.com/naluneotlichno/FP-GO-API/api"
	"github.com/naluneotlichno/FP-GO-API/database"
)

func main() {
	log.Println("✅ 🔥 Запускаем нашего монстра!")

	// ✅ Инициализация базы данных
	if err := database.InitDB(database.GetDBPath()); err != nil {
		log.Fatalf("❌ Ошибка инициализации БД: %v", err)
	}

	// ✅ Регистрация хендлеров
	registerHandlers()

	// ✅ Запуск сервера
	startServer()
}

// 🔥 registerHandlers регистрирует все хендлеры
func registerHandlers() {
	// Добавление задачи (POST)
	http.HandleFunc("/api/task", api.AddTaskHandler) // Для добавления новой задачи (по примеру твоего предыдущего кода)

	// Получение задач (GET)
	http.HandleFunc("/api/task", api.GetTaskHandler) // Для получения задачи по ID (GET запрос)

	// Обновление задачи (PUT)
	http.HandleFunc("/api/task", api.UpdateTaskHandler) // Для обновления задачи (PUT запрос)

	// Обработчик для следующих запросов (например, обработка даты)
	http.HandleFunc("/api/nextdate", api.HandleNextDate)

	// Получение всех задач (GET)
	http.HandleFunc("/api/tasks", api.GetTasksHandler)

	// Для отдачи статики (веб-страницы)
	http.Handle("/", http.FileServer(http.Dir("./web"))) 
}

// 🔥 startServer запускает сервер
func startServer() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("✅ 🚀 Сервер выезжает на порт %s. Подрубаемся!", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Ой-ой, сервер упал: %v", err)
	}
}
