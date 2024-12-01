package handlers

import (
    "net/http"
    "time"
    "knovel/internal/models"
    "knovel/internal/repository"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
    userRepo  repository.UserRepository
    secretKey string
}

type RegisterRequest struct {
    Name     string      `json:"name" binding:"required"`
    Email    string      `json:"email" binding:"required,email"`
    Password string      `json:"password" binding:"required,min=6"`
    Role     models.Role `json:"role" binding:"required,oneof=employee employer"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func NewAuthHandler(userRepo repository.UserRepository, secretKey string) *AuthHandler {
    return &AuthHandler{
        userRepo:  userRepo,
        secretKey: secretKey,
    }
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     req.Role,
    }

    if err := h.userRepo.Create(c, user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userRepo.GetByEmail(c, req.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString([]byte(h.secretKey))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
} 