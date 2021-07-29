package codeexchange

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"github.com/gtank/cryptopasta"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/suite"
)

// AnyTime is used to expect arbitrary time value
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type RepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository ExchangeRepository
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) SetupTest() {
	encryptionKey := "ASBzXLpOI0GOorN41dKF47gcFnaILuVh"
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	repo, _ := NewRepository(db, encryptionKey)
	s.repository = repo
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *RepositoryTestSuite) TestExchange() {
	s.Run("should insert access token if workspace does not exist", func() {
		var oldCryptopastaEncryptor = cryptopasta.Encrypt
		defer func() {
			cryptopastaEncryptor = oldCryptopastaEncryptor
		}()
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "access_tokens" WHERE workspace = 'test'`)
		insertQuery := regexp.QuoteMeta(`INSERT INTO "access_tokens" ("created_at","updated_at","access_token","workspace","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)

		inputToken := &AccessToken{
			ID:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AccessToken: "test-token",
			Workspace:   "test",
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertQuery).WithArgs(inputToken.CreatedAt,
			inputToken.UpdatedAt, base64.StdEncoding.EncodeToString([]byte(inputToken.AccessToken)),
			inputToken.Workspace, inputToken.ID).
			WillReturnRows(sqlmock.NewRows(nil))

		cryptopastaEncryptor = func(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
			return []byte("test-token"), nil
		}
		err := s.repository.Upsert(inputToken)
		s.Nil(err)
	})

	s.Run("should update exchange code if workspace exists", func() {
		var oldCryptopastaEncryptor = cryptopasta.Encrypt
		defer func() {
			cryptopastaEncryptor = oldCryptopastaEncryptor
		}()
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "access_tokens" WHERE workspace = 'test'`)
		updateQuery := regexp.QuoteMeta(`UPDATE "access_tokens" SET "created_at"=$1,"updated_at"=$2,"access_token"=$3,"workspace"=$4 WHERE id = $5 AND "id" = $6`)
		timeNow := time.Now()

		accessToken1 := &AccessToken{
			ID:          10,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
			AccessToken: "test-token",
			Workspace:   "test",
		}

		accessToken2 := &AccessToken{
			ID:          10,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
			AccessToken: "test-token-4",
			Workspace:   "test",
		}

		expectedRows1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "access_token", "workspace"}).
			AddRow(accessToken1.ID, accessToken1.CreatedAt, accessToken1.UpdatedAt,
				accessToken1.AccessToken, accessToken1.Workspace)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)

		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{}, AnyTime{}, base64.StdEncoding.EncodeToString([]byte(accessToken2.AccessToken)), accessToken2.Workspace,
			accessToken2.ID, accessToken2.ID).WillReturnResult(sqlmock.NewResult(10, 1))

		cryptopastaEncryptor = func(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
			return []byte("test-token-4"), nil
		}
		err := s.repository.Upsert(accessToken2)
		s.Nil(err)
		err = s.dbmock.ExpectationsWereMet()
		s.Nil(err)
	})

	s.Run("should return error in updating access token", func() {
		var oldCryptopastaEncryptor = cryptopasta.Encrypt
		defer func() {
			cryptopastaEncryptor = oldCryptopastaEncryptor
		}()
		expectedErrorMessage := "random error"
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "access_tokens" WHERE workspace = 'test'`)
		updateQuery := regexp.QuoteMeta(`UPDATE "access_tokens" SET "created_at"=$1,"updated_at"=$2,"access_token"=$3,"workspace"=$4 WHERE id = $5 AND "id" = $6`)
		timeNow := time.Now()

		accessToken1 := &AccessToken{
			ID:          10,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
			AccessToken: "test-token",
			Workspace:   "test",
		}

		accessToken2 := &AccessToken{
			ID:          10,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
			AccessToken: "test-token-4",
			Workspace:   "test",
		}

		expectedRows1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "access_token", "workspace"}).
			AddRow(accessToken1.ID, accessToken1.CreatedAt, accessToken1.UpdatedAt,
				accessToken1.AccessToken, accessToken1.Workspace)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)

		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{}, AnyTime{}, base64.StdEncoding.EncodeToString([]byte(accessToken2.AccessToken)), accessToken2.Workspace,
			accessToken2.ID, accessToken2.ID).WillReturnError(errors.New(expectedErrorMessage))

		cryptopastaEncryptor = func(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
			return []byte("test-token-4"), nil
		}
		err := s.repository.Upsert(accessToken2)
		s.EqualError(err, expectedErrorMessage)
		err = s.dbmock.ExpectationsWereMet()
		s.Nil(err)
	})

	s.Run("should return error in checking if workspace exists", func() {
		var oldCryptopastaEncryptor = cryptopasta.Encrypt
		defer func() {
			cryptopastaEncryptor = oldCryptopastaEncryptor
		}()
		expectedErrorMessage := "random error"
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "access_tokens" WHERE workspace = 'test'`)
		timeNow := time.Now()

		accessToken2 := &AccessToken{
			ID:          10,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
			AccessToken: "test-token-4",
			Workspace:   "test",
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnError(errors.New(expectedErrorMessage))

		cryptopastaEncryptor = func(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
			return []byte("test-token-4"), nil
		}
		err := s.repository.Upsert(accessToken2)
		s.EqualError(err, expectedErrorMessage)
		err = s.dbmock.ExpectationsWereMet()
		s.Nil(err)
	})

	s.Run("should return empty string if cryptopasta fails to encrypt", func() {
		var oldCryptopastaEncryptor = cryptopasta.Encrypt
		defer func() {
			cryptopastaEncryptor = oldCryptopastaEncryptor
		}()

		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "access_tokens" WHERE workspace = 'test'`)
		inputToken := &AccessToken{
			ID:          1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AccessToken: "test-token",
			Workspace:   "test",
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))
		cryptopastaEncryptor = func(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
			return []byte(""), errors.New("failed to encrypt")
		}
		err := s.repository.Upsert(inputToken)
		s.EqualError(err, "encryption failed: failed to encrypt")

	})
}

func (s *RepositoryTestSuite) TestNewRepository() {
	s.Run("should through error if encryption key is less then 32 char", func(){
		encryptionKey := "ASBzXLpOI0GOorN41dKF47gcFnaI"
		repo, err := NewRepository(nil, encryptionKey)
		s.Nil(repo)
		s.EqualError(err, "random hash should be 32 chars in length")
	})

}