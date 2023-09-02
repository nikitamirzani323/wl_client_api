package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikitamirzani323/wl_agen_backend_api/configs"
	"github.com/nikitamirzani323/wl_agen_backend_api/db"
	"github.com/nikitamirzani323/wl_agen_backend_api/helpers"
	"github.com/nleeper/goment"
)

func Login_Model(idmasteragen, username, password, ipaddress, timezone string) (bool, string, string, string, error) {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	tglnow, _ := goment.New()

	tbl_mst_member, _, _, _ := Get_mappingdatabase(idmasteragen)

	var idmaster_db, idmember_DB, password_member_DB string
	sql_select := `
			SELECT
			A.idmember, B.idmaster, A.password_member      
			FROM ` + tbl_mst_member + ` as A 
			JOIN ` + configs.DB_tbl_mst_master_agen + ` as B ON B.idmasteragen = A.idmasteragen  
			WHERE A.username_member=$1 
			AND A.idmasteragen=$2 
			AND A.status_member='Y'
		`

	row := con.QueryRowContext(ctx, sql_select, username, idmasteragen)
	switch e := row.Scan(&idmember_DB, &idmaster_db, &password_member_DB); e {
	case sql.ErrNoRows:
		return false, "", "", "", errors.New("Username and Password Not Found")
	case nil:
		flag = true
	default:
		return false, "", "", "", errors.New("Username and Password Not Found")
	}

	hashpass := helpers.HashPasswordMD5(password)

	if hashpass != password_member_DB {
		return false, "", "", "", nil
	}

	if flag {

		sql_update := `
			UPDATE ` + tbl_mst_member + ` 
			SET lastlogin_member=$1, ipaddress_member=$2, timezone_member=$3,    
			update_member=$4,  updatedate_member=$5    
			WHERE idmember=$6  
			AND idmasteragen=$7   
			AND status_member='Y' 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_master_agen_member, "UPDATE",
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), ipaddress, timezone,
			idmember_DB, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idmember_DB, idmasteragen)

		if flag_update {
			flag = true
		} else {
			fmt.Println(msg_update)
		}
	}

	return true, idmaster_db, idmasteragen, idmember_DB, nil
}
