package handler

func TestHandlerIsSingleton(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	ucHandler1 := HealthCheckHandler{
		HealthCheckUsecase: mockUsecase,
	}
	ucHandler2 := HealthCheckHandler{
		HealthCheckUsecase: mockUsecase,
	}
	assert.Equal(t, ucHandler1, ucHandler2)
}

func TestUsecaseCheckCalled(t *testing.T) {
	mockUsecase := new(mocks.Usecase)
	mockUsecase.On("GetDBTimestamp").Return(time.Now())
	ucHandler := HealthCheckHandler{
		HealthCheckUsecase: mockUsecase,
	}

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	ucHandler.Check(c)
	mockUsecase.AssertExpectations(t)
}
