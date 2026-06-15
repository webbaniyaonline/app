package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Functions for FaceBook Login
// Struct for FaceBook Login
type FBRawData struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	ID        string `json:"id"`
	LastName  string `json:"last_name"`
	Name      string `json:"name"`
}

func init() {

	// Initialize Facebook provider
	goth.UseProviders(
		facebook.New(
			os.Getenv("FACEBOOK_APP_ID"), // Facebook App ID
			os.Getenv("FACEBOOK_APP_SECRET"),
			os.Getenv("CommonURL")+"/auth/facebook/callback", // Facebook App Secret
		),
	)

}

// Facebook Login initiates the login process
func FacebookLogin(c *fiber.Ctx) error {

	c.Request().URI().SetQueryString("provider=facebook") // Add provider query parameter

	fasthttpadaptor.NewFastHTTPHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r) // Start Facebook OAuth flow
	}))(c.Context())

	return nil
}

// Facebook Callback handles the callback from Google
func FacebookCallback(c *fiber.Ctx) error {
	s, _ := store.Get(c) // start session

	var user goth.User
	var err error

	// Add provider query parameter
	c.Request().URI().SetQueryString("provider=facebook")

	fasthttpadaptor.NewFastHTTPHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Complete the OAuth flow and get user info
		user, err = gothic.CompleteUserAuth(w, r)
	}))(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error completing Facebook OAuth: " + err.Error())
	}

	//fmt.Println(user)
	//email := user.Email
	//fmt.Println(email)
	//name := user.Name
	//fmt.Println(name)

	if strings.TrimSpace(user.Name) != "" && strings.TrimSpace(user.Email) != "" {
		s.Set("SocialMerchantName", strings.TrimSpace(user.Name))
		s.Set("SocialMerchantEmail", strings.TrimSpace(user.Email))
		s.Set("SocialType", "Facebook")
		if err := s.Save(); err != nil {
			return err
		}
		err = MerchantSocialLogin(c)
		if err != nil {
			fmt.Println("issue in Social Media Account")
		}

	}

	//return c.JSON(user)
	return nil

}
