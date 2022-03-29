package controllers

import "github.com/guigoebel/user-order-api/user-api/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Orders routes
	s.Router.HandleFunc("/orders", middlewares.SetMiddlewareJSON(s.CreateOrder)).Methods("POST")
	s.Router.HandleFunc("/orders", middlewares.SetMiddlewareJSON(s.GetOrders)).Methods("GET")
	s.Router.HandleFunc("/orders/{id}", middlewares.SetMiddlewareJSON(s.GetOrder)).Methods("GET")
	s.Router.HandleFunc("/orders/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateOrder))).Methods("PUT")
	s.Router.HandleFunc("/orders/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteOrder)).Methods("DELETE")
}
