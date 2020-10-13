package wallet

import (
	"errors"
	"os"
	"log"
	"strconv"
	"strings"
	"io"
	"github.com/google/uuid"
	"github.com/s-zer0/wallet/pkg/types"
)

var ErrPhoneRegistred = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")
var ErrFileNotFound = errors.New("file not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	var targetpayment *types.Payment
		for _, payment := range s.payments{
			if payment.ID == paymentID {
				targetpayment=payment
			}
		}
		if targetpayment ==nil {
			return nil, ErrAccountNotFound
		}

		pay,err:=s.Pay(targetpayment.AccountID,targetpayment.Amount,targetpayment.Category)
		if err!=nil {
			return nil,err
		}
		
		return pay,nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	defer func() {
		if ferr := file.Close(); ferr != nil {
			log.Print(ferr)
		}
	}()

	str := ""

	for _, data := range s.accounts {
			str += strconv.Itoa(int(data.ID)) + ";"
			str += string(data.Phone) + ";"
			str += strconv.Itoa(int(data.Balance)) + "|"
	}
	
	_, err = file.Write([]byte(str))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if ferr := file.Close(); ferr != nil {
			return
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)

	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}

		content = append(content, buf[:read]...)
	}

	data := strings.Split(string(content), "|")
	data = data[:len(data)-1]
	for _, dt := range data {
		splits := strings.Split(dt, ";")

		id, err := strconv.Atoi(splits[0])
		if err != nil {
			return err
		}

		balance, err := strconv.Atoi(splits[2])
		if err != nil {
			return err
		}

		account := &types.Account{
			ID:      int64(id),
			Phone:   types.Phone(splits[1]),
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, account)
	}
	return nil
}