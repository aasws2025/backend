package event

import (
	"api/model"
	"api/repository/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"time"
)

// TambahEvent godoc
// @Summary Tambah event baru
// @Description Menambahkan data event ke database
// @Tags event
// @Accept json
// @Produce json
// @Param event body model.DataEvent true "Event data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /protected/event [post]
func TambahEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access!",
		})
	}

	var requestData model.DataEvent
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Trim spaces
	requestData.Judul = strings.TrimSpace(requestData.Judul)
	requestData.Tanggal = strings.TrimSpace(requestData.Tanggal)
	requestData.Harga = strings.TrimSpace(requestData.Harga)
	requestData.Lokasi = strings.TrimSpace(requestData.Lokasi)
	requestData.Deskripsi = strings.TrimSpace(requestData.Deskripsi)
	requestData.Kategori = strings.ToLower(strings.TrimSpace(requestData.Kategori))

	// Cek field wajib
	if requestData.Judul == "" ||
		requestData.Tanggal == "" ||
		requestData.Harga == "" ||
		requestData.Lokasi == "" ||
		requestData.Deskripsi == "" ||
		requestData.Kategori == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required (judul, tanggal, harga, lokasi, deskripsi, kategori)",
		})
	}

	// Validasi harga (harus angka)
	if _, err := strconv.Atoi(requestData.Harga); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Harga harus berupa angka",
		})
	}

	// Validasi tanggal (format YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", requestData.Tanggal); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tanggal harus dalam format YYYY-MM-DD",
		})
	}

	// Validasi deskripsi minimal 10 karakter
	if len(requestData.Deskripsi) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Deskripsi minimal harus 10 karakter",
		})
	}

	// Generate ID dan simpan
	requestData.ID = uuid.New().String()
	requestData.AuthorID = userID.(string)

	if err := db.InsertDataEvent(requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data ke database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event berhasil ditambahkan",
		"data":    requestData,
	})
}

// GetAllEvent godoc
// @Summary Ambil semua event
// @Description Mengambil seluruh data event dari database
// @Tags event
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /event [get]
func GetAllEvent(c *fiber.Ctx) error {
	filter := bson.M{}

	requestData, err := db.GetDataEventFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

// GetEventID godoc
// @Summary Ambil detail event
// @Description Mengambil data event berdasarkan ID
// @Tags event
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /event/{id} [get]
func GetEventID(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"id": id}
	requestData, err := db.GetOneDataEventFilter(filter)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

// EditEvent godoc
// @Summary Edit event
// @Description Mengubah data event berdasarkan ID
// @Tags event
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param event body model.DataEvent true "Event data baru"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /protected/event/{id} [put]
func EditEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"id": id}

	var requestData model.DataEvent
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	requestData.ID = id

	if err := db.EditEvent(filter, requestData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create profile!",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    requestData,
	})
}

// DeleteEvent godoc
// @Summary Hapus event
// @Description Menghapus event berdasarkan ID
// @Tags event
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /protected/event/{id} [delete]
func DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"id": id}

	_, err := db.DeleteEvent(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal hapus data!",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data terhapus!",
	})
}
