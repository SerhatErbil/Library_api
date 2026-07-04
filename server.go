package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Book struct {
	Name     string  `json:"name" bson:"name"`
	ID       int     `json:"id" bson:"id"`
	Author   string  `json:"author" bson:"author"`
	Borrower *string `json:"borrower,omitempty" bson:"borrower,omitempty"`
}

var userCollection *mongo.Collection
var bookCollection *mongo.Collection

func connectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Mongo bağlantı hatası:", err)
	}

	db := client.Database("mydb")
	userCollection = db.Collection("users")
	bookCollection = db.Collection("books")
	log.Println("MongoDB'ye bağlandı.")
}

var users = map[string]string{}
var books = map[string]string{}

func AuthMiddleware(c *fiber.Ctx) error {
	type Body struct {
		Username string `json:"username"`
	}

	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz JSON formatı"})
	}

	if body.Username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Username eksik"})
	}

	if _, ok := users[body.Username]; !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Yetkisiz erişim"})
	}

	c.Locals("username", body.Username)

	return c.Next()
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing User
	err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existing)
	if err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Kullanıcı zaten var"})
	}

	hashed, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Şifre hashlenemedi"})
	}
	user.Password = hashed

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kayıt başarısız"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Kullanıcı başarıyla kaydedildi"})
}

func Login(c *fiber.Ctx) error {
	var input User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := userCollection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
	}

	if !checkPasswordHash(input.Password, user.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Şifre yanlış"})
	}

	users[input.Username] = user.Password

	return c.Status(200).JSON(fiber.Map{"message": "Giriş başarılı"})
}

func Profile(c *fiber.Ctx) error {
	type RequestBody struct {
		Username string `json:"username"`
	}
	var body RequestBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := userCollection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
	}

	cursor, err := bookCollection.Find(ctx, bson.M{"borrower": body.Username})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kitaplar getirilemedi"})
	}
	defer cursor.Close(ctx)

	var borrowedBooks []Book
	if err := cursor.All(ctx, &borrowedBooks); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kitaplar çözümlenemedi"})
	}

	return c.Status(200).JSON(fiber.Map{
		"username": user.Username,
		"books":    borrowedBooks,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	username := c.Params("username")
	if _, exists := users[username]; !exists {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	delete(users, username)
	return c.Status(200).JSON(fiber.Map{"message": "User deleted"})
}

func AddBook(c *fiber.Ctx) error {
	var book Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz JSON"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := bookCollection.InsertOne(ctx, book)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kitap eklenemedi"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Kitap başarıyla eklendi"})
}

func GetLibrary(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bookCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kitaplar getirilemedi"})
	}
	defer cursor.Close(ctx)

	var books []Book
	if err := cursor.All(ctx, &books); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kitaplar çözümlenemedi"})
	}

	return c.JSON(books)
}
func BorrowBook(c *fiber.Ctx) error {
	type RequestBody struct {
		BookID   int    `json:"book_id"`
		Username string `json:"username"`
	}
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz JSON"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := userCollection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kullanıcı bulunamadı"})
	}

	count, err := bookCollection.CountDocuments(ctx, bson.M{"borrower": body.Username})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sayma hatası"})
	}
	if count >= 2 {
		return c.Status(403).JSON(fiber.Map{"error": "Bu kullanıcı en fazla 2 kitap alabilir"})
	}

	var book Book
	err = bookCollection.FindOne(ctx, bson.M{"id": body.BookID}).Decode(&book)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kitap bulunamadı"})
	}
	if book.Borrower != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Bu kitap zaten başka biri tarafından alınmış"})
	}

	_, err = bookCollection.UpdateOne(ctx, bson.M{"id": body.BookID}, bson.M{
		"$set": bson.M{"borrower": body.Username},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kitap alma başarısız"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Kitap başarıyla alındı"})
}

func main() {
	app := fiber.New()
	connectMongo()

	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Post("/books", AuthMiddleware, AddBook)
	app.Get("/library", AuthMiddleware, GetLibrary)
	app.Get("/profile", AuthMiddleware, Profile)
	app.Post("/borrow", AuthMiddleware, BorrowBook)

	log.Fatal(app.Listen(":3000"))
}
