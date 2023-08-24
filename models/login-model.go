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
	var idmaster_db, idagenmember_DB, password_agenmember_DB string
	sql_select := `
			SELECT
			A.idagenmember, B.idmaster, A.password_agenmember      
			FROM ` + configs.DB_tbl_mst_master_agen_member + ` as A 
			JOIN ` + configs.DB_tbl_mst_master_agen + ` as B ON B.idmasteragen = A.idmasteragen  
			WHERE A.username_agenmember=$1 
			AND A.idmasteragen=$2 
			AND A.status_agenmember='Y'
		`

	row := con.QueryRowContext(ctx, sql_select, username, idmasteragen)
	switch e := row.Scan(&idagenmember_DB, &idmaster_db, &password_agenmember_DB); e {
	case sql.ErrNoRows:
		return false, "", "", "", errors.New("Username and Password Not Found")
	case nil:
		flag = true
	default:
		return false, "", "", "", errors.New("Username and Password Not Found")
	}

	hashpass := helpers.HashPasswordMD5(password)

	if hashpass != password_agenmember_DB {
		return false, "", "", "", nil
	}

	if flag {

		sql_update := `
			UPDATE ` + configs.DB_tbl_mst_master_agen_member + ` 
			SET lastlogin_agenmember=$1, ipaddress_agenmember=$2, timezone_agenmember=$3,    
			update_agenmember=$4,  updatedate_agenmember=$5    
			WHERE idagenmember=$6  
			AND idmasteragen=$7   
			AND status_agenmember='Y' 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_master_agen_member, "UPDATE",
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), ipaddress, timezone,
			idagenmember_DB, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idagenmember_DB, idmasteragen)

		if flag_update {
			flag = true
		} else {
			fmt.Println(msg_update)
		}
	}

	return true, idmaster_db, idmasteragen, idagenmember_DB, nil
}
