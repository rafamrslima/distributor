package db

import (
	"context"
	"log"

	"github.com/rafamrslima/distributor/internal/domain"
)

func GetReportInfo(clientEmail string, reportName string) ([]domain.Report, error) {
	pool, err := connect()
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	ctx := context.Background()

	rows, err := pool.Query(ctx,
		`SELECT client_email, report_name, gains, losses, info_date 
		FROM investment_results WHERE info_date::date = CURRENT_DATE AND client_email = $1 AND report_name = $2`,
		clientEmail, reportName)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Report

	for rows.Next() {
		var res domain.Report
		if err := rows.Scan(&res.ClientEmail, &res.ReportName, &res.Gains, &res.Losses, &res.InfoDate); err != nil {
			log.Fatal(err)
		}
		results = append(results, res)
	}

	return results, nil
}
