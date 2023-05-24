package services

// ! basic crud with getting data from request body and params
import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/godgodwinter/go-fiber-1-keuangan/app/configs"
	"github.com/godgodwinter/go-fiber-1-keuangan/app/models"
	"github.com/godgodwinter/go-fiber-1-keuangan/app/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultService struct{}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func (s *DefaultService) GetBasic(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	data := types.ExampleStruct{
		Name: c.Params("userId"),
		Desc: "Ini userId " + userId,
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

var validate = validator.New()

func (s *DefaultService) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.UserModel
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    err.Error(),
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)
	// Mengiterasi setiap dokumen dalam cursor dan menambahkannya ke slice users
	for cursor.Next(ctx) {
		var user models.UserModel

		// Mendekode dokumen ke struct user
		if err := cursor.Decode(&user); err != nil {
			// Mengembalikan respons kesalahan jika terjadi kesalahan saat mendekode dokumen
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"data":    err.Error(),
				"error":   err.Error(),
			})
		}

		// Menambahkan pengguna ke slice users
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		// Mengembalikan respons kesalahan jika terjadi kesalahan pada cursor
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

func (s *DefaultService) GetUserWhereId(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.UserModel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    err.Error(),
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func (s *DefaultService) CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.UserModel
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    err.Error(),
			"error":   err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    validationErr.Error(),
			"error":   validationErr.Error(),
		})
	}

	newUser := models.UserModel{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    "Data berhasil dibuat!",
		"userId":  result.InsertedID,
	})
}

func (s *DefaultService) UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.UserModel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    err.Error(),
			"error":   err.Error(),
		})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    validationErr.Error(),
			"error":   validationErr.Error(),
		})
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
	fmt.Println(objId)
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    err.Error(),
			"error":   err.Error(),
		})
	}
	//get updated user details
	var updatedUser models.UserModel
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"data":    err.Error(),
				"error":   err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    updatedUser,
	})

}
func (s *DefaultService) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    err.Error(),
			"error":   err.Error(),
		})

	}

	if result.DeletedCount < 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    "User with specified ID not found!",
			"error":   "User with specified ID not found!",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    "User successfully deleted!",
	})
}
