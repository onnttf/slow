package dal

import (
	"context"
	"slow/dal/model"

	"gorm.io/gorm"
)

// Entity interface that any model must implement to be used with the DAO.
type Entity interface {
	model.Petrol
}

// DAO interface defines methods for data access operations including create, update, querying single entity, querying a list of entities, and counting.
type DAO[T Entity] interface {
	Create(ctx context.Context, newValue T) (*T, error)
	Update(ctx context.Context, newValue T, funcs ...func(*gorm.DB) *gorm.DB) error
	QueryOne(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) (*T, error)
	QueryList(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) ([]T, error)
	Count(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) (int64, error)
}

// Dao struct holds a reference to the gorm.DB, providing base for performing database operations.
type Dao[T Entity] struct {
	database *gorm.DB
}

// NewDao function initializes a Dao instance.
func NewDao[T Entity](database *gorm.DB) *Dao[T] {
	return &Dao[T]{database: database}
}

// Create method inserts a new record into the database.
func (dao *Dao[T]) Create(ctx context.Context, newValue T) (*T, error) {
	if err := dao.database.WithContext(ctx).Create(&newValue).Error; err != nil {
		return nil, err
	}
	return &newValue, nil
}

// Update method updates an existing record in the database.
func (dao *Dao[T]) Update(ctx context.Context, newValue T, funcs ...func(*gorm.DB) *gorm.DB) error {
	return dao.database.WithContext(ctx).Model(&newValue).Scopes(funcs...).Updates(&newValue).Error
}

// QueryOne method retrieves a single record from the database.
func (dao *Dao[T]) QueryOne(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) (*T, error) {
	var record T
	result := dao.database.WithContext(ctx).Scopes(funcs...).Limit(1).Find(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	return &record, nil
}

// QueryList method retrieves a list of records from the database.
func (dao *Dao[T]) QueryList(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	var recordList []T
	err := dao.database.WithContext(ctx).Scopes(funcs...).Find(&recordList).Error
	return recordList, err
}

// Count method counts the records in the database that meet the criteria.
func (dao *Dao[T]) Count(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	result := dao.database.WithContext(ctx).Model(new(T)).Scopes(funcs...).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// Paginate function provides implementation for pagination.
func Paginate(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo <= 0 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
