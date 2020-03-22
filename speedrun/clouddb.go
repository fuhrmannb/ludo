package speedrun

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const GameCollectionID = "games"
const CategoryCollectionID = "categories"
const RecordCollectionID = "records"

type Record struct {
	Player   string        `firestore:"player"`
	Duration time.Duration `firestore:"duration"`
	Date     time.Time     `firestore:"date,serverTimestamp"`
	Event    string        `firestore:"event,omitempty"`
}

// FirestoreCredentialsFile is the name of the Firebase credential file to provide to have
// a Cloud DB access
const FirestoreCredentialsFile = "firebase_credentials.json"

// CloudDB is an access to remote database to store and show speedrun records
type CloudDB struct {
	db *firestore.Client
}

// NewCloudDB crates a CloudDB instance based on speedrun directory configuration file
// A file with needed credentials must be specified (see `FirestoreCredentialsFile`)
func NewCloudDB(speedrunDir string) (*CloudDB, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(filepath.Join(speedrunDir, FirestoreCredentialsFile))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, fmt.Errorf("cannot create CloudDB firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create CloudDB firestore client: %w", err)
	}

	return &CloudDB{
		db: client,
	}, nil
}

func (c *CloudDB) getRecordCollection(gameID string, categoryID string) *firestore.CollectionRef {
	return c.db.Collection(GameCollectionID).Doc(gameID).Collection(CategoryCollectionID).Doc(categoryID).Collection(RecordCollectionID)
}

// ListRecords list all speedrun records for a specific game/category
func (c *CloudDB) ListRecords(gameID string, categoryID string) ([]*Record, error) {
	var result []*Record

	players := make(map[string]bool)
	recordCol := c.getRecordCollection(gameID, categoryID).OrderBy("duration", firestore.Asc)
	iter := recordCol.Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get record document from CloudDB: %w", err)
		}
		var r Record
		err = doc.DataTo(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal record document from CloudDB: %w", err)
		}

		if !players[r.Player] {
			result = append(result, &r)
			players[r.Player] = true
		}

	}
	return result, nil
}

// AddRecord create a new speedrun record for a game in a specific category
func (c *CloudDB) AddRecord(gameID string, categoryID string, record Record) error {
	_, _, err := c.getRecordCollection(gameID, categoryID).Add(context.Background(), record)
	if err != nil {
		return fmt.Errorf("failed to add record into CloudDB: %w", err)
	}
	return nil
}

// Close terminate the Cloud database connection
func (c *CloudDB) Close() {
	c.db.Close()
}
