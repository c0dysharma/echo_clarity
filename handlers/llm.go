package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/c0dysharma/echo_clarity/helpers"
	"github.com/c0dysharma/echo_clarity/structs"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OCRResponse struct{
	Data []structs.CalendarEvent `json:"data"`
}

func OCREventData(c echo.Context) error {
	dstFilePath := fmt.Sprintf("uploads/%s", uuid.New().String())

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "file not found"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid file"})
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(dstFilePath)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to create file"})
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to copy file to destination"})
	}

	// call llm and return response
	resp, err := helpers.CallLLM(dstFilePath)
	if(err != nil){
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to process image"})
	}
	log.Info(resp)

	return c.JSON(http.StatusOK, OCRResponse{Data: resp})
}