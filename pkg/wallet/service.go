package wallet

import (
	"errors"

	"github.com/google/uuid"
	"github.com/s-zer0/wallet/pkg/types"
)

var ErrPhoneRegistred = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistred
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платёж.
	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (s *Service) Reject(paymentID string) error {
	var targetPaymet *types.Payment
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			targetPaymet = payment
			break
		}
	}
	if targetPaymet == nil {
		return ErrPaymentNotFound
	}

	var targetAccount *types.Account
	for _, account := range s.accounts {
		if account.ID == targetPaymet.AccountID {
			targetAccount = account
			break
		}
	}
	if targetAccount == nil {
		return ErrAccountNotFound
	}

	targetPaymet.Status = types.PaymentStatusFail
	targetAccount.Balance += targetPaymet.Amount
	return nil
}