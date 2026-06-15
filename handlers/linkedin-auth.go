package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

type LiUserProfile struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Locale        struct {
		Country  string `json:"country"`
		Language string `json:"language"`
	} `json:"locale"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

var oauth2Config *oauth2.Config
var oauth2StateString = "AQU-xTZQPM_NMwciMH_x1XhmF2drsqxUWWP1CsHUDjtToUMEmDgJyGfpVOSeRxgvnMTv2e6a2aZRJdsvopJEQDHO9WGebd8rp3-dtMByJAwrwkSKnpLkqG_tMZo_CzIVblyUl9MRxg0qyKGHukBNYs0dfvcAO0X5x0gAURxByGPhxbWSCiqQh1c0MdYL5Th7lAGnobr6YKURY90-gqcrlYmk7Bwgd_EUXSYqdT6-yE9wYNoRyLwyHwi3iR17Erxd9ZVBVW7_NItL9Vn3XfY11oidmYuF-0lT-75VHQMVymmD8bWCv3gtyunSXv8YG5d1gcF0BOCQamKCVFwyETOwWpvUevcMtg" // Use a randomly generated string in production
func init() {

	// Initialize LinkedIn OAuth2 configuration
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
		ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("CommonURL") + "/auth/linkedin/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     linkedin.Endpoint,
	}

}

// Linked Login initiates the login process
func LinkedInLogin(c *fiber.Ctx) error {

	// Generate LinkedIn login URL with state parameter
	authURL := oauth2Config.AuthCodeURL(oauth2StateString, oauth2.AccessTypeOffline)
	return c.Redirect(authURL)
}

// Linked Callback handles the callback from Google
func LinkedInCallback(c *fiber.Ctx) error {
	s, _ := store.Get(c) // start session
	// Check the state parameter to prevent CSRF attacks
	if c.Query("state") != oauth2StateString {
		return c.Status(fiber.StatusBadRequest).SendString("State mismatch")
	}

	// Get the authorization code from the query string
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Code not found")
	}

	// Exchange the authorization code for an access token
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get token: " + err.Error())
	}

	// Use the access token to fetch the user's LinkedIn profile
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://api.linkedin.com/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch user profile: " + err.Error())
	}
	defer resp.Body.Close()

	// Read and process the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response: " + err.Error())
	}

	// Parse JSON response into struct
	var userProfile LiUserProfile
	err = json.Unmarshal(body, &userProfile)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// Access specific fields
	fmt.Printf("User Profile: %+v\n", userProfile)
	fmt.Println("Name:", userProfile.Name)
	fmt.Println("Email:", userProfile.Email)
	fmt.Println("Picture:", userProfile.Picture)

	if strings.TrimSpace(userProfile.Name) != "" && strings.TrimSpace(userProfile.Email) != "" {
		s.Set("SocialMerchantName", strings.TrimSpace(userProfile.Name))
		s.Set("SocialMerchantEmail", strings.TrimSpace(userProfile.Email))
		s.Set("SocialType", "LinkedIn")
		if err := s.Save(); err != nil {
			return err
		}
		err = MerchantSocialLogin(c)
		if err != nil {
			fmt.Println("issue in Social Media Account")
		}

	}

	// Output user profile
	return nil

}
