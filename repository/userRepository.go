package repository

import (
	"backend/model"
	"context"
	"errors"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (model.User, error)
	InsertUser(ctx context.Context, user model.User) error
	FindUserById(ctx context.Context, userId string) (model.User, error)
	UpdateUserById(ctx context.Context, userId string, user model.User) error
	DeleteUserById(ctx context.Context, userId string) error
}

type userRepositoryImpl struct {
	db *mongo.Database
}

func (u userRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	filter := bson.M{"email": email}
	timeoutContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var user model.User
	err := u.db.Collection("users").FindOne(timeoutContext, filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u userRepositoryImpl) InsertUser(ctx context.Context, user model.User) error {
	timeoutContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := u.db.Collection("users").InsertOne(timeoutContext, user)
	return err
}

func (u userRepositoryImpl) FindUserById(ctx context.Context, userId string) (model.User, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return model.User{}, errors.New("userId not valid")
	}
	filter := bson.M{"_id": id}
	var user model.User
	timeoutContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	err = u.db.Collection("users").FindOne(timeoutContext, filter).Decode(&user)
	if err != nil {
		return model.User{}, errors.New("user not found with user id: " + userId)
	}
	return user, nil
}

func (u userRepositoryImpl) UpdateUserById(ctx context.Context, userId string, user model.User) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errors.New("userId not valid")
	}
	updates := bson.D{}
	typeData := reflect.TypeOf(user)
	values := reflect.ValueOf(user)
	for i := 1; i < typeData.NumField(); i++ {
		field := typeData.Field(i)
		val := values.Field(i)
		tag := field.Tag.Get("bson")
		if !isZeroType(val) {
			update := bson.E{Key: tag, Value: val.Interface()}
			updates = append(updates, update)
		}
	}
	updateFilter := bson.M{"$set": updates}
	timeoutContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = u.db.Collection("users").UpdateByID(timeoutContext, id, updateFilter)
	return err
}

func isZeroType(value reflect.Value) bool {
	zero := reflect.Zero(value.Type()).Interface()

	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map:
		return value.Len() == 0
	default:
		return reflect.DeepEqual(zero, value.Interface())
	}
}

func (u userRepositoryImpl) DeleteUserById(ctx context.Context, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errors.New("userId not valid")
	}
	filter := bson.M{"_id": id}
	timeoutContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = u.db.Collection("users").DeleteOne(timeoutContext, filter)
	return err
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return userRepositoryImpl{
		db: db,
	}
}
