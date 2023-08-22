package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/entities"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

const database_provider_local = configs.DB_tbl_mst_provider

func Fetch_providerHome() (helpers.Response, error) {
	var obj entities.Model_provider
	var arraobj []entities.Model_provider
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idprovider , nmprovider, nmownerprovider,  emailprovider, phone1provider, phone2provider, urlwebsite, statusprovider, 
			createprovider, to_char(COALESCE(createdateprovider,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatecateprovider, to_char(COALESCE(updatedateprovider,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_provider_local + `  
			ORDER BY createdateprovider DESC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idprovider_db                                                                                                               int
			nmprovider_db, nmownerprovider_db, emailprovider_db, phone1provider_db, phone2provider_db, urlwebsite_db, statusprovider_db string
			createprovider_db, createdateprovider_db, updatecateprovider_db, updatedateprovider_db                                      string
		)

		err = row.Scan(&idprovider_db, &nmprovider_db, &nmownerprovider_db, &emailprovider_db, &phone1provider_db, &phone2provider_db,
			&urlwebsite_db, &statusprovider_db,
			&createprovider_db, &createdateprovider_db, &updatecateprovider_db, &updatedateprovider_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		status_css := configs.STATUS_CANCEL
		if createprovider_db != "" {
			create = createprovider_db + ", " + createdateprovider_db
		}
		if updatecateprovider_db != "" {
			update = updatecateprovider_db + ", " + updatedateprovider_db
		}
		if statusprovider_db == "Y" {
			status_css = configs.STATUS_COMPLETE
		}

		obj.Provider_id = idprovider_db
		obj.Provider_name = nmprovider_db
		obj.Provider_owner = nmownerprovider_db
		obj.Provider_email = emailprovider_db
		obj.Provider_phone1 = phone1provider_db
		obj.Provider_phone2 = phone2provider_db
		obj.Provider_urlwebsite = urlwebsite_db
		obj.Provider_status = statusprovider_db
		obj.Provider_status_css = status_css
		obj.Provider_create = create
		obj.Provider_update = update
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
func Save_provider(admin, name, owner, email, phone1, phone2, urlwebsite, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + database_provider_local + ` (
					idprovider , nmprovider, nmownerprovider,
					emailprovider, phone1provider, phone2provider, urlwebsite, statusprovider,
					createprovider, createdateprovider  
				) values (
					$1, $2, $3,   
					$4, $5, $6, $7, $8,
					$9, $10
				)
			`
		field_column := database_provider_local + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, database_provider_local, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), name, owner, email, phone1, phone2, urlwebsite, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_provider_local + `  
				SET nmprovider=$1, nmownerprovider=$2, 
				emailprovider=$3, phone1provider=$4, phone2provider=$5, urlwebsite=$6, statusprovider=$7,    
				updatecateprovider=$8, updatedateprovider=$9     
				WHERE idprovider=$10   
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_provider_local, "UPDATE",
			name, owner, email, phone1, phone2, urlwebsite, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
