package main

import (
	"fmt"
)

type Payment struct {
	ID         string
	PayerEmail string `json:"payerEmail"`
	PayerPhone string `json:"payerPhone"`
	PayerName  string `json:"payerName"`
	Token      string `json:"id"`
	ClientIP   string `json:"client_ip"`
	Created    int64  `json:"created"`
	Error      string
}

func paymentFind() []Payment {
	rows, e := db.Query("SELECT * FROM payments")
	res := []Payment{}
	if err(e) {
		return res
	}
	defer rows.Close()
	for rows.Next() {
		prod := Payment{}
		e := rows.Scan(&prod.ID, &prod.PayerEmail, &prod.PayerPhone, &prod.PayerName, &prod.Token, &prod.ClientIP, &prod.Created, &prod.Error)
		if err(e) {
			return res
		}
		res = append(res, prod)
	}
	return res
}

func paymentFindByID(id string) Payment {
	sql := fmt.Sprintf("SELECT * FROM payments WHERE id = '%s'", id)
	rows, e := db.Query(sql)
	prod := Payment{}
	if err(e) {
		return prod
	}

	defer rows.Close()
	if rows.Next() {
		e := rows.Scan(&prod.ID, &prod.PayerEmail, &prod.PayerPhone, &prod.PayerName, &prod.Token, &prod.ClientIP, &prod.Created, &prod.Error)
		if err(e) {
			return prod
		}
	}
	return prod
}

func initPaymentTables() {
	_, e := db.Exec(`CREATE TABLE IF NOT EXISTS payments (
							id VARCHAR(255) PRIMARY KEY, 
							payeremail VARCHAR(255),
							payerphone VARCHAR(255),
							payername VARCHAR(255),
							token VARCHAR(255),
							clientip VARCHAR(255), 
							created INT,
							error VARCHAR(255))
	`)
	crash(e)
}

func (payment Payment) save() {
	sql := fmt.Sprintf(`INSERT INTO payments 
						(id, payeremail, payerphone, payername, token, clientip, created, error) 
						VALUES 
						('%s', '%s', '%s', '%s', '%s', '%s', %d, '%s')
	`, payment.ID, payment.PayerEmail, payment.PayerPhone, payment.PayerName, payment.Token, payment.ClientIP, payment.Created, payment.Error)
	_, e := db.Exec(sql)
	err(e)
}

func (payment Payment) delete() {
	sql := fmt.Sprintf("DELETE from payments WHERE id = '%s'", payment.ID)
	_, e := db.Exec(sql)
	err(e)
}
