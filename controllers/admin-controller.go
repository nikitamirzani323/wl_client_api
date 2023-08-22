package controllers

import (
	"fmt"
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

const Fieldadmin_home_redis = "LISTADMIN_AGEN"

func Adminhome(c *fiber.Ctx) error {
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	_, client_idmasteragen, _, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	var obj entities.Model_admin
	var arraobj []entities.Model_admin
	var obj_listruleadmin entities.Model_adminrule
	var arraobj_listruleadmin []entities.Model_adminrule
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldadmin_home_redis + "_" + client_idmasteragen)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	listruleadmin_RD, _, _, _ := jsonparser.Get(jsonredis, "listruleadmin")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		admin_id, _ := jsonparser.GetString(value, "admin_id")
		admin_idrule, _ := jsonparser.GetInt(value, "admin_idrule")
		admin_tipe, _ := jsonparser.GetString(value, "admin_tipe")
		admin_nmrule, _ := jsonparser.GetString(value, "admin_nmrule")
		admin_username, _ := jsonparser.GetString(value, "admin_username")
		admin_nama, _ := jsonparser.GetString(value, "admin_nama")
		admin_phone1, _ := jsonparser.GetString(value, "admin_phone1")
		admin_phone2, _ := jsonparser.GetString(value, "admin_phone2")
		admin_lastipaddres, _ := jsonparser.GetString(value, "admin_lastipaddres")
		admin_lastlogin, _ := jsonparser.GetString(value, "admin_lastlogin")
		admin_status, _ := jsonparser.GetString(value, "admin_status")
		admin_status_css, _ := jsonparser.GetString(value, "admin_status_css")
		admin_create, _ := jsonparser.GetString(value, "admin_create")
		admin_update, _ := jsonparser.GetString(value, "admin_update")

		obj.Admin_id = admin_id
		obj.Admin_idrule = int(admin_idrule)
		obj.Admin_tipe = admin_tipe
		obj.Admin_nmrule = admin_nmrule
		obj.Admin_username = admin_username
		obj.Admin_nama = admin_nama
		obj.Admin_phone1 = admin_phone1
		obj.Admin_phone2 = admin_phone2
		obj.Admin_lastlogin = admin_lastlogin
		obj.Admin_lastipaddres = admin_lastipaddres
		obj.Admin_status = admin_status
		obj.Admin_status_css = admin_status_css
		obj.Admin_create = admin_create
		obj.Admin_update = admin_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listruleadmin_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		adminrule_idruleadmin, _ := jsonparser.GetInt(value, "adminrule_idruleadmin")
		adminrule_nmruleadmin, _ := jsonparser.GetString(value, "adminrule_nmruleadmin")

		obj_listruleadmin.Adminrule_idruleadmin = int(adminrule_idruleadmin)
		obj_listruleadmin.Adminrule_nmruleadmin = adminrule_nmruleadmin
		arraobj_listruleadmin = append(arraobj_listruleadmin, obj_listruleadmin)
	})
	if !flag {
		result, err := models.Fetch_adminHome(client_idmasteragen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldadmin_home_redis+"_"+client_idmasteragen, result, 60*time.Minute)
		fmt.Println("ADMIN AGEN MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("ADMIN AGEN CACHE")
		return c.JSON(fiber.Map{
			"status":        fiber.StatusOK,
			"message":       "Success",
			"record":        arraobj,
			"listruleadmin": arraobj_listruleadmin,
			"time":          time.Since(render_page).String(),
		})
	}
}

func AdminSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_adminsave)
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
	_, client_idmasteragen, client_admin, _, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_adminHome(
		client_admin,
		client.Admin_id, client_idmasteragen, client.Admin_username, client.Admin_password,
		client.Admin_nama, client.Admin_phone1, client.Admin_phone2, client.Admin_status,
		client.Sdata, client.Admin_idrule)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_admin(client_idmasteragen)
	return c.JSON(result)
}
func _deleteredis_admin(idmasteragen string) {
	val_master := helpers.DeleteRedis(Fieldadmin_home_redis + "_" + idmasteragen)
	log.Printf("Redis Delete AGEN ADMIN : %d", val_master)

}
