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

const Fieldcatebank_home_redis = "LISTCATEBANK_BACKEND"
const Fieldcatebank_home_client_redis = "LISTCATEBANK_FRONTEND"

func CateBankhome(c *fiber.Ctx) error {
	var obj entities.Model_catebank
	var arraobj []entities.Model_catebank
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcatebank_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		catebank_id, _ := jsonparser.GetInt(value, "catebank_id")
		catebank_name, _ := jsonparser.GetString(value, "catebank_name")
		catebank_status, _ := jsonparser.GetString(value, "catebank_status")
		catebank_status_css, _ := jsonparser.GetString(value, "catebank_status_css")
		catebank_create, _ := jsonparser.GetString(value, "catebank_create")
		catebank_update, _ := jsonparser.GetString(value, "catebank_update")

		var objbanktype entities.Model_bankType
		var arraobjbanktype []entities.Model_bankType
		record_banktype_RD, _, _, _ := jsonparser.Get(value, "catebank_list")
		jsonparser.ArrayEach(record_banktype_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			banktype_id, _ := jsonparser.GetString(value, "banktype_id")
			banktype_name, _ := jsonparser.GetString(value, "banktype_name")
			banktype_img, _ := jsonparser.GetString(value, "banktype_img")
			banktype_status, _ := jsonparser.GetString(value, "banktype_status")
			Banktype_status_css, _ := jsonparser.GetString(value, "Banktype_status_css")
			banktype_create, _ := jsonparser.GetString(value, "banktype_create")
			banktype_update, _ := jsonparser.GetString(value, "banktype_update")

			objbanktype.Banktype_id = banktype_id
			objbanktype.Banktype_name = banktype_name
			objbanktype.Banktype_img = banktype_img
			objbanktype.Banktype_status = banktype_status
			objbanktype.Banktype_status_css = Banktype_status_css
			objbanktype.Banktype_create = banktype_create
			objbanktype.Banktype_update = banktype_update
			arraobjbanktype = append(arraobjbanktype, objbanktype)
		})

		obj.Catebank_id = int(catebank_id)
		obj.Catebank_name = catebank_name
		obj.Catebank_status = catebank_status
		obj.Catebank_status_css = catebank_status_css
		obj.Catebank_list = arraobjbanktype
		obj.Catebank_create = catebank_create
		obj.Catebank_update = catebank_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_catebankHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcatebank_home_redis, result, 60*time.Minute)
		fmt.Println("CATEBANK MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("CATEBANK CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func CateBankSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_catebanksave)
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

	// admin, name, status, sData string, idrecord int
	result, err := models.Save_catebank(
		client_admin,
		client.Catebank_name, client.Catebank_status, client.Sdata, client.Catebank_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_catebank()
	return c.JSON(result)
}
func BankTypeSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_banktypesave)
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

	// admin, idrecord, name, img, status, sData string, idcatebank int
	result, err := models.Save_banktype(
		client_admin,
		client.Banktype_id, client.Banktype_name, client.Banktype_img, client.Banktype_status, client.Sdata, client.Banktype_idcatebank)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_catebank()
	return c.JSON(result)
}
func _deleteredis_catebank() {
	val_master := helpers.DeleteRedis(Fieldcatebank_home_redis)
	fmt.Printf("Redis Delete BACKEND CATEBANK : %d", val_master)
}
