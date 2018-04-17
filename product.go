package main

import (
	"fmt"
	"strings"
)

type ProductOption struct {
	Name    string
	Lang    string
	Choices []string
	Stock   []int
}

type Product struct {
	ID            string
	NameEn        string
	NameEs        string
	Price         int
	DescriptionEn string
	DescriptionEs string
	Image         []string
	Options       []ProductOption
}

func productFind() []Product {
	rows, e := db.Query("SELECT * FROM products")
	res := []Product{}
	if err(e) {
		return res
	}
	var images string
	var options string
	defer rows.Close()
	for rows.Next() {
		prod := Product{}
		e := rows.Scan(&prod.ID, &prod.NameEn, &prod.NameEs, &prod.Price, &images, &prod.DescriptionEn, &prod.DescriptionEs, &options)
		if err(e) {
			return res
		}
		prod.Image = strings.Split(images, "::")
		prod.Options = optionsFromJSON(options)
		fmt.Println("raw options", options)
		fmt.Println("optionsFromJSON", prod.Options)
		res = append(res, prod)
	}
	return res
}

func productFindByID(id string) Product {
	sql := fmt.Sprintf("SELECT * FROM products WHERE id = '%s'", id)
	rows, e := db.Query(sql)
	prod := Product{}
	if err(e) {
		return prod
	}

	var images string
	var options string

	defer rows.Close()
	if rows.Next() {
		e := rows.Scan(&prod.ID, &prod.NameEn, &prod.NameEs, &prod.Price, &images, &prod.DescriptionEn, &prod.DescriptionEs, &options)
		if err(e) {
			return prod
		}
		prod.Image = strings.Split(images, "::")
		prod.Options = optionsFromJSON(options)
	}
	return prod
}

func initProductTables() {
	_, e := db.Exec(`CREATE TABLE IF NOT EXISTS products (
							id VARCHAR (255) PRIMARY KEY, 
							name_en VARCHAR (255), 
							name_es VARCHAR (255), 
							price INT, 
							image VARCHAR (255), 
							description_en VARCHAR (2000),
							description_es VARCHAR (2000),
							options VARCHAR (2000))`)
	crash(e)
}

func (opt ProductOption) restStock(name string, choice string) ProductOption {
	print("restStock", opt)
	for i, v := range opt.Choices {
		print("restStock: Compare names", opt.Name, name)
		print("restStock: Compare choices", v, choice)
		if opt.Name == name && v == choice {
			print("Rest stock")
			opt.Stock[i]--
		}
	}
	return opt
}

func (product Product) save() {
	p := productFindByID(product.ID)
	if p.NameEs == "" {
		if product.ID == "" {
			product.ID = randStringRunes(6)
		}
		images := strings.Join(product.Image, "::")
		options := optionsToJSON(product.Options)
		sql := fmt.Sprintf(`INSERT INTO products (id, name_en, name_es, price, image, description_en, description_es, options) 
							VALUES ('%s', '%s', '%s', %d, '%s', '%s', '%s', '%s')`,
			product.ID, product.NameEn, product.NameEs, product.Price, images, product.DescriptionEn, product.DescriptionEs, options)
		print("Product Save New SQL", sql)
		_, e := db.Exec(sql)
		err(e)
	} else {
		sql := `UPDATE products SET `
		sql += fmt.Sprintf(` name_en = '%s', `, product.NameEn)
		sql += fmt.Sprintf(` name_es = '%s', `, product.NameEs)
		sql += fmt.Sprintf(` price = %d, `, product.Price)
		if len(product.Image) != 0 {
			images := strings.Join(product.Image, "::")
			sql += fmt.Sprintf(` image = '%s', `, images)
		}
		if len(product.Options) != 0 {
			options := optionsToJSON(product.Options)
			sql += fmt.Sprintf(` options = '%s', `, options)
		}
		sql += fmt.Sprintf(` description_en = '%s', `, product.DescriptionEn)
		sql += fmt.Sprintf(` description_es = '%s' `, product.DescriptionEs)
		sql += fmt.Sprintf(` WHERE id = '%s'`, product.ID)
		print("Product Save Update SQL", sql)
		_, e := db.Exec(sql)
		err(e)
	}
}

func (product Product) delete() {
	sql := fmt.Sprintf("DELETE from products WHERE id = '%s'", product.ID)
	_, e := db.Exec(sql)
	err(e)
}
