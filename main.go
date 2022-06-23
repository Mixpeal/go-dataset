package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mixpeal/go-dataset/models"
	"github.com/morkid/paginate"

	"github.com/mixpeal/go-dataset/storage"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type User struct {
	Name     string `json:"name" validate:"required,min=3,max=40"`
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password,omitempty" validate:"required"`
	Date     string `json:"date" validate:"required"`
	Company  string `json:"company" validate:"required,min=3,max=40"`
}

type Repository struct {
	DB *gorm.DB
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(user User) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func returnNewUser(user User) User {
	newUser := User{
		Name:    user.Name,
		Email:   user.Email,
		Date:    user.Date,
		Company: user.Company,
	}
	return newUser
}
func (r *Repository) CreateUser(context *fiber.Ctx) error {
	user := User{}
	err := context.BodyParser(&user)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request failed"})

		return err
	}
	errors := ValidateStruct(user)
	if errors != nil {
		return context.Status(fiber.StatusBadRequest).JSON(errors)
	}
	hash, err := hashPassword(user.Password)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})

	}
	user.Password = hash
	if err := r.DB.Create(&user).Error; err != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User has been added", "data": returnNewUser(user)})
	return nil
}
func (r *Repository) UpdateUser(context *fiber.Ctx) error {
	type UpdateUserInput struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	var uui UpdateUserInput
	if err := context.BodyParser(&uui); err != nil {
		return context.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	db := r.DB
	userModel := &models.Users{}

	userModel.Name = &uui.Name
	userModel.Email = &uui.Email
	if db.Model(&userModel).Where("id = ?", id).Updates(&userModel).RowsAffected == 0 {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get User with given id"})
	}

	return context.JSON(fiber.Map{"status": "success", "message": "User successfully updated"})
}

func (r *Repository) DeleteUser(context *fiber.Ctx) error {
	userModel := models.Users{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	err := r.DB.Delete(userModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete boo"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User delete successfully"})
	return nil
}

func (r *Repository) GetUsers(context *fiber.Ctx) error {
	db := r.DB
	model := db.Model(&models.Users{})

	pg := paginate.New(&paginate.Config{
		DefaultSize: 20,
	})

	page := pg.With(model).Request(context.Request()).Response(&[]models.Users{})

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"data": page,
	})
	return nil
}

func (r *Repository) GetUserByID(context *fiber.Ctx) error {
	id := context.Params("id")
	userModel := &models.Users{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(userModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get the user"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User id fetched successfully", "data": userModel})
	return nil
}

func (repo *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/users", repo.GetUsers)
	api.Post("/users", repo.CreateUser)
	api.Patch("/users/:id", repo.UpdateUser)
	api.Delete("/users/:id", repo.DeleteUser)
	api.Get("/users/:id", repo.GetUserByID)
}

func main() {
	_, ok := os.LookupEnv("APP_ENV")

	if !ok {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load the database")
	}

	err = models.MigrateUsers(db)

	if err != nil {
		log.Fatal("Could not migrate db")
	}

	repo := Repository{
		DB: db,
	}
	app := fiber.New()
	repo.SetupRoutes(app)
	app.Listen(":9091")
}
