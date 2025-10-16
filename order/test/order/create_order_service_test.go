package order_test

import (
	"fmt"
	"synapsis/order/app/repository"
	"synapsis/order/app/service"
	"synapsis/order/test/mocks"
	"synapsis/proto-definitions/inventory"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	. "gopkg.in/check.v1"
)

func (s *S) TestSuccessfulCreateOrder(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)

	// call check stock
	mockInventoryClient.EXPECT().
		CheckStock(gomock.Any(), &inventory.CheckStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
		}).
		Return(&inventory.CheckStockResponse{
			Product: &inventory.Product{
				Id:          dummy.ProductId,
				Name:        dummy.ProductName,
				Sku:         dummy.ProductSKU,
				Description: dummy.ProductDesc,
				Price:       float32(dummy.ProductPrice),
			},
			Stock: &inventory.Stock{
				Id:             1,
				ProductId:      dummy.ProductId,
				AvailableStock: 100,
				ReservedStock:  0,
				TotalStock:     100,
			},
			Quantity:    dummy.Quantity,
			IsAvailable: true,
		}, nil)

	s.mock.ExpectBegin()

	expectedQuery := `INSERT INTO "orders" ("product_id","product_name","quantity","price","final_price","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`
	s.mock.ExpectQuery(expectedQuery).
		WithArgs(
			dummy.ProductId,
			dummy.ProductName,
			dummy.Quantity,
			dummy.ProductPrice,
			float64(dummy.Quantity)*dummy.ProductPrice,
			"confirmed",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))

	mockInventoryClient.EXPECT().
		ReserveStock(gomock.Any(), &inventory.ReserveStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
			OrderId:   1,
		}).
		Return(&inventory.ReserveStockResponse{
			Success: true,
			Message: "success reserve stock",
		}, nil)

	s.mock.ExpectCommit()

	// prepare test execution
	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCreateOrderService(s.instance, mockInventoryClient, orderRepo)

	req := new(service.CreateOrderRequest)
	req.Product.Id = dummy.ProductId
	req.Product.Quantity = dummy.Quantity
	// execution
	output, err := createOrderService.Serve(req, s.ctx)
	c.Assert(err, IsNil)
	c.Assert(output, NotNil)
	c.Assert(output.Order, NotNil)
	c.Assert(output.Order.Status, Equals, "confirmed")

	s.TearDownTest(c)
}

func (s *S) TestFailedCreateOrderOnUnavailableStock(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)

	// call check stock
	mockInventoryClient.EXPECT().
		CheckStock(gomock.Any(), &inventory.CheckStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
		}).
		Return(&inventory.CheckStockResponse{
			Product: &inventory.Product{
				Id:          dummy.ProductId,
				Name:        dummy.ProductName,
				Sku:         dummy.ProductSKU,
				Description: dummy.ProductDesc,
				Price:       float32(dummy.ProductPrice),
			},
			Stock: &inventory.Stock{
				Id:             1,
				ProductId:      dummy.ProductId,
				AvailableStock: 0,
				ReservedStock:  100,
				TotalStock:     100,
			},
			Quantity:    0,
			IsAvailable: false,
		}, nil)

	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCreateOrderService(s.instance, mockInventoryClient, orderRepo)

	req := new(service.CreateOrderRequest)
	req.Product.Id = dummy.ProductId
	req.Product.Quantity = dummy.Quantity
	// execution
	output, err := createOrderService.Serve(req, s.ctx)

	c.Assert(err, NotNil)
	c.Assert(output, IsNil)

	s.TearDownTest(c)
}

func (s *S) TestFailedCreateOrderOnGrpcUnavailable(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)

	// call check stock
	mockInventoryClient.EXPECT().
		CheckStock(gomock.Any(), &inventory.CheckStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
		}).
		Return(nil, fmt.Errorf("not available"))

	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCreateOrderService(s.instance, mockInventoryClient, orderRepo)

	req := new(service.CreateOrderRequest)
	req.Product.Id = dummy.ProductId
	req.Product.Quantity = dummy.Quantity
	// execution
	output, err := createOrderService.Serve(req, s.ctx)

	c.Assert(err, NotNil)
	c.Assert(output, IsNil)

	s.TearDownTest(c)
}

func (s *S) TestFailedCreateOrderOnInvalidStock(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)

	// call check stock
	mockInventoryClient.EXPECT().
		CheckStock(gomock.Any(), &inventory.CheckStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
		}).
		Return(&inventory.CheckStockResponse{
			Product: &inventory.Product{
				Id:          dummy.ProductId,
				Name:        dummy.ProductName,
				Sku:         dummy.ProductSKU,
				Description: dummy.ProductDesc,
				Price:       float32(dummy.ProductPrice),
			},
			Stock: &inventory.Stock{
				Id:             1,
				ProductId:      dummy.ProductId,
				AvailableStock: 100,
				ReservedStock:  0,
				TotalStock:     100,
			},
			Quantity:    dummy.Quantity,
			IsAvailable: true,
		}, nil)

	s.mock.ExpectBegin()

	expectedQuery := `INSERT INTO "orders" ("product_id","product_name","quantity","price","final_price","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`
	s.mock.ExpectQuery(expectedQuery).
		WithArgs(
			dummy.ProductId,
			dummy.ProductName,
			dummy.Quantity,
			dummy.ProductPrice,
			float64(dummy.Quantity)*dummy.ProductPrice,
			"confirmed",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))

	mockInventoryClient.EXPECT().
		ReserveStock(gomock.Any(), &inventory.ReserveStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
			OrderId:   1,
		}).
		Return(nil, fmt.Errorf("stock not available"))

	s.mock.ExpectRollback()

	// prepare test execution
	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCreateOrderService(s.instance, mockInventoryClient, orderRepo)

	req := new(service.CreateOrderRequest)
	req.Product.Id = dummy.ProductId
	req.Product.Quantity = dummy.Quantity
	// execution
	output, err := createOrderService.Serve(req, s.ctx)
	c.Assert(err, NotNil)
	c.Assert(output, IsNil)

	s.TearDownTest(c)
}
