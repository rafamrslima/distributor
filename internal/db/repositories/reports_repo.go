package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/rafamrslima/distributor/internal/db"
	"github.com/rafamrslima/distributor/internal/domain"
)

func GetReportInfo(ctx context.Context, clientEmail string, reportName string) ([]domain.Report, error) {
	pool, err := db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database pool: %w", err)
	}

	// Create context with timeout for the database operation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := pool.Query(ctx,
		`SELECT client_email, report_name, gains, losses, info_date FROM reports WHERE client_email = $1 AND report_name = $2`,
		clientEmail, reportName)

	if err != nil {
		return nil, fmt.Errorf("failed to query reports: %w", err)
	}
	defer rows.Close()

	var results []domain.Report

	for rows.Next() {
		var res domain.Report
		if err := rows.Scan(&res.ClientEmail, &res.ReportName, &res.Gains, &res.Losses, &res.InfoDate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, res)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}
