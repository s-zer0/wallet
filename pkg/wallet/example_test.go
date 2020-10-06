package wallet

import (
	"reflect"
	"fmt"
	"testing"
	"github.com/s-zer0/wallet/pkg/types"
)

func TestService_RegisterAccount_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_FindAccoundByID_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(2)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}

}

func TestService_Reject_success(t *testing.T) {
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

func TestService_FindPaymentByID_success(t *testing.T) {
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
// func TestService_FindPaymentByID_Success(t *testing.T) {
// 	svc := &Service{}
// 	svc.RegisterAccount("+9920000001")

// 	payment, err := svc.Pay(2, types.Money(110), "auto")

// 	if err != nil {
// 		switch err {
// 			case ErrAmountMustBePositive:
// 				fmt.Println("Сумма должна быть положительной")
// 			case ErrAccountNotFound:
// 				fmt.Println("Аккаунт пользователя не найден")
// 			case ErrNotEnoughBalance:
// 				fmt.Println("Недостаточно средств на балансе")
// 			}
// 			return 
// 	}

// 	result, err := svc.FindPaymentByID(payment.ID)
// 	if err != nil{
// 		fmt.Println("Платёж не найден")
// 	}

// 	if !reflect.DeepEqual(payment, result) {
// 		t.Errorf("invalid result, expected: %v, actual: %v", payment, result)
// 	}
// }

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, error=%v", err)
	}

	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	return account, nil
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