package usecases

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/gateway"
	paymentuc "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/usecases"
	productentity "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	productuc "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
	productorderentity "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
	productorderuc "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type UseCases struct {
	orderGateway        *gateway.Gateway
	productUseCase      productuc.UseCases
	productOrderUseCase productorderuc.UseCases
	paymentUseCase      *paymentuc.UseCases
}

func Build(
	orderGateway *gateway.Gateway,
	productUseCase productuc.UseCases,
	productOrderUseCase productorderuc.UseCases,
	paymentUseCase paymentuc.UseCases,
) *UseCases {
	return &UseCases{
		orderGateway:        orderGateway,
		productUseCase:      productUseCase,
		productOrderUseCase: productOrderUseCase,
		paymentUseCase:      &paymentUseCase,
	}
}

func (u *UseCases) CreateCompleteOrder(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	var productIds []string
	for _, item := range orderDTO.Products {
		productIds = append(productIds, item.ProductID)
	}

	products, findErr := u.productUseCase.FindByIDs(ctx, productIds)
	if findErr != nil {
		return "", findErr
	}
	if len(products) != len(orderDTO.Products) {
		return "", &apperror.NotFoundError{
			Msg: "some products not found",
		}
	}

	populatedOrder := generateOrderByProducts(orderDTO, products)
	createdOrder, createErr := u.orderGateway.Create(ctx, populatedOrder.Build())
	if createErr != nil {
		return "", createErr
	}

	productOrders, _ := generateProductOrderFromOrderAndProducts(createdOrder.ID, orderDTO.Products, products)
	_, createBulkErr := u.productOrderUseCase.CreateBulk(ctx, productOrders)
	if createBulkErr != nil {
		return "", createBulkErr
	}

	payment, paymentErr := u.paymentUseCase.CreateByOrderID(ctx, createdOrder.ID)
	if paymentErr != nil {
		return "", paymentErr
	}

	return payment.QrCode, nil
}

// todo verify if we can move this function to other package
func generateOrderByProducts(orderDTO dto.CreateOrderDTO, products []productentity.Product) entity.Order {
	totalPrice, preparingTime := getOrderInfoFromProducts(products, orderDTO)

	return entity.Order{
		CustomerID:    orderDTO.CustomerID,
		Status:        enum.OrderStatusAwaitingPayment,
		Price:         totalPrice,
		PreparingTime: preparingTime,
	}
}

func getOrderInfoFromProducts(products []productentity.Product, orderDTO dto.CreateOrderDTO) (float64, uint) {
	var totalPrice float64
	var preparingTime uint

	for _, item := range orderDTO.Products {
		for _, product := range products {
			if product.Id == item.ProductID {
				totalPrice += product.Price * float64(item.Quantity)
				preparingTime += product.PreparingTime
			}
		}
	}

	return totalPrice, preparingTime
}

func generateProductOrderFromOrderAndProducts(
	orderID string,
	orderProductInfo []dto.OrderProductInfo,
	products []productentity.Product,
) ([]productorderentity.ProductOrder, error) {
	var res []productorderentity.ProductOrder

	for _, product := range products {
		for _, item := range orderProductInfo {
			if product.Id == item.ProductID {
				res = append(res, productorderentity.ProductOrder{
					OrderID:   orderID,
					ProductID: product.Id,
					Quantity:  item.Quantity,
					UnitPrice: product.Price,
				})
			}
		}
	}

	return res, nil
}

func (u *UseCases) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.orderGateway.Create(ctx, order)
}

func (u *UseCases) GetAll(ctx context.Context) ([]entity.Order, error) {
	return u.orderGateway.GetAll(ctx)
}

func (u *UseCases) GetPanel(ctx context.Context, status []string) ([]entity.Order, error) {
	return u.orderGateway.GetPanel(ctx, status)
}

func (u *UseCases) FindByID(ctx context.Context, id string) (entity.Order, error) {
	return u.orderGateway.FindByID(ctx, id)
}

func (u *UseCases) Update(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.orderGateway.Update(ctx, order)
}
