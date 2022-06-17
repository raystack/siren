package receiver

//go:generate mockery --name=StrategyService -r --case underscore --with-expecter --structname StrategyService --filename receiver_strategy_service.go --output=./mocks
type StrategyService interface {
	Encrypt(r *Receiver) error
	Decrypt(r *Receiver) error
	PopulateReceiver(rcv *Receiver) (*Receiver, error)
	Notify(rcv *Receiver, payloadMessage NotificationMessage) error
}
