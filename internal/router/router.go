package router

import (
	"github.com/ZnNr/GopherVault/internal/handler"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

// Middleware для установки заголовка Content-Type "application/json"
func setJSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// New создает новый маршрутизатор с настройками и обработчиками
func New(db models.Storage, log *zap.SugaredLogger) *chi.Mux {
	// Создаем обработчик HTTP запросов
	httpHandler := handler.New(db, log)

	// Инициализируем новый маршрутизатор Chi
	r := chi.NewRouter()
	// Применение middleware для установки Content-Type заголовка на всем маршрутизаторе
	r.Use(setJSONContentTypeMiddleware)

	// Группа маршрутов для авторизации и регистрации
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", httpHandler.RegisterHandler)
		r.Post("/auth/login", httpHandler.LoginHandler)
	})

	// Группа маршрутов для управления данными пользователя
	r.Group(func(r chi.Router) {
		r.Use(httpHandler.CheckAuthorization)

		// Маршруты для управления учетными данными
		r.Post("/save/credentials", httpHandler.SaveUserCredentialsHandler)
		r.Post("/delete/credentials", httpHandler.DeleteUserCredentialsHandler)
		r.Post("/get/credentials", httpHandler.GetUserCredentialsHandler)
		r.Post("/update/credentials", httpHandler.UpdateUserCredentialsHandler)

		// Маршруты для управления заметками пользователя
		r.Post("/save/note", httpHandler.SaveUserNoteHandler)
		r.Post("/delete/note", httpHandler.DeleteUserNotesHandler)
		r.Post("/get/note", httpHandler.GetUserNoteHandler)
		r.Post("/update/note", httpHandler.UpdateUserNoteHandler)

		// Маршруты для управления банковскими картами
		// Обновление карт не предусмотрено
		r.Post("/save/card", httpHandler.SaveCardHandler)
		r.Post("/delete/card", httpHandler.DeleteCardHandler)
		r.Post("/get/card", httpHandler.GetCardHandler)
	})

	// Возвращаем итоговый маршрутизатор
	return r
}
