package database

import (
	"context"
	"testing"

	"main/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestGorm(t *testing.T) {
	suite.Run(t, new(GormSuite))
}

type GormSuite struct {
	suite.Suite
	ctx context.Context
}

func (su *GormSuite) SetupSuite() {
	su.ctx = context.Background()
}

func (su *GormSuite) mockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func (su *GormSuite) TestTxSuccess() {
	db, mock, err := su.mockDB()
	su.Require().NoError(err)
	su.Require().NotNil(db)
	su.Require().NotNil(mock)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `models` (`value`,`id`) VALUES (?,?)").
		WithArgs("Hello", 1).
		WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectCommit()

	repo := Repository{db: db}
	su.Require().NoError(
		repo.Tx(su.ctx, func(txCtx context.Context) error {
			return repo.getDB(txCtx).Debug().Model(&entity.Model{}).Create(entity.Model{ID: 1, Value: "Hello"}).Error
		}),
	)

	su.NoError(
		mock.ExpectationsWereMet(),
	)

	mock.ExpectQuery("SELECT * FROM `models` LIMIT ?").
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "value"}).
				AddRow(1, "Hello"),
		)

	result := entity.Model{}
	dbResult := repo.getDB(su.ctx).Debug().Model(&entity.Model{}).Take(&result)
	su.NoError(dbResult.Error)

	su.NoError(
		mock.ExpectationsWereMet(),
	)
}

func (su *GormSuite) TestTxFailed() {
	db, mock, err := su.mockDB()
	su.Require().NoError(err)
	su.Require().NotNil(db)
	su.Require().NotNil(mock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `models` (`value`,`id`) VALUES (?,?)").
		WithArgs("Hello", 1).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectExec("INSERT INTO `models` (`value`,`id`) VALUES (?,?)").
		WithArgs("Hello", 1).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	repo := Repository{db: db}
	su.Require().Error(
		repo.Tx(su.ctx, func(txCtx context.Context) error {
			if err := repo.getDB(txCtx).Debug().Model(&entity.Model{}).Create(entity.Model{ID: 1, Value: "Hello"}).Error; err != nil {
				return err
			}

			return repo.getDB(txCtx).Debug().Model(&entity.Model{}).Create(entity.Model{ID: 1, Value: "Hello"}).Error
		}),
	)

	su.NoError(
		mock.ExpectationsWereMet(),
	)

	mock.ExpectQuery("SELECT * FROM `models` LIMIT ?").
		WithArgs(1).
		WillReturnError(gorm.ErrRecordNotFound)

	result := entity.Model{}
	su.Error(
		repo.getDB(su.ctx).Debug().Model(&entity.Model{}).Take(&result).Error,
	)

	su.NoError(
		mock.ExpectationsWereMet(),
	)
}
