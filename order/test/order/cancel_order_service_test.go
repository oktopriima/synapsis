package order_test

import (
	"fmt"
	"synapsis/order/app/repository"
	"synapsis/order/app/service"
	"synapsis/order/test/mocks"
	inventoryProto "synapsis/proto-definitions/inventory"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	. "gopkg.in/check.v1"
)

func (s *S) TestSuccessfulCancelOrder(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)
	mockInventoryClient.EXPECT().
		ReleaseStock(gomock.Any(), &inventoryProto.ReleaseStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
			OrderId:   1,
		}).
		Return(&inventoryProto.ReleaseStockResponse{
			Success: true,
			Message: "success release stock",
		}, nil)

	s.mock.ExpectBegin()
	expectedRows := s.mock.NewRows([]string{"id", "product_id", "product_name", "quantity", "price", "final_price", "status", "created_at", "updated_at"})
	expectedRows.AddRow(1, dummy.ProductId, dummy.ProductName, dummy.Quantity, dummy.ProductPrice, float64(dummy.Quantity)*dummy.ProductPrice, "confirmed", time.Now(), time.Now())

	expectedSelect := `SELECT * FROM "orders" WHERE id = $1 ORDER BY "orders"."id" LIMIT $2`
	s.mock.ExpectQuery(expectedSelect).
		WithArgs(1, 1).
		WillReturnRows(expectedRows)

	expectedUpdateQuery := `UPDATE "orders" SET "product_id"=$1,"product_name"=$2,"quantity"=$3,"price"=$4,"final_price"=$5,"status"=$6,"created_at"=$7,"updated_at"=$8 WHERE "id" = $9`
	s.mock.ExpectExec(expectedUpdateQuery).
		WithArgs(
			dummy.ProductId,
			dummy.ProductName,
			dummy.Quantity,
			dummy.ProductPrice,
			float64(dummy.Quantity)*dummy.ProductPrice,
			"cancelled",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	// prepare test execution
	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCancelOrderService(s.instance, mockInventoryClient, orderRepo)

	output, err := createOrderService.Execute(s.ctx, service.CancelOrderRequest{
		OrderId: 1,
	})
	c.Assert(err, IsNil)
	c.Assert(output, NotNil)

	s.TearDownTest(c)

}

func (s *S) TestFailedCancelOrderOnNotFoundOrder(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)

	expectedSelect := `SELECT * FROM "orders" WHERE id = $1 ORDER BY "orders"."id" LIMIT $2`
	s.mock.ExpectQuery(expectedSelect).
		WithArgs(1, 1).
		WillReturnError(fmt.Errorf("order not found"))

	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCancelOrderService(s.instance, mockInventoryClient, orderRepo)

	output, err := createOrderService.Execute(s.ctx, service.CancelOrderRequest{
		OrderId: 1,
	})
	c.Assert(err, NotNil)
	c.Assert(output, IsNil)

	s.TearDownTest(c)
}

func (s *S) TestFailedCancelOrderOnAlreadyCancelled(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)
	expectedRows := s.mock.NewRows([]string{"id", "product_id", "product_name", "quantity", "price", "final_price", "status", "created_at", "updated_at"})
	expectedRows.AddRow(1, dummy.ProductId, dummy.ProductName, dummy.Quantity, dummy.ProductPrice, float64(dummy.Quantity)*dummy.ProductPrice, "cancelled", time.Now(), time.Now())

	expectedSelect := `SELECT * FROM "orders" WHERE id = $1 ORDER BY "orders"."id" LIMIT $2`
	s.mock.ExpectQuery(expectedSelect).
		WithArgs(1, 1).
		WillReturnRows(expectedRows)

	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCancelOrderService(s.instance, mockInventoryClient, orderRepo)

	output, err := createOrderService.Execute(s.ctx, service.CancelOrderRequest{
		OrderId: 1,
	})

	c.Assert(err, NotNil)
	c.Assert(output, IsNil)
	c.Assert(err.Error(), Equals, "order already cancelled")

	s.TearDownTest(c)
}

func (s *S) TestFailedCancelOrderOnUnexpectedReleaseStock(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)
	mockInventoryClient.EXPECT().
		ReleaseStock(gomock.Any(), &inventoryProto.ReleaseStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
			OrderId:   1,
		}).
		Return(nil, fmt.Errorf("release stock error"))

	expectedRows := s.mock.NewRows([]string{"id", "product_id", "product_name", "quantity", "price", "final_price", "status", "created_at", "updated_at"})
	expectedRows.AddRow(1, dummy.ProductId, dummy.ProductName, dummy.Quantity, dummy.ProductPrice, float64(dummy.Quantity)*dummy.ProductPrice, "confirmed", time.Now(), time.Now())

	expectedSelect := `SELECT * FROM "orders" WHERE id = $1 ORDER BY "orders"."id" LIMIT $2`
	s.mock.ExpectQuery(expectedSelect).
		WithArgs(1, 1).
		WillReturnRows(expectedRows)

	// prepare test execution
	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCancelOrderService(s.instance, mockInventoryClient, orderRepo)

	output, err := createOrderService.Execute(s.ctx, service.CancelOrderRequest{
		OrderId: 1,
	})
	c.Assert(err, NotNil)
	c.Assert(output, IsNil)

	s.TearDownTest(c)
}

func (s *S) TestFailedCancelOrderOnUnexpectedCancelOrder(c *C) {
	s.SetUpTest(c)

	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockInventoryClient := mocks.NewMockInventoryServiceClient(ctrl)
	mockInventoryClient.EXPECT().
		ReleaseStock(gomock.Any(), &inventoryProto.ReleaseStockRequest{
			ProductId: dummy.ProductId,
			Quantity:  dummy.Quantity,
			OrderId:   1,
		}).
		Return(&inventoryProto.ReleaseStockResponse{
			Success: true,
			Message: "success release stock",
		}, nil)

	s.mock.ExpectBegin()
	expectedRows := s.mock.NewRows([]string{"id", "product_id", "product_name", "quantity", "price", "final_price", "status", "created_at", "updated_at"})
	expectedRows.AddRow(1, dummy.ProductId, dummy.ProductName, dummy.Quantity, dummy.ProductPrice, float64(dummy.Quantity)*dummy.ProductPrice, "confirmed", time.Now(), time.Now())

	expectedSelect := `SELECT * FROM "orders" WHERE id = $1 ORDER BY "orders"."id" LIMIT $2`
	s.mock.ExpectQuery(expectedSelect).
		WithArgs(1, 1).
		WillReturnRows(expectedRows)

	expectedUpdateQuery := `UPDATE "orders" SET "product_id"=$1,"product_name"=$2,"quantity"=$3,"price"=$4,"final_price"=$5,"status"=$6,"created_at"=$7,"updated_at"=$8 WHERE "id" = $9`
	s.mock.ExpectExec(expectedUpdateQuery).
		WithArgs(
			dummy.ProductId,
			dummy.ProductName,
			dummy.Quantity,
			dummy.ProductPrice,
			float64(dummy.Quantity)*dummy.ProductPrice,
			"cancelled",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			1,
		).WillReturnError(fmt.Errorf("failed to cancel order"))

	s.mock.ExpectRollback()

	// prepare test execution
	orderRepo := repository.NewOrderRepository(s.instance)
	createOrderService := service.NewCancelOrderService(s.instance, mockInventoryClient, orderRepo)

	output, err := createOrderService.Execute(s.ctx, service.CancelOrderRequest{
		OrderId: 1,
	})
	c.Assert(err, NotNil)
	c.Assert(output, IsNil)

	s.TearDownTest(c)
}
