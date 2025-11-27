package repository

import (
	"database/sql"
	"time"

	"employee-app/internal/models"
)

type EmployeeRepository interface {
	Create(name, password, designation string) (models.Employee, error)
	GetByCredentials(name, password string) (models.Employee, error)
	GetByID(id int64) (models.Employee, error)
}

type PostgresEmployeeRepository struct {
	DB *sql.DB
}

func NewPostgresEmployeeRepository(db *sql.DB) *PostgresEmployeeRepository {
	return &PostgresEmployeeRepository{DB: db}
}

func (r *PostgresEmployeeRepository) Create(name, password, designation string) (models.Employee, error) {
	query := `
		INSERT INTO employees (name, password, designation)
		VALUES ($1, $2, $3)
		RETURNING id, name, designation, created_at;
	`

	var emp models.Employee
	var createdAt time.Time

	err := r.DB.QueryRow(query, name, password, designation).
		Scan(&emp.ID, &emp.Name, &emp.Designation, &createdAt)
	if err != nil {
		return models.Employee{}, err
	}
	emp.CreatedAt = createdAt
	return emp, nil
}

func (r *PostgresEmployeeRepository) GetByCredentials(name, password string) (models.Employee, error) {
	query := `
		SELECT id, name, password, designation, created_at
		FROM employees
		WHERE name=$1 AND password=$2
		LIMIT 1;
	`

	var emp models.Employee
	var createdAt time.Time

	err := r.DB.QueryRow(query, name, password).
		Scan(&emp.ID, &emp.Name, &emp.Password, &emp.Designation, &createdAt)
	if err != nil {
		return models.Employee{}, err
	}
	emp.CreatedAt = createdAt
	return emp, nil
}

func (r *PostgresEmployeeRepository) GetByID(id int64) (models.Employee, error) {
	query := `
		SELECT id, name, designation, created_at
		FROM employees
		WHERE id=$1;
	`

	var emp models.Employee
	var createdAt time.Time

	err := r.DB.QueryRow(query, id).
		Scan(&emp.ID, &emp.Name, &emp.Designation, &createdAt)
	if err != nil {
		return models.Employee{}, err
	}
	emp.CreatedAt = createdAt
	return emp, nil
}
