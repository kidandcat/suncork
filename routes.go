package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

func setupRoutes(server *http.Server) {
	// Static
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.Handle("/products/", http.StripPrefix("/products/", http.FileServer(http.Dir("./products"))))

	// Routes
	http.Handle("/", mw(mwSession(handlerIndex)))
	http.Handle("/language", mw(mwSession(handlerLanguage)))
	http.Handle("/logout", mw(mwSession(handlerLogout)))
	http.Handle("/shestakova/login", mw(mwSession(handlerLogin)))
	http.Handle("/shestakova/loginform", mw(mwSession(handlerLoginForm)))
	http.Handle("/cookies", mw(mwSession(handlerCookies)))
	http.Handle("/about", mw(mwSession(handlerAbout)))
	http.Handle("/product", mw(mwSession(handlerProduct)))
	http.Handle("/product/cart", mw(mwSession(handlerProductToCart)))
	http.Handle("/product/remove", mw(mwSession(handlerProductRemoveFromCart)))
	http.Handle("/cart", mw(mwSession(handlerCart)))
	http.Handle("/buy", mw(mwSession(handlerPayment)))
	http.Handle("/product/new", mw(mwSession(handlerProductNew)))
	http.Handle("/product/delete", mw(mwSession(handlerProductDelete)))
	http.Handle("/shestakova", mw(mwSession(handlerAdmin)))
	http.Handle("/shestakova/translations", mw(mwSession(handlerTranslations)))
	http.Handle("/shestakova/translations/new", mw(mwSession(handlerTranslationsPost)))
}

// Handlers
func handlerCookies(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerCookies")
	data["cart"] = s.cart
	renderTemplate(w, "cookies.html", data, s)
}

func handlerAbout(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerAbout")
	data["cart"] = s.cart
	renderTemplate(w, "about-"+s.lang+".html", data, s)
}

func handlerLogout(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerLogout")
	s.admin = false
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerLanguage(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerLanguage")
	lang, ok := r.URL.Query()["lang"]
	if ok {
		s.lang = lang[0]
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlerTranslations(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerTranslations")
	if !s.admin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "trans.html", data, s)
}

func handlerTranslationsPost(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerTranslationsPost")
	if !s.admin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	trans := map[string]string{}
	for key, value := range r.PostForm {
		trans[key] = value[0]
	}
	setTrans(s.lang, trans)
	http.Redirect(w, r, "/shestakova/translations", http.StatusSeeOther)
}

func handlerLoginForm(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerLoginForm")
	r.ParseForm()
	user := r.Form["User"][0]
	pass := r.Form["Pass"][0]
	if user == "olga" && pass == "te.quiero.mucho." {
		s.admin = true
	}
	http.Redirect(w, r, "/shestakova/login", http.StatusSeeOther)
}

func handlerLogin(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerLogin")
	if s.admin {
		http.Redirect(w, r, "/shestakova", http.StatusSeeOther)
	}
	renderTemplate(w, "login.html", data, s)
}

func handlerIndex(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerIndex")
	data["products"] = productFind()
	data["cart"] = s.cart
	renderTemplate(w, "index.html", data, s)
}

func handlerProduct(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerProduct")
	productID, ok := r.URL.Query()["id"]
	if ok {
		data["product"] = productFindByID(productID[0])
		data["cart"] = s.cart
		renderTemplate(w, "product.html", data, s)
	} else {
		w.WriteHeader(404)
	}
}

func handlerProductRemoveFromCart(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerProductRemoveFromCart")
	productID, ok := r.URL.Query()["id"]
	if ok {
		data["product"] = productFindByID(productID[0])
		for i, v := range s.cart {
			if productID[0] == v.Product.ID {
				s.cart = append(s.cart[:i], s.cart[i+1:]...)
				break
			}
		}
		data["cart"] = s.cart
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
	} else {
		w.WriteHeader(404)
	}
}

func handlerProductToCart(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerProductToCart")
	productID, ok := r.URL.Query()["id"]
	optName, ok1 := r.URL.Query()["optName"]
	optVal, ok2 := r.URL.Query()["optVal"]
	if ok {
		product := productFindByID(productID[0])
		pwo := ProductWithOption{
			Product: product,
		}
		if ok1 && ok2 {
			for i, v := range optName {
				pwo.Options = append(pwo.Options, Option{
					Name:  v,
					Value: optVal[i],
				})
			}
		}
		s.cart = append(s.cart, pwo)

		http.Redirect(w, r, "/product?id="+product.ID, http.StatusSeeOther)
	} else {
		w.WriteHeader(404)
	}
}

func handlerProductDelete(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerProductDelete")
	if !s.admin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	productID, ok := r.URL.Query()["id"]
	if ok {
		product := productFindByID(productID[0])
		product.delete()
	}
	http.Redirect(w, r, "/shestakova", http.StatusSeeOther)
}

func handlerCart(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerCart")
	data["products"] = s.cart
	data["cart"] = s.cart
	renderTemplate(w, "cart.html", data, s)
}

func handlerAdmin(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerAdmin")
	if !s.admin {
		http.Redirect(w, r, "/shestakova/login", http.StatusSeeOther)
		return
	}
	productID, ok := r.URL.Query()["edit"]
	if ok {
		product := productFindByID(productID[0])
		data["product"] = product
	} else {
		data["product"] = Product{}
	}
	data["products"] = productFind()
	renderTemplate(w, "admin.html", data, s)
}

func handlerProductNew(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerProductNew")
	if !s.admin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	reader, err := r.MultipartReader()

	filenames := []string{}

	if err != nil {
		panic(err)
	}

	product := Product{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FormName() != "" {
			switch part.FormName() {
			case "ID":
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				product.ID = buf.String()
			case "NameEn":
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				product.NameEn = buf.String()
			case "Options[]":
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				str := buf.String()
				fmt.Println("OPTIONS", str)
				choices := strings.Split(strings.Split(str, "::")[1], ",")
				stocksString := strings.Split(strings.Split(str, "::")[3], ",")
				fmt.Println("::::::: stocksString", stocksString)
				stocksInt := []int{}
				for _, v := range stocksString {
					i, e := strconv.Atoi(v)
					if e == nil {
						stocksInt = append(stocksInt, i)
					} else {
						stocksInt = append(stocksInt, 0)
					}
				}
				for len(stocksInt) < len(choices) {
					stocksInt = append(stocksInt, 0)
				}
				fmt.Println("::::::: stocksInt", stocksInt)
				product.Options = append(product.Options, ProductOption{
					Name:    strings.Split(str, "::")[0],
					Choices: choices,
					Stock:   stocksInt,
					Lang:    strings.Split(str, "::")[2],
				})
			case "NameEs":
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				product.NameEs = buf.String()
			case "DescriptionEn":
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				product.DescriptionEn = buf.String()
			case "DescriptionEs":
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				product.DescriptionEs = buf.String()
			case "Price":
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				p, err := strconv.ParseFloat(buf.String(), 64)
				if err == nil {
					product.Price = int(p * 100)
				}
			}
		}

		if part.FileName() == "" {
			continue
		}
		randname := randStringRunes(8) + filepath.Ext(part.FileName())
		dst, err := os.Create("./products/" + randname)
		defer dst.Close()

		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dst, part); err != nil {
			panic(err)
		}

		filenames = append(filenames, "/products/"+randname)
	}

	product.Image = filenames
	product.save()
	http.Redirect(w, r, "/shestakova", http.StatusSeeOther)
}

func handlerPayment(w http.ResponseWriter, r *http.Request, data Data, s *Session) {
	print("CONTROLLER handlerPayment")
	r.ParseForm()
	shippingID, ok := r.Form["shippingOptionId"]
	if !ok {
		fmt.Println("error reading shippingOptionId")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	shipping, ok := getShippingOptionByID(shippingID[0])
	if !ok {
		fmt.Println("error reading getShippingOptionByID", shippingID, shipping)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	token := r.Form["token"][0]

	shippingAddress := ShippingAddress{
		AddressLine: r.Form["address"][0],
		City:        r.Form["city"][0],
		Country:     r.Form["country"][0],
		PostalCode:  r.Form["postalCode"][0],
	}

	stripe.Key = SK_STRIPE_KEY

	amount := 0
	for _, p := range s.cart {
		amount += p.Product.Price
	}
	amount += shipping.Amount

	// Charge ready
	params := &stripe.ChargeParams{
		Amount:   uint64(amount),
		Currency: "eur",
		Desc:     "123",
	}
	params.SetSource(token)

	// Create Payment
	createdTimestamp := time.Now().Unix()
	var clientIp string
	if len(r.Form["clientIp"]) == 0 {
		clientIp = ""
	} else {
		clientIp = r.Form["clientIp"][0]
	}
	payment := Payment{
		ID:         IDGeneratorPayment(),
		ClientIP:   clientIp,
		Created:    createdTimestamp,
		PayerEmail: r.Form["payerEmail"][0],
		PayerName:  r.Form["payerName"][0],
		PayerPhone: r.Form["payerPhone"][0],
		Token:      r.Form["token"][0],
	}

	// Create Order
	order := Order{
		ID:              IDGeneratorOrder(),
		PaymentID:       payment.ID,
		Products:        s.cart,
		ShippingOption:  shipping,
		ShippingAddress: shippingAddress,
	}

	_, err := charge.New(params)

	s.cart = []ProductWithOption{}
	http.Redirect(w, r, "/", http.StatusSeeOther)

	if err != nil {
		payment.Error = err.Error()
		order.Error = err.Error()
		mail("olga@suncork.net", "Payment ERROR!", mailTemplatePayment(&payment, &order))
		mail("jairo@suncork.net", "Payment ERROR!", mailTemplatePayment(&payment, &order))
	} else {
		mail("olga@suncork.net", "Payment SUCCESS", mailTemplatePayment(&payment, &order))
		mail("jairo@suncork.net", "Payment SUCCESS", mailTemplatePayment(&payment, &order))
	}

	order.save()
	payment.save()

	print("Order saved")
	print("Payment saved")

	// Rest Stock
	print("Stock rest start")
	for _, v := range order.Products {
		print("Stock rest: product")
		for _, vv := range v.Options {
			print("Stock rest: options")
			for iii, _ := range v.Product.Options {
				print("OPT", "-"+v.Product.Options[iii].Name+"-", "::", "-"+vv.Name+"-")
				if v.Product.Options[iii].Name == vv.Name {
					v.Product.Options[iii] = v.Product.Options[iii].restStock(vv.Name, vv.Value)
					v.Product.save()
				}
			}
		}
	}
	print("Stock rested")
}
