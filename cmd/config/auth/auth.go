package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	// "github.com/rchmachina/grpc/dto/model"
	userPb "github.com/rchmachina/grpc/dto/authpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CustomClaims defines the structure of JWT claims with custom fields
type CustomClaims struct {
	Username string `json:"username"` // Your custom claim (e.g., username)
	jwt.StandardClaims
	Roles string `json:"roles"`
}

// JWT secret
var jwtSecret = []byte("secretK3yz") // Use a strong secret key

// JWTClaims represents the structure of JWT claims
type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// AuthInterceptor selectively applies JWT authentication to certain methods
func AuthInterceptor(authRequiredMethods []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Check if the method requires authentication
		for _, method := range authRequiredMethods {
			if strings.HasPrefix(info.FullMethod, method) {
				// Extract the metadata from the context
				md, ok := metadata.FromIncomingContext(ctx)
				if !ok {
					return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
				}

				// Get the 'authorization' value from metadata
				authHeader, ok := md["authorization"]
				if !ok || len(authHeader) == 0 {
					return nil, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
				}

				// Check if the token is in the format "Bearer <token>"
				tokenString := authHeader[0]
				if !strings.HasPrefix(tokenString, "Bearer ") {
					return nil, status.Errorf(codes.Unauthenticated, "Invalid token format")
				}

				// Extract the actual token by removing "Bearer "
				tokenString = strings.TrimPrefix(tokenString, "Bearer ")

				// Parse and validate the JWT token
				claims := &JWTClaims{}
				token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					// Ensure the signing method is correct
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, status.Errorf(codes.Unauthenticated, "Unexpected signing method")
					}
					return jwtSecret, nil
				})

				if err != nil || !token.Valid {
					return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
				}

				// Save the claims or user info in the context for further use
				ctx = context.WithValue(ctx, "username", claims.Username)
				break
			}
		}

		// Call the handler
		return handler(ctx, req)
	}
}


// Helper function to extract token from context metadata
func extractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "No metadata found")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "No authorization header found")
	}

	return values[0], nil
}

func GenerateToken(user *userPb.LoginResponse) (string, error) {
	// Token expiration time: 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create custom claims
	claims := &CustomClaims{
		Username: user.UserName,
		Roles: user.Roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create a new JWT token with the claims and the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// Return the generated token
	return tokenString, nil
}
