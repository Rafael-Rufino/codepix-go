package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending string = "pending"
	TransactionCompleted string = "completed"
	TransactionError string = "error"
	TransactionConfirmed string = "confirmed"
)

type TransactionsRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []*Transaction
}

type Transaction struct {
	Base 											`valid:"required"`
	AccountFrom *Account 			`valid:"-"`
	Amount float64 						`json:"amount" valid:"notnull"`
	PixKeyTo *Pixkey 					`valid:"-"`
	Status string 						`json:"status" valid:"notnull"`
	Description string 				`json:"description" valid:"notnull"`
	CancelDescription string 	`json:"cancel_description" valid:"notnull"`
}

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)

}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if t.Amount <= 0 {
		return errors.New("The amount must be greater than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("Invalid status for the transaction")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("The source and destination account cannot be the same")
	}

	if err != nil {
		return err
	}
	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *Pixkey, description string, id string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount: amount,
		PixKeyTo: pixKeyTo,
		Status: TransactionPending,
		Description: description,
	}


	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}


func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
	t.UpdatedAt = time.Now()
	t.CancelDescription = description
	err := t.isValid()
	return err
}