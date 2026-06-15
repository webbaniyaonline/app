package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func vals() (int, int) {
	return 3, 7
}

// Handler to get balance by address
func TestAPIS3(c *fiber.Ctx) error {

	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	for j := 0; j < 3; j++ {
		fmt.Println(j)
	}

	for i := range 3 {
		fmt.Println("range", i)
	}

	for {
		fmt.Println("loop")
		break
	}

	for n := range 6 {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}

	return nil

}
