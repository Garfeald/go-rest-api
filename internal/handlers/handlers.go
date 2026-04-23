package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api/internal/database"
	"rest-api/internal/models"
	"strconv"
	"strings"
)

// Handler — это структура, которая содержит ссылку на TaskStore для доступа к данным.
// Она инкапсулирует логику обработки HTTP-запросов, связанных с задачами.
// Методы Handler будут использовать TaskStore для выполнения операций над задачами
// и формировать HTTP-ответы в формате JSON.
type Handlers struct {
	store *database.TaskStore
}

// NewHandler создаёт новый экземпляр Handler, принимая TaskStore в качестве аргумента.
// Это позволяет отделить логику доступа к данным от логики обработки HTTP-запросов,
// что улучшает тестируемость и поддерживаемость кода. Возвращает готовый к использованию Handler.
func NewHandlers(store *database.TaskStore) *Handlers {
	// Создаём новый экземпляр Handler, передавая ему ссылку на TaskStore для доступа к данным
	return &Handlers{store: store}
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

// getIDFromURL извлекает id задачи из URL, удаляя префикс "/tasks/" и разделяя оставшуюся часть по "/".
// Если URL не соответствует ожидаемому формату (например, содержит дополнительные сегменты), отправляет ответ с кодом 400 и сообщением об ошибке.
// Конвертирует извлечённый id из строки в целое число. Если id не является числом, отправляет ответ с кодом 400 и сообщением об ошибке.
// Возвращает извлечённый id в виде целого числа и ошибку, если id не удалось извлечь или конвертировать.
func getIDFromURL(w http.ResponseWriter, path string) (int, error) {
	// Извлекаем id из URL, удаляя префикс "/tasks/" и разделяя оставшуюся часть по "/"
	pathParts := strings.Split(strings.TrimPrefix(path, "/tasks/"), "/")
	if len(pathParts) != 1 {
		// Если URL не соответствует ожидаемому формату (например, содержит дополнительные сегменты), отправляем ответ с кодом 400 и сообщением об ошибке
		respondWithError(w, http.StatusBadRequest, "Неверный URL")
		return 0, fmt.Errorf("invalid URL")
	}
	// Конвертируем извлечённый id из строки в целое число
	idStr := pathParts[0]
	// Если id не является числом, отправляем ответ с кодом 400 и сообщением об ошибке
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "Неверный ID задачи")
		return 0, fmt.Errorf("invalid task ID")
	}
	// Конвертируем извлечённый id из строки в целое число
	id, err := strconv.Atoi(idStr)
	// Если id не является числом, отправляем ответ с кодом 400 и сообщением об ошибке
	if err != nil {
		// Если id не является числом, отправляем ответ с кодом 400 и сообщением об ошибке
		respondWithError(w, http.StatusBadRequest, "Неверный ID задачи")
		return 0, fmt.Errorf("invalid task ID")
	}
	return id, nil
}

// GetAllTasks обрабатывает HTTP-запросы на получение всех задач.
// Она вызывает метод GetAll у TaskStore для получения списка задач из базы данных.
// Если возникает ошибка при получении данных, отправляет ответ с кодом 500 и сообщением об ошибке.
// В случае успешного получения задач, отправляет их в формате JSON с кодом 200 OK.
func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
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
func (h *Handlers) GetTaskById(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(w, r.URL.Path)
	if err != nil {
		// Если id не удалось извлечь или конвертировать, завершаем обработку запроса, так как уже был отправлен ответ с ошибкой
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

// CreateTask обрабатывает HTTP-запросы на создание новой задачи.
// Она декодирует данные новой задачи из тела запроса в структуру CreateTaskInput.
// Если возникает ошибка при декодировании (например, неверный формат данных), отправляет ответ с кодом 400 и сообщением об ошибке.
// В случае успешного декодирования, дальнейшая логика по сохранению задачи в базе данных будет реализована внутри этой функции.
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Объявляем переменную для хранения данных новой задачи, которые будут декодированы из тела запроса
	var input models.CreateTaskInput
	// Декодируем данные новой задачи из тела запроса в структуру CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// Если возникает ошибка при декодировании (например, неверный формат данных), отправляем ответ с кодом 400 и сообщением об ошибке
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Поле title не может быть пустым")
		return
	}

	// В случае успешного декодирования, дальнейшая логика по сохранению задачи в базе данных будет реализована внутри этой функции.
	task, err := h.store.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// В случае успешного создания задачи, отправляем её в формате JSON с кодом 201 Created
	respondWithJSON(w, http.StatusCreated, task)
}

// UpdateTask обрабатывает HTTP-запросы на обновление существующей задачи по её id.
// Она извлекает id из URL, конвертирует его в целое число и декодирует данные для обновления из тела запроса в структуру UpdateTaskInput.
// Если id не является числом или возникает ошибка при декодировании, отправляет ответ с кодом 400 и сообщением об ошибке.
// В случае успешного извлечения id и декодирования данных, дальнейшая логика по обновлению задачи в базе данных будет реализована внутри этой функции.
func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Извлекаем id задачи из URL и конвертируем его в целое число
	id, err := getIDFromURL(w, r.URL.Path)
	if err != nil {
		// Если id не удалось извлечь или конвертировать, завершаем обработку запроса, так как уже был отправлен ответ с ошибкой
		return
	}
	// Объявляем переменную для хранения данных для обновления задачи, которые будут декодированы из тела запроса
	var input models.UpdateTaskInput
	// Декодируем данные для обновления задачи из тела запроса в структуру UpdateTaskInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}
	// Если в input есть поле Title, проверяем, что оно не пустое (после удаления пробелов)
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Поле title не может быть пустым")
		return
	}
	// В случае успешного извлечения id и декодирования данных, дальнейшая логика по обновлению задачи в базе данных будет реализована внутри этой функции.
	task, err := h.store.Update(id, input)
	if err != nil {
		// Если возникает ошибка при обновлении задачи (например, задача не найдена), отправляем ответ с кодом 404 и сообщением об ошибке
		if strings.Contains(err.Error(), "not found") {
			// Если задача не найдена, отправляем ответ с кодом 404 и сообщением об ошибке
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			// В случае других ошибок при обновлении задачи, отправляем ответ с кодом 500 и сообщением об ошибке
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// В случае успешного обновления задачи, отправляем её в формате JSON с кодом 200 OK
	respondWithJSON(w, http.StatusOK, task)
}

// DeleteTask обрабатывает HTTP-запросы на удаление задачи по её id.
// Она извлекает id из URL, конвертирует его в целое число и вызывает метод Delete у TaskStore.
// Если id не является числом или возникает ошибка при удалении задачи, отправляет ответ с кодом 400 и сообщением об ошибке.
// В случае успешного удаления задачи, отправляет JSON-ответ с сообщением об успешном удалении и кодом 200 OK.
func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Извлекаем id задачи из URL и конвертируем его в целое число
	id, err := getIDFromURL(w, r.URL.Path)
	if err != nil {
		// Если id не удалось извлечь или конвертировать, завершаем обработку запроса, так как уже был отправлен ответ с ошибкой
		return
	}
	// Вызываем метод Delete у TaskStore для удаления задачи по её id из базы данных
	err = h.store.Delete(id)
	// Если возникает ошибка при удалении задачи (например, задача не найдена), отправляем ответ с кодом 400 и сообщением об ошибке
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// Если задача не найдена, отправляем ответ с кодом 404 и сообщением об ошибке
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			// В случае других ошибок при обновлении задачи, отправляем ответ с кодом 500 и сообщением об ошибке
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Задача успешно удалена"})
}
