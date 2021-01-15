package usecase

// func TestUsecaseIsSingleton(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	usecase1 := NewHealthCheckUsecase(mockRepo)
// 	usecase2 := NewHealthCheckUsecase(mockRepo)
// 	assert.Equal(t, usecase1, usecase2)
// }

// func TestUseCaseIsCallingRepository(t *testing.T) {
// 	mockRepo := new(mocks.Repository)
// 	mockRepo.On("GetDBTimestamp").Return(time.Now())
// 	usecase := NewHealthCheckUsecase(mockRepo)
// 	healthcheck := usecase.GetDBTimestamp()
// 	assert.NotEmpty(t, healthcheck)
// 	mockRepo.AssertExpectations(t)
// }
