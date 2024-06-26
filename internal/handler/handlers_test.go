package handler

import (
	"errors"
	"fmt"
	"github.com/ZnNr/GopherVault/internal/database"
	"github.com/ZnNr/GopherVault/internal/models"
	"github.com/ZnNr/GopherVault/internal/models/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Login(t *testing.T) {
	userName := "jaime"
	password := "cersei"
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	testCases := []struct {
		name            string
		storageResponse error
		expectedCode    int
		cookies         bool
	}{
		{
			name:            "positive: success login",
			storageResponse: nil,
			expectedCode:    http.StatusOK,
			cookies:         true,
		},
		{
			name:            "negative: no such user",
			storageResponse: database.ErrNoSuchUser,
			expectedCode:    http.StatusUnauthorized,
			cookies:         false,
		},
		{
			name:            "negative: invalid credentials",
			storageResponse: database.ErrInvalidCredentials,
			expectedCode:    http.StatusUnauthorized,
			cookies:         false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Login", mock.Anything, userName, password).Return(tt.storageResponse)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/login", h.LoginHandler)
			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, userName, password)).
				Post(fmt.Sprintf("%s/auth/login", srv.URL))
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			if tt.cookies {
				assert.True(t, len(h.cookies) == 1)
				assert.True(t, h.cookies[userName] != "")
				assert.True(t, len(resp.Header().Get("Authorization")) > 1)
				assert.True(t, len(resp.Cookies()) == 1)
			}
		})
	}
	t.Run("negative: login is not provided", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/login", h.LoginHandler)
		srv := httptest.NewServer(r)
		defer srv.Close()

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"password": %q}`, password)).
			Post(fmt.Sprintf("%s/auth/login", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
	t.Run("negative: invalid body", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/login", h.LoginHandler)
		srv := httptest.NewServer(r)
		defer srv.Close()

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(`"password"`).
			Post(fmt.Sprintf("%s/auth/login", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
}

func TestHandler_Register(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	userName := "joffrey"
	password := "iamabadguy"

	testCases := []struct {
		name            string
		storageResponse error
		expectedCode    int
		cookies         bool
	}{
		{
			name:            "positive: success register",
			storageResponse: nil,
			expectedCode:    http.StatusOK,
			cookies:         true,
		},
		{
			name:            "negative: user exists",
			storageResponse: database.ErrUserAlreadyExists,
			expectedCode:    http.StatusConflict,
			cookies:         false,
		},
		{
			name:            "negative: db error",
			storageResponse: errors.New("some db error"),
			expectedCode:    http.StatusInternalServerError,
			cookies:         false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, userName, password).Return(tt.storageResponse)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			srv := httptest.NewServer(r)
			defer srv.Close()

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, userName, password)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			if tt.cookies {
				assert.True(t, len(h.cookies) == 1)
				assert.True(t, h.cookies[userName] != "")
				assert.True(t, len(resp.Header().Get("Authorization")) > 1)
				assert.True(t, len(resp.Cookies()) == 1)
			}
		})
	}
}

func TestHandler_GetUserCredentials(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	userName := "margaery"
	systemPassword := "rose"
	password := "ihatecersei"
	login := "queen"

	testCases := []struct {
		name                 string
		storageResponse      []models.Credentials
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name: "positive: success getting credentials",
			storageResponse: []models.Credentials{
				{
					UserName: userName,
					Login:    &login,
					Password: &password,
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: fmt.Sprintf(`[{"user_name":%q,"login":%q,"password":%q}]`, userName, login, password),
		},
		{
			name:                 "negative: no data for user",
			storageResponse:      []models.Credentials{},
			storageResponseError: database.ErrNoData,
			expectedCode:         http.StatusNoContent,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, userName, systemPassword).Return(nil)
			mockedStorage.On("GetCredentials", mock.Anything, models.Credentials{UserName: userName}).Return(tt.storageResponse, tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/get/credentials", h.GetUserCredentialsHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, userName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q}`, userName)).
				Post(fmt.Sprintf("%s/get/credentials", srv.URL))
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}

	t.Run("negative: invalid json", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, userName, password).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/get/credentials", h.GetUserCredentialsHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, userName, password)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q`, userName)).
			Post(fmt.Sprintf("%s/get/credentials", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
	t.Run("negative: unauthorized user", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, userName, password).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/get/credentials", h.GetUserCredentialsHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, userName, password)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(`{"user_name": "other"}`).
			Post(fmt.Sprintf("%s/get/credentials", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusUnauthorized)
	})
}

func TestHandler_SaveUserCredentials(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "shae"
	systemPassword := "mylion"
	loginName := "badgirl"
	password := "money"
	metadata := "capital"

	testCases := []struct {
		name                 string
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name:         "positive: success saving credentials",
			expectedCode: http.StatusOK,
			expectedBody: `Учетные данные для пользователя "shae" сохранены`,
		},
		{
			name:                 "negative: saving error",
			expectedCode:         http.StatusInternalServerError,
			storageResponseError: errors.New("save error"),
			expectedBody:         `ошибка запроса пользователя "shae": save error`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
			mockedStorage.On("SaveCredentials", mock.Anything, models.Credentials{UserName: systemName, Login: &loginName, Password: &password, Metadata: &metadata}).Return(tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/save/credentials", h.SaveUserCredentialsHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q, "login": %q, "password": %q, "metadata": %q}`, systemName, loginName, password, metadata)).
				Post(fmt.Sprintf("%s/save/credentials", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}
	t.Run("negative: bad json", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/save/credentials", h.SaveUserCredentialsHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q, "login1": %q}`, systemName, loginName)).
			Post(fmt.Sprintf("%s/save/credentials", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
}

func TestHandler_DeleteUserCredentials(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "missandei"
	systemPassword := "qwerty12"
	login := "blackgirl"

	t.Run("positive: with login", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
		mockedStorage.On("DeleteCredentials", mock.Anything, models.Credentials{UserName: systemName, Login: &login}).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/delete/credentials", h.DeleteUserCredentialsHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		request := fmt.Sprintf(`{"login": %q, "user_name": %q}`, login, systemName)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(request).
			Post(fmt.Sprintf("%s/delete/credentials", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusOK)
		assert.Equal(t, resp.String(), "Учетные данные для пользователя \"missandei\" с логином \"blackgirl\" были успешно удалены")
	})
	t.Run("positive: with no login", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
		mockedStorage.On("DeleteCredentials", mock.Anything, models.Credentials{UserName: systemName}).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/delete/credentials", h.DeleteUserCredentialsHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		request := fmt.Sprintf(`{"user_name": %q}`, systemName)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(request).
			Post(fmt.Sprintf("%s/delete/credentials", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusOK)
		assert.Equal(t, resp.String(), "Учетные данные для пользователя \"missandei\" были успешно удалены")
	})
}

func TestHandler_UpdateUserCredentials(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "shae"
	systemPassword := "mylion"
	loginName := "badgirl"
	password := "money"
	metadata := "capital"

	testCases := []struct {
		name                 string
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name:         "positive: success updating credentials",
			expectedCode: http.StatusOK,
			expectedBody: `Учетные данные для пользователя "shae" были успешно обновлены`,
		},
		{
			name:                 "negative: updating error",
			expectedCode:         http.StatusInternalServerError,
			storageResponseError: errors.New("update error"),
			expectedBody:         `ошибка запроса пользователя "shae": update error`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
			mockedStorage.On("UpdateCredentials", mock.Anything, models.Credentials{UserName: systemName, Login: &loginName, Password: &password, Metadata: &metadata}).Return(tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/update/credentials", h.UpdateUserCredentialsHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q, "login": %q, "password": %q, "metadata": %q}`, systemName, loginName, password, metadata)).
				Post(fmt.Sprintf("%s/update/credentials", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}
	t.Run("negative: bad json", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/update/credentials", h.UpdateUserCredentialsHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q, "login1": %q}`, systemName, loginName)).
			Post(fmt.Sprintf("%s/update/credentials", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
}

func TestHandler_SaveUserNote(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "hound"
	systemPassword := "ihavebadbrother"
	loginName := "puppy"
	metadata := "no fire please"
	title := "some title"
	content := "some content"

	testCases := []struct {
		name                 string
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name:         "positive: success saving notest",
			expectedCode: http.StatusOK,
			expectedBody: `Заметка для пользователя "hound" была успешно сохранена`,
		},
		{
			name:                 "negative: saving error",
			expectedCode:         http.StatusInternalServerError,
			storageResponseError: errors.New("save error"),
			expectedBody:         `ошибка запроса пользователя "hound": save error`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
			mockedStorage.On("SaveNote", mock.Anything, models.Note{UserName: systemName, Title: &title, Content: &content, Metadata: &metadata}).Return(tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/save/note", h.SaveUserNoteHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q, "title": %q, "content": %q, "metadata": %q}`, systemName, title, content, metadata)).
				Post(fmt.Sprintf("%s/save/note", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}
	t.Run("negative: bad json", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/save/note", h.SaveUserNoteHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q, "login1": %q}`, systemName, loginName)).
			Post(fmt.Sprintf("%s/save/note", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
}

func TestHandler_GetUserNote(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "hound"
	systemPassword := "ihavebadbrother"
	metadata := "no fire please"
	title := "some title"
	content := "some content"

	testCases := []struct {
		name                 string
		storageResponse      []models.Note
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name: "positive: success getting credentials",
			storageResponse: []models.Note{
				{
					UserName: systemName,
					Title:    &title,
					Content:  &content,
					Metadata: &metadata,
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: fmt.Sprintf(`[{"user_name":%q,"title":%q,"content":%q,"metadata":%q}]`, systemName, title, content, metadata),
		},
		{
			name:                 "negative: no data for user",
			storageResponse:      []models.Note{},
			storageResponseError: database.ErrNoData,
			expectedCode:         http.StatusNoContent,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
			mockedStorage.On("GetNotes", mock.Anything, models.Note{UserName: systemName}).Return(tt.storageResponse, tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/get/note", h.GetUserNoteHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q}`, systemName)).
				Post(fmt.Sprintf("%s/get/note", srv.URL))
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}

	t.Run("negative: invalid json", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/get/note", h.GetUserNoteHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q`, systemName)).
			Post(fmt.Sprintf("%s/get/note", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
	t.Run("negative: unauthorized user", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/get/note", h.GetUserNoteHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(`{"user_name": "other"}`).
			Post(fmt.Sprintf("%s/get/note", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusUnauthorized)
	})
}

func TestHandler_DeleteUserNotes(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "missandei"
	systemPassword := "qwerty12"
	title := "some title"

	t.Run("positive: with title", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
		mockedStorage.On("DeleteNotes", mock.Anything, models.Note{UserName: systemName, Title: &title}).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/delete/note", h.DeleteUserNotesHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q, "title":%q}`, systemName, title)).
			Post(fmt.Sprintf("%s/delete/note", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusOK)
		assert.Equal(t, resp.String(), `Заметки для пользователя "missandei" с заголовком "some title" были успешно удалены`)
	})
	t.Run("positive: with no title", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
		mockedStorage.On("DeleteNotes", mock.Anything, models.Note{UserName: systemName}).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/delete/note", h.DeleteUserNotesHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q}`, systemName)).
			Post(fmt.Sprintf("%s/delete/note", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusOK)
		assert.Equal(t, resp.String(), "Заметки для пользователя \"missandei\" были успешно удалены")
	})
}

func TestHandler_UpdateUserNote(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "shae"
	systemPassword := "mylion"
	title := "some title"
	content := "some content"
	metadata := "capital"

	testCases := []struct {
		name                 string
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name:         "positive: success updating notes",
			expectedCode: http.StatusOK,
			expectedBody: `Заметка для пользователя "shae" успешно обновлена`,
		},
		{
			name:                 "negative: updating error",
			expectedCode:         http.StatusInternalServerError,
			storageResponseError: errors.New("update error"),
			expectedBody:         `ошибка запроса пользователя "shae": update error`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
			mockedStorage.On("UpdateNote", mock.Anything, models.Note{UserName: systemName, Title: &title, Content: &content, Metadata: &metadata}).Return(tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/update/note", h.UpdateUserNoteHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q, "title": %q, "content": %q, "metadata": %q}`, systemName, title, content, metadata)).
				Post(fmt.Sprintf("%s/update/note", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}
	t.Run("negative: bad json", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/update/note", h.UpdateUserNoteHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q, "login1": %q}`, systemName, title)).
			Post(fmt.Sprintf("%s/update/note", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
}

func TestHandler_SaveCard(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "hound"
	systemPassword := "ihavebadbrother"
	cv := "724"
	metadata := "green bank"
	bankName := "sber"
	password := "strong"
	number := "1111000033338888"

	testCases := []struct {
		name                 string
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name:         "positive: success saving",
			expectedCode: http.StatusOK,
			expectedBody: `Карточка для пользователя "hound" успешно сохранена`,
		},
		{
			name:                 "negative: saving error",
			expectedCode:         http.StatusInternalServerError,
			storageResponseError: errors.New("save error"),
			expectedBody:         `ошибка запроса пользователя "hound": save error`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
			mockedStorage.On("SaveCard", mock.Anything, models.Card{UserName: systemName, BankName: &bankName, Number: &number, CV: &cv, Password: &password, Metadata: &metadata}).Return(tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/save/card", h.SaveCardHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q, "bank_name": %q, "number": %q,"cv":%q,"password":%q,"metadata": %q}`, systemName, bankName, number, cv, password, metadata)).
				Post(fmt.Sprintf("%s/save/card", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}
}

func TestHandler_GetCard(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	systemName := "hound"
	systemPassword := "ihavebadbrother"
	metadata := "no fire please"
	bankName := "alpha"
	number := "0000888822227777"
	cv := "728"
	password := "strong"

	testCases := []struct {
		name                 string
		storageResponse      []models.Card
		storageResponseError error
		expectedCode         int
		expectedBody         string
	}{
		{
			name: "positive: success getting credentials",
			storageResponse: []models.Card{
				{
					UserName: systemName,
					BankName: &bankName,
					Number:   &number,
					CV:       &cv,
					Password: &password,
					Metadata: &metadata,
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: fmt.Sprintf(`[{"user_name":%q,"bank_name":%q,"number":%q,"cv":%q,"password":%q,"metadata":%q}]`, systemName, bankName, number, cv, password, metadata),
		},
		{
			name:                 "negative: no data for user",
			storageResponse:      []models.Card{},
			storageResponseError: database.ErrNoData,
			expectedCode:         http.StatusNoContent,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockedStorage := mocks.NewStorage(t)
			mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
			mockedStorage.On("GetCard", mock.Anything, models.Card{UserName: systemName}).Return(tt.storageResponse, tt.storageResponseError)

			r := chi.NewRouter()
			h := New(mockedStorage, log)
			r.Post("/auth/register", h.RegisterHandler)
			r.Group(func(r chi.Router) {
				r.Post("/auth/register", h.RegisterHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(h.CheckAuthorization)
				r.Post("/get/card", h.GetCardHandler)
			})
			srv := httptest.NewServer(r)
			defer srv.Close()

			_, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
				Post(fmt.Sprintf("%s/auth/register", srv.URL))
			assert.NoError(t, err)

			resp, err := resty.New().R().
				SetHeader("content-type", "application/json").
				SetBody(fmt.Sprintf(`{"user_name": %q}`, systemName)).
				Post(fmt.Sprintf("%s/get/card", srv.URL))
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode(), tt.expectedCode)
			assert.Equal(t, resp.String(), tt.expectedBody)
		})
	}

	t.Run("negative: invalid json", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/get/card", h.GetCardHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q`, systemName)).
			Post(fmt.Sprintf("%s/get/card", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusBadRequest)
	})
	t.Run("negative: unauthorized user", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/get/card", h.GetCardHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(`{"user_name": "other"}`).
			Post(fmt.Sprintf("%s/get/card", srv.URL))
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusUnauthorized)
	})
}

func TestHandler_DeleteCard(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	systemName := "hound"
	systemPassword := "ihavebadbrother"
	bankName := "alpha"
	number := "0000888822227777"

	t.Run("positive: with bank name", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
		mockedStorage.On("DeleteCards", mock.Anything, models.Card{UserName: systemName, BankName: &bankName}).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/delete/card", h.DeleteCardHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q, "bank_name":%q}`, systemName, bankName)).
			Post(fmt.Sprintf("%s/delete/card", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusOK)
		assert.Equal(t, resp.String(), `Карты "alpha" банка, принадлежащие пользователю "hound" были успешно удалены`)
	})
	t.Run("positive: with no number", func(t *testing.T) {
		mockedStorage := mocks.NewStorage(t)
		mockedStorage.On("Register", mock.Anything, systemName, systemPassword).Return(nil)
		mockedStorage.On("DeleteCards", mock.Anything, models.Card{UserName: systemName, Number: &number}).Return(nil)

		r := chi.NewRouter()
		h := New(mockedStorage, log)
		r.Post("/auth/register", h.RegisterHandler)
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", h.RegisterHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.CheckAuthorization)
			r.Post("/delete/card", h.DeleteCardHandler)
		})
		srv := httptest.NewServer(r)
		defer srv.Close()

		_, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"login": %q, "password": %q}`, systemName, systemPassword)).
			Post(fmt.Sprintf("%s/auth/register", srv.URL))
		assert.NoError(t, err)

		resp, err := resty.New().R().
			SetHeader("content-type", "application/json").
			SetBody(fmt.Sprintf(`{"user_name": %q, "number":%q}`, systemName, number)).
			Post(fmt.Sprintf("%s/delete/card", srv.URL))

		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode(), http.StatusOK)
		assert.Equal(t, resp.String(), "Карты с номером \"0000888822227777\" принадлежащие пользователю \"hound\" были успешно удалены")
	})
}
