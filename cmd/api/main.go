package main

import (
	"log"
	"net/http"
	"os"
	"rest-api/internal/database"
	"rest-api/internal/handlers"
)

// main — это точка входа в приложение. Она отвечает за настройку и запуск HTTP-сервера.
func main() {
	// Читаем URL подключения к базе данных из переменной окружения DATABASE_URL.
	databaseURL := os.Getenv("DATABASE_URL")
	// Если переменная окружения DATABASE_URL не установлена, используем значение по умолчанию для локального подключения к PostgreSQL.
	if databaseURL == "" {
		// Устанавливаем значение по умолчанию для URL подключения к базе данных PostgreSQL, если переменная окружения не задана
		databaseURL = "postgres://taskuser:taskpass@localhost:5433/tasksdb?sslmode=disable"
	}
	// Читаем порт для сервера из переменной окружения SERVER_PORT.
	serverPort := os.Getenv("SERVER_PORT")
	// Если переменная окружения SERVER_PORT не установлена, используем значение по умолчанию 8080.
	if serverPort == "" {
		serverPort = "8080"
	}
	// Логируем информацию о запуске сервера, включая порт, на котором он будет работать
	log.Printf("Начинаем запуск вервера %s", serverPort)
	// Подключаемся к базе данных, используя URL подключения. Если подключение не удалось, логируем ошибку и завершаем программу.
	db, err := database.Connect(databaseURL)
	// Если возникает ошибка при подключении к базе данных, логируем её и завершаем программу
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// Гарантируем, что соединение с базой данных будет закрыто при завершении работы программы
	defer db.Close()
	// Логируем успешное подключение к базе данных
	log.Printf("Успешно подключились к базе данных")
	// Создаём новый экземпляр TaskStore, передавая ему подключение к базе данных. Это позволяет TaskStore выполнять операции с базой данных.
	taskStore := database.NewTaskStore(db)
	// Создаём новый экземпляр Handler, передавая ему TaskStore для доступа к данным. Handler будет обрабатывать HTTP-запросы, используя методы TaskStore.
	handler := handlers.NewHandlers(taskStore)
	// Создаём новый HTTP-сервер, используя стандартную библиотеку net/http. Устанавливаем маршруты и обработчики для каждого пути.
	mux := http.NewServeMux()
	// Регистрируем обработчики для маршрутов /tasks и /tasks/create, используя методHandler для проверки HTTP-метода запроса.
	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTasks, "GET"))
	// Регистрируем обработчик для маршрута /tasks/create, который будет обрабатывать POST-запросы для создания новых задач.
	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, "POST"))
	// Регистрируем обработчик для маршрута /tasks/{id}, который будет обрабатывать GET, PUT и DELETE запросы для получения, обновления и удаления задач по id.
	mux.HandleFunc("/tasks/", taskIDHandler(handler))
	// Оборачиваем основной маршрутизатор в промежуточное ПО для логирования входящих запросов. Это позволит нам видеть информацию о каждом запросе в логах.
	loggedMux := LoggingMiddleware(mux)
	// Запускаем HTTP-сервер на указанном порту, используя loggedMux в качестве обработчика. Если возникает ошибка при запуске сервера, логируем её и завершаем программу.
	serverAddr := ":" + serverPort
	// Запускаем HTTP-сервер на указанном порту, используя loggedMux в качестве обработчика. Если возникает ошибка при запуске сервера, логируем её и завершаем программу.
	err = http.ListenAndServe(serverAddr, loggedMux)
	// Если возникает ошибка при запуске сервера, логируем её и завершаем программу
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
	// Логируем информацию о том, что сервер успешно запущен и на каком адресе он работает
	log.Printf("Сервер запущен на %s", serverAddr)

}

// methodHandler — это вспомогательная функция, которая оборачивает обработчик HTTP-запросов,
// проверяя, что HTTP-метод запроса соответствует ожидаемому. Если метод не совпадает,
// она возвращает ошибку 405 Method Not Allowed. Это позволяет централизованно управлять
// поддерживаемыми методами для каждого маршрута и улучшает читаемость кода.
func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	// Возвращаем новую функцию-обёртку, которая будет проверять HTTP-метод запроса перед вызовом основного обработчика
	return func(w http.ResponseWriter, r *http.Request) {
		// Если HTTP-метод запроса не совпадает с разрешённым, отправляем ответ с кодом 405 Method Not Allowed и сообщением об ошибке
		if r.Method != allowedMethod {
			// Отправляем ответ с кодом 405 Method Not Allowed и сообщением об ошибке, если метод запроса не совпадает с разрешённым
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		handlerFunc(w, r)
	}
}

// taskIDHandler — это функция, которая возвращает обработчик HTTP-запросов для маршрута /tasks/{id}.
// Она проверяет HTTP-метод запроса и вызывает соответствующий метод Handler для обработки GET, PUT или DELETE запросов.
// Если HTTP-метод запроса не поддерживается, она возвращает ошибку 405 Method Not Allowed.
func taskIDHandler(handler *handlers.Handlers) http.HandlerFunc {
	// Возвращаем функцию-обёртку, которая будет обрабатывать запросы к маршруту /tasks/{id} и проверять HTTP-метод запроса
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handler.GetTaskById(w, r)
		case "PUT":
			handler.UpdateTask(w, r)
		case "DELETE":
			handler.DeleteTask(w, r)
		default:
			// Если HTTP-метод запроса не поддерживается, отправляем ответ с кодом 405 Method Not Allowed и сообщением об ошибке
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

// LoggingMiddleware — это промежуточное ПО (middleware) для HTTP-сервера, которое логирует информацию о каждом входящем запросе.
// Она принимает следующий обработчик (next) и возвращает новый обработчик, который выполняет логирование перед вызовом следующего обработчика.
// В данном случае, она логирует удалённый адрес клиента, HTTP-метод и URL запроса.
func LoggingMiddleware(next http.Handler) http.Handler {
	// Возвращаем новый обработчик, который будет логировать информацию о каждом входящем запросе перед вызовом следующего обработчика
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем удалённый адрес клиента, HTTP-метод и URL запроса для каждого входящего запроса
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
