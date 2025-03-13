package repositories

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gistsapp/api/types"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Abstraction for database related operations
type Database interface {
	//Bootstrap the database by running the migrations for example
	Bootstrap() error
	CreateUser(user *types.User) (*types.User, error)
	GetUserByID(id string) (*types.User, error)
	DeleteUser(id string) error
	UpdateUser(user *types.User) (*types.User, error)
	GetUserThroughFederatedIdentity(federated_id string) (*types.User, error)
	CreateFederatedIdentity(federated_identity *types.FederatedIdentity) (*types.FederatedIdentity, error)
	GetFederatedIdentityByID(id string) (*types.FederatedIdentity, error)
	DeleteFederatedIdentity(id string) error
	CreateOpaqueToken(opaque_token *types.OpaqueToken) (*types.OpaqueToken, error)
	GetOpaqueTokenByID(id string) (*types.OpaqueToken, error)
	GetOpaqueTokenByUserEmail(email string) (*types.OpaqueToken, error)
	GetOpaqueTokenByToken(token string) (*types.OpaqueToken, error)
	DeleteOpaqueToken(id string) error
	CreateVerificationToken(verification_token *types.VerificationToken) (*types.VerificationToken, error)
	GetVerificationTokenByEmail(email string) (*types.VerificationToken, error)
	DeleteVerificationToken(email string, value string) error
}

type PgDatabase struct {
	db *sqlx.DB
	username string
	password string
	host string
	port int
	dbname string
}

func NewPgDatabase(username string, password string, host string, port int, dbname string) (*PgDatabase, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, username, dbname, password))
	if err != nil {
		return nil, err
	}
	database := PgDatabase{
		db: db,
		username: username,
		password: password,
		host: host,
		port: port,
		dbname: dbname,
	}

	return &database, nil
}

func (db *PgDatabase) Bootstrap() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	migrationsPath := filepath.Join(filepath.Dir(ex), "migrations")

	fmt.Println(migrationsPath)
	fmt.Println(db)

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.username, db.password, db.host, db.port, db.dbname))

	if err != nil {
		return err
	}

	return m.Up()
}

func (db *PgDatabase) CreateUser(user *types.User) (*types.User, error) {
	var created_user types.User
	err := db.db.Get(&created_user, "INSERT INTO user_entity (username, email, picture) VALUES ($1, $2, $3) RETURNING *", user.Username, user.Email, user.Picture)
	if err != nil {
		return nil, err
	}
	return &created_user, nil
}

func (db *PgDatabase) GetUserByID(id string) (*types.User, error) {
	var user types.User
	err := db.db.Get(&user, "SELECT * FROM user_entity WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *PgDatabase) DeleteUser(id string) error {
	_, err := db.db.Exec("DELETE FROM user_entity WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *PgDatabase) UpdateUser(user *types.User) (*types.User, error) {
	var updated_user types.User
	err := db.db.Get(&updated_user, "UPDATE user_entity SET username = $1, email = $2, picture = $3 WHERE id = $4", user.Username, user.Email, user.Picture, user.ID)
	if err != nil {
		return nil, err
	}
	return &updated_user, nil
}

func (db *PgDatabase) CreateFederatedIdentity(federated_identity *types.FederatedIdentity) (*types.FederatedIdentity, error) {
	var created_federated_identity types.FederatedIdentity // TODO: change to federated_identity_id
	err := db.db.Get(&created_federated_identity, "INSERT INTO federated_identity_entity (id, user_id, provider, data) VALUES ($1, $2, $3, $4) RETURNING *", federated_identity.ID, federated_identity.UserID, federated_identity.Provider, federated_identity.Data)
	return &created_federated_identity, err
}

func (db *PgDatabase) GetFederatedIdentityByID(id string) (*types.FederatedIdentity, error) {
	var federated_identity types.FederatedIdentity
	err := db.db.Get(&federated_identity, "SELECT * FROM federated_identity_entity WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &federated_identity, nil
}

func (db *PgDatabase) DeleteFederatedIdentity(id string) error {
	_, err := db.db.Exec("DELETE FROM federated_identity_entity WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *PgDatabase) CreateOpaqueToken(opaque_token *types.OpaqueToken) (*types.OpaqueToken, error) {
	var created_opaque_token types.OpaqueToken
	err := db.db.Get(&created_opaque_token, "INSERT INTO opaque_token_entity (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4) RETURNING *", opaque_token.ID, opaque_token.UserID, opaque_token.Token, opaque_token.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &created_opaque_token, nil
}

func (db *PgDatabase) GetOpaqueTokenByID(id string) (*types.OpaqueToken, error) {
	var opaque_token types.OpaqueToken
	err := db.db.Get(&opaque_token, "SELECT * FROM opaque_token_entity WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &opaque_token, nil
}

func (db *PgDatabase) DeleteOpaqueToken(id string) error {
	_, err := db.db.Exec("DELETE FROM opaque_token_entity WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *PgDatabase) GetOpaqueTokenByUserEmail(email string) (*types.OpaqueToken, error) {
	var opaque_token types.OpaqueToken
	err := db.db.Get(&opaque_token, "SELECT * FROM opaque_token_entity WHERE user_id = (SELECT id FROM user_entity WHERE email = $1)", email)
	if err != nil {
		return nil, err
	}
	return &opaque_token, nil
}

func (db *PgDatabase) GetOpaqueTokenByToken(token string) (*types.OpaqueToken, error) {
    var opaqueToken types.OpaqueToken
    err := db.db.Get(&opaqueToken, "SELECT * FROM opaque_token_entity WHERE token = $1", token)
    if err != nil {
        return nil, err
    }
    return &opaqueToken, nil
}

func (db *PgDatabase) GetUserThroughFederatedIdentity(federated_id string) (*types.User, error) {
	var user types.User
	err := db.db.Get(&user, "SELECT * FROM user_entity WHERE user_id = (SELECT user_id FROM federated_identity WHERE federated_identity_id = $1)", federated_id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *PgDatabase) CreateVerificationToken(verification_token *types.VerificationToken) (*types.VerificationToken, error) {
	var created_verification_token types.VerificationToken
	err := db.db.Get(&created_verification_token, "INSERT INTO verification_token (token, email) VALUES ($1, $2) RETURNING *", verification_token.Token, verification_token.Email)
	if err != nil {
		return nil, err
	}
	return &created_verification_token, nil
}

func (db *PgDatabase) GetVerificationTokenByEmail(email string) (*types.VerificationToken, error) {
	var verification_token types.VerificationToken
	err := db.db.Get(&verification_token, "SELECT * FROM verification_token WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &verification_token, nil
}


func (db *PgDatabase) DeleteVerificationToken(email string, value string) error {
	_, err := db.db.Exec("DELETE FROM verification_token WHERE email = $1 AND token = $2", email, value)
	if err != nil {
		return err
	}
	return nil
}

