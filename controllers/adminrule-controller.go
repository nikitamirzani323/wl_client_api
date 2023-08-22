package controllers

import (
	"log"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nikitamirzani323/wl_agen_backend_api/models"
)

const Fieldadminrule_home_redis = "LISTADMINRULE_AGEN"

func Adminrulehome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_idmasteragen, _, _, _ := helpers.Parsing_Decry(temp_decp, "==")
	log.Println(client_idmasteragen)
	var obj entities.Model_agenadminrule
	var arraobj []entities.Model_agenadminrule
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldadminrule_home_redis + "_" + client_idmasteragen)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		agenadminrule_id, _ := jsonparser.GetInt(value, "agenadminrule_id")
		agenadminrule_name, _ := jsonparser.GetString(value, "agenadminrule_name")
		agenadminrule_rule, _ := jsonparser.GetString(value, "agenadminrule_rule")
		agenadminrule_create, _ := jsonparser.GetString(value, "agenadminrule_create")
		agenadminrule_update, _ := jsonparser.GetString(value, "agenadminrule_update")

		obj.Agenadminrule_id = int(agenadminrule_id)
		obj.Agenadminrule_name = agenadminrule_name
		obj.Agenadminrule_rule = agenadminrule_rule
		obj.Agenadminrule_create = agenadminrule_create
		obj.Agenadminrule_update = agenadminrule_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_adminruleHome(client_idmasteragen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldadminrule_home_redis+"_"+client_idmasteragen, result, 60*time.Minute)
		log.Println("ADMIN RULE MYSQL")
		return c.JSON(result)
	} else {
		log.Println("ADMIN RULE CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func AdminruleSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_agenadminrulesave)
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
	_, client_idmasteragen, client_idagenadmin, _, _ := helpers.Parsing_Decry(temp_decp, "==")
	//admin, idmasteragen, nmrule, rule, sData string, idrecord int
	result, err := models.Save_adminrule(client_idagenadmin, client_idmasteragen,
		client.Agenadminrule_name, client.Agenadminrule_rule, client.Sdata, client.Agenadminrule_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_adminrule(client_idmasteragen)
	return c.JSON(result)
}
func _deleteredis_adminrule(idmasteragen string) {
	val_master := helpers.DeleteRedis(Fieldadminrule_home_redis + "_" + idmasteragen)
	log.Printf("Redis Delete AGEN ADMIN RULE : %d", val_master)

}
