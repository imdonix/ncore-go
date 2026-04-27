package rest

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imdonix/ncore-go/pkg/ncore"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {

	gin.DefaultWriter = os.Stdout
	gin.DefaultErrorWriter = os.Stderr

	s := &Server{
		router: gin.Default(),
	}

	s.router.POST("/login", s.handleLogin)
	s.router.POST("/search", s.handleSearch)
	s.router.GET("/torrent/:id", s.handleGetTorrent)
	s.router.GET("/torrent/:id/download", s.handleDownload)
	s.router.GET("/activity", s.handleGetByActivity)
	s.router.GET("/recommended", s.handleGetRecommended)
	s.router.POST("/logout", s.handleLogout)

	return s
}

func (s *Server) Start(addr string) error {
	fmt.Printf("\n🚀 Server is running on http://localhost%s\n\n", addr)
	return s.router.Run(addr)
}

func (s *Server) getClient(c *gin.Context) (*ncore.Client, error) {
	token := c.GetHeader("X-Ncore-Auth")
	if token == "" {
		// Fallback to Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return nil, fmt.Errorf("missing authentication token in headers (X-Ncore-Auth or Authorization)")
	}

	return ncore.NewClientFromToken(15*time.Second, token)
}

func (s *Server) handleLogin(c *gin.Context) {
	var loginReq struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		TwoFactor string `json:"2factor"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := ncore.NewClient(15 * time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := client.Login(loginReq.Username, loginReq.Password, loginReq.TwoFactor)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) handleSearch(c *gin.Context) {
	client, err := s.getClient(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var searchReq struct {
		Pattern   string                 `json:"pattern"`
		Type      ncore.SearchParamType  `json:"type"`
		Where     ncore.SearchParamWhere `json:"where"`
		SortBy    ncore.ParamSort        `json:"sort_by"`
		SortOrder ncore.ParamSeq         `json:"sort_order"`
		Page      int                    `json:"page"`
	}

	if err := c.ShouldBindJSON(&searchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if searchReq.Page == 0 {
		searchReq.Page = 1
	}

	res, err := client.Search(searchReq.Pattern, searchReq.Type, searchReq.Where, searchReq.SortBy, searchReq.SortOrder, searchReq.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleGetTorrent(c *gin.Context) {
	client, err := s.getClient(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	res, err := client.GetTorrent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleGetByActivity(c *gin.Context) {
	client, err := s.getClient(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	res, err := client.GetByActivity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleGetRecommended(c *gin.Context) {
	client, err := s.getClient(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tType := ncore.SearchParamType(c.Query("type"))
	res, err := client.GetRecommended(tType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleDownload(c *gin.Context) {
	client, err := s.getClient(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	torrent, err := client.GetTorrent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	body, filename, err := client.Download(torrent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer body.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/x-bittorrent")
	c.DataFromReader(http.StatusOK, -1, "application/x-bittorrent", body, nil)
}

func (s *Server) handleLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out (stateless)"})
}
