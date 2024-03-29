package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/GuiTadeu/mercado-fresh-panic/internal/purchaseOrders"

	"github.com/GuiTadeu/mercado-fresh-panic/cmd/server/controller"
	db "github.com/GuiTadeu/mercado-fresh-panic/cmd/server/database"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/buyers"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/carries"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/employees"
	inboundorders "github.com/GuiTadeu/mercado-fresh-panic/internal/inboundOrders"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/localities"
	productrecords "github.com/GuiTadeu/mercado-fresh-panic/internal/product_records"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/products"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/products/batches"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/sections"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/sellers"
	"github.com/GuiTadeu/mercado-fresh-panic/internal/warehouses"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error to load .env", err)
	}

	storageDB := db.Init()
	server := gin.Default()

	sellerRepository, warehouseRepository, sectionRepository, productRepository, buyerRepository, employeeRepository, inboundOrderRepository, localityRepository, carrieRepository, batchesRepository, productRecordsRepository, purchaseOrdersRepository := buildRepositories(storageDB)

	sellersHandlers(sellerRepository, server)
	warehousesHandlers(warehouseRepository, server)
	sectionHandlers(sectionRepository, server)
	productHandlers(productRepository, server)
	buyerHandlers(buyerRepository, server)
	employeeHandlers(employeeRepository, server)
	inboundOrderHandlers(inboundOrderRepository, employeeRepository, warehouseRepository, server)
	localitiesHandlers(localityRepository, server)
	carriersHandlers(carrieRepository, server)
	productBatchesHandlers(batchesRepository, sectionRepository, productRepository, server)
	productRecordsHandlers(productRecordsRepository, productRepository, server)
	purchaseOrdersHandlers(purchaseOrdersRepository, server)

	port := os.Getenv("MERCADO_FRESH_HOST_PORT")
	server.Run(port)
}

func carriersHandlers(carrierRepository carries.CarrierRepository, server *gin.Engine) {
	carrierService := carries.NewCarrierService(carrierRepository)
	carrierController := controller.NewCarrierController(carrierService)

	CarrierGroup := server.Group("/api/v1/carries")
	CarrierGroup.GET("/reportCarries", carrierController.GetAllCarrierInfo())
	CarrierGroup.POST("/", carrierController.Create())
}

func sellersHandlers(sellerRepository sellers.Repository, server *gin.Engine) {
	sellerService := sellers.NewService(sellerRepository)
	sellerController := controller.NewSeller(sellerService)

	sellerGroup := server.Group("/api/v1/sellers")
	sellerGroup.GET("/", sellerController.FindAll())
	sellerGroup.GET("/:id", sellerController.FindOne())
	sellerGroup.POST("/", sellerController.Create())
	sellerGroup.PATCH("/:id", sellerController.Update())
	sellerGroup.DELETE("/:id", sellerController.Delete())
}

func warehousesHandlers(warehouseRepository warehouses.WarehouseRepository, server *gin.Engine) {
	warehouseService := warehouses.NewService(warehouseRepository)
	warehouseController := controller.NewWarehouseController(warehouseService)

	warehouseGroup := server.Group("/api/v1/warehouses")
	warehouseGroup.GET("/", warehouseController.GetAll())
	warehouseGroup.GET("/:id", warehouseController.Get())
	warehouseGroup.POST("/", warehouseController.Create())
	warehouseGroup.PATCH("/:id", warehouseController.Update())
	warehouseGroup.DELETE("/:id", warehouseController.Delete())
}

func productHandlers(productRepository products.ProductRepository, server *gin.Engine) {

	productService := products.NewProductService(productRepository)
	productHandler := controller.NewProductController(productService)

	productRoutes := server.Group("/api/v1/products")

	productRoutes.GET("/", productHandler.GetAll())
	productRoutes.GET("/:id", productHandler.Get())
	productRoutes.GET("/reportrecords", productHandler.GetAllReportRecords())
	productRoutes.POST("/", productHandler.Create())
	productRoutes.PATCH("/:id", productHandler.Update())
	productRoutes.DELETE("/:id", productHandler.Delete())
}

func productRecordsHandlers(productRecordsRepository productrecords.ProductRecordsRepository, productRepository products.ProductRepository, server *gin.Engine) {

	productRecordsService := productrecords.NewProductRecordsService(productRecordsRepository, productRepository)
	productRecordsHandler := controller.NewProductRecordsController(productRecordsService)

	productRecordsRoutes := server.Group("/api/v1/productRecords")

	productRecordsRoutes.POST("/", productRecordsHandler.Create())
}

func sectionHandlers(sectionRepository sections.SectionRepository, server *gin.Engine) {

	sectionService := sections.NewService(sectionRepository)
	sectionHandler := controller.NewSectionController(sectionService)

	sectionRoutes := server.Group("/api/v1/sections")

	sectionRoutes.GET("/", sectionHandler.GetAll())
	sectionRoutes.GET("/:id", sectionHandler.Get())
	sectionRoutes.POST("/", sectionHandler.Create())
	sectionRoutes.PATCH("/:id", sectionHandler.Update())
	sectionRoutes.DELETE("/:id", sectionHandler.Delete())
}

func employeeHandlers(employeeRepository employees.EmployeeRepository, server *gin.Engine) {

	employeeService := employees.NewEmployeeService(employeeRepository)
	employeeHandler := controller.NewEmployeeController(employeeService)

	employeeRoutes := server.Group("/api/v1/employees")

	employeeRoutes.GET("/", employeeHandler.GetAll())
	employeeRoutes.POST("/", employeeHandler.Create())
	employeeRoutes.DELETE("/:id", employeeHandler.Delete())
	employeeRoutes.GET("/:id", employeeHandler.Get())
	employeeRoutes.PATCH("/:id", employeeHandler.Update())
	employeeRoutes.GET("/reportInboundOrders", employeeHandler.CountInboundOrders())
}

func inboundOrderHandlers(inboundOrderRepository inboundorders.InboundOrderRepository, employeeRepository employees.EmployeeRepository, warehouseRepository warehouses.WarehouseRepository, server *gin.Engine) {

	inboundOrderService := inboundorders.NewInboundOrderService(employeeRepository, warehouseRepository, inboundOrderRepository)

	cInboundOrders := controller.NewInboundOrderController(inboundOrderService)

	inboundOrderRoutes := server.Group("/api/v1/inboundOrders")

	inboundOrderRoutes.POST("/", cInboundOrders.Create())
}

func buyerHandlers(buyerRepository buyers.BuyerRepository, server *gin.Engine) {

	sBuyers := buyers.NewBuyerService(buyerRepository)
	cBuyers := controller.NewBuyerController(sBuyers)

	buyerRoutes := server.Group("/api/v1/buyers")

	buyerRoutes.GET("/reportPurchaseOrders", cBuyers.CountPurchaseOrdersByBuyers())
	buyerRoutes.GET("/", cBuyers.GetAll())
	buyerRoutes.GET("/:id", cBuyers.Get())
	buyerRoutes.POST("/", cBuyers.Create())
	buyerRoutes.PATCH("/:id", cBuyers.Update())
	buyerRoutes.DELETE("/:id", cBuyers.Delete())
}

func localitiesHandlers(localityRepository localities.Repository, server *gin.Engine) {
	localityService := localities.NewService(localityRepository)
	localityController := controller.NewLocality(localityService)

	localityGroup := server.Group("/api/v1/localities")
	localityGroup.POST("/", localityController.Create())
	localityGroup.GET("/reportSellers", localityController.GetLocalityInfo())
}

func productBatchesHandlers(
	pbr batches.ProductBatchRepository,
	sr sections.SectionRepository,
	pr products.ProductRepository,
	server *gin.Engine,
) {
	batchesService := batches.NewProductBatchesService(pbr, sr, pr)
	batchesController := controller.NewProductBatchController(batchesService)

	batchesGroup := server.Group("/api/v1/productBatches")
	batchesGroup.POST("/", batchesController.Create())

	server.GET("/api/v1/sections/reportProducts", batchesController.CountProductsBySections())
}

func buildRepositories(storageDB *sql.DB) (
	sellers.Repository,
	warehouses.WarehouseRepository,
	sections.SectionRepository,
	products.ProductRepository,
	buyers.BuyerRepository,
	employees.EmployeeRepository,
	inboundorders.InboundOrderRepository,
	localities.Repository,
	carries.CarrierRepository,
	batches.ProductBatchRepository,
	productrecords.ProductRecordsRepository,
	purchaseOrders.PurchaseOrdersRepository) {

	sellerRepository := sellers.NewRepository(storageDB)
	warehouseRepository := warehouses.NewRepository(storageDB)
	sectionRepository := sections.NewRepository(storageDB)
	productRepository := products.NewProductRepository(storageDB)
	buyerRepository := buyers.NewBuyerRepository(storageDB)
	employeeRepository := employees.NewRepository(storageDB)
	inboundOrderRepository := inboundorders.NewRepository(storageDB)
	localityRepository := localities.NewRepository(storageDB)
	carrieRepository := carries.NewCarrierRepository(storageDB)
	productRecordsRepository := productrecords.NewProductRecordsRepository(storageDB)
	productBatchesRepository := batches.NewProductBatchRepository(storageDB)
	purchaseOrdersRepository := purchaseOrders.NewPurchaseOrdersRepository(storageDB)

	return sellerRepository, warehouseRepository, sectionRepository, productRepository, buyerRepository, employeeRepository, inboundOrderRepository, localityRepository, carrieRepository, productBatchesRepository, productRecordsRepository, purchaseOrdersRepository
}

func purchaseOrdersHandlers(purchaseOrdersRepository purchaseOrders.PurchaseOrdersRepository, server *gin.Engine) {
	purchaseOrderService := purchaseOrders.NewPurchaseOrdersService(purchaseOrdersRepository)
	purchaseOrderHandler := controller.NewPurchaseOrderController(purchaseOrderService)

	purchaseOrderRoutes := server.Group("/api/v1/")

	purchaseOrderRoutes.POST("/purchaseOrders", purchaseOrderHandler.Create())

}
