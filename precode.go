package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

// Тесты
func TestMainHandlerCorrectRequest(t *testing.T) {
	// Создание запроса с корректными параметрами
	req, err := http.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	require.NoError(t, err)

	// Создаём Recorder для записи ответа
	responseRecorder := httptest.NewRecorder()

	// Выполняем запрос
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка, что ответ имеет статус 200
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Проверка, что тело ответа не пустое
	assert.NotEmpty(t, responseRecorder.Body.String())

	// Проверка содержимого ответа (должно быть два кафе)
	assert.Equal(t, "Мир кофе,Сладкоежка", responseRecorder.Body.String())
}

func TestMainHandlerWrongCity(t *testing.T) {
	// Создание запроса с неправильным городом
	req, err := http.NewRequest("GET", "/cafe?count=2&city=ny", nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверяем код ошибки 400
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	// Проверяем сообщение об ошибке
	assert.Equal(t, "wrong city value", responseRecorder.Body.String())
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// здесь нужно добавить необходимые проверки
	// Проверяем, что код ответа 200
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Проверяем, что вернулись все доступные кафе (totalCount кафе)
	response := responseRecorder.Body.String()
	cafeList := strings.Split(response, ",")
	assert.Len(t, cafeList, totalCount)

	// Проверяем, что вернулись правильные названия кафе
	assert.Equal(t, "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент", response)
}
