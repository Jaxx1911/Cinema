package service

import (
	"TTCS/src/common/log"
	"TTCS/src/common/mail"
	"TTCS/src/core/domain"
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
		mailService:    mailService,
	}
}

func (p *PaymentService) HandleCallback(ctx context.Context, callback request.PaymentCallback) (*domain.Payment, error) {
	caller := "PaymentService.HandleCallback"
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
	for _, ticket := range tickets {
		ticket.Status = "success"
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
		ShowTime:    showtime.StartTime.Format("15:04"),
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

func (p *PaymentService) GetPaymentsByCinemaId(ctx context.Context, cinemaID uuid.UUID, req request.GetPaymentsByCinemaRequest) ([]response.PaymentDetail, error) {
	// Set end date to end of day
	endDate := req.EndDate.Add(24*time.Hour - time.Second)

	payments, err := p.paymentRepo.GetByCinemaIDAndDateRange(ctx, cinemaID, req.StartDate, endDate)
	if err != nil {
		return nil, err
	}

	var paymentDetails []response.PaymentDetail
	for _, payment := range payments {
		if payment.OrderID == nil {
			continue
		}

		// Lấy order
		order, err := p.orderRepo.GetByID(ctx, *payment.OrderID)
		if err != nil {
			continue
		}

		// Lấy user
		user, err := p.userRepo.GetById(ctx, order.UserID)
		if err != nil {
			continue
		}

		// Lấy showtime
		showtime, err := p.showtimeRepo.GetById(ctx, order.ShowtimeID)
		if err != nil {
			continue
		}

		// Lấy movie
		movie, err := p.movieRepo.GetById(ctx, showtime.MovieID)
		if err != nil {
			continue
		}

		// Lấy room
		room, err := p.roomRepo.GetById(ctx, showtime.RoomID)
		if err != nil {
			continue
		}

		// Lấy tickets và tạo danh sách ghế
		tickets, err := p.ticketRepo.FindByOrderID(ctx, order.ID)
		if err != nil {
			continue
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

		// Tính tổng combo amount
		orderCombos, err := p.orderComboRepo.GetByOrderID(ctx, order.ID)
		totalComboAmount := float64(0)
		if err == nil {
			for _, combo := range orderCombos {
				totalComboAmount += combo.TotalPrice
			}
		}

		paymentDetail := response.PaymentDetail{
			Date:             payment.PaymentTime,
			UserName:         user.Name,
			MovieName:        movie.Title,
			RoomName:         room.Name,
			Tickets:          seatNames,
			TotalComboAmount: totalComboAmount,
			TotalAmount:      payment.Amount,
			Status:           payment.Status,
		}

		paymentDetails = append(paymentDetails, paymentDetail)
	}

	return paymentDetails, nil
}

func (p *PaymentService) AcceptAll(ctx context.Context) error {
	orders, err := p.orderRepo.GetAllPendingOrders(ctx)
	if err != nil {
		return err
	}
	for _, order := range orders {
		_, err := p.HandleCallback(ctx, request.PaymentCallback{
			Payment: request.Payment{
				TransactionId:   "tp-bank" + strconv.Itoa(rand.Intn(100000000)),
				Content:         order.ID.String(),
				Amount:          order.TotalPrice,
				Date:            time.Now(),
				Gate:            "tp-bank",
				AccountReceiver: "10001259387",
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
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
