package models

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
)

const database_curr_local = configs.DB_tbl_mst_curr

func Fetch_currHome() (helpers.Response, error) {
	var obj entities.Model_curr
	var arraobj []entities.Model_curr
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcurr , multipliercurr 
			FROM ` + database_curr_local + `  
			ORDER BY idcurr ASC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcurr_db         string
			multipliercurr_db float32
		)

		err = row.Scan(&idcurr_db, &multipliercurr_db)

		helpers.ErrorCheck(err)

		obj.Curr_id = idcurr_db
		obj.Curr_multiplier = multipliercurr_db
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
