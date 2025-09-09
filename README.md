# 🚀 Todo List API на Go

RESTful API для управления списком задач с использованием MongoDB в качестве базы данных.

## ✨ Возможности

- ✅ Создание, чтение, обновление и удаление задач (CRUD)
- 🗄️ Хранение данных в MongoDB
- 📋 Валидация входных данных
- ⚡ Высокая производительность благодаря Go
- 🛡️ Graceful shutdown для корректного завершения
- 📊 Логирование запросов с помощью middleware

## 🏗️ Архитектура

- **Маршрутизация**: Chi Router
- **Рендеринг ответов**: Renderer
- **База данных**: MongoDB с драйвером mgo.v2
- **Обработка ошибок**: Единая функция проверки ошибок

## 📦 Зависимости

- `github.com/go-chi/chi` - маршрутизатор
- `github.com/go-chi/chi/middleware` - middleware для логирования
- `github.com/thedevsaddam/renderer` - рендеринг JSON ответов
- `gopkg.in/mgo.v2` - драйвер MongoDB
- `gopkg.in/mgo.v2/bson` - работа с BSON форматом

## 🚀 Быстрый старт

### 1. Установите MongoDB

Убедитесь, что MongoDB установлена и запущена на `localhost:27017`

### 2. Установите зависимости

```bash
go mod init todo-app
go get github.com/go-chi/chi
go get github.com/go-chi/chi/middleware
go get github.com/thedevsaddam/renderer
go get gopkg.in/mgo.v2
3. Запустите приложение
bash
go run main.go
Сервер запустится на порту :9000

📡 API Endpoints
GET /
Главная страница (возвращает HTML шаблон)

GET /todo/
Получить все задачи

bash
curl http://localhost:9000/todo/
POST /todo/
Создать новую задачу

bash
curl -X POST http://localhost:9000/todo/ \
  -H "Content-Type: application/json" \
  -d '{"title":"Новая задача"}'
PUT /todo/{id}
Обновить задачу

bash
curl -X PUT http://localhost:9000/todo/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{"title":"Обновленная задача", "completed":true}'
DELETE /todo/{id}
Удалить задачу

bash
curl -X DELETE http://localhost:9000/todo/507f1f77bcf86cd799439011
📊 Структуры данных
Запрос (Request)
json
{
  "title": "string",
  "completed": "boolean"
}
Ответ (Response)
json
{
  "id": "string",
  "title": "string", 
  "completed": "boolean",
  "created_at": "timestamp"
}
⚙️ Конфигурация
Константы в коде для настройки:

go
const (
  hostName       = "localhost:27017"  // Адрес MongoDB
  dbName         = "demo_todo"        // Имя базы данных
  collectionName = "todo"             // Имя коллекции
  port           = ":9000"            // Порт сервера
)
## 🎯 Пример использования
1. Создание задачи
bash
curl -X POST http://localhost:9000/todo/ \
  -H "Content-Type: application/json" \
  -d '{"title":"Купить молоко"}'
Ответ:

json
{
  "message": "Todo created successfully",
  "todo_id": "507f1f77bcf86cd799439011"
}
## 2. Получение всех задач
bash
curl http://localhost:9000/todo/
Ответ:

json
{
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Купить молоко",
      "completed": false,
      "created_at": "2023-12-07T10:30:00Z"
    }
  ]
}
## 🔧 Функции обработчиков
homeHandler
Обрабатывает главную страницу, рендерит HTML шаблон

createTodo
Создает новую задачу с валидацией заголовка

updateTodo
Обновляет существующую задачу

fetchTodos
Возвращает список всех задач

deleteTodo
Удаляет задачу по ID

## 🛡️ Обработка ошибок
Валидация ObjectID

Проверка обязательных полей

Обработка ошибок MongoDB

HTTP статус коды для всех сценариев

## 📋 Требования
Go 1.16+
MongoDB 4.4+


MongoDB 4.4+

Интернет соединение для загрузки зависимостей
