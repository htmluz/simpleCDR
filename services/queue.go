package services

import (
	"database/sql"
	"fmt"
	"radiusgo/models"
	"sync"
)

var (
	callQueue = make(map[string]*models.Bilhete)
	queueLock = sync.RWMutex{}
)

func AddCall(call *models.Bilhete) {
	queueLock.Lock()
	defer queueLock.Unlock()
	callQueue[call.CallID] = call
}

func UpdateCall(updatedCall *models.Bilhete) {
	queueLock.Lock()
	defer queueLock.Unlock()

	if existingCall, exists := callQueue[updatedCall.CallID]; exists {
		*existingCall = *updatedCall
	} else {
		callQueue[updatedCall.CallID] = updatedCall
	}
}

// tem que mudar o tratamento pra identificar as duas pernas e gerar uma chamada relacionada
// leg a e leg b

func RemoveCall(db *sql.DB, call *models.Bilhete) {
	queueLock.Lock()
	defer queueLock.Unlock()
	if existingCall, exists := callQueue[call.CallID]; exists {
		*existingCall = *call
	} else {
		fmt.Print("tentou deletar chamada q n tinha")
	}
	InsertBilhete(db, call)
	delete(callQueue, call.CallID)
}

func GetActiveCalls() ([]*models.Bilhete, int) {
	queueLock.RLock()
	defer queueLock.RUnlock()
	activeCalls := []*models.Bilhete{}
	for _, call := range callQueue {
		activeCalls = append(activeCalls, call)
	}
	return activeCalls, len(activeCalls)
}

func QueueHandleBilhete(db *sql.DB, call *models.Bilhete) {
	_, exists := callQueue[call.CallID]
	fmt.Println(call.AcctStatusType)
	switch call.AcctStatusType {
	case "Start":
		if !exists {
			AddCall(call)
		} else {
			UpdateCall(call)
		}
	case "Stop":
		RemoveCall(db, call)
	}
}
