package user

import (
	"api/model"
	jwtoken "api/package/token"
	"api/repository/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"strings"
)

func CreateUser(c *fiber.Ctx) error {
	var requestData model.UserAccount
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal memproses data!",
		})
	}

	// Trim all input
	requestData.Nama = strings.TrimSpace(requestData.Nama)
	requestData.Email = strings.TrimSpace(requestData.Email)
	requestData.Telfon = strings.TrimSpace(requestData.Telfon)
	requestData.Alamat = strings.TrimSpace(requestData.Alamat)
	requestData.Password = strings.TrimSpace(requestData.Password)

	// Cek field wajib
	if requestData.Nama == "" ||
		requestData.Email == "" ||
		requestData.Telfon == "" ||
		requestData.Alamat == "" ||
		requestData.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semua field harus diisi (nama, email, telfon, alamat, password)",
		})
	}

	// Validasi nama minimal 2 karakter
	if len(requestData.Nama) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama terlalu pendek",
		})
	}

	// Validasi email (regex sederhana)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(requestData.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format email tidak valid",
		})
	}

	// Validasi nomor telepon (hanya angka dan minimal 10 digit)
	phoneRegex := regexp.MustCompile(`^[0-9]{10,15}$`)
	if !phoneRegex.MatchString(requestData.Telfon) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nomor telepon tidak valid. Hanya angka dan minimal 10 digit",
		})
	}

	// Validasi password minimal 6 karakter
	if len(requestData.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password minimal 6 karakter",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengenkripsi password",
		})
	}
	requestData.Password = string(hashedPassword)
	requestData.ID = uuid.New().String()

	// Cek duplikat email/telfon
	checkUserData, err := db.CheckEmailOrTelfonExists(requestData.Email, requestData.Telfon)
	if err != nil {
		log.Println("Terjadi kesalahan:", err)
	} else {
		switch checkUserData {
		case 1:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email sudah digunakan",
			})
		case 2:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Nomor telepon sudah digunakan",
			})
		case 3:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email dan nomor telepon sudah digunakan",
			})
		}
	}

	// Simpan data ke database
	if err := db.InsertUserSata(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data ke database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registrasi akun berhasil!",
	})
}

func GetOneUser(c *fiber.Ctx) error {
	email := c.Params("email")
	filter := bson.M{"email": email}
	requestData, err := db.GetOneUserFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

func Authorize(c *fiber.Ctx) error {
	var requestData model.AuthRequest
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	filter := bson.M{"email": requestData.Email}

	accountData, err := db.GetOneUserFilter(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve account",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(accountData.Password), []byte(requestData.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := jwtoken.GenerateJWT(accountData.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	response := fiber.Map{
		"message": "Login successful",
		"token":   token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
