package seeds

import (
	"context"

	"github.com/rs/zerolog/log"
)

// deleteAll modified to include `notification_templates`.
func (s *Seed) deleteAll() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Add `notification_templates` to the list of tables to clear
	tables := []string{"notification_templates"}
	for _, table := range tables {
		query := `DELETE FROM ` + table
		_, err = tx.Exec(query)
		if err != nil {
			log.Error().Err(err).Msgf("Error deleting table %s", table)
			return
		}
	}

	log.Info().Msg("All tables deleted successfully")
}
