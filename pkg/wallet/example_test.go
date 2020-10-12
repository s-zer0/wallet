package wallet

import (
	"reflect"
	"testing"
	"github.com/s-zer0/wallet/pkg/types"
)

func TestService_RegisterAccount_Success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_FindAccoundByID_NotFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(2)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}

}

func TestService_Reject_Success(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+9920000001")
	account, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Reject(): can't register account, error = %v", err)
		return
	}

	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("Reject(): can't deposit account, error = %v", err)
		return
	}

	payment, err := s.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can't create payment, error = %v", err)
		return
	}

	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
}

func TestService_FindPaymentByID_Success(t *testing.T) {
	//создаём сервис
	s := &Service{}
	
	// регистрируем там пользователя
	phone := types.Phone("+9920000001")
	account, err :=s.RegisterAccount(phone)
	if err !=nil {
		t.Errorf("FindPaymentByID(): can't register account, error = %v", err)
		return
	}

	//пополняем его счёт
	err = s.Deposit(account.ID, 10_000_00)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't deposit account, error = %v", err)
		return
	}

	//осуществляем платёж на его счёт
	payment, err := s.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
		return
	}

	//пробуем найти платёж
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_Repeat_Success(t *testing.T) {
	s := &Service{}
	s.RegisterAccount("+9920000001")

	account, err := s.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = s.Deposit(account.ID, 1000_00)
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