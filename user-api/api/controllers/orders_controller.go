package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/guigoebel/user-order-api/user-api/api/auth"
	"github.com/guigoebel/user-order-api/user-api/api/models"
	"github.com/guigoebel/user-order-api/user-api/api/responses"
	"github.com/guigoebel/user-order-api/user-api/api/utils/formaterror"
)

func (server *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Order := models.Order{}
	err = json.Unmarshal(body, &Order)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Order.Prepare()
	err = Order.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != Order.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	OrderCreated, err := Order.SaveOrder(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, OrderCreated.ID))
	responses.JSON(w, http.StatusCreated, OrderCreated)
}

func (server *Server) GetOrders(w http.ResponseWriter, r *http.Request) {

	order := models.Order{}

	orders, err := order.FindAllOrders(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, orders)
}

func (server *Server) GetOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	oid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	order := models.Order{}

	orderReceived, err := order.FindOrderByID(server.DB, oid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, orderReceived)
}

func (server *Server) UpdateOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the Order id is valid
	oid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Order exist
	Order := models.Order{}
	err = server.DB.Debug().Model(models.Order{}).Where("id = ?", oid).Take(&Order).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Order not found"))
		return
	}

	// If a user attempt to update a Order not belonging to him
	if uid != Order.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data Ordered
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	OrderUpdate := models.Order{}
	err = json.Unmarshal(body, &OrderUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != OrderUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	OrderUpdate.Prepare()
	err = OrderUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	OrderUpdate.ID = Order.ID //this is important to tell the model the Order id to update, the other update field are set above

	OrderUpdated, err := OrderUpdate.UpdateAOrder(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, OrderUpdated)
}

func (server *Server) DeleteOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid Order id given to us?
	oid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Order exist
	Order := models.Order{}
	err = server.DB.Debug().Model(models.Order{}).Where("id = ?", oid).Take(&Order).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this Order?
	if uid != Order.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = Order.DeleteAOrder(server.DB, oid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", oid))
	responses.JSON(w, http.StatusNoContent, "")
}
