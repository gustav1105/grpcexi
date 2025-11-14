package repo

import (
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	pb "grpcexi/protos"
	"log"
)

type ContactRepo interface {
	Create(contact *pb.Contact) error
	GetAll() ([]*pb.Contact, error)
}

type SQLLContactRepo struct {
	db *sql.DB
}

func NewSQLiteContactRepository(dbPath string) *SQLLContactRepo {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open sqlite DB: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS contacts (
		id TEXT PRIMARY KEY,
		name TEXT
		);

		CREATE TABLE IF NOT EXISTS phones (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		contact_id TEXT,
		digits TEXT,
		type INTEGER,
		FOREIGN KEY(contact_id) REFERENCES contacts(id)
		);
		`)

	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	return &SQLLContactRepo{db: db}
}

func (r *SQLLContactRepo) Create(contact *pb.Contact) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	contactId := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO contacts (id, name) VALUES (?, ?)
		`, contactId, contact.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, p := range contact.Phones {
		_, err = tx.Exec(`
			INSERT INTO phones (contact_id, digits, type)
			VALUES (?, ?, ?)
			`, contactId, p.Digits, int(p.Type))

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLLContactRepo) GetAll() ([]*pb.Contact, error) {
	rows, err := r.db.Query(`
		SELECT 
		c.id,
		c.name,
		p.digits,
		p.type
		FROM contacts c
		LEFT JOIN phones p ON p.contact_id = c.id
		ORDER BY c.id
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contactMap := make(map[string]*pb.Contact)

	for rows.Next() {
		var (
			id        string
			name      string
			digits    sql.NullString
			phoneType sql.NullInt64
		)

		if err := rows.Scan(&id, &name, &digits, &phoneType); err != nil {
			return nil, err
		}

		if _, exists := contactMap[id]; !exists {
			contactMap[id] = &pb.Contact{
				Name:   name,
				Phones: []*pb.PhoneNumber{},
			}
		}

		if digits.Valid {
			contactMap[id].Phones = append(contactMap[id].Phones, &pb.PhoneNumber{
				Digits: digits.String,
				Type:   pb.PhoneType(phoneType.Int64),
			})
		}
	}

	contacts := make([]*pb.Contact, 0, len(contactMap))
	for _, c := range contactMap {
		contacts = append(contacts, c)
	}

	return contacts, nil
}
