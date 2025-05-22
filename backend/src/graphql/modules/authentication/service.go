package authentication

import (
	"github.com/samber/lo"

	"database/sql"
	"errors"

	"backend/src/.gen/cdr-intelligence/public/model"
	"backend/src/.gen/cdr-intelligence/public/table"
	error_utils "backend/src/error"
	"backend/src/jwt"
	"backend/src/utils"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
	"github.com/thanhpk/randstr"
)

type AuthenticationService struct {
	Db    *sql.DB
	redis *redis.Client
}

func NewAuthenticationService(params NewAuthenticationResolverParams) AuthenticationService {
	return AuthenticationService{
		Db:    params.Db,
		redis: params.Redis,
	}
}

func (service AuthenticationService) Login(input LoginData) (Authentication, error) {
	getAccountStmt := table.Accounts.
		SELECT(table.Accounts.AllColumns).
		WHERE(
			table.Accounts.Username.EQ(pg.String(input.Username)),
		).
		LIMIT(1)

	account := model.Accounts{}
	err := getAccountStmt.Query(service.Db, &account)

	if err != nil && error_utils.HasNoRow(err) {
		return Authentication{}, errors.New("username or password is incorrect")
	}

	if err != nil {
		return Authentication{}, error_utils.InternalServerError
	}

	match, err := utils.ComparePasswordAndHash(input.Password, account.Password)
	if err != nil {
		return Authentication{}, errors.New("username or password is incorrect")
	}

	if !match {
		return Authentication{}, errors.New("username or password is incorrect")
	}

	token, err := jwt.SignToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
		Nounce:    randstr.Hex(16),
	})

	if err != nil {
		return Authentication{}, error_utils.SignTokenFailed
	}

	refreshToken, err := jwt.SignRefreshToken(jwt.SignedTokenParams{
		AccountId: account.ID.String(),
		Nounce:    randstr.Hex(16),
	})

	if err != nil {
		return Authentication{}, error_utils.SignTokenFailed
	}

	insertSessionTokenStmt := table.SessionToken.INSERT(
		table.SessionToken.Token,
		table.SessionToken.AccountId,
		table.SessionToken.Revoke,
	).MODEL(model.SessionToken{
		Token:     token,
		AccountId: account.ID,
		Revoke:    false,
	})

	_, err = insertSessionTokenStmt.Exec(service.Db)
	if err != nil {
		return Authentication{}, error_utils.InternalServerError
	}

	projects := []model.Projects{}
	projectsStmt := table.ProjectAccounts.
		SELECT(table.ProjectAccounts.ProjectId).
		FROM(table.ProjectAccounts).
		WHERE(table.ProjectAccounts.AccountId.EQ(pg.UUID(account.ID)))
	err = projectsStmt.Query(service.Db, &projects)

	if err != nil {
		return Authentication{}, error_utils.InternalServerError
	}

	projectIds := lo.Map(projects, func(item model.Projects, index int) graphql.ID {
		return graphql.ID(item.ID.String())
	})

	return Authentication{
		Token:        token,
		RefreshToken: refreshToken,
		AccountId:    graphql.ID(account.ID.String()),
		ProjectIds:   &projectIds,
	}, nil
}

func (service AuthenticationService) Logout(token string) (Logout, error) {
	if token == "" {
		return Logout{
			Success: false,
			Message: "Invalid token",
		}, nil
	}

	updateSessionTokenStmt := table.SessionToken.
		UPDATE(table.SessionToken.Revoke).
		MODEL(model.SessionToken{
			Revoke: true,
		}).
		WHERE(table.SessionToken.Token.EQ(pg.String(token)))

	result, err := updateSessionTokenStmt.Exec(service.Db)

	if err != nil {
		return Logout{}, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return Logout{
			Success: false,
			Message: "Invalid token",
		}, nil
	}

	return Logout{
		Success: true,
		Message: "Logout success",
	}, nil
}
