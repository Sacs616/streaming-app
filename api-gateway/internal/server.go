package internal

import (
    "context"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    authpb "github.com/Sacs616/streaming-app/services/auth/proto/auth"
)

type Gateway struct {
    authClient authpb.AuthServiceClient
}

func NewGateway() (*Gateway, error) {
    authURL := os.Getenv("AUTH_SERVICE_URL")
    if authURL == "" {
        authURL = "localhost:50051"
    }

    authConn, err := grpc.Dial(
        authURL,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, err
    }

    return &Gateway{
        authClient: authpb.NewAuthServiceClient(authConn),
    }, nil
}

func (g *Gateway) SetupRouter() *gin.Engine {
    router := gin.Default()

    // CORS
    router.Use(g.corsMiddleware())

    // Public routes
    auth := router.Group("/api/v1/auth")
    {
        auth.POST("/register", g.handleRegister)
        auth.POST("/login", g.handleLogin)
        auth.POST("/refresh", g.handleRefreshToken)
    }

    // Protected routes
    protected := router.Group("/api/v1")
    protected.Use(g.authMiddleware())
    {
        protected.GET("/me", g.handleGetMe)
        protected.POST("/logout", g.handleLogout)
    }

    return router
}

func (g *Gateway) corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func (g *Gateway) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
        defer cancel()

        resp, err := g.authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
            Token: token,
        })

        if err != nil || !resp.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Set("user", resp.User)
        c.Next()
    }
}

func (g *Gateway) handleRegister(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required,email"`
        Username string `json:"username" binding:"required,min=3"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    resp, err := g.authClient.Register(ctx, &authpb.RegisterRequest{
        Email:    req.Email,
        Username: req.Username,
        Password: req.Password,
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "user":          resp.User,
        "token":         resp.Token,
        "refresh_token": resp.RefreshToken,
        "expires_at":    resp.ExpiresAt,
    })
}

func (g *Gateway) handleLogin(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    resp, err := g.authClient.Login(ctx, &authpb.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user":          resp.User,
        "token":         resp.Token,
        "refresh_token": resp.RefreshToken,
        "expires_at":    resp.ExpiresAt,
    })
}

func (g *Gateway) handleGetMe(c *gin.Context) {
    user, _ := c.Get("user")
    c.JSON(http.StatusOK, gin.H{"user": user})
}

func (g *Gateway) handleLogout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    token := strings.TrimPrefix(authHeader, "Bearer ")

    ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

    _, err := g.authClient.Logout(ctx, &authpb.LogoutRequest{
        Token: token,
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (g *Gateway) handleRefreshToken(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()

    resp, err := g.authClient.RefreshToken(ctx, &authpb.RefreshTokenRequest{
        RefreshToken: req.RefreshToken,
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token":         resp.Token,
        "refresh_token": resp.RefreshToken,
        "expires_at":    resp.ExpiresAt,
    })
}
