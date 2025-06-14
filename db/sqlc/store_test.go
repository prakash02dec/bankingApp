package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)
	fromAccount1 := createRandomAccount(t)
	toAccount2 := createRandomAccount(t)

	// Run n concurrent transfers
	amount := int64(10)
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i:= 0 ; i < n ; i++ {
		go func() {
			// Perform the transfer
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount1.ID,
				ToAccountID:   toAccount2.ID,
				Amount:        amount,
			})
			// If there is an error, log it
			if err != nil {
				t.Errorf("Transfer failed: %v", err)
			}
			// Check the result
			errs <- err
			results <- result
		}()
	}
	
	for i:= 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		

		// Check the transfer
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, fromAccount1.ID, result.Transfer.FromAccountID)
		require.Equal(t, toAccount2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)
		// ensure the transfer is created in the database
		transfer, err := testQueries.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)


		// Check entries
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, fromAccount1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)
		fromEntry, err := testQueries.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, toAccount2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)
		toEntry, err := testQueries.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)

		// Check accounts
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, fromAccount1.ID, result.FromAccount.ID)
		require.Equal(t, fromAccount1.Owner, result.FromAccount.Owner)
		require.Equal(t, fromAccount1.Currency, result.FromAccount.Currency)
		require.NotZero(t, result.FromAccount.CreatedAt)
		fromAccount, err := testQueries.GetAccount(context.Background(), result.FromAccount.ID)
		require.NoError(t, err)
		require.NotEmpty(t, fromAccount)


		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, toAccount2.ID, result.ToAccount.ID)
		require.Equal(t, toAccount2.Owner, result.ToAccount.Owner)
		require.Equal(t, toAccount2.Currency, result.ToAccount.Currency)
		require.NotZero(t, result.ToAccount.CreatedAt)
		toAccount, err := testQueries.GetAccount(context.Background(), result.ToAccount.ID)
		require.NoError(t, err)
		require.NotEmpty(t, toAccount)
	}



}
