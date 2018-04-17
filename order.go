package main

import (
	"fmt"
	"strconv"
)

type ShippingAddress struct {
	AddressLine string `json:"addressLine"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PostalCode  string `json:"postalCode"`
}

type Order struct {
	ID              string
	PaymentID       string
	Products        []ProductWithOption
	ShippingOption  ShippingOption  `json:"shippingOption"`
	ShippingAddress ShippingAddress `json:"shippingAddress"`
	Error           string
}

func orderFind() []Order {
	print("Find orders")
	rows, e := db.Query("SELECT * FROM orders")
	res := []Order{}
	if err(e) {
		return res
	}
	defer rows.Close()
	var shippingOptionID string
	var shippingAddressID string
	for rows.Next() {
		prod := Order{}
		e := rows.Scan(&prod.ID, &prod.PaymentID, &prod.Products, shippingOptionID, shippingAddressID, &prod.Error)
		if err(e) {
			return res
		}
		prod.ShippingOption, _ = getShippingOptionByID(shippingOptionID)
		prod.ShippingAddress, _ = getShippingAddressByID(shippingAddressID)
		res = append(res, prod)
	}
	print("Orders found", strconv.Itoa(len(res)))
	return res
}

func orderFindByID(id string) Order {
	print("Find order by id", id)
	prod := Order{}
	sql := fmt.Sprintf("SELECT * FROM orders WHERE id = '%s'", id)
	rows, e := db.Query(sql)
	if err(e) {
		return prod
	}

	defer rows.Close()
	var shippingOptionID string
	var shippingAddressID string
	if rows.Next() {
		e := rows.Scan(&prod.ID, &prod.PaymentID, &prod.Products, shippingOptionID, shippingAddressID, &prod.Error)
		if err(e) {
			return prod
		}
		prod.ShippingOption, _ = getShippingOptionByID(shippingOptionID)
		prod.ShippingAddress, _ = getShippingAddressByID(shippingAddressID)
	}
	return prod
}

func initOrderTables() {
	_, e := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
							id VARCHAR(255) PRIMARY KEY, 
							payment VARCHAR(255), 
							products VARCHAR(255),
							shippingoption VARCHAR(255),
							shippingaddress VARCHAR(255),
							error VARCHAR(255))
	`)
	crash(e)
}

func (order Order) save() {
	products := ""
	for i, v := range order.Products {
		if i != len(order.Products)-1 {
			products += v.Product.ID + "::"
		} else {
			products += v.Product.ID
		}
	}
	sql := fmt.Sprintf("INSERT INTO orders (id, payment, products, shippingoption, shippingaddress, error) VALUES ('%s', '%s', '%s', '%s', '%s', '%s')", order.ID, order.PaymentID, products, order.ShippingOption.ID, order.ShippingAddress.toString(), order.Error)
	_, e := db.Exec(sql)
	err(e)
}

func (order Order) delete() {
	sql := fmt.Sprintf("DELETE from orders WHERE id = '%s'", order.ID)
	_, e := db.Exec(sql)
	err(e)
}
