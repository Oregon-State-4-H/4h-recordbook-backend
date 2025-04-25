package middleware

import (
	"context"
	"log"
	"net/url"
	"net/http"
	"strings"
	"time"
	"os"
	"io"

	"github.com/gin-gonic/gin"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gwatts/gin-adapter"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
	Name string `json:"nickname"`
}

// Userinfo data
type UserInfo struct {
	Sub string `json:"sub"`
	Name string `json:"name"`
	Email string `json:"email"`
	Verified bool `json:"email_verified"`
	Picture string `json:"picture"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() gin.HandlerFunc {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return adapter.Wrap(middleware.CheckJWT)
}

func GetToken() gin.HandlerFunc {
	return func(c* gin.Context){
		token := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		customClaims := token.CustomClaims.(*CustomClaims)
		// Get the user info with the JWT
		req, _ := http.NewRequest("GET", token.RegisteredClaims.Audience[1], nil)
		req.Header.Set("Authorization", c.GetHeader("Authorization"))

		userinfo := UserInfo{};

		client := &http.Client{};
		res, _ := client.Do(req)
		buf := new(strings.Builder);
		if(res.StatusCode == 200){
			io.Copy(buf, res.Body);
			json.Unmarshall([]byte(buf.String()), &userinfo);
		}
		
		c.Set("user_id", userinfo.Sub)
		c.Set("user_name", userinfo.Name)
		c.Set("user_email", userinfo.Email)
		c.Set("user_picture", userinfo.Picture)
		c.Set("user_verified", userinfo.Verified)
	}
}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
