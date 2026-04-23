package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect устанавливает соединение с базой данных PostgreSQL, используя предоставленный URL подключения.
// Возвращает указатель на sqlx.DB для выполнения запросов и ошибку, если соединение не удалось установить.
// Внутри функции используется sqlx.Connect для создания подключения, а также устанавливаются параметры
// для управления количеством открытых и неиспользуемых соединений в пуле.
func Connect(databaseURL string) (*sqlx.DB, error) {
	// Implement database connection logic here
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		// Возвращаем обёрнутую ошибку с дополнительным контекстом, если не удалось подключиться к базе данных
		return nil, fmt.Errorf("Connect db error: %w", err)
	}
	// Устанавливаем максимальное количество открытых соединений в пуле (например, 25)
	db.SetMaxOpenConns(25)
	// Устанавливаем максимальное количество неиспользуемых соединений в пуле (например, 5)
	db.SetMaxIdleConns(5)

	return db, nil
}
