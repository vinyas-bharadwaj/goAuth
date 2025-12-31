package mongodb

import (
	"context"
	"fmt"
	"goAuth/internal/models"
	"goAuth/pkg/utils"
	pb "goAuth/proto/gen"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error connecting to the database")
	}
	defer client.Disconnect(ctx)

	filter := bson.M{"username": username}

	var user models.User
	err = client.Database("auth").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrorHandler(err, "User not found. Incorrect username or password ")
		}
		return nil, utils.ErrorHandler(err, "Internal error")
	}
	return &user, nil
}

func AddUserToDB(ctx context.Context, userFromRequest *pb.RegisterRequest) (*pb.User, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error connecting to mongodb")
	}
	defer client.Disconnect(ctx)

	// Create a new user model from the registration request
	modelUser := &models.User{
		Username: userFromRequest.Username,
		Email:    userFromRequest.Email,
		Password: userFromRequest.Password,
		Role:     "user", // Auto-set default role
	}

	// Hash the password before storing
	if modelUser.Password != "" {
		hashedPassword, err := utils.HashPassword(modelUser.Password)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error hashing password")
		}
		modelUser.Password = hashedPassword
	}

	res, err := client.Database("auth").Collection("users").InsertOne(ctx, modelUser)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error inserting data into mongodb")
	}

	objectId, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		modelUser.Id = objectId.Hex()
	}

	pbUser := MapModelUserToPbUser(modelUser)

	return pbUser, nil
}

func ModifyUserRoleInDB(ctx context.Context, userIdFromReq, updatedRole string) error {
	client, err := CreateMongoClient()
	if err != nil {
		return utils.ErrorHandler(err, "Error connecting to mongodb")
	}
	defer client.Disconnect(ctx)

	objId, err := primitive.ObjectIDFromHex(userIdFromReq)
	if err != nil {
		return utils.ErrorHandler(err, "Invalid ID")
	}

	_, err = client.Database("auth").Collection("users").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"role": updatedRole}})
	if err != nil {
		return utils.ErrorHandler(err, fmt.Sprintf("Error updating user with ID: %s", userIdFromReq))
	}

	return nil
}
