package handlers

import (
	"customer-service/config"
	"customer-service/internal/core/domain"
	"customer-service/internal/mock"
	"customer-service/pkg/dto"
	"customer-service/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type RestHandlerTestSuite struct {
	suite.Suite
	MockService *mock.CustomerService
	TestHandler *HTTPHandler
	TestRouter  *gin.Engine
	Cfg         *config.Config
	TestData    struct {
		User     domain.User
		Customer domain.Customer
	}
}

func (suite *RestHandlerTestSuite) SetupSuite() {
	cfgPath := "../../test/customer.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	logger := logging.MockLogger{}

	mockService := new(mock.CustomerService)

	router := gin.New()
	gin.SetMode(gin.TestMode)

	deliveryHandler := NewRest(mockService, router, logger, cfg)
	deliveryHandler.SetupEndpoints()

	suite.Cfg = cfg
	suite.MockService = mockService
	suite.TestRouter = router
	suite.TestHandler = deliveryHandler
	suite.TestData = struct {
		User     domain.User
		Customer domain.Customer
	}{
		User: domain.User{
			ID:       "test-id-2",
			Name:     "test-name-2",
			LastName: "test-lastname-2",
		},
		Customer: domain.Customer{
			UserID: "test-id",
			User: domain.User{
				ID:       "test-id",
				Name:     "test-name",
				LastName: "test-lastname",
			},
			ServiceArea: 1,
		},
	}
}

func (suite *RestHandlerTestSuite) SetupTest() {
	suite.MockService.ExpectedCalls = nil
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll() {
	suite.MockService.On("GetAll").Return([]domain.Customer{suite.TestData.Customer}, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/customers", nil)
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.CustomerListResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.Len(responseObject, 1)

	suite.EqualValues(suite.TestData.Customer.User.Name, responseObject[0].Name)
	suite.EqualValues(suite.TestData.Customer.UserID, responseObject[0].UserID)
	suite.EqualValues(suite.TestData.Customer.ServiceArea, responseObject[0].ServiceArea)
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll_NoneFound() {
	suite.MockService.On("GetAll").Return([]domain.Customer{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/customers", nil)
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get() {
	suite.MockService.On("Get", suite.TestData.Customer.UserID).Return(suite.TestData.Customer, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", suite.TestData.Customer.UserID), nil)
	request.Header.Set("X-User-Id", suite.TestData.Customer.UserID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.CustomerResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Customer.User.Name, responseObject.Name)
	suite.EqualValues(suite.TestData.Customer.UserID, responseObject.UserID)
	suite.EqualValues(suite.TestData.Customer.ServiceArea, responseObject.ServiceArea)
	suite.EqualValues(suite.TestData.Customer.User.LastName, responseObject.LastName)
}

func (suite *RestHandlerTestSuite) TestHandler_Get_BadID() {
	suite.MockService.On("Get", "test").Return(domain.Customer{}, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", "test"), nil)
	request.Header.Set("X-User-Id", suite.TestData.Customer.UserID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get_NotFound() {
	suite.MockService.On("Get", suite.TestData.Customer.UserID).Return(domain.Customer{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", suite.TestData.Customer.UserID), nil)
	request.Header.Set("X-User-Id", suite.TestData.Customer.UserID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create() {
	suite.MockService.On("Create", suite.TestData.Customer.UserID, suite.TestData.Customer.ServiceArea).Return(suite.TestData.Customer, nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateCustomer{
		ServiceArea: suite.TestData.Customer.ServiceArea,
		ID:          suite.TestData.Customer.UserID,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)
	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.CustomerResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Customer.User.Name, responseObject.Name)
	suite.EqualValues(suite.TestData.Customer.UserID, responseObject.UserID)
	suite.EqualValues(suite.TestData.Customer.ServiceArea, responseObject.ServiceArea)
	suite.EqualValues(suite.TestData.Customer.User.LastName, responseObject.LastName)
}

func (suite *RestHandlerTestSuite) TestHandler_Create_BadInput() {
	rr := httptest.NewRecorder()

	data, err := json.Marshal(struct {
		Test string
	}{
		Test: suite.TestData.Customer.UserID,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create_CouldNotCreate() {
	suite.MockService.On("Create", suite.TestData.Customer.UserID, suite.TestData.Customer.ServiceArea).Return(domain.Customer{}, errors.New("could not create"))

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateCustomer{
		ServiceArea: suite.TestData.Customer.ServiceArea,
		ID:          suite.TestData.Customer.UserID,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/customers", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_UpdateServiceArea() {
	suite.MockService.On("UpdateServiceArea", suite.TestData.Customer.UserID, suite.TestData.Customer.ServiceArea).Return(suite.TestData.Customer, nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyUpdateServiceArea{
		ServiceArea: suite.TestData.Customer.ServiceArea,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s/service-area", suite.TestData.Customer.UserID), strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)
	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.CustomerResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Customer.UserID, responseObject.UserID)
	suite.EqualValues(suite.TestData.Customer.ServiceArea, responseObject.ServiceArea)
}

func (suite *RestHandlerTestSuite) TestHandler_Update_CouldNotCreate() {
	suite.MockService.On("UpdateServiceArea", suite.TestData.Customer.UserID, suite.TestData.Customer.ServiceArea).Return(domain.Customer{}, errors.New("could not update"))

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateCustomer{
		ServiceArea: suite.TestData.Customer.ServiceArea,
		ID:          suite.TestData.Customer.UserID,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/customers/%s/service-area", suite.TestData.Customer.UserID), strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func TestIntegration_RestHandlerTestSuite(t *testing.T) {
	repoSuite := new(RestHandlerTestSuite)
	suite.Run(t, repoSuite)
}
