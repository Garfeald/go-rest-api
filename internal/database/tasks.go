package database

import (
	"database/sql"
	"fmt"
	"rest-api/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// TaskStore инкапсулирует логику работы с таблицей tasks в базе данных.
// Содержит методы для выполнения CRUD-операций над задачами.
// Поле db — экземпляр sqlx.DB, предоставляющий расширенный интерфейс
// над стандартным database/sql с поддержкой маппинга строк в структуры.
type TaskStore struct {
	db *sqlx.DB
}

// NewTaskStore создаёт новый экземпляр TaskStore.
// Принимает подключение к базе данных (*sqlx.DB) и возвращает
// готовый к использованию объект хранилища задач.
func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

// GetAll возвращает список всех задач из таблицы tasks,
// отсортированных по дате создания в порядке убывания (новые — первыми).
// Использует sqlx.Select для маппинга всех строк результата в срез моделей Task.
// Возвращает ошибку, если запрос к базе не удался.
func (s *TaskStore) GetAll() ([]models.Task, error) {

	// Объявляем срез для хранения найденных задач
	var tasks []models.Task

	// Формируем SQL-запрос: выбираем все задачи, сортируем по дате создания (новые первые)
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks order by created_at desc`

	// Выполняем запрос и маппим результат в срез tasks
	err := s.db.Select(&tasks, query)
	if err != nil {
		// Возвращаем ошибку при неудачном запросе к БД
		return nil, err
	}
	// Возвращаем найденные задачи
	return tasks, nil
}

// GetByID ищет задачу по её уникальному идентификатору (id).
// Выполняет параметризованный запрос с плейсхолдером $1 для защиты от SQL-инъекций.
// Если задача с указанным id не найдена, возвращает sql.ErrNoRows,
// который перехватывается и оборачивается в понятную ошибку с номером id.
// Возвращает указатель на модель Task, чтобы вызывающий код мог отличить
// «не найдено» (nil, error) от пустого значения.
func (s *TaskStore) GetByID(id int) (*models.Task, error) {

	// Объявляем переменную для хранения найденной задачи
	var task models.Task

	// Формируем SQL-запрос: ищем задачу по id с параметризованным плейсхолдером $1
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = $1`

	// Выполняем запрос и маппим одну строку в структуру task
	err := s.db.Get(&task, query, id)

	if err == sql.ErrNoRows {
		// Задача не найдена — возвращаем понятную ошибку с указанием id
		return nil, fmt.Errorf("task with id %d not found", id)
	}

	if err != nil {
		// Произошла другая ошибка при запросе к БД
		return nil, err
	}
	// Возвращаем найденную задачу
	return &task, nil

}

// Create вставляет новую задачу в таблицу tasks.
// Получает на вход модель CreateTaskInput с полями title, description и completed.
// Время created_at и updated_at устанавливается автоматически в текущий момент (time.Now).
// Использует конструкцию RETURNING в INSERT-запросе, чтобы PostgreSQL вернула
// все поля только что созданной строки — это позволяет сразу получить
// сгенерированный id и timestamps без дополнительного SELECT.
// StructScan маппит возвращённые столбцы в структуру Task по совпадению имён.
func (s *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {

	// Объявляем переменную для хранения созданной задачи
	var task models.Task

	// Формируем SQL-запрос: вставляем новую строку и возвращаем все поля через RETURNING
	query := `
	INSERT INTO tasks (title, description, completed, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, title, description, completed, created_at, updated_at
	`

	// Фиксируем текущее время для полей created_at и updated_at
	now := time.Now()

	// Выполняем INSERT-запрос с входными данными и маппим возвращённую строку в структуру
	err := s.db.QueryRowx(query, input.Title, input.Description, input.Completed, now, now).StructScan(&task)

	if err != nil {
		// Возвращаем ошибку при неудачной вставке в БД
		return nil, err
	}
	// Возвращаем созданную задачу со всеми заполненными полями
	return &task, nil
}

// Update изменяет существующую задачу по её id.
// Получает на вход id задачи и модель UpdateTaskInput, в которой поля могут быть nil,
// что означает «не обновлять это поле».
// Сначала получает текущую задачу по id, чтобы знать её текущее состояние.
// Затем обновляет только те поля, для которых в input есть новые значения (не nil).
// Обновляет поле updated_at на текущее время.
// Выполняет UPDATE-запрос с новыми данными и возвращает обновлённую задачу через RETURNING.
// Если задача с указанным id не найдена, возвращает ошибку.
func (s *TaskStore) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	// Сначала получаем текущую задачу по id, чтобы знать её текущее состояние
	task, err := s.GetByID(id)

	if err != nil {
		// Если задача не найдена, возвращаем ошибку
		return nil, err
	}

	// Обновляем поля задачи только если в input есть новые значения (не nil)
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	// Обновляем поле updated_at на текущее время
	task.UpdatedAt = time.Now()

	// Формируем SQL-запрос для обновления задачи по id
	query := `
	UPDATE tasks 
	SET title = $1, description = $2, completed = $3, updated_at = $4 
	WHERE id = $5 
	RETURNING id, title, description, completed, created_at, updated_at
	`

	var updatedTask models.Task

	// Выполняем UPDATE-запрос с новыми данными и маппим возвращённую строку в структуру updatedTask
	err = s.db.QueryRowx(query, task.Title, task.Description, task.Completed, task.UpdatedAt, id).StructScan(&updatedTask)

	if err != nil {
		// Возвращаем ошибку при неудачном обновлении в БД
		return nil, err
	}
	// Возвращаем обновлённую задачу
	return &updatedTask, nil
}

// Delete удаляет задачу по её id.
// Выполняет параметризованный DELETE-запрос с плейсхолдером $1 для защиты от SQL-инъекций.
// Если задача с указанным id не найдена, запрос не удалит ничего, и метод вернёт nil,
// что позволяет вызывающему коду не беспокоиться о том, была ли задача удалена или её не было.
func (s *TaskStore) Delete(id int) error {
	// Формируем SQL-запрос для удаления задачи по id
	query := `DELETE FROM tasks WHERE id = $1`

	// Выполняем DELETE-запрос с указанным id
	result, err := s.db.Exec(query, id)
	if err != nil {
		// Возвращаем ошибку при неудачном удалении в БД
		return err
	}

	// Проверяем, была ли удалена хотя бы одна строка (задача)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Если ни одна строка не была удалена, значит задача с таким id не найдена
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	return nil

}
