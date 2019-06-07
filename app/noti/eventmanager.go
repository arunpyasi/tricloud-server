package noti

type EventManager struct {
	cs            CredentialStore
	eventProvider Provider
	chanLogConsume chan []byte
	//consumeStat chan *StatType
}

func NewEventManager()*EventManager{
	return &EventManager{}
}

func (e *EventManager) ConsumeLog(data []byte) {

}

func (e *EventManager) SaveKey() {

}

func (e *EventManager) SaveToken() {

}




/*
fegnEF0AXtY:APA91bG4f6R6S0I1vtAkf7ngd0z6Vo3aaUiMnCMpy7pmgDZF0aplQ41tt4F4ww0FRhK1BEkZFnEk1nEa79D0hFeGk5ydYldwjSX67P17a71sbCT9iwiJ5JLmXizEOz9xVGzA9i8Ux3M9

func (e *EventManager)Run(){

}

// type StatType struct {
// 	username string
// 	agentname string
// 	data []byte
// }


*/
