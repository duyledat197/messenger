package postgres

import (
	"context"
	"fmt"
	"strings"

	"openmyth/messgener/internal/user/entity"
	"openmyth/messgener/internal/user/repository"
	"openmyth/messgener/util/database"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

// Create inserts a new user into the database.
//
// ctx: the context.Context object for the function.
// db: the database.Executor object for executing the SQL statement.
// data: the entity.User object containing the user data to be inserted.
// Returns an error if there was a problem executing the SQL statement.
func (u *userRepository) Create(ctx context.Context, db database.Executor, data *entity.User) error {
	fields, values := database.FieldMap(data)
	placeHolders := database.GetPlaceholders(len(fields))
	stmt := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s)`,
		data.TableName(),
		strings.Join(fields, ","),
		placeHolders,
	)

	if _, err := db.ExecContext(ctx, stmt, values...); err != nil {
		return err
	}

	return nil
}

// UpdateInfoByID updates user information by ID.
//
// ctx is the context.Context for the operation, db is the database.Executor for database operations,
// id is the identifier of the user, and data is the entity.User containing updated information.
// Returns an error if the update operation fails.
func (u *userRepository) UpdateInfoByID(ctx context.Context, db database.Executor, id string, data *entity.User) error {
	_, values := database.FieldMap(data)
	stmt := fmt.Sprintf(`
	UPDATE FROM %s 
	SET
		hashed_password = COALESCE($4, hashed_password),
		enable_2fa = COALESCE($5, enable_2fa),
		email = COALESCE($6, email),
		google = COALESCE($7, google),
		facebook = COALESCE($8, facebook),
		discord = COALESCE($9, discord),
		github = COALESCE($10, github)
	WHERE
		id = $1
	`,
		data.TableName(),
	)

	if _, err := db.ExecContext(ctx, stmt, values...); err != nil {
		return err
	}

	return nil
}

// RetrieveByUserName retrieves a user from the userRepository based on the provided username.
func (u *userRepository) RetrieveByUserName(ctx context.Context, db database.Executor, userName string) (*entity.User, error) {
	e := &entity.User{}
	fields, values := database.FieldMap(e)
	stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE username = $1`,
		strings.Join(fields, ","),
		e.TableName(),
	)
	if err := db.QueryRowContext(ctx, stmt, &userName).Scan(values...); err != nil {
		return nil, err
	}

	return e, nil
}
