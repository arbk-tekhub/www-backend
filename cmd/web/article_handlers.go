package main

import (
	"errors"
	"net/http"

	"github.com/benk-techworld/www-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func (app *application) createArticleHandler(c *gin.Context) {

	var input service.CreateArticleInput

	err := c.BindJSON(&input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	err = app.service.CreateArticle(&input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": input.ValidationErrors,
			})
		default:
			app.internalServerErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "article successfully created",
	})

}

func (app *application) fetchArticleHandler(c *gin.Context) {

	idString := c.Param("id")

	article, err := app.service.GetArticleByID(idString)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrFailedValidation):
			app.notFoundResponse(c)
		default:
			app.internalServerErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"article": article,
	})
}

func (app *application) deleteArticleHandler(c *gin.Context) {

	idString := c.Param("id")

	err := app.service.DeleteArticle(idString)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrFailedValidation):
			app.notFoundResponse(c)
		default:
			app.internalServerErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "article successfully deleted",
	})
}

func (app *application) fetchArticlesHandler(c *gin.Context) {

	var input service.FilterArticlesInput

	qs := c.Request.URL.Query()

	input.Title = readString(qs, "title", "")
	input.Tags = readCsv(qs, "tags", []string{})
	input.Page = readInt(qs, "page", 1)
	input.PageSize = readInt(qs, "page_size", 20)
	input.Sort = readString(qs, "sort", "-published")

	articles, err := app.service.GetArticles(&input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFailedValidation):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": input.ValidationErrors,
			})
		default:
			app.internalServerErrorResponse(c, err)

		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}
