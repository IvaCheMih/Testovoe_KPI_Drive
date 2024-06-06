package dto

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func GetRequestFact(c *fiber.Ctx) (GetFactRequest, error) {
	body := c.Body()

	var request GetFactRequest

	err := json.Unmarshal(body, &request)

	return request, err
}

func CreateRequestFact(c *fiber.Ctx) (CreateFactRequest, error) {
	body := c.Body()

	var request CreateFactRequest

	err := json.Unmarshal(body, &request)

	return request, err
}
