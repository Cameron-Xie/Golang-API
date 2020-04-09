package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Storage struct {
	db   *gorm.DB
	conn *DBConn
}

type DBConn struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
	RootCert string
	MaxConn  int
	ConnLife time.Duration
}

func New(conn *DBConn) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) Open() (*Storage, error) {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			// NOSONAR
			"host=%v port=%v dbname=%v user=%v password=%s sslrootcert=%s sslmode=%s",
			s.conn.Host,
			s.conn.Port,
			s.conn.Database,
			s.conn.Username,
			s.conn.Password,
			s.conn.RootCert,
			s.conn.SSLMode,
		),
	)
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(s.conn.ConnLife)
	db.DB().SetMaxIdleConns(s.conn.MaxConn)
	db.DB().SetMaxOpenConns(s.conn.MaxConn)

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
