package service

import (
    "context"
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/Sacs616/streaming-app/services/auth/internal/domain"
    "github.com/Sacs616/streaming-app/services/auth/internal/repository"
)

var (
    ErrUserExists         = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
)

type Claims struct {
    UserID uint64 `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

type AuthService struct {
    userRepo  repository.UserRepository
    jwtSecret []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        jwtSecret: []byte(jwtSecret),
    }
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*domain.User, string, string, time.Time, error) {
    // Check if user exists
    existingUser, _ := s.userRepo.FindByEmail(ctx, email)
    if existingUser != nil {
        return nil, "", "", time.Time{}, ErrUserExists
    }

    // Create user
    user := &domain.User{
        Email:    email,
        Username: username,
        Role:     "user",
    }

    if err := user.HashPassword(password); err != nil {
        return nil, "", "", time.Time{}, err
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, "", "", time.Time{}, err
    }

    // Generate tokens
    token, expiresAt, err := s.generateToken(user)
    if err != nil {
        return nil, "", "", time.Time{}, err
    }

    refreshToken, _, err := s.generateRefreshToken(user)
    if err != nil {
        return nil, "", "", time.Time{}, err
    }

    return user, token, refreshToken, expiresAt, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, string, time.Time, error) {
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, "", "", time.Time{}, ErrInvalidCredentials
    }

    if !user.CheckPassword(password) {
        return nil, "", "", time.Time{}, ErrInvalidCredentials
    }

    token, expiresAt, err := s.generateToken(user)
    if err != nil {
        return nil, "", "", time.Time{}, err
    }

    refreshToken, _, err := s.generateRefreshToken(user)
    if err != nil {
        return nil, "", "", time.Time{}, err
    }

    return user, token, refreshToken, expiresAt, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return s.jwtSecret, nil
    })

    if err != nil {
        return nil, ErrInvalidToken
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        user, err := s.userRepo.FindByID(ctx, claims.UserID)
        if err != nil {
            return nil, ErrInvalidToken
        }
        return user, nil
    }

    return nil, ErrInvalidToken
}

func (s *AuthService) generateToken(user *domain.User) (string, time.Time, error) {
    expiresAt := time.Now().Add(24 * time.Hour)
    
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(s.jwtSecret)
    
    return tokenString, expiresAt, err
}

func (s *AuthService) generateRefreshToken(user *domain.User) (string, time.Time, error) {
    expiresAt := time.Now().Add(7 * 24 * time.Hour)
    
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(s.jwtSecret)
    
    return tokenString, expiresAt, err
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, time.Time, error) {
    token, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return s.jwtSecret, nil
    })

    if err != nil {
        return "", "", time.Time{}, ErrInvalidToken
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return "", "", time.Time{}, ErrInvalidToken
    }

    user, err := s.userRepo.FindByID(ctx, claims.UserID)
    if err != nil {
        return "", "", time.Time{}, ErrInvalidToken
    }

    // Generate new token pair
    newToken, expiresAt, err := s.generateToken(user)
    if err != nil {
        return "", "", time.Time{}, err
    }

    newRefreshToken, _, err := s.generateRefreshToken(user)
    if err != nil {
        return "", "", time.Time{}, err
    }

    return newToken, newRefreshToken, expiresAt, nil
}

func (s *AuthService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return user, nil
}