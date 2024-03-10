package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/marcojulian/go-jwt/initializers"
	"github.com/marcojulian/go-jwt/models"
	"gorm.io/gorm"
)

func AddWord(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to be logged in to perform this action"})
		return
	}
	usr, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user data"})
		return
	}

	err := c.Request.ParseMultipartForm(10 << 20) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Upload request could not be processed"})
		return
	}

	text := c.PostForm("text")
	translation1 := c.PostForm("translation1")
	translation2 := c.PostForm("translation2")
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if mkDirErr := os.MkdirAll(uploadsDir, os.ModePerm); mkDirErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
			return
		}
	}
	filePath := filepath.Join("uploads", newFileName)
	if err := c.SaveUploadedFile(header, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	} else {
		fmt.Println("File saved successfully")
	
	}

	word := models.Word{
		Text:         text,
		Translation1: translation1,
		Translation2: translation2,
		ImagePath:    filePath,
		UserID:       usr.ID,
		CreatedAt:    time.Now(),
	}
	result := initializers.DB.Create(&word)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save word"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Word added successfully", "wordId": word.ID})
}

func GetWords(c *gin.Context) {
    var words []models.Word
    result := initializers.DB.Find(&words)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve words"})
        return
    }
    for i := range words {
        words[i].ImagePath = fmt.Sprintf("http://%s/uploads/%s", c.Request.Host, filepath.Base(words[i].ImagePath))
    }
    c.JSON(http.StatusOK, gin.H{"words": words})
}

func DeleteWord(c *gin.Context) {
    wordID := c.Param("id")
    if wordID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: ID is required"})
        return
    }

    var word models.Word
    result := initializers.DB.First(&word, wordID)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query the database"})
        }
        return
    }

    deleteResult := initializers.DB.Delete(&word)
    if deleteResult.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the word"})
        return
    }

    if err := os.Remove(word.ImagePath); err != nil {
        fmt.Printf("Warning: Failed to delete image file: %s\n", err)
    }

    c.JSON(http.StatusOK, gin.H{"message": "Word deleted successfully"})
}

func UpdateWord(c *gin.Context) {
    wordID := c.Param("id")
    if wordID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: ID is required"})
        return
    }

    var updateRequest models.Word
    if err := c.ShouldBindJSON(&updateRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    var word models.Word
    result := initializers.DB.First(&word, wordID)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query the database"})
        }
        return
    }

    word.Text = updateRequest.Text
    word.Translation1 = updateRequest.Translation1
    word.Translation2 = updateRequest.Translation2

    saveResult := initializers.DB.Save(&word)
    if saveResult.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the word"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Word updated successfully", "word": word})
}