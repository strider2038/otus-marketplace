package messaging

import (
	"context"
	"encoding/json"

	"billing-service/internal/billing"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/strider2038/pkg/persistence"
)

type CreatePayment struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

type PaymentSucceeded struct {
	ID uuid.UUID `json:"id"`
}

func (p PaymentSucceeded) Name() string {
	return "Billing/PaymentSucceeded"
}

type PaymentDeclined struct {
	ID     uuid.UUID `json:"id"`
	Reason string    `json:"reason"`
}

func (p PaymentDeclined) Name() string {
	return "Billing/PaymentDeclined"
}

type CreatePaymentProcessor struct {
	accounts           billing.AccountRepository
	operations         billing.OperationRepository
	transactionManager persistence.TransactionManager
	dispatcher         Dispatcher
}

func NewCreatePaymentProcessor(
	accounts billing.AccountRepository,
	operations billing.OperationRepository,
	transactionManager persistence.TransactionManager,
	dispatcher Dispatcher,
) *CreatePaymentProcessor {
	return &CreatePaymentProcessor{
		accounts:           accounts,
		operations:         operations,
		transactionManager: transactionManager,
		dispatcher:         dispatcher,
	}
}

func (processor *CreatePaymentProcessor) Process(ctx context.Context, message []byte) error {
	var command CreatePayment
	err := json.Unmarshal(message, &command)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal CreatePayment command")
	}

	payment := billing.NewPayment(command.ID, command.UserID, command.Amount, command.Description)
	err = processor.transactionManager.DoTransactionally(ctx, processor.addPayment(payment))
	if err != nil {
		return err
	}

	return nil
}

func (processor *CreatePaymentProcessor) addPayment(payment *billing.Operation) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		account, err := processor.accounts.FindByIDForUpdate(ctx, payment.AccountID)
		if err != nil {
			return errors.WithMessagef(err, `failed to find user account by id "%s"`, payment.AccountID)
		}

		account.Amount -= payment.Amount
		if account.Amount < 0 {
			err = processor.dispatcher.Dispatch(ctx, PaymentDeclined{
				ID:     payment.ID,
				Reason: "not enough money",
			})
			if err != nil {
				return errors.WithMessage(err, "failed to dispatch PaymentDeclined event")
			}

			return nil
		}

		err = processor.operations.Add(ctx, payment)
		if err != nil {
			return errors.WithMessagef(err, `failed to add payment "%s" for user "%s"`, payment.ID, payment.AccountID)
		}

		err = processor.accounts.Save(ctx, account)
		if err != nil {
			return errors.WithMessagef(err, `failed to save user account "%s"`, account.ID)
		}

		err = processor.dispatcher.Dispatch(ctx, PaymentSucceeded{ID: payment.ID})
		if err != nil {
			return errors.WithMessage(err, "failed to dispatch PaymentSucceeded event")
		}

		return nil
	}
}
