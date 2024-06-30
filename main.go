package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName       string = "localhost:27017" // Адрес MongoDB
	dbName         string = "demo_todo"       // Имя базы данных
	collectionName string = "todo"            // Имя коллекции (таблицы)
	port           string = ":9000"           // Порт для сервера
)

// todoModel - это структура для данных в MongoDB.
type todoModel struct {
	ID        bson.ObjectId `bson:"_id,omitempty"` // Уникальный идентификатор
	Title     string        `bson:"title"`         // Заголовок задачи
	Completed bool          `bson:"completed"`     // Статус выполнения
	CreatedAt time.Time     `bson:"createdAt"`     // Время создания
}

// todo - это структура для ответа клиенту.
type todo struct {
	ID        string    `json:"id"`         // Уникальный идентификатор
	Title     string    `json:"title"`      // Заголовок задачи
	Completed bool      `json:"completed"`  // Статус выполнения
	CreatedAt time.Time `json:"created_at"` // Время создания
}

// init - функция для начальной настройки.
func init() {
	rnd = renderer.New() // Создание объекта для отрисовки страниц
	uri := fmt.Sprintf("mongodb://%s/%s", hostName, dbName)
	sess, err := mgo.Dial(uri) // Подключение к MongoDB
	checkErr(err)
	sess.SetMode(mgo.Monotonic, true) // Установка режима работы MongoDB
	db = sess.DB(dbName)              // Выбор базы данных
}

// homeHandler - функция для главной страницы.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil) // Отрисовка шаблона
	checkErr(err)
}

// createTodo - функция для создания новой задачи.
func createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo // Создание переменной для новой задачи

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil { // Чтение данных из запроса
		rnd.JSON(w, http.StatusProcessing, err) // Отправка ошибки
		return
	}

	// Проверка, что заголовок задачи не пустой
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	// Создание задачи в MongoDB
	tm := todoModel{
		ID:        bson.NewObjectId(),
		Title:     t.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	if err := db.C(collectionName).Insert(&tm); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"todo_id": tm.ID.Hex(),
	})
}

// updateTodo - функция для обновления задачи.
func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id")) // Получение идентификатора задачи из URL

	if !bson.IsObjectIdHex(id) { // Проверка, что идентификатор в правильном формате
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var t todo // Создание переменной для обновления задачи

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil { // Чтение данных из запроса
		rnd.JSON(w, http.StatusProcessing, err) // Отправка ошибки
		return
	}

	// Проверка, что заголовок задачи не пустой
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is required",
		})
		return
	}

	// Обновление задачи в MongoDB
	if err := db.C(collectionName).
		Update(
			bson.M{"_id": bson.ObjectIdHex(id)},
			bson.M{"title": t.Title, "completed": t.Completed},
		); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
	})
}

// fetchTodos - функция для получения списка задач.
func fetchTodos(w http.ResponseWriter, r *http.Request) {
	todos := []todoModel{} // Создание переменной для списка задач

	if err := db.C(collectionName).
		Find(bson.M{}).
		All(&todos); err != nil { // Получение всех задач из MongoDB
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"error":   err,
		})
		return
	}

	todoList := []todo{} // Создание переменной для ответа клиенту
	for _, t := range todos {
		todoList = append(todoList, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoList,
	})
}

// deleteTodo - функция для удаления задачи.
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id")) // Получение идентификатора задачи из URL

	if !bson.IsObjectIdHex(id) { // Проверка, что идентификатор в правильном формате
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(id)); err != nil { // Удаление задачи из MongoDB
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete todo",
			"error":   err,
		})
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo deleted successfully",
	})
}

// main - главная функция приложения.
func main() {
	stopChan := make(chan os.Signal)      // Создание канала для остановки сервера
	signal.Notify(stopChan, os.Interrupt) // Оповещение о сигнале остановки

	r := chi.NewRouter()     // Создание маршрутизатора
	r.Use(middleware.Logger) // Использование middleware для логирования
	r.Get("/", homeHandler)  // Маршрут для главной страницы

	r.Mount("/todo", todoHandlers()) // Маршруты для задач

	srv := &http.Server{ // Настройка сервера
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() { // Запуск сервера в отдельной горутине
		log.Println("Listening on port ", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // Ожидание сигнала остановки
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx) // Граутина для корректной остановки сервера
	defer cancel()
	log.Println("Server gracefully stopped!")
}

// todoHandlers - функция для создания обработчиков задач.
func todoHandlers() http.Handler {
	rg := chi.NewRouter() // Создание маршрутизатора для задач
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)        // Маршрут для получения списка задач
		r.Post("/", createTodo)       // Маршрут для создания задачи
		r.Put("/{id}", updateTodo)    // Маршрут для обновления задачи
		r.Delete("/{id}", deleteTodo) // Маршрут для удаления задачи
	})
	return rg
}

// checkErr - функция для обработки ошибок.
func checkErr(err error) {
	if err != nil {
		log.Fatal(err) // Вывод ошибки и завершение работы
	}
}
