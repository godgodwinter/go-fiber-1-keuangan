package services

import (
	"context"
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

func (s *DefaultService) GetAll(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	data := types.ExampleStruct{
		Name: c.Params("userId"),
		Desc: "Ini userId " + userId,
	}
	return c.JSON(data)
}

var validate = validator.New()

func (s *DefaultService) GetAllUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(err.Error())
	}
	defer cursor.Close(ctx)
	// Mengiterasi setiap dokumen dalam cursor dan menambahkannya ke slice users
	for cursor.Next(ctx) {
		var user models.User

		// Mendekode dokumen ke struct user
		if err := cursor.Decode(&user); err != nil {
			// Mengembalikan respons kesalahan jika terjadi kesalahan saat mendekode dokumen
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
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
	return c.JSON(users)
}

func (s *DefaultService) GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.JSON(err.Error())
	}

	return c.JSON(user)
}

func (s *DefaultService) CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(err.Error())
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(validationErr.Error())
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(err.Error())
	}

	return c.JSON(result)
}
