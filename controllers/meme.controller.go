package controllers

import (
	"context"
	"fmt"
	// "log"
	"net/http"
	"time"

	"meme/config"
	"meme/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "context"
	// "fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	// "log"
	// "net/http"
	"os"
	"path/filepath"
	// "time"

	// "meme/config"
	// "meme/models"
	// "image"
	// "image/color"
	// "image/draw"
	"image/png"
	// "os"

	// "github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/image/draw"
)

// UploadMeme handles the uploading of memes with translations and localized images


// Upload and process a meme image with text extraction
func UploadMemeWithTextExtraction(c *gin.Context) {
	// Save uploaded file
	filePath, err := saveUploadedFile(c, "image", "./uploads/memes")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(filePath) // Cleanup file

	// Extract text
	extractedText, err := extractTextFromFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Text extraction failed", "details": err.Error()})
		return
	}

	// Get user ID from form (assuming it's passed)
	userIDStr := c.PostForm("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Convert UserID to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prepare meme for insertion
	meme := models.Meme{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		OriginalImage: filePath,
		ExtractedText: extractedText,
		CreatedAt:     time.Now(),
	}

	// Insert meme into database
	collection := config.DB.Collection("memes")
	_, err = collection.InsertOne(context.TODO(), meme)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload meme"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Meme uploaded successfully",
		"meme_id":         meme.ID,
		"extracted_text":  extractedText,
		"original_image":  filePath,
	})
}

// Extract text from an image using Tesseract OCR
func extractTextFromFile(imagePath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	// Preprocess image for better OCR
	preprocessedPath, err := preprocessImage(imagePath)
	if err == nil {
		client.SetImage(preprocessedPath)
		defer os.Remove(preprocessedPath) // Cleanup preprocessed image
	} else {
		client.SetImage(imagePath)
	}

	// Set OCR language
	client.SetLanguage("eng")

	// Extract text
	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("OCR extraction error: %v", err)
	}

	return text, nil
}

// Preprocess the image to improve OCR accuracy
func preprocessImage(imagePath string) (string, error) {
	// Open the source image
	srcFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// Decode the image
	src, _, err := image.Decode(srcFile)
	if err != nil {
		return "", err
	}

	// Create a grayscale version of the image
	bounds := src.Bounds()
	grayImage := image.NewGray(bounds)

	// Draw the grayscale image
	draw.Draw(grayImage, bounds, src, bounds.Min, draw.Src)

	// Save the preprocessed image
	preprocessedPath := imagePath + "_preprocessed.png"
	preprocessedFile, err := os.Create(preprocessedPath)
	if err != nil {
		return "", err
	}
	defer preprocessedFile.Close()

	// Encode the preprocessed image as PNG
	err = png.Encode(preprocessedFile, grayImage)
	if err != nil {
		return "", err
	}

	return preprocessedPath, nil
}

// Save uploaded file to a specified directory
func saveUploadedFile(c *gin.Context, formField string, uploadDir string) (string, error) {
	file, err := c.FormFile(formField)
	if err != nil {
		return "", fmt.Errorf("No file uploaded")
	}

	os.MkdirAll(uploadDir, os.ModePerm)
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("Failed to save file")
	}

	return filePath, nil
}














// GetMeme retrieves a single meme by ID
func GetMeme(c *gin.Context) {
	memeIDStr := c.Param("id")

	// Convert string ID to ObjectID
	memeID, err := primitive.ObjectIDFromHex(memeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meme ID"})
		return
	}

	var meme models.Meme
	collection := config.DB.Collection("memes")
	
	err = collection.FindOne(context.TODO(), bson.M{"_id": memeID}).Decode(&meme)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meme not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meme"})
		return
	}

	c.JSON(http.StatusOK, meme)
}

// ListMemes retrieves memes with advanced filtering and pagination
func ListMemes(c *gin.Context) {
	collection := config.DB.Collection("memes")

	// Pagination parameters
	limit := int64(10) // Default limit
	skip := int64(0)   // Default skip

	// Parse limit
	if l := c.Query("limit"); l != "" {
		if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	// Parse skip
	if s := c.Query("skip"); s != "" {
		if _, err := fmt.Sscanf(s, "%d", &skip); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skip parameter"})
			return
		}
	}

	// Optional language filter
	language := c.Query("language")
	filter := bson.M{}
	if language != "" {
		// Filter memes with translations or localized images in the specified language
		filter = bson.M{
			"$or": []bson.M{
				{"translations.language": language},
				{"localized_images.language": language},
			},
		}
	}

	// Optional user ID filter
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		filter["user_id"] = userID
	}

	// Find memes with options
	findOptions := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.M{"created_at": -1}) // Sort by most recent first

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch memes"})
		return
	}
	defer cursor.Close(context.TODO())

	// Decode results
	var memes []models.Meme
	if err = cursor.All(context.TODO(), &memes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process memes"})
		return
	}

	// If no memes found
	if len(memes) == 0 {
		c.JSON(http.StatusOK, []models.Meme{})
		return
	}

	c.JSON(http.StatusOK, memes)
}

// UpdateMeme allows partial updates to a meme
func UpdateMeme(c *gin.Context) {
	memeIDStr := c.Param("id")

	// Convert string ID to ObjectID
	memeID, err := primitive.ObjectIDFromHex(memeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meme ID"})
		return
	}

	// Parse update request
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	// Remove _id if accidentally included
	delete(updateData, "_id")
	delete(updateData, "user_id")
	delete(updateData, "created_at")

	// Prepare update
	update := bson.M{"$set": updateData}

	collection := config.DB.Collection("memes")
	result, err := collection.UpdateOne(
		context.TODO(), 
		bson.M{"_id": memeID}, 
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meme"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meme not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Meme updated successfully",
		"modified_count": result.ModifiedCount,
	})
}

// DeleteMeme removes a meme by ID
func DeleteMeme(c *gin.Context) {
	memeIDStr := c.Param("id")

	// Convert string ID to ObjectID
	memeID, err := primitive.ObjectIDFromHex(memeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meme ID"})
		return
	}

	collection := config.DB.Collection("memes")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": memeID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meme"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meme not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Meme deleted successfully",
		"deleted_count": result.DeletedCount,
	})
}