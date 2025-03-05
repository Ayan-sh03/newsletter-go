package controllers

import (
	"fmt"
	"net/http"
	"newsletter/internal/api/models"
	"newsletter/internal/storage/db"
	"newsletter/internal/utils"
)

func CreateNewsletter(w http.ResponseWriter, r *http.Request) {

	// Parse the request body
	var newsletter models.Newsletter

	err := utils.ParseRequestBody(r, &newsletter)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Create the newsletter

	createNewsLetterParams := db.CreateNewsletterParams{
		Title:       newsletter.Title,
		Description: newsletter.Description,
		// Author: newsletter.Author,
	}

	fmt.Println(createNewsLetterParams)

	q := db.New(db.Db)

	_, err = q.CreateNewsletter(r.Context(), createNewsLetterParams)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create newsletter")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Newsletter created"})
}
