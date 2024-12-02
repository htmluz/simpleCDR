package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func AutoClean(db *sql.DB, interval time.Duration) {
	go func() {
		for {
			var days int
			q := `SELECT cleanup_days FROM cleanup_config LIMIT 1;`
			if err := db.QueryRow(q).Scan(&days); err != nil {
				log.Printf("Erro ao obter o cleanup_days %v\n", err)
				time.Sleep(interval)
				continue
			}
			var count int
			countQ := fmt.Sprintf(`SELECT COUNT(*) FROM call_records WHERE created_at < now() - INTERVAL '%d day'`, days)
			if err := db.QueryRow(countQ).Scan(&count); err != nil {
				log.Printf("Erro ao contar a quantia de tickets a serem apagadas %v\n", err)
				time.Sleep(interval)
				continue
			}

			delQ := fmt.Sprintf(`DELETE FROM call_records WHERE created_at < NOW() - INTERVAL '%d day'`, days)
			if _, err := db.Exec(delQ); err != nil {
				log.Printf("Erro ao executar a limpeza %v\n", err)
			} else {
				log.Printf("Limpeza concluída, %d registros com mais de %d dias apagados", count, days)
			}
			time.Sleep(interval)
			// TODO gravar historico de limpezas
		}
	}()
}
