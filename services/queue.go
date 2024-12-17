package services

import (
	"database/sql"
	"fmt"
	"radiusgo/models"
	"sync"
	"time"
)

type CallQueue struct {
	mu    sync.RWMutex
	bilhs map[string]*models.BilheteFull
	db    *sql.DB
}

func generateUniqueKey(bilh *models.Bilhete) string {
	layout := "15:04:05.000 -0700 Mon Jan 2 2006"
	setupTime, err := time.Parse(layout, bilh.H323SetupTime)
	if err != nil {
		fmt.Printf("Erro parseando time %v\n", err)
		return ""
	}
	roundedTime := setupTime.Truncate(10 * time.Second)
	return fmt.Sprintf("%s|%s|%s|%s",
		bilh.CallingStationID,
		bilh.UserName,
		roundedTime.Format("2006-01-02T15:04:05"),
		"2sOffset")
}

func NewCallQueue(db *sql.DB) *CallQueue {
	return &CallQueue{
		bilhs: make(map[string]*models.BilheteFull),
		db:    db,
	}
}

func (q *CallQueue) Add(bilh *models.Bilhete) {
	q.mu.Lock()
	defer q.mu.Unlock()
	bilhKey := generateUniqueKey(bilh)
	fmt.Print(bilhKey, "\n")
	if _, ok := q.bilhs[bilhKey]; !ok {
		q.bilhs[bilhKey] = &models.BilheteFull{
			Bid:  bilhKey,
			LegA: &models.Bilhete{},
			LegB: &models.Bilhete{},
		}
	}

	if bilh.AcctStatusType == "Start" {
		if bilh.H323CallOrigin == "answer" {
			q.bilhs[bilhKey].LegA = bilh
		} else if bilh.H323CallOrigin == "originate" {
			q.bilhs[bilhKey].LegB = bilh
		}
	} else if bilh.AcctStatusType == "Stop" {
		if bilh.H323CallOrigin == "answer" {
			q.bilhs[bilhKey].LegA = bilh
			InsertBilhete(q.db, bilh)
			InsertBid(q.db, *q.bilhs[bilhKey])
		} else {
			q.bilhs[bilhKey].LegB = bilh
			InsertBilhete(q.db, bilh)
			InsertBid(q.db, *q.bilhs[bilhKey])
		}
		if q.bilhs[bilhKey].LegA.AcctStatusType == "Stop" && q.bilhs[bilhKey].LegB.AcctStatusType == "Stop" {
			InsertBid(q.db, *q.bilhs[bilhKey])
			delete(q.bilhs, bilhKey)
		}
	}
}

// funcao de rotina pra limpar a queue
// trato os start e stops desordenados por aqui
func (q *CallQueue) QueueCleanup(interval time.Duration) {
	go func() {
		for {
			q.mu.Lock()
			copyMap := make(map[string]*models.BilheteFull, len(q.bilhs))
			for k, v := range q.bilhs {
				copyMap[k] = v
			}
			q.mu.Unlock()

			for k := range copyMap {
				if BidExists(q.db, copyMap[k].Bid) {
					// valida se o id do bilhete ja ta no banco, ocorre quando os stop vem antes do start
					q.mu.Lock()
					delete(q.bilhs, copyMap[k].Bid)
					q.mu.Unlock()
				}
			}
			time.Sleep(interval)
		}
	}()
}

func (q *CallQueue) GetQueueSize() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.bilhs)
}

func (q *CallQueue) GetAllCalls() []models.BilheteFull {
	q.mu.RLock()
	defer q.mu.RUnlock()
	var calls []models.BilheteFull
	for _, bilhete := range q.bilhs {
		calls = append(calls, *bilhete)
	}
	return calls
}
