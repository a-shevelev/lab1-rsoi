package handlers

import (
	"fmt"
	"lab1-rsoi/internal/dto"
	"lab1-rsoi/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	service *service.PersonService
}

func New(service *service.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func (h *PersonHandler) RegisterRoutes(rg *gin.RouterGroup) {
	personRoutes := rg.Group("/persons")
	{
		personRoutes.POST("", h.CreatePerson)
		personRoutes.GET("", h.ListPersons)
		personRoutes.GET("/:id", h.GetPerson)
		personRoutes.PATCH("/:id", h.UpdatePerson)
		personRoutes.DELETE("/:id", h.DeletePerson)
	}
}

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := &dto.CreatePersonRequest{
		Name:    req.Name,
		Age:     req.Age,
		Address: req.Address,
		Work:    req.Work,
	}

	id, err := h.service.Create(c.Request.Context(), person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", fmt.Sprintf("/persons/%d", id))
	c.Status(http.StatusCreated)
}

func (h *PersonHandler) ListPersons(c *gin.Context) {
	persons, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]dto.PersonResponse, len(persons))
	for i, p := range persons {
		resp[i] = dto.PersonResponse{
			ID:      p.ID,
			Name:    p.Name,
			Age:     p.Age,
			Address: p.Address,
			Work:    p.Work,
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PersonHandler) GetPerson(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	person, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if person == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	c.JSON(http.StatusOK, dto.PersonResponse{
		ID:      person.ID,
		Name:    person.Name,
		Age:     person.Age,
		Address: person.Address,
		Work:    person.Work,
	})
}

func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := h.service.Get(c.Request.Context(), id)
	if err != nil || person == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	if req.Age != nil {
		person.Age = req.Age
	}
	if req.Address != nil {
		person.Address = req.Address
	}
	if req.Work != nil {
		person.Work = req.Work
	}

	personResponse, err := h.service.Update(c.Request.Context(), id, *person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, personResponse)
}

func (h *PersonHandler) DeletePerson(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
