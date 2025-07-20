package handler

// HTTP handler layer for the Product domain.
// All comments are in English, as required.
// The handler exposes a classic CRUD REST API and is intentionally thin:
// * validation / parsing of the request
// * mapping to service layer
// * marshalling of the response or error

import (
	"errors"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/model"
	"github.com/ADRPUR/event-driven-marketplace/internal/product/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

// Handler wires the HTTP endpoints to ProductService.
// Prefer using the RegisterHTTPRoutes helper for readability in main().

type Handler struct {
	svc *service.ProductService
}

// New creates a new HTTP handler.
func New(svc *service.ProductService) *Handler {
	return &Handler{svc: svc}
}

// RegisterHTTPRoutes is a convenience wrapper used by main.go.
// It instantiates a Handler and mounts `/products` under the supplied Gin engine.
func RegisterHTTPRoutes(r *gin.Engine, svc *service.ProductService) {
	New(svc).RegisterRoutes(r)
}

// RegisterRoutes attaches all product routes under the provided router group.
// Example tree:
//
// POST   /products          → create product
// GET    /products/:id      → get product by id
// GET    /products          → list products (page, pageSize)
// PUT    /products/:id      → full update product
// DELETE /products/:id      → delete product
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/products")
	{
		g.POST("", h.create)
		g.GET(":id", h.get)
		g.GET("", h.list)
		g.PUT(":id", h.update)
		g.DELETE(":id", h.delete)
	}
}

// ---- request / response DTOs ----

type createProductReq struct {
	Name        string  `json:"name" binding:"required,min=2,max=255"`
	Description string  `json:"description" binding:"max=1024"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

type updateProductReq struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
}

// ---- handlers ----

func (h *Handler) create(c *gin.Context) {
	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prod, err := h.svc.Create(c, req.Name, req.Description, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toJSON(prod))
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")
	prod, err := h.svc.Get(c, id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toJSON(prod))
}

func (h *Handler) list(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	products, err := h.svc.List(c, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]gin.H, len(products))
	for i, p := range products {
		resp[i] = toJSON(&p)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) update(c *gin.Context) {
	var req updateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")

	prod, err := h.svc.Update(c, id, req.Name, req.Description, req.Price)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toJSON(prod))
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := h.svc.Delete(c, id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ---- helpers ----

func toJSON(m *model.Product) gin.H {
	return gin.H{
		"id":          m.ID,
		"name":        m.Name,
		"description": m.Description,
		"price":       m.Price,
		"createdAt":   m.CreatedAt,
	}
}
