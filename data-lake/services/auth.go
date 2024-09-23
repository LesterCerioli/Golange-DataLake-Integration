package services

import (
    "crypto/sha256"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
    "github.com/gofiber/fiber/v2"
    "time"
    
)


var secretKey = "your_secret_key_here"


type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}


func Login(c *fiber.Ctx) error {
    user := new(User)

    if err := c.BodyParser(user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Cannot parse request",
        })
    }

    
    hashedPassword := sha256.Sum256([]byte(user.Password))
    if user.Username == "admin" && fmt.Sprintf("%x", hashedPassword) == fmt.Sprintf("%x", sha256.Sum256([]byte("password"))) {
        
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "username": user.Username,
            "exp":      time.Now().Add(time.Hour * 72).Unix(),
        })

        
        tokenString, err := token.SignedString([]byte(secretKey))
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to generate token",
            })
        }

        return c.JSON(fiber.Map{
            "token": tokenString,
        })
    }

    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "Invalid credentials",
    })
}


func SecureData(c *fiber.Ctx) error {
    
    tokenString := c.Get("Authorization")

    if tokenString == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Missing or malformed token",
        })
    }

    
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // Access secure data with valid token
        return c.JSON(fiber.Map{
            "message": "This is a protected route",
            "user":    claims["username"],
        })
    }

    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "Invalid token",
    })
}
