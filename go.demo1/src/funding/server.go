package funding

type Transactor func(fund *Fund)

type TransactionCommand struct {
        Transactor Transactor
        Done chan bool
}


type FundServer struct {
        // commands chan interface{}
        commands chan TransactionCommand
        fund *Fund
}

type WithdrawCommand struct {
        Amount int
}

type BalanceCommand struct {
        Response chan int
}

func (s *FundServer) Transact(transactor Transactor) {
        command := TransactionCommand {
                Transactor: transactor,
                Done: make(chan bool),
        }
        s.commands <- command
        <- command.Done
}

func NewFundServer(initialBalance int) *FundServer {
        server := &FundServer {
                // make() creates builtins like channels, maps, and slices
                //commands: make(chan interface{}),
                commands: make(chan TransactionCommand),
                fund: NewFund(initialBalance),
        }

        // Spawn off the server's main loop imediatley
        go server.loop()
        return server
}

func (s *FundServer) Balance() int {
        var balance int
        s.Transact(func(f *Fund) {
                balance = f.Balance()
        })
        return balance

        //responseChan := make(chan int)
        //s.commands <- BalanceCommand{ Response: responseChan }
        //return <- responseChan
}

func (s *FundServer) Withdraw(amount int) {
        s.Transact(func (f *Fund) {
                f.Withdraw(amount)
        })
        //s.commands <-WithdrawCommand { Amount: amount }
}

func (s *FundServer) loop() {
        for transaction := range s.commands {
                transaction.Transactor(s.fund)
                transaction.Done <- true
        }
}

