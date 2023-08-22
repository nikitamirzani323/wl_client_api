package controllers

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

func CheckLogin(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Login)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}

	result, idmaster, idmasteragen, idagenadmin, ruleadmin, tipeadmin, err := models.Login_Model(client.Username, client.Password, client.Ipaddress)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	if !result {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Username or Password Not Found",
			})

	} else {
		dataclient := idmaster + "==" + idmasteragen + "==" + idagenadmin + "==" + ruleadmin + "==" + tipeadmin
		dataclient_encr, keymap := helpers.Encryption(dataclient)
		dataclient_encr_final := dataclient_encr + "|" + strconv.Itoa(keymap)
		t, err := helpers.GenerateNewAccessToken(dataclient_encr_final)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{
			"status": fiber.StatusOK,
			"token":  t,
		})

	}
}
func Home(c *fiber.Ctx) error {
	client := new(entities.Home)
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_idmaster, client_idmasteragen, client_idagenadmin, idruleadmin, tipeadmin := helpers.Parsing_Decry(temp_decp, "==")
	fmt.Println("IDMASTER:", client_idmaster)
	fmt.Println("IDMASTER_AGEN:", client_idmasteragen)
	fmt.Println("IDAGENADMIN:", client_idagenadmin)
	fmt.Println("IDAGENADMINRULE:", idruleadmin)
	fmt.Println("TIPE:", tipeadmin)
	fmt.Println(client.Page)

	if tipeadmin == "ADMIN" {
		ruleadmin := models.Get_AdminRule("ruleadmingroup", client_idmasteragen, idruleadmin)
		flag := models.Get_listitemsearch(ruleadmin, ",", client.Page)
		if !flag {
			c.Status(fiber.StatusForbidden)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusForbidden,
				"message": "Anda tidak bisa akses halaman ini",
				"record":  nil,
			})
		} else {
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusOK,
				"message": "ADMIN",
				"record":  nil,
			})
		}
	} else {
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "ADMIN",
			"record":  nil,
		})
	}

}
