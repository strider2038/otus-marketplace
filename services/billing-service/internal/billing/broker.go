package billing

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type BrokerAccount struct {
	brokerID   uuid.UUID
	accounts   AccountRepository
	operations OperationRepository
}

func NewBrokerAccount(brokerID uuid.UUID, accounts AccountRepository, operations OperationRepository) *BrokerAccount {
	return &BrokerAccount{brokerID: brokerID, accounts: accounts, operations: operations}
}

func (broker *BrokerAccount) ID() uuid.UUID {
	return broker.brokerID
}

func (broker *BrokerAccount) ChargeCommission(ctx context.Context, commission float64, comment string) error {
	account, err := broker.accounts.FindByIDForUpdate(ctx, broker.brokerID)
	if err != nil {
		return errors.WithMessagef(err, `failed to find broker account by id "%s"`, broker.brokerID)
	}

	account.Amount += commission

	err = broker.operations.Add(ctx, NewCommission(broker.brokerID, commission, comment))
	if err != nil {
		return errors.WithMessagef(err, `failed to charge commission for broker "%s"`, broker.brokerID)
	}

	err = broker.accounts.Save(ctx, account)
	if err != nil {
		return errors.WithMessagef(err, `failed to save broker account "%s"`, account.ID)
	}

	return nil
}
