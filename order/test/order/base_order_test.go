package order_test

import (
	"context"
	"synapsis/order/database/connection"
	"synapsis/order/test/mocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "gopkg.in/check.v1"
)

type S struct {
	instance connection.DBInstance
	mock     sqlmock.Sqlmock
	ctx      context.Context
}

var dummy = struct {
	ProductId    int64
	ProductName  string
	ProductSKU   string
	ProductDesc  string
	ProductPrice float64
	Quantity     int64
}{
	ProductId:    1,
	ProductName:  "Product 1",
	ProductSKU:   "Product SKU 1",
	ProductDesc:  "Product Desc 1",
	ProductPrice: 10,
	Quantity:     1,
}

func Test(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&S{
	mock:     nil,
	instance: nil,
	ctx:      context.Background(),
})

func (s *S) SetUpTest(c *C) {
	s.instance, s.mock = mocks.Instance() // fresh DB per test
}

func (s *S) TearDownTest(c *C) {
	err := s.mock.ExpectationsWereMet()
	c.Assert(err, IsNil)
}
