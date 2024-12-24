package services

import (
	"database/sql"
	"fmt"
	"radiusgo/models"
	"radiusgo/utils"
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
	last := ""
	if len(bilh.CalledStationID) < 4 {
		last = "6669"
	} else {
		last = bilh.CalledStationID[len(bilh.CalledStationID)-4:]
	}
	roundedTime := setupTime.Truncate(3 * time.Second)
	return fmt.Sprintf("%s|%s|%s|%s|%s",
		bilh.CallingStationID,
		last,
		bilh.UserName,
		roundedTime.Format("2006-01-02T15:04:05"),
		"3sOffset")
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
	if _, ok := q.bilhs[bilhKey]; !ok {
		q.bilhs[bilhKey] = &models.BilheteFull{
			Bid:  bilhKey,
			LegA: nil,
			LegB: nil,
		}
	}
	q.bilhs[bilhKey].Lock()
	defer q.bilhs[bilhKey].Unlock()

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
			InsertBid(q.db, q.bilhs[bilhKey])
		} else {
			q.bilhs[bilhKey].LegB = bilh
			InsertBilhete(q.db, bilh)
			InsertBid(q.db, q.bilhs[bilhKey])
		}
		if q.bilhs[bilhKey].LegA != nil && q.bilhs[bilhKey].LegB != nil &&
			q.bilhs[bilhKey].LegA.AcctStatusType == "Stop" && q.bilhs[bilhKey].LegB.AcctStatusType == "Stop" {
			InsertBid(q.db, q.bilhs[bilhKey])
			delete(q.bilhs, bilhKey)
		}
	}
}

func (q *CallQueue) WriteAndRemove(bilhKey string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.bilhs[bilhKey].LegA != nil {
		InsertBilhete(q.db, q.bilhs[bilhKey].LegA)
	}
	if q.bilhs[bilhKey].LegB != nil {
		InsertBilhete(q.db, q.bilhs[bilhKey].LegB)
	}
	InsertBid(q.db, q.bilhs[bilhKey])
	delete(q.bilhs, bilhKey)
}

// funcao de rotina pra limpar a queue
// trato os start e stops desordenados por aqui
func (q *CallQueue) QueueCleanup(interval time.Duration) {
	go func() {
		for {
			// cria uma copia do map pra nao dar lock enquanto faz a limpeza, vai ser sempre um pouco atrasado
			q.mu.Lock()
			copyMap := make(map[string]*models.BilheteFull, len(q.bilhs))
			for k, v := range q.bilhs {
				copyMap[k] = v
			}
			q.mu.Unlock()

			now := time.Now()
			for _, bilh := range copyMap {
				bilh.Lock()
				if bilh.LegA == nil && bilh.LegB != nil || bilh.LegA != nil && bilh.LegB == nil {
					// valida se so tem uma perna, se tiver e for mais velha que 5 minutos apaga
					var setupTime string
					if bilh.LegA != nil {
						setupTime = bilh.LegA.H323SetupTime
					}
					if bilh.LegB != nil {
						setupTime = bilh.LegB.H323SetupTime
					}
					timeBilh, err := utils.ConvertToTimestamp(setupTime)
					if err != nil {
						fmt.Println("Erro convertendo para tempo: ", err)
						break
					}
					parsedTime, e := time.Parse("2006-01-02 15:04:05.000-07:00", timeBilh)
					if e != nil {
						fmt.Println("Erro parseando tempo: ", err)
						break
					}
					if now.Sub(parsedTime) > 5*time.Minute {
						q.WriteAndRemove(bilh.Bid)
					}
				} else if bilh.LegA != nil && bilh.LegB != nil {
					parsedTime, err := time.Parse("15:04:05.000 -0700 Mon Jan 02 2006", bilh.LegA.H323SetupTime)
					if err != nil {
						fmt.Println("Erro parseando tempo: ", err)
						break
					}
					if bilh.LegA.AcctStatusType == "Start" && bilh.LegB.AcctStatusType == "Stop" ||
						bilh.LegA.AcctStatusType == "Stop" && bilh.LegB.AcctStatusType == "Start" {
						if now.Sub(parsedTime) > 5*time.Minute {
							q.WriteAndRemove(bilh.Bid)
						}
					} else {
						if now.Sub(parsedTime) > 4*time.Hour {
							q.WriteAndRemove(bilh.Bid)
						}
					}
				} else {
					exists := false
					if bilh.LegA != nil {
						exists = CallIDExists(q.db, bilh.LegA.CallID)
					} else if bilh.LegB != nil {
						exists = CallIDExists(q.db, bilh.LegA.CallID)
					}
					if exists {
						q.WriteAndRemove(bilh.Bid)
					}
				}
				bilh.Unlock()
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

func (q *CallQueue) GetAllCalls() []*models.BilheteFull {
	q.mu.RLock()
	defer q.mu.RUnlock()
	var calls []*models.BilheteFull
	for _, bilhete := range q.bilhs {
		calls = append(calls, bilhete)
	}
	return calls
}
