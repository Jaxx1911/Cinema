package service

import (
	"TTCS/src/common/log"
	"TTCS/src/common/mail"
	"TTCS/src/core/domain"
	"TTCS/src/infra/cache"
	"TTCS/src/present/httpui/request"
	"TTCS/src/present/httpui/response"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	paymentRepo    domain.PaymentRepo
	orderRepo      domain.OrderRepo
	ticketRepo     domain.TicketRepo
	userRepo       domain.UserRepo
	movieRepo      domain.MovieRepo
	cinemaRepo     domain.CinemaRepo
	roomRepo       domain.RoomRepo
	seatRepo       domain.SeatRepo
	showtimeRepo   domain.ShowtimeRepo
	comboRepo      domain.ComboRepository
	orderComboRepo domain.OrderComboRepository
	cache          *cache.RedisCache
	mailService    *mail.GmailService
}

func NewPaymentService(
	paymentRepo domain.PaymentRepo,
	orderRepo domain.OrderRepo,
	ticketRepo domain.TicketRepo,
	userRepo domain.UserRepo,
	movieRepo domain.MovieRepo,
	cinemaRepo domain.CinemaRepo,
	roomRepo domain.RoomRepo,
	seatRepo domain.SeatRepo,
	showtimeRepo domain.ShowtimeRepo,
	comboRepo domain.ComboRepository,
	orderComboRepo domain.OrderComboRepository,
	cache *cache.RedisCache,
	mailService *mail.GmailService,
) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		orderRepo:      orderRepo,
		ticketRepo:     ticketRepo,
		userRepo:       userRepo,
		movieRepo:      movieRepo,
		cinemaRepo:     cinemaRepo,
		roomRepo:       roomRepo,
		seatRepo:       seatRepo,
		showtimeRepo:   showtimeRepo,
		comboRepo:      comboRepo,
		orderComboRepo: orderComboRepo,
		cache:          cache,
		mailService:    mailService,
	}
}

func (p *PaymentService) HandleCallback(ctx context.Context, callback request.PaymentCallback) (*domain.Payment, error) {
	caller := "PaymentService.HandleCallback"
	if p.cache.Exists(callback.Payment.TransactionId) {
		err := fmt.Errorf("[%v] block duplicate transaction id %v", caller, callback.Payment.TransactionId)
		return nil, err
	}
	_ = p.cache.Set(callback.Payment.TransactionId, true, 60)
	oid, err := uuid.Parse(callback.Payment.Content)
	if err != nil {
		_, err = p.paymentRepo.Create(ctx, &domain.Payment{
			UserID:        nil,
			OrderID:       nil,
			TransactionID: callback.Payment.TransactionId,
			Status:        "failed",
			Amount:        callback.Payment.Amount,
			PaymentTime:   time.Now(),
		})
		log.Error(ctx, "[%v] invalid content %+v", caller, err)
		return nil, err
	}
	order, err := p.orderRepo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if order.TotalPrice != callback.Payment.Amount {
		err := errors.New("invalid amount")
		_, err = p.paymentRepo.Create(ctx, &domain.Payment{
			UserID:        &order.UserID,
			OrderID:       &order.ID,
			TransactionID: callback.Payment.TransactionId,
			Status:        "failed",
			Amount:        callback.Payment.Amount,
			PaymentTime:   time.Now(),
		})
		log.Error(ctx, "[%v] invalid amount %+v", caller, err)
		return nil, err
	}

	payment, err := p.paymentRepo.Create(ctx, &domain.Payment{
		UserID:        &order.UserID,
		OrderID:       &order.ID,
		TransactionID: callback.Payment.TransactionId,
		Status:        "success",
		Amount:        callback.Payment.Amount,
		PaymentTime:   time.Now(),
	})
	order.Status = "success"

	order, err = p.orderRepo.Update(ctx, order)

	if err != nil {
		return nil, err
	}

	tickets, err := p.ticketRepo.FindByOrderID(ctx, order.ID)

	if err != nil {
		return nil, err
	}
	for i := range tickets {
		tickets[i].Status = "success"
	}
	tickets, err = p.ticketRepo.UpdateBatch(ctx, tickets)
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.GetById(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	emailData, err := p.prepareEmailData(ctx, order, user)
	if err != nil {
		log.Error(ctx, "[%v] failed to prepare email data %+v", caller, err)
	} else {
		err = p.mailService.SendEmailOAuth2("Đặt vé thành công", user.Email, emailData, "booking-success.txt")
		if err != nil {
			log.Error(ctx, "[%v] failed to send email %+v", caller, err)
		}
	}

	return payment, nil
}

func (p *PaymentService) prepareEmailData(ctx context.Context, order *domain.Order, user *domain.User) (*emailData, error) {
	tickets, err := p.ticketRepo.FindByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	if len(tickets) == 0 {
		return nil, errors.New("no tickets found")
	}
	showtime, err := p.showtimeRepo.GetById(ctx, order.ShowtimeID)
	if err != nil {
		return nil, err
	}
	movie, err := p.movieRepo.GetById(ctx, showtime.MovieID)
	if err != nil {
		return nil, err
	}
	room, err := p.roomRepo.GetById(ctx, showtime.RoomID)
	if err != nil {
		return nil, err
	}
	cinema, err := p.cinemaRepo.FindByID(ctx, room.CinemaID)
	if err != nil {
		return nil, err
	}
	var seatNames []string
	for _, ticket := range tickets {
		seat, err := p.seatRepo.GetById(ctx, ticket.SeatID)
		if err != nil {
			continue
		}
		seatName := fmt.Sprintf("%s%d", seat.RowNumber, seat.SeatNumber)
		seatNames = append(seatNames, seatName)
	}
	orderCombos, err := p.getOrderCombos(ctx, order.ID)
	if err != nil {
		log.Error(ctx, "failed to get order combos: %v", err)
		orderCombos = []comboInfo{} // Không có combo thì để empty
	}
	emailData := &emailData{
		UserName:    user.Name,
		OrderID:     order.ID.String(),
		MovieTitle:  movie.Title,
		CinemaName:  cinema.Name,
		RoomName:    room.Name,
		ShowDate:    showtime.StartTime.Format("02/01/2006"),
		ShowTime:    showtime.StartTime.Add(7 * time.Hour).Format("15:04"),
		Seats:       strings.Join(seatNames, ", "),
		Combos:      orderCombos,
		TotalAmount: fmt.Sprintf("%.0f", order.TotalPrice),
	}
	return emailData, nil
}

func (p *PaymentService) getOrderCombos(ctx context.Context, orderID uuid.UUID) ([]comboInfo, error) {
	orderCombos, err := p.orderComboRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var comboInfos []comboInfo
	for _, orderCombo := range orderCombos {
		combo, err := p.comboRepo.FindByID(ctx, orderCombo.ComboID)
		if err != nil {
			continue
		}

		comboInfos = append(comboInfos, comboInfo{
			Name:     combo.Name,
			Quantity: orderCombo.Quantity,
			Price:    fmt.Sprintf("%.0f", orderCombo.TotalPrice),
		})
	}

	return comboInfos, nil
}

func (p *PaymentService) GetPaymentsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Payment, error) {
	payments, err := p.paymentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (p *PaymentService) GetPaymentsByCinemaId(ctx context.Context, cinemaID uuid.UUID, req request.GetPaymentsByCinemaRequest) ([]response.PaymentCinemaDetail, error) {
	// Set end date to end of day
	endDate := req.EndDate.Add(24*time.Hour - time.Second)

	payments, err := p.paymentRepo.GetByCinemaIDAndDateRange(ctx, cinemaID, req.StartDate, endDate)
	if err != nil {
		return nil, err
	}

	var paymentDetails []response.PaymentCinemaDetail
	for _, payment := range payments {
		if payment.OrderID == nil {
			continue
		}

		// Lấy order
		order, err := p.orderRepo.GetByID(ctx, *payment.OrderID)
		if err != nil {
			continue
		}

		// Tính tổng combo amount
		orderCombos, err := p.orderComboRepo.GetByOrderID(ctx, order.ID)
		totalComboPrice := float64(0)
		if err == nil {
			for _, combo := range orderCombos {
				totalComboPrice += combo.TotalPrice
			}
		}

		// Lấy thông tin discount nếu có
		var discountInfo *response.DiscountInfo
		if order.DiscountID != nil && payment.Order.Discount != nil {
			discountInfo = &response.DiscountInfo{
				ID:         payment.Order.Discount.ID,
				Code:       payment.Order.Discount.Code,
				Percentage: payment.Order.Discount.Percentage,
			}
		}

		paymentDetail := response.ToPaymentCinemaDetailResponse(payment, *order, totalComboPrice, discountInfo)
		paymentDetails = append(paymentDetails, paymentDetail)
	}

	return paymentDetails, nil
}

func (p *PaymentService) AcceptAll(ctx context.Context) ([]domain.Payment, error) {
	orders, err := p.orderRepo.GetAllPendingOrders(ctx)
	if err != nil {
		return nil, err
	}
	var payments []domain.Payment
	for _, order := range orders {
		payment, err := p.HandleCallback(ctx, request.PaymentCallback{
			Payment: request.Payment{
				TransactionId:   "tbbank" + strconv.Itoa(rand.Intn(100000000000)),
				Content:         order.ID.String(),
				Amount:          order.TotalPrice,
				Date:            time.Now(),
				Gate:            "tp-bank",
				AccountReceiver: "10001259387",
			},
		})
		if err != nil {
			return nil, err
		}
		payments = append(payments, *payment)
	}
	return payments, nil
}

func (p *PaymentService) GetList(ctx context.Context, req request.GetListPaymentRequest) ([]response.PaymentWithCustomerAndDiscount, int64, error) {
	// Set default values for pagination
	req.Page.SetDefaults()

	payments, total, err := p.paymentRepo.GetList(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	// Tính tổng combo price cho tất cả orders
	orderCombos := make(map[uuid.UUID]float64)
	for _, payment := range payments {
		if payment.OrderID != nil {
			combos, err := p.orderComboRepo.GetByOrderID(ctx, *payment.OrderID)
			if err == nil {
				totalComboPrice := float64(0)
				for _, combo := range combos {
					totalComboPrice += combo.TotalPrice
				}
				orderCombos[*payment.OrderID] = totalComboPrice
			}
		}
	}

	return response.ToPaymentsWithCustomerAndDiscountResponse(payments, orderCombos), total, nil
}

type emailData struct {
	UserName    string
	OrderID     string
	MovieTitle  string
	CinemaName  string
	RoomName    string
	ShowDate    string
	ShowTime    string
	Seats       string
	Combos      []comboInfo
	TotalAmount string
}

type comboInfo struct {
	Name     string
	Quantity int
	Price    string
}
