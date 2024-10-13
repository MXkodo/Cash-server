package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MXkodo/cash-server/internal/service"
	"github.com/gin-gonic/gin"
)

type DocHandler struct {
	docService *service.DocService
}

func NewDocHandler(docService *service.DocService) *DocHandler {
	return &DocHandler{docService}
}

func (h *DocHandler) UploadDocument(c *gin.Context) {
	userId := c.GetInt("userID")
	fmt.Println("UserId: ", userId)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный файл"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
		return
	}
	defer file.Close()

	response, err := h.docService.UploadDocument(c, userId, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), false, file)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки документа"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"json": response, "file": fileHeader.Filename}})
}

func (h *DocHandler) GetDocuments(c *gin.Context) {
	userId := c.GetInt("userID")

	login := c.Query("login") 
	key := c.Query("key")     
	value := c.Query("value") 
	limit := c.Query("limit") 

	if limit == "" {
		limit = "10"
	}

	docs, err := h.docService.GetDocuments(c, userId, login, key, value, limit)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения документов"})
		return
	}

	var formattedDocs []struct {
		ID      int      `json:"id"`
		Name    string   `json:"name"`
		Mime    string   `json:"mime"`
		File    bool     `json:"file"`
		Public  bool     `json:"public"`
		Created string   `json:"created"`
		Grant   []string `json:"grant"`
	}

	for _, doc := range docs {
		formattedDocs = append(formattedDocs, struct {
			ID      int      `json:"id"`
			Name    string   `json:"name"`
			Mime    string   `json:"mime"`
			File    bool     `json:"file"`
			Public  bool     `json:"public"`
			Created string   `json:"created"`
			Grant   []string `json:"grant"`
		}{
			ID:      doc.ID,
			Name:    doc.Name,
			Mime:    doc.Mime,
			File:    true, 
			Public:  doc.IsPublic,
			Created: doc.CreatedAt.Format("2006-01-02 15:04:05"), 
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"docs": formattedDocs}})
}

func (h *DocHandler) GetDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный идентификатор документа"})
		return
	}

	doc, err := h.docService.GetDocument(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Документ не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": doc})
}

func (h *DocHandler) DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный идентификатор документа"})
		return
	}

	if err := h.docService.DeleteDocument(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления документа"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": gin.H{"success": true}})
}
