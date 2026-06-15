package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
)

func init() {
	// Set up Google OAuth2 configuration
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),                    // Get Data from ENV file
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),                // Get Data from ENV file
		RedirectURL:  os.Getenv("CommonURL") + "/auth/google/callback", // Get Data from ENV file
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

}

// Google Login initiates the login process
func GoogleLogin(c *fiber.Ctx) error {
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

// Google Callback handles the callback from Google
func GoogleCallback(c *fiber.Ctx) error {
	s, _ := store.Get(c) // start session
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Code not found")
	}

	// Exchange code for a token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	// Fetch user info
	client := oauthConfig.Client(context.Background(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer userInfo.Body.Close()

	var user struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(userInfo.Body).Decode(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to decode user info")
	}

	// Check if user exists in database

	// Find Duplicate Email in DB

	//Alerts := ""
	//getEmail := user.Email
	//getName := user.Name
	if strings.TrimSpace(user.Name) != "" && strings.TrimSpace(user.Email) != "" {
		s.Set("SocialMerchantName", strings.TrimSpace(user.Name))
		s.Set("SocialMerchantEmail", strings.TrimSpace(user.Email))
		s.Set("SocialType", "Google")
		if err := s.Save(); err != nil {
			return err
		}
		err = MerchantSocialLogin(c)
		if err != nil {
			fmt.Println("issue in Social Media Account")
		}
	}

	return c.SendString(fmt.Sprintf("Welcome, %s!", user.Name))
}
