package billing

import (
	"avito-test/internal/models"
)

func (r *repository) SubmitDecision(bidId string, decision models.Decision) error {
	query :=
		`
		INSERT INTO bid_decisions(bid_id, value)
		VALUES($1,$2)
	`
	_, err := r.db.Exec(query, &bidId, decision.String())
	if err != nil {
		return err
	}
	return nil
}
