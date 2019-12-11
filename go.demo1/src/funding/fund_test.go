package funding

import (
        "sync"
        "testing"
)

const WORKERS=10

func BenchmarkWithdrawals(b *testing.B) {
        // Skip N -1
        if b.N < WORKERS {
                return
        }

        server := NewFundServer(b.N)

        // Add as many dollars as we have iterations this run
        // fund := NewFund(b.N)

        // Casually assume b.N divides cleanly
        dollarsPerFounder := b.N / WORKERS

        // WaitGroup struct don't need to be initialized
        // (their "zero value" is ready to use).
        // So, we just declare one and then use it.
        var wg sync.WaitGroup

        for i := 0; i < WORKERS; i++ {
                wg.Add(1)

                go func() {
                        // Mark this worker done when the function finishes
                        defer wg.Done()

                        pizzaTime := false
                        for i := 0; i < dollarsPerFounder; i++ {

                                server.Transact( func(fund *Fund) {

                                        // Stop when we're down to pizza money
                                        if fund.Balance() <= 10 {
                                                pizzaTime = true
                                                return
                                        }
                                        fund.Withdraw(1)
                                        // server.Commands <- WithdrawCommand{ Amount: 1}
                                })
                                if pizzaTime {
                                        break
                                }
                        }
                }() // Remember to call the closure.
        }
        wg.Wait()


        balance := server.Balance()
        // balanceResponseChan := make(chan int)
        // server.Commands <- BalanceCommand { Response: balanceResponseChan }
        // balance := <- balanceResponseChan

        if balance != 10 {
                b.Error("Balance wasn't 10 dollars: ", balance)
        }

}

