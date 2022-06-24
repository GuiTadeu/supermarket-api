package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/GuiTadeu/mercado-fresh-panic/cmd/server/database"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/sellers"
	"github.com/GuiTadeu/mercado-fresh-panic/pkg/web"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_Create_201(t *testing.T) {

	validSeller := db.Seller{
		Id:          1,
		Cid:         1,
		CompanyName: "Nike",
		Address:     "Avenida Paulista, 202",
		Telephone:   "13997780890",
	}

	jsonValue, _ := json.Marshal(validSeller)
	requestBody := bytes.NewBuffer(jsonValue)

	mockService := mockSellerService{
		result: validSeller,
		err:    nil,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v1/sellers", requestBody)
	router.ServeHTTP(response, request)

	responseData := db.Seller{}
	decodeWebResponse(response, &responseData)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, validSeller, responseData)
}

func Test_Create_422(t *testing.T) {

	invalidSeller := db.Seller{}
	jsonValue, _ := json.Marshal(invalidSeller)
	requestBody := bytes.NewBuffer(jsonValue)

	mockService := mockSellerService{}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v1/sellers", requestBody)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
}

func Test_Create_409(t *testing.T) {

	validSeller := db.Seller{
		Id:          1,
		Cid:         1,
		CompanyName: "Nike",
		Address:     "Avenida Paulista, 202",
		Telephone:   "13997780890",
	}

	jsonValue, _ := json.Marshal(validSeller)
	requestBody := bytes.NewBuffer(jsonValue)

	mockService := mockSellerService{
		result: db.Seller{},
		err:    sellers.ExistsSellerCodeError,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/v1/sellers", requestBody)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusConflict, response.Code)
}

func Test_GetAll_200(t *testing.T) {

	sellersList := []db.Seller{
		{
			Id:          1,
			Cid:         1,
			CompanyName: "Nike",
			Address:     "Avenida Paulista, 202",
			Telephone:   "13997780890",
		},
		{
			Id:          2,
			Cid:         2,
			CompanyName: "adidas",
			Address:     "Avenida Mineira, 202",
			Telephone:   "13927180890",
		},
		{
			Id:          3,
			Cid:         3,
			CompanyName: "Puma",
			Address:     "Avenida Goiás, 202",
			Telephone:   "13997780112",
		},
	}

	mockService := mockSellerService{
		result: sellersList,
		err:    nil,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/sellers", nil)
	router.ServeHTTP(response, request)

	responseData := []db.Seller{}
	decodeWebResponse(response, &responseData)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, sellersList, responseData)
}

func Test_Get_200(t *testing.T) {

	foundSeller := db.Seller{
		Id:          1,
		Cid:         1,
		CompanyName: "Nike",
		Address:     "Avenida Paulista, 202",
		Telephone:   "13997780890",
	}

	mockService := mockSellerService{
		result: foundSeller,
		err:    nil,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/sellers/1", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

func Test_Get_404(t *testing.T) {

	mockService := mockSellerService{
		result: db.Seller{},
		err:    sellers.SellerNotFoundError,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/Sellers/666", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func Test_Update_200(t *testing.T) {

	sellerToUpdate := db.Seller{
		Id:          1,
		Cid:         1,
		CompanyName: "Nike",
		Address:     "Avenida Paulista, 202",
		Telephone:   "13997780890",
	}

	jsonValue, _ := json.Marshal(sellerToUpdate)
	requestBody := bytes.NewBuffer(jsonValue)

	updatedSeller := db.Seller{
		Id:          1,
		Cid:         42,
		CompanyName: "Nike",
		Address:     "Avenida Paulista, 202",
		Telephone:   "13997780890",
	}

	mockService := mockSellerService{
		result: updatedSeller,
		err:    nil,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PATCH", "/api/v1/sellers/1", requestBody)
	router.ServeHTTP(response, request)

	responseData := db.Seller{}
	decodeWebResponse(response, &responseData)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, updatedSeller, responseData)
}

func Test_Update_404(t *testing.T) {

	SellerToUpdate := db.Seller{
		Id:          1,
		Cid:         1,
		CompanyName: "Nike",
		Address:     "Avenida Paulista, 202",
		Telephone:   "13997780890",
	}

	jsonValue, _ := json.Marshal(SellerToUpdate)
	requestBody := bytes.NewBuffer(jsonValue)

	mockService := mockSellerService{
		result: db.Seller{},
		err:    sellers.SellerNotFoundError,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PATCH", "/api/v1/sellers/1", requestBody)
	router.ServeHTTP(response, request)

	responseData := db.Seller{}
	decodeWebResponse(response, &responseData)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func Test_Delete_204(t *testing.T) {

	mockService := mockSellerService{
		result: db.Seller{},
		err:    nil,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/v1/sellers/1", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Code)
}

func Test_Delete_404(t *testing.T) {

	mockService := mockSellerService{
		result: db.Seller{},
		err:    sellers.SellerNotFoundError,
	}

	router := setupRouter(mockService)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/v1/sellers/1", nil)
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func decodeWebResponse(response *httptest.ResponseRecorder, responseData any) {
	responseStruct := web.Response{}
	json.Unmarshal(response.Body.Bytes(), &responseStruct)

	jsonData, _ := json.Marshal(responseStruct.Data)
	json.Unmarshal(jsonData, &responseData)
}

func setupRouter(mockService mockSellerService) *gin.Engine {
	controller := NewSeller(mockService)

	router := gin.Default()
	router.POST("/api/v1/sellers", controller.Create())
	router.GET("/api/v1/sellers", controller.FindAll())
	router.GET("/api/v1/sellers/:id", controller.FindOne())
	router.PATCH("/api/v1/sellers/:id", controller.Update())
	router.DELETE("/api/v1/sellers/:id", controller.Delete())

	return router
}