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

func exportDataToCSV(db *sql.DB, filePath string) error {
	query := `
		SELECT c.full_name, c.email, p.amount, pm.method_name, p.status, p.payment_date
		FROM payments p
		JOIN customers c ON p.customer_id = c.id
		JOIN payment_methods pm ON p.payment_method_id = pm.id;
	`
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

	// Write header
	writer.Write([]string{"Customer Name", "Email", "Amount", "Payment Method", "Status", "Payment Date"})

	// Write rows
	for rows.Next() {
		var fullName, email, methodName, status string
		var amount float64
		var paymentDate time.Time

		if err := rows.Scan(&fullName, &email, &amount, &methodName, &status, &paymentDate); err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		record := []string{fullName, email, fmt.Sprintf("%.2f", amount), methodName, status, paymentDate.Format(time.RFC3339)}
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadFileToAzure(filePath string) error {
	// Get Azure Storage account details from environment variables
	accountName := os.Getenv("AZURE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_ACCOUNT_KEY")
	containerName := os.Getenv("AZURE_CONTAINER_NAME")
	blobName := os.Getenv("AZURE_BLOB_NAME")

	if accountName == "" || accountKey == "" || containerName == "" || blobName == "" {
		return fmt.Errorf("missing Azure environment variables")
	}

	// Service client using Azure Storage account details
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return fmt.Errorf("error creating Azure credentials: %v", err)
	}

	// Create the service client (for the storage account)
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	serviceClient, err := azblob.NewServiceClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		return fmt.Errorf("error creating Azure service client: %v", err)
	}

	// Get a reference to the container using the service client
	containerClient := serviceClient.NewContainerClient(containerName)

	// Open the file for upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file for upload: %v", err)
	}
	defer file.Close()

	// Upload the file
	blobClient := containerClient.NewBlockBlobClient(blobName)
	_, err = blobClient.UploadFile(context.Background(), file, nil)
	if err != nil {
		return fmt.Errorf("error uploading file to Azure Data Lake: %v", err)
	}

	return nil
}

func ExportAndUpload() error {
	// Get PostgreSQL connection string from environment variables
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		return fmt.Errorf("missing PostgreSQL DSN environment variable")
	}

	// Connect to PostgreSQL
	db, err := sql.Open("pgx", postgresDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	filePath := "data.csv"
	// Export data to CSV
	err = exportDataToCSV(db, filePath)
	if err != nil {
		return fmt.Errorf("error exporting data to CSV: %v", err)
	}

	// Upload the CSV to Azure Data Lake
	err = uploadFileToAzure(filePath)
	if err != nil {
		return fmt.Errorf("error uploading file to Azure Data Lake: %v", err)
	}

	log.Println("Data export and upload to Azure Data Lake completed successfully.")
	return nil
}
