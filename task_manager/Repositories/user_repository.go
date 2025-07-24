package repositories

import (
	"context"
	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type userRepository struct {
	database mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) domain.UserRepository{
	return &userRepository{
		database: db,
		collection: collection,
	}
}

func (r *userRepository) Register(c context.Context, user domain.User) error{
	collection := r.database.Collection(r.collection)
	_, err := collection.InsertOne(c, user)
	return err
}

func (r *userRepository) GetUser(c context.Context, username string) (*domain.User, error){
	collection := r.database.Collection(r.collection)

	filter := bson.D{{Key: "username", Value: username}}
	var user domain.User
	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		return nil, err
	} 
	return &user, nil
}

func (r *userRepository) GetNewUserID(c context.Context) (int, error){
	collection := r.database.Collection(r.collection)

	opts:= options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	var lastUser domain.User
	err := collection.FindOne(c, bson.M{}, opts).Decode(&lastUser)
	if err != nil && err != mongo.ErrNoDocuments{
		return -1, err
	}

	if err == mongo.ErrNoDocuments{
		return 1, nil
	} 
	return lastUser.ID+1, nil
}