package service

import (
	"database/sql"
	"fmt"
	"merch_shop/internal/entity"
	"merch_shop/internal/model"
	"merch_shop/internal/repository"
)

type TransactionService struct {
	uow repository.UnitOfWork
}

func NewTransactionService(uow repository.UnitOfWork) *TransactionService {
	return &TransactionService{uow: uow}
}

func (t TransactionService) GetInfo(userId uint) (model.InfoResponse, error) {
	tx, err := t.uow.BeginTransaction(&sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return model.InfoResponse{}, fmt.Errorf("failed to begin transaction")
	}
	defer tx.Commit()

	userRepository := tx.UserRepository()
	transactionRepository := tx.TransactionRepository()
	user, err := userRepository.FindUserById(userId)
	if err != nil {
		return model.InfoResponse{}, fmt.Errorf("error getting user")
	}
	if user == nil {
		return model.InfoResponse{}, fmt.Errorf("user not found")
	}
	outcome, err := transactionRepository.GetOutcomeTransactions(userId)
	if err != nil {
		return model.InfoResponse{}, fmt.Errorf("error getting outcome transactions")
	}
	income, err := transactionRepository.GetIncomeTransactions(userId)
	if err != nil {
		return model.InfoResponse{}, fmt.Errorf("error getting income transactions")
	}
	inventory, err := transactionRepository.GetUserInventory(userId)
	if err != nil {
		return model.InfoResponse{}, fmt.Errorf("error getting inventory")
	}
	inventoryModel := make([]model.Inventory, 0, len(inventory))
	for _, v := range inventory {
		inventoryModel = append(inventoryModel, model.Inventory{
			Name:     v.Item.Name,
			Quantity: v.Quantity,
		})
	}
	outcomeModel := make([]model.CoinHistorySent, 0, len(outcome))
	for _, v := range outcome {
		outcomeModel = append(outcomeModel, model.CoinHistorySent{
			ToUser: v.ToUser.Name,
			Amount: v.Amount,
		})
	}
	incomeModel := make([]model.CoinHistoryReceived, 0, len(income))
	for _, v := range income {
		incomeModel = append(incomeModel, model.CoinHistoryReceived{
			FromUser: v.FromUser.Name,
			Amount:   v.Amount,
		})
	}
	coinHistoryModel := model.CoinHistory{
		Received: incomeModel,
		Sent:     outcomeModel,
	}
	infoResponse := model.InfoResponse{
		Coins:       user.Balance,
		Inventory:   inventoryModel,
		CoinHistory: coinHistoryModel,
	}
	return infoResponse, nil
}

func (t TransactionService) SendCoin(userId uint, toUserName string, amount uint) error {
	tx, err := t.uow.BeginTransaction(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}

	userRepository := tx.UserRepository()
	transactionRepository := tx.TransactionRepository()

	fromUser, err := userRepository.FindUserById(userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find user")
	}
	if fromUser == nil {
		tx.Rollback()
		return fmt.Errorf("user not found")
	}
	toUser, err := userRepository.FindUserByName(toUserName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find user")
	}
	if toUser == nil {
		tx.Rollback()
		return fmt.Errorf("user not found")
	}

	if toUser.ID == fromUser.ID {
		tx.Rollback()
		return fmt.Errorf("cannot send coin to yourself")
	}
	if fromUser.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient balance")
	}
	fromUser.Balance -= amount
	toUser.Balance += amount
	err = userRepository.UpdateUser(fromUser)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user")
	}
	err = userRepository.UpdateUser(toUser)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user")
	}
	transaction := entity.Transaction{
		FromId: fromUser.ID,
		ToId:   toUser.ID,
		Amount: amount,
	}
	err = transactionRepository.CreateTransaction(&transaction)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create transaction")
	}
	tx.Commit()
	return nil

}

func (t TransactionService) BuyItem(userId uint, name string) error {
	tx, err := t.uow.BeginTransaction(&sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	userRepository := tx.UserRepository()
	transactionRepository := tx.TransactionRepository()
	user, err := userRepository.FindUserById(userId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find user")
	}
	if user == nil {
		tx.Rollback()
		return fmt.Errorf("user not found")
	}
	item, err := transactionRepository.GetItemByName(name)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to find item")
	}
	if item == nil {
		tx.Rollback()
		return fmt.Errorf("item not found")
	}
	if user.Balance < item.Price {
		tx.Rollback()
		return fmt.Errorf("insufficient balance")
	}
	user.Balance -= item.Price
	err = userRepository.UpdateUser(user)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user")
	}
	err = transactionRepository.AddItem(user.ID, item.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add item to inventory: %v", err)
	}

	tx.Commit()
	return nil
}
