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

func TestService_FindPaymentByID_Success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+9920000001")

	payment, err := svc.Pay(2, types.Money(110), "auto")

	if err != nil {
		switch err {
			case ErrAmountMustBePositive:
				fmt.Println("Сумма должна быть положительной")
			case ErrAccountNotFound:
				fmt.Println("Аккаунт пользователя не найден")
			case ErrNotEnoughBalance:
				fmt.Println("Недостаточно средств на балансе")
			}
			return 
	}

	result, err := svc.FindPaymentByID(payment.ID)
	if err != nil{
		fmt.Println("Платёж не найден")
	}

	if !reflect.DeepEqual(payment, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", payment, result)
	}
}