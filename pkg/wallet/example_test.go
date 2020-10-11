package wallet

import (
	"testing"
	"github.com/s-zer0/wallet/pkg/types"
)

func TestService_RegisterAccount_success(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	s.RegisterAccount(phone)

	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_FindAccoundByIdmethod_notFound(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	s.RegisterAccount(phone)

	account, err := s.FindAccountByID(2)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestDeposit(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	s.RegisterAccount(phone)

	err := s.Deposit(1, 100_00)
	if err != nil {
		t.Error("Не удалось пополнить счёт")
	}

	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_Reject_success(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	s.RegisterAccount(phone)

	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := s.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = s.Reject(pay.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}

func TestService_Reject_fail(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	s.RegisterAccount(phone)

	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := s.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payID := pay.ID + " "
	err = s.Reject(payID)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	s.RegisterAccount(phone)

	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := s.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err = s.Repeat(pay.ID)
	if err != nil {
		t.Errorf("Repeat(): Error(): can't pay for an account(%v): %v", pay.ID, err)
	}
}

func TestService_Favorite_Success(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	
	account, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("RegisterAccount returned not nil error, account => %v", account)
	}

	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("Deposit returned not nil error, error => %v", err)
	}

	payment, err := s.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Pay() Error() can't pay for an account(%v): %v", account, err)
	}

	favorite, err := s.FavoritePayment(payment.ID, "megafon")
	if err != nil {
		t.Errorf("Error() can't for an favorite(%v): %v", favorite, err)
	}

	paymentFavorite, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("Error() can't for an favorite(%v): %v", paymentFavorite, err)
	}
}