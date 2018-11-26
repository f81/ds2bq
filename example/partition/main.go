package example

import (
	"net/http"
	"time"

	"github.com/f81/ds2bq"
)

const (
	apiDeleteBackup   = "/api/datastore-management/delete-old-backups"
	tqDeleteOldBackup = "/tq/datastore-management/delete-old-backups"
	tqDeleteBackup    = "/tq/datastore-management/delete-backup"
	apiReceiveOCN     = "/api/gcs/object-change-notification"
	tqImportBigQuery  = "/tq/gcs/object-to-bq"
)

func init() {
	// Datastore backup management
	queueName := "exec-rm-old-datastore-backups"
	expireAfter := 24 * time.Hour * 30
	http.HandleFunc(apiDeleteBackup, ds2bq.DeleteOldBackupAPIHandlerFunc(queueName, tqDeleteOldBackup))
	http.HandleFunc(tqDeleteOldBackup, ds2bq.DeleteOldBackupTaskHandlerFunc(queueName, tqDeleteBackup, expireAfter))
	http.HandleFunc(tqDeleteBackup, ds2bq.DeleteBackupTaskHandlerFunc(queueName))

	// import GCS to BigQuery
	queueName = "datastore-to-bq"
	bucketName := "ds2bqexample-nethttp"
	datasetID := "datastore_imports"
	targetKinds := []string{"Article", "User"}
	http.HandleFunc(apiReceiveOCN, ds2bq.ReceivePubSubNotificationHandlerFunc(bucketName, queueName, tqImportBigQuery, targetKinds)) // from GCS, via pubsub

	bqPartitionConf := &ds2bq.BQConfigMap{
		"user": {Kind: "user", TimePartitioningField: "time", ClusteringFields: []string{"department", "name"}},
	}
	http.HandleFunc(tqImportBigQuery, ds2bq.ImportBigQueryWithConfHandlerFunc(datasetID, bqPartitionConf))
}
