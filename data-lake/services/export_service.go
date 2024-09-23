package services

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
)

const (
	accountName   = "your-account-name"
	accountKey    = "your-account-key"
	containerName = "your-container-name"  // Equivalent to filesystem in Data Lake
	blobName      = "exports/data.csv"     // Path in Azure Data Lake where the file will be uploaded
)

const (
	postgresDSN = "postgres://user:password@localhost:5432/payments_db"
)


func exportDataCSV(db *sql.DB, filePath string) error {
	query := '
		SELECT c.full_name, c.email, p.amount, pm.method_name, p.status, p.payment_date
		FROM payments p
		JOIN customers c ON p.customer_id = c.id
		JOIN payment_methods pm ON p.payment_method_id = pm.id;
	'
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying the database: %v", err)
	}
	defer rows.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	
	writer.Write([]string{"Customer Name", "Email", "Amount", "Payment Method", "Status", "Payment Date"})

	
	for rows.Next() {
		var fullName, email, methodName, status string
		var amount float64
		var paymentDate time.Time

		if err := rows.Scan(&fullName, &email, &amount, &methodName, &status, &paymentDate); err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		record := []string{fullName, email, fmt.Sprintf("%.2f", amount), methodName, status, paymentDate.Format(time.RFC3339)}
		writer.Write(record)
	}
	return nil
}

func ExportDataAndUploadToAzure() error {
	db, err := sql.Open("pgx", postgresDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	filePath := "data.csv"
	err = exportDataToCSV(db, filePath)
	if err != nil {
		return fmt.Errorf("error exporting data to CSV: %v", err)
	}

	err = uploadFileToAzure(filePath)
	if err != nil {
		return fmt.Errorf("error uploading file to azure data lake: %v", err)
	}

	fmt.Println("Data export and upload to Azure Data Lake completed successfully.")
	return nil
}