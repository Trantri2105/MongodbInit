package transport

import (
	requestDto "backend/dto/request"
	"backend/endpoints"
	"backend/middleware"
	"backend/repository"
	"backend/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.mongodb.org/mongo-driver/mongo"
)

func decodeSignupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request requestDto.SignupDto
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request requestDto.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request requestDto.UpdateUserDto
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		return nil, errors.New("user id is missing")
	}
	return userId, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func NewRouter(db *mongo.Database) *gin.Engine {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	userRepository := repository.NewUserRepository(db)

	jwtService := service.NewJwtService()
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userRepository, jwtService)

	authEndpoint := endpoints.NewAuthEndpoint(authService)
	userEndpoint := endpoints.NewUserEndpoint(userService)
	middleware := middleware.NewMiddleware(jwtService)

	signupHandler := httptransport.NewServer(
		authEndpoint.Signup(),
		decodeSignupRequest,
		encodeResponse,
		options...,
	)

	loginHandler := httptransport.NewServer(
		authEndpoint.Login(),
		decodeLoginRequest,
		encodeResponse,
		options...,
	)

	updateHandler := httptransport.NewServer(
		userEndpoint.UpdateUserInfo(),
		decodeUpdateRequest,
		encodeResponse,
		options...,
	)

	getInfoHander := httptransport.NewServer(
		userEndpoint.GetUserInfo(),
		decodeGetRequest,
		encodeResponse,
		options...,
	)

	deleteHandler := httptransport.NewServer(
		userEndpoint.DeleteUser(),
		decodeDeleteRequest,
		encodeResponse,
		options...,
	)

	r := gin.Default()

	authRoute := r.Group("/auth")
	authRoute.POST("/signup", gin.WrapH(signupHandler))
	authRoute.POST("/login", gin.WrapH(loginHandler))

	userRoute := r.Group("/user", middleware.RequireAuth)
	userRoute.GET("/info", gin.WrapH(getInfoHander))
	userRoute.PATCH("/update", gin.WrapH(updateHandler))
	userRoute.DELETE("/delete", middleware.RequireAdminRole, gin.WrapH(deleteHandler))

	return r
}
