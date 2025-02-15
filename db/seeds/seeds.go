package seeds

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Seed struct {
	db *sqlx.DB
}

// NewSeed initializes a new Seed instance with a database connection.
func newSeed(db *sqlx.DB) Seed {
	return Seed{db: db}
}

// Execute runs the seeder for the specified table with the given number of entries.
func Execute(db *sqlx.DB, table string) {
	seed := newSeed(db)
	seed.run(table)
}

// NotificationTemplatesSeed seeds the `notification_templates` table.
func (s *Seed) notificationTemplatesSeed() {
	templates := []map[string]any{
		{
			"template_name": "register",
			"subject":       "Welcome to Dzikra App",
			"body":          "Hi welcome to our service! We're glad to have you.",
		},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}
	defer rollbackOrCommit(tx, &err)

	// Check if `notification_templates` table already contains data
	var count int
	err = tx.Get(&count, `SELECT COUNT(*) FROM notification_templates`)
	if err != nil {
		log.Error().Err(err).Msg("Error checking notification_templates table")
		return
	}
	if count > 0 {
		log.Info().Msg("Notification templates already seeded")
		return
	}

	// Insert data into `notification_templates`
	insertTemplateQuery := `
		INSERT INTO notification_templates 
		(template_name, subject, body) 
		VALUES (:template_name, :subject, :body)`
	for _, template := range templates {
		template["id"] = uuid.New().String()
		_, err = tx.NamedExec(insertTemplateQuery, template)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting notification template")
			return
		}
	}

	log.Info().Msg("notification_templates table seeded successfully")
}

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

// rollbackOrCommit handles transaction rollback or commit based on error state.
func rollbackOrCommit(tx *sqlx.Tx, err *error) {
	if *err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Error().Err(rbErr).Msg("Error rolling back transaction")
		}
	} else {
		cmErr := tx.Commit()
		if cmErr != nil {
			log.Error().Err(cmErr).Msg("Error committing transaction")
		}
	}
}

// Update the `run` function to include `notification_templatesSeed`.
func (s *Seed) run(table string) {
	switch table {
	case "notification_templates":
		s.notificationTemplatesSeed()
	case "all":
		s.notificationTemplatesSeed()
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}
}
