package hotelbooking

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

/* Types */

type Invoice struct {
	RoomID     *int    `json:"room_id"`
	CustomerID *int    `json:"customer_id"`
	CheckIn    *string `json:"in"`
	CheckOut   *string `json:"out"`
}

type CreateInvoiceResponseData struct {
	ID int `json:"id"`
}

/* API */

func CreateInvoice(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var invoice Invoice
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("CreateInvoice :", err)

		SendBadRequestWithData(w)
		return
	}

	err = json.Unmarshal(body, &invoice)

	if err != nil {
		log.Println("CreateInvoice :", err)

		SendBadRequestWithData(w)
		return
	}

	if invoice.RoomID == nil || invoice.CustomerID == nil || invoice.CheckIn == nil || invoice.CheckOut == nil {
		SendBadRequestWithData(w)
		return
	}

	checkIn, err := time.Parse("02-01-2006", *invoice.CheckIn)

	if err != nil {
		SendBadRequestWithData(w)
		return
	}

	checkOut, err := time.Parse("02-01-2006", *invoice.CheckOut)

	if err != nil {
		SendBadRequestWithData(w)
		return
	}

	if !checkIn.Before(checkOut) {
		SendBadRequestWithData(w)
		return
	}

	var price, totalPrice, dummy int
	statement, err := db.Prepare("SELECT price FROM room WHERE id = ?")

	if err != nil {
		log.Println("CreateInvoice :", err)
		return
	}

	defer statement.Close()

	err = statement.QueryRow(*invoice.RoomID).Scan(&price)

	if err != nil {
		log.Println("CreateInvoice :", err)

		SendBadRequestWithData(w)
		return
	}

	statement, err = db.Prepare("SELECT id FROM customer WHERE id = ?")

	if err != nil {
		log.Println("CreateInvoice :", err)
		return
	}

	err = statement.QueryRow(*invoice.CustomerID).Scan(&dummy)

	if err != nil {
		log.Println("CreateInvoice :", err)

		SendBadRequestWithData(w)
		return
	}

	statement, err = db.Prepare("INSERT INTO invoice VALUES (0, ?, ?, ?, ?, ?, 0, 0)")

	if err != nil {
		log.Println("CreateInvoice :", err)
		return
	}

	totalPrice = (int(checkOut.Sub(checkIn).Hours()) / 24) * price
	res, err := statement.Exec(*invoice.RoomID, *invoice.CustomerID, checkIn, checkOut, totalPrice)

	if err != nil {
		log.Println("CreateInvoice :", err)
		return
	}

	id, _ := res.LastInsertId()

	data := CreateInvoiceResponseData{
		ID: int(id),
	}

	SendOKWithData(w, data)
}
