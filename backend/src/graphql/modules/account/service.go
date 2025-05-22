package account

import (
	"database/sql"
	"errors"

	"backend/src/.gen/cdr-intelligence/public/model"
	"backend/src/.gen/cdr-intelligence/public/table"
	error_utils "backend/src/error"
	"backend/src/utils"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type AccountService struct {
	Db    *sql.DB
	Redis *redis.Client
}

func NewAccountService(params NewAccountParams) AccountService {

	return AccountService{
		Db:    params.Db,
		Redis: params.Redis,
	}
}

func (this AccountService) GetById(id uuid.UUID) (Account, error) {
	stmt := table.Accounts.
		SELECT(table.Accounts.AllColumns).
		WHERE(table.Accounts.ID.EQ(pg.UUID(id))).
		LIMIT(1)

	account := model.Accounts{}
	err := stmt.Query(this.Db, &account)

	if err != nil && error_utils.HasNoRow(err) {
		return Account{}, errors.New("account not found")
	}

	if err != nil {
		log.Printf("get-account-profile-by-id-error: %v", err.Error())
		return Account{}, err
	}

	return transformToGraphql(account), nil
}

func (this AccountService) Create(input CreateAccountData) (Account, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Account{}, err
	}

	hashPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return Account{}, err
	}

	createAccountStmt := table.Accounts.
		INSERT(
			table.Accounts.ID,
			table.Accounts.Username,
			table.Accounts.Password,
			table.Accounts.CreatedBy,
		).
		MODEL(model.Accounts{
			ID:        id,
			Username:  input.Username,
			Password:  hashPassword,
			CreatedBy: "System",
		}).
		RETURNING(table.Accounts.AllColumns)

	account := model.Accounts{}
	err = createAccountStmt.Query(this.Db, &account)
	if err != nil && error_utils.IsDuplicate(err) {
		return Account{}, errors.New("username is already exist")
	}

	if err != nil {
		return Account{}, err
	}

	return transformToGraphql(account), nil
}

func transformToGraphql(data model.Accounts) Account {
	return Account{
		Id:        graphql.ID(data.ID.String()),
		Username:  data.Username,
		CreatedAt: graphql.Time{Time: data.CreatedAt},
	}
}
