package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api/internal/database"
	"strconv"
	"strings"
)

// Handler — это структура, которая содержит ссылку на TaskStore для доступа к данным.
// Она инкапсулирует логику обработки HTTP-запросов, связанных с задачами.
// Методы Handler будут использовать TaskStore для выполнения операций над задачами
// и формировать HTTP-ответы в формате JSON.
type Handler struct {
	store *database.TaskStore
}

// NewHandler создаёт новый экземпляр Handler, принимая TaskStore в качестве аргумента.
// Это позволяет отделить логику доступа к данным от логики обработки HTTP-запросов,
// что улучшает тестируемость и поддерживаемость кода. Возвращает готовый к использованию Handler.
func NewHandler(store *database.TaskStore) *Handler {
	// Создаём новый экземпляр Handler, передавая ему ссылку на TaskStore для доступа к данным
	return &Handler{store: store}
}

//	respondWithJSON — это вспомогательная функция для отправки JSON-ответов клиенту.
//
// Она устанавливает заголовок Content-Type в application/json, задаёт статус ответа
// и кодирует переданный payload в JSON-формат, отправляя его в тело ответа.
// Это позволяет стандартизировать формат ответов и упростить обработку данных на клиентской стороне.
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	// Устанавливаем заголовок Content-Type для указания формата данных в ответе
	w.Header().Set("Content-Type", "application/json")
	// Устанавливаем статус ответа (например, 200 OK, 404 Not Found и т.д.)
	w.WriteHeader(statusCode)
	// Кодируем payload в JSON и отправляем его в тело ответа
	json.NewEncoder(w).Encode(payload)
}

// respondWithError — это вспомогательная функция для отправки JSON-ответов с сообщением об ошибке.
// Она вызывает respondWithJSON, передавая статус ошибки и объект с полем "error", содержащим сообщение об ошибке.
// Это позволяет стандартизировать формат ошибок в ответах и упростить обработку ошибок на клиентской стороне.
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// GetAllTasks обрабатывает HTTP-запросы на получение всех задач.
// Она вызывает метод GetAll у TaskStore для получения списка задач из базы данных.
// Если возникает ошибка при получении данных, отправляет ответ с кодом 500 и сообщением об ошибке.
// В случае успешного получения задач, отправляет их в формате JSON с кодом 200 OK.
func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Вызываем метод GetAll у TaskStore для получения списка всех задач из базы данных
	tasks, err := h.store.GetAll()
	if err != nil {
		// Если возникает ошибка при получении данных, отправляем ответ с кодом 500 и сообщением об ошибке
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения задач")
		// Завершаем обработку запроса, так как произошла ошибка
		return
	}
	// В случае успешного получения задач, отправляем их в формате JSON с кодом 200 OK
	respondWithJSON(w, http.StatusOK, tasks)
}

// GetTaskById обрабатывает HTTP-запросы на получение задачи по её уникальному идентификатору (id).
// Она извлекает id из URL, конвертирует его в целое число и вызывает метод GetByID у TaskStore.
// Если id не является числом или возникает ошибка при получении задачи, отправляет ответ с кодом 400 и сообщением об ошибке.
// В случае успешного получения задачи, отправляет её в формате JSON с кодом 200 OK.
func (h *Handler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	// Извлекаем id из URL, удаляя префикс "/tasks/" и разделяя оставшуюся часть по "/"
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	if len(pathParts) != 1 {
		// Если URL не соответствует ожидаемому формату (например, содержит дополнительные сегменты), отправляем ответ с кодом 400 и сообщением об ошибке
		respondWithError(w, http.StatusBadRequest, "Неверный URL")
		return
	}
	// Конвертируем извлечённый id из строки в целое число
	idStr := pathParts[0]
	// Если id не является числом, отправляем ответ с кодом 400 и сообщением об ошибке
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Если id не является числом, отправляем ответ с кодом 400 и сообщением об ошибке
		respondWithError(w, http.StatusBadRequest, "Неверный ID задачи")
		return
	}
	// Вызываем метод GetByID у TaskStore для получения задачи по её id из базы данных
	task, err := h.store.GetByID(id)
	if err != nil {
		// Если возникает ошибка при получении задачи (например, задача не найдена), отправляем ответ с кодом 400 и сообщением об ошибке
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// В случае успешного получения задачи, отправляем её в формате JSON с кодом 200 OK
	respondWithJSON(w, http.StatusOK, task)
}
