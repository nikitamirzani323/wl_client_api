package controllers

import (
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

const Fieldprovider_home_redis = "LISTPROVIDER_BACKEND"
const Fieldprovider_home_client_redis = "LISTPROVIDER_FRONTEND"

func Providerhome(c *fiber.Ctx) error {
	var obj entities.Model_provider
	var arraobj []entities.Model_provider
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldprovider_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		provider_id, _ := jsonparser.GetInt(value, "provider_id")
		provider_name, _ := jsonparser.GetString(value, "provider_name")
		provider_owner, _ := jsonparser.GetString(value, "provider_owner")
		provider_email, _ := jsonparser.GetString(value, "provider_email")
		provider_phone1, _ := jsonparser.GetString(value, "provider_phone1")
		provider_phone2, _ := jsonparser.GetString(value, "provider_phone2")
		provider_urlwebsite, _ := jsonparser.GetString(value, "provider_urlwebsite")
		provider_status, _ := jsonparser.GetString(value, "provider_status")
		provider_status_css, _ := jsonparser.GetString(value, "provider_status_css")
		provider_create, _ := jsonparser.GetString(value, "provider_create")
		provider_update, _ := jsonparser.GetString(value, "provider_update")

		obj.Provider_id = int(provider_id)
		obj.Provider_name = provider_name
		obj.Provider_owner = provider_owner
		obj.Provider_email = provider_email
		obj.Provider_phone1 = provider_phone1
		obj.Provider_phone2 = provider_phone2
		obj.Provider_urlwebsite = provider_urlwebsite
		obj.Provider_status = provider_status
		obj.Provider_status_css = provider_status_css
		obj.Provider_create = provider_create
		obj.Provider_update = provider_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_providerHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldprovider_home_redis, result, 60*time.Minute)
		fmt.Println("PROVIDER MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("PROVIDER CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func ProviderSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_providersave)
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
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, _, client_admin, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_provider(
		client_admin,
		client.Provider_name, client.Provider_owner,
		client.Provider_email, client.Provider_phone1, client.Provider_phone2, client.Provider_urlwebsite, client.Provider_status, client.Sdata, client.Provider_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_provider()
	return c.JSON(result)
}

func _deleteredis_provider() {
	val_master := helpers.DeleteRedis(Fieldprovider_home_redis)
	fmt.Printf("Redis Delete BACKEND PROVIDER : %d", val_master)

	val_client := helpers.DeleteRedis(Fieldprovider_home_client_redis)
	fmt.Printf("Redis Delete CLIENT PROVIDER : %d", val_client)

}
