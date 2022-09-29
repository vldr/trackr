package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"trackr/src/forms/requests"
	"trackr/src/forms/responses"
	"trackr/src/models"
	"trackr/tests"
)

func TestGetValuesRoute(t *testing.T) {
	suite := tests.StartupWithRouter()
	method, path := "GET", "/api/values/"

	//
	// Test not logged in path.
	//

	response, _ := json.Marshal(responses.Error{
		Error: "Not authorized to access this resource.",
	})
	httpRecorder := httptest.NewRecorder()
	httpRequest, _ := http.NewRequest(method, path, nil)
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusForbidden, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid request parameters path.
	//

	response, _ = json.Marshal(responses.Error{
		Error: "Invalid request parameters provided.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, nil)
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid order parameter path.
	//

	request, _ := json.Marshal(requests.GetValues{
		Order: "invalid",
	})
	response, _ = json.Marshal(responses.Error{
		Error: "Invalid order parameter provided.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid offset parameter path.
	//

	request, _ = json.Marshal(requests.GetValues{
		Order:  "asc",
		Offset: -1,
	})
	response, _ = json.Marshal(responses.Error{
		Error: "Invalid offset parameter provided.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid limit parameter path.
	//

	request, _ = json.Marshal(requests.GetValues{
		Order:  "asc",
		Offset: 0,
		Limit:  -1,
	})
	response, _ = json.Marshal(responses.Error{
		Error: "Invalid limit parameter provided.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid field id parameter path.
	//

	request, _ = json.Marshal(requests.GetValues{
		Order:   "asc",
		Offset:  0,
		Limit:   0,
		FieldID: models.Field{}.ID,
	})
	response, _ = json.Marshal(responses.Error{
		Error: "Failed to find field.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test successful path.
	//

	expectedValues := []models.Value{suite.Value}
	numberOfExpectedValues := 10 + len(expectedValues)

	for i := 2; i <= numberOfExpectedValues; i++ {
		newValue := suite.Value
		newValue.ID = uint(i)
		newValue.CreatedAt = time.Now()
		newValue.Value = fmt.Sprintf("%d.00", i)

		expectedValues = append(expectedValues, newValue)

		suite.Service.GetValueService().AddValue(newValue)
	}

	for offset := 0; offset < numberOfExpectedValues; offset++ {
		for limit := 0; limit <= numberOfExpectedValues; limit++ {
			for _, order := range []string{"desc", "asc"} {
				expectedValuesSlice := expectedValues[:]

				sort.Slice(expectedValuesSlice, func(i, j int) bool {
					if order == "desc" {
						return expectedValuesSlice[i].CreatedAt.After(
							expectedValuesSlice[j].CreatedAt,
						)
					} else {
						return expectedValuesSlice[i].CreatedAt.Before(
							expectedValuesSlice[j].CreatedAt,
						)
					}
				})

				if offset == 0 {
					expectedValuesSlice = expectedValuesSlice[:]
				} else {
					expectedValuesSlice = expectedValuesSlice[offset:]
				}

				request, _ = json.Marshal(requests.GetValues{
					Order:   order,
					Offset:  offset,
					Limit:   limit,
					FieldID: suite.Field.ID,
				})

				httpRecorder = httptest.NewRecorder()
				httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
				httpRequest.Header.Add("Cookie", "Session=SessionID")
				suite.Router.ServeHTTP(httpRecorder, httpRequest)

				assert.Equal(t, http.StatusOK, httpRecorder.Code)

				var response responses.ValueList
				err := json.Unmarshal(httpRecorder.Body.Bytes(), &response)
				assert.Nil(t, err)

				if limit > 0 {
					assert.LessOrEqual(t, len(response.Values), limit)
				} else {
					assert.Equal(t, len(expectedValuesSlice), len(response.Values))
				}

				for i := 0; i < len(response.Values); i++ {
					assert.Equal(t, expectedValuesSlice[i].Value, response.Values[i].Value)
				}
			}
		}
	}
}

func TestDeleteValuesRoute(t *testing.T) {
	suite := tests.StartupWithRouter()
	method, path := "DELETE", "/api/values/"

	//
	// Test not logged in path.
	//

	response, _ := json.Marshal(responses.Error{
		Error: "Not authorized to access this resource.",
	})
	httpRecorder := httptest.NewRecorder()
	httpRequest, _ := http.NewRequest(method, path+"invalid", nil)
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusForbidden, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid field id parameter path.
	//

	response, _ = json.Marshal(responses.Error{
		Error: "Invalid :fieldId parameter provided.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path+"invalid", nil)
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test non existant field parameter path.
	//

	response, _ = json.Marshal(responses.Error{
		Error: "Failed to find field.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path+"0", nil)
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test successful path.
	//

	value, err := suite.Service.GetValueService().GetValue(suite.Value.ID, suite.User)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, suite.Value.ID, value.ID)

	response, _ = json.Marshal(responses.Empty{})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path+fmt.Sprintf("%d", suite.Field.ID), nil)
	httpRequest.Header.Add("Cookie", "Session=SessionID")
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusOK, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	value, err = suite.Service.GetValueService().GetValue(suite.Value.ID, suite.User)
	assert.NotNil(t, err)
	assert.Nil(t, value)
}

func TestAddValueRoute(t *testing.T) {
	suite := tests.StartupWithRouter()
	method, path := "POST", "/api/values/"

	//
	// Test invalid request parameters path.
	//

	response, _ := json.Marshal(responses.Error{
		Error: "Invalid request parameters provided.",
	})
	httpRecorder := httptest.NewRecorder()
	httpRequest, _ := http.NewRequest(method, path, nil)
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test empty value path.
	//

	request, _ := json.Marshal(requests.AddValue{})
	response, _ = json.Marshal(responses.Error{
		Error: "The value cannot be empty.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid api key path.
	//

	request, _ = json.Marshal(requests.AddValue{
		Value:  "2.00",
		APIKey: "invalid",
	})
	response, _ = json.Marshal(responses.Error{
		Error: "Failed to find project, invalid API key.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test invalid field id parameter path.
	//

	request, _ = json.Marshal(requests.AddValue{
		Value:   "2.00",
		APIKey:  suite.Project.APIKey,
		FieldID: models.Field{}.ID,
	})
	response, _ = json.Marshal(responses.Error{
		Error: "Failed to find field.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test maximum values path.
	//

	request, _ = json.Marshal(requests.AddValue{
		Value:   "2.00",
		APIKey:  suite.Project.APIKey,
		FieldID: suite.Field.ID,
	})
	response, _ = json.Marshal(responses.Error{
		Error: "You have exceeded your max values limit.",
	})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	//
	// Test successful path.
	//

	suite.User.MaxValues++
	suite.Service.GetUserService().UpdateUser(suite.User)

	value, err := suite.Service.GetValueService().GetValue(2, suite.User)
	assert.NotNil(t, err)
	assert.Nil(t, value)

	request, _ = json.Marshal(requests.AddValue{
		Value:   "2.00",
		APIKey:  suite.Project.APIKey,
		FieldID: suite.Field.ID,
	})
	response, _ = json.Marshal(responses.Empty{})
	httpRecorder = httptest.NewRecorder()
	httpRequest, _ = http.NewRequest(method, path, bytes.NewReader(request))
	suite.Router.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusOK, httpRecorder.Code)
	assert.Equal(t, response, httpRecorder.Body.Bytes())

	value, err = suite.Service.GetValueService().GetValue(2, suite.User)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, uint(2), value.ID)
	assert.Equal(t, "2.00", value.Value)
}