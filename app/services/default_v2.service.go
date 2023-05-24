package services

// ! v2
// ! basic crud with data model from controller
// ! dont use request.body / request.params on services
import (
	"context"
	"time"

	"github.com/godgodwinter/go-fiber-1-keuangan/app/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *DefaultService) V2_GetAll(c *fiber.Ctx) error {
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

func (s *DefaultService) V2_GetUserWhereId(c *fiber.Ctx, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func (s *DefaultService) V2_CreateUser(c *fiber.Ctx, user models.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := models.UserModel{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"userId":  result.InsertedID,
	})
}

func (s *DefaultService) V2_UpdateUser(c *fiber.Ctx, userId string, user models.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
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
func (s *DefaultService) V2_DeleteUser(c *fiber.Ctx, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
