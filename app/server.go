package app

import (
	"fmt"
	"frappuccino/internal/dal"
	"frappuccino/internal/handler"
	"frappuccino/internal/service"
	"frappuccino/tools"
	"log"
	"net/http"
)

func StartServer() {
	err := tools.CheckJsonFiles()
	if err != nil {
		log.Fatal(err)
		return
	}
	tools.InitLogger()

	mux := http.NewServeMux()

	invRepo := dal.InventoryFilePath()
	menuRepo := dal.MenuFilePath()
	orderRepo := dal.NewOrderRepoImpl()

	invSvc := service.NewInventoryService(invRepo, menuRepo)
	invHandle := handler.NewInventoryHandler(invSvc)

	menuSvc := service.NewMenuService(menuRepo, invRepo)
	menuHandle := handler.NewMenuHandler(menuSvc)

	orderSvc := service.NewOrderService(orderRepo, menuRepo, invRepo)
	orderHandle := handler.NewOrderHandler(orderSvc)

	mux.HandleFunc("POST /inventory", invHandle.AddInventoryOfHandle())
	mux.HandleFunc("GET /inventory", invHandle.ReadInventoryOfHandle())
	mux.HandleFunc("GET /inventory/{id}", invHandle.ReadInventoryOfHandleByID())
	mux.HandleFunc("PUT /inventory/{id}", invHandle.UpdateInventoryOfHandle())
	mux.HandleFunc("DELETE /inventory/{id}", invHandle.DeleteInventoryOfHandle())

	mux.HandleFunc("POST /menu", menuHandle.AddMenuOfHandle())
	mux.HandleFunc("GET /menu", menuHandle.ReadMenuOfHandle())
	mux.HandleFunc("GET /menu/{id}", menuHandle.ReadMenuOfHandleByID())
	mux.HandleFunc("PUT /menu/{id}", menuHandle.UpdateMenuOfHandle())
	mux.HandleFunc("DELETE /menu/{id}", menuHandle.DeleteMenuOfHandle())

	mux.HandleFunc("POST /orders", orderHandle.AddOrderOfHandle())
	mux.HandleFunc("GET /orders", orderHandle.ReadOrderOfHandle())
	mux.HandleFunc("GET /orders/{id}", orderHandle.ReadOrderOfHandleByID())
	mux.HandleFunc("PUT /orders/{id}", orderHandle.UpdateOrderOfHandle())
	mux.HandleFunc("DELETE /orders/{id}", orderHandle.DeleteOrderOfHandle())
	mux.HandleFunc("POST /orders/{id}/close", orderHandle.CloseOrderOfHandle())

	mux.HandleFunc("GET /reports/total-sales", orderHandle.TotalSalesOfHandle())
	mux.HandleFunc("GET /reports/popular-items", orderHandle.PopularItemsOfHandle())

	fmt.Println("Server listening on port", *tools.Port)
	log.Fatal(http.ListenAndServe(":"+*tools.Port, mux))
}
