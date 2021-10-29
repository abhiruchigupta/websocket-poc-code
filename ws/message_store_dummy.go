package ws



import (
"time"
)

type DummyMessageStore struct {
	db string
}


func (p *DummyMessageStore) Initialize() error {
 	return nil
}

func (p *DummyMessageStore) StoreMessage(userID string, message string, receivedAt time.Time) (int64, error) {
	return  int64(0), nil
}
func (p *DummyMessageStore) RetrieveUnsentMessages(userID string) []DbMessage {
	return make([]DbMessage, 0)
}

func (p *DummyMessageStore) ConfirmSentMessage(id int64) error {
	return nil
}

func (p *DummyMessageStore) Close() {

}

