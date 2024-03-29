package employees

import (
	"database/sql"
	models "github.com/GuiTadeu/mercado-fresh-panic/cmd/server/database"
)

type EmployeeRepository interface {
	GetAll() ([]models.Employee, error)
	Create(cardNumberId string, firstName string, lastName string, wareHouseId uint64) (models.Employee, error)
	Get(id uint64) (models.Employee, error)
	Update(updatedEmployee models.Employee) (models.Employee, error)
	Delete(id uint64) error
	ExistsEmployeeCardNumberId(cardNumberId string) (bool, error)
	CountInboundOrdersByEmployeeId(id uint64) (models.ReportInboundOrders, error)
	CountInboundOrders() ([]models.ReportInboundOrders, error)
	ExistsEmployee(uint64) (bool)
}

type employeeRepository struct {
	db *sql.DB
}

func NewRepository(employees *sql.DB) EmployeeRepository {
	return &employeeRepository{
		db: employees,
	}
}

func (r *employeeRepository) Create(
	cardNumberId string, firstName string, lastName string, wareHouseId uint64,
) (models.Employee, error) {

	stmt, err := r.db.Prepare(`
		INSERT INTO employees(
			id_card_number, 
			first_name, 
			last_name, 
			warehouse_id
		) VALUES(?, ?, ?, ?)
	`)

	if err != nil {
		return models.Employee{}, err
	}

	defer stmt.Close()
	var result sql.Result
	result, err = stmt.Exec(
		cardNumberId,
		firstName,
		lastName,
		wareHouseId,
	)
	if err != nil {
		return models.Employee{}, err
	}

	insertId, _ := result.LastInsertId()
	employees := models.Employee{
		Id:           uint64(insertId),
		CardNumberId: cardNumberId,
		FirstName:    firstName,
		LastName:     lastName,
		WarehouseId:  wareHouseId,
	}

	return employees, nil

}
func (r *employeeRepository) Get(id uint64) (models.Employee, error) {

	var myEmployee models.Employee
	stmt := r.db.QueryRow(`
		SELECT 
		    id, 
		    id_card_number, 
		    first_name, 
		    last_name, 
		    warehouse_id 
		FROM employees 
		WHERE id = ? 
		`, id)

	err := stmt.Scan(
		&myEmployee.Id,
		&myEmployee.CardNumberId,
		&myEmployee.FirstName,
		&myEmployee.LastName,
		&myEmployee.WarehouseId,
	)
	if err != nil {
		return models.Employee{}, err
	}

	return myEmployee, nil
}

func (r *employeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee

	stmt, err := r.db.Query(`
		SELECT 
		    id, 
		    id_card_number, 
		    first_name, 
		    last_name, 
		    warehouse_id 
		FROM employees
		`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for stmt.Next() {
		var oneEmployee models.Employee
		err := stmt.Scan(
			&oneEmployee.Id,
			&oneEmployee.CardNumberId,
			&oneEmployee.FirstName,
			&oneEmployee.LastName,
			&oneEmployee.WarehouseId,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, oneEmployee)
	}

	return employees, nil
}

func (r *employeeRepository) Update(updatedEmployee models.Employee) (models.Employee, error) {

	stmt, err := r.db.Prepare(`
		UPDATE employees 
		SET id_card_number = ?, 
		    first_name = ?, 
		    last_name = ?, 
		    warehouse_id = ? 
		WHERE id = ?
		`)
	if err != nil {
		return models.Employee{}, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		updatedEmployee.CardNumberId,
		updatedEmployee.FirstName,
		updatedEmployee.LastName,
		updatedEmployee.WarehouseId,
		updatedEmployee.Id,
	)
	if err != nil {
		return models.Employee{}, err
	}

	return updatedEmployee, nil
}

func (r *employeeRepository) Delete(id uint64) error {

	stmt, err := r.db.Prepare(`
		DELETE FROM employees 
		       WHERE id = ?
		       `)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (r *employeeRepository) ExistsEmployeeCardNumberId(cardNumberId string) (bool, error) {
	var employee models.Employee
	rows, err := r.db.Query(`
		SELECT id, 
		       id_card_number, 
		       first_name, 
		       last_name, 
		       warehouse_id 
		FROM employees 
		WHERE id_card_number = ?
		`, cardNumberId)

	if err != nil {
		return false, err
	}

	for rows.Next() {

		// Fields must be in the same order as in the models
		err := rows.Scan(
			&employee.Id,
			&employee.CardNumberId,
			&employee.FirstName,
			&employee.LastName,
			&employee.WarehouseId,
		)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (r *employeeRepository) CountInboundOrdersByEmployeeId(id uint64) (models.ReportInboundOrders, error) {
	var report models.ReportInboundOrders

	stmt := r.db.QueryRow(`
	SELECT 	employees.id, 
			id_card_number, 
			first_name, 
			last_name, 
			employees.warehouse_id, 
	COUNT(inbound_orders.id) AS inbound_orders_count
	FROM employees
	LEFT JOIN inbound_orders
	ON employees.id = inbound_orders.employee_id
	WHERE employees.id = ?
	GROUP BY employees.id;
		`, id)


	err := stmt.Scan(
		&report.Id,
		&report.CardNumberId,
		&report.FirstName,
		&report.LastName,
		&report.WarehouseId,
		&report.InboundOrdersCount,
	)
	if err != nil {
		return models.ReportInboundOrders{}, err
	}

	return report, nil
}

func (r *employeeRepository) CountInboundOrders() ([]models.ReportInboundOrders, error) {
	var reports []models.ReportInboundOrders

	stmt, err := r.db.Query(`
	SELECT 	employees.id, 
			id_card_number, 
			first_name, 
			last_name, 
			employees.warehouse_id, 
	COUNT(inbound_orders.id) AS inbound_orders_count
	FROM employees
	LEFT JOIN inbound_orders
	ON employees.id = inbound_orders.employee_id
	GROUP BY employees.id;
		`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for stmt.Next() {
		var oneReport models.ReportInboundOrders
		err := stmt.Scan(
			&oneReport.Id,
			&oneReport.CardNumberId,
			&oneReport.FirstName,
			&oneReport.LastName,
			&oneReport.WarehouseId,
			&oneReport.InboundOrdersCount,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, oneReport)
	}

	return reports, nil
}

func (r *employeeRepository) ExistsEmployee(id uint64) (bool) {

	stmt := r.db.QueryRow(`SELECT EXISTS(SELECT ID from employees WHERE id = ?)`, id)
	var result uint64
	stmt.Scan(
		&result,
	)

	if result == 0 {
        return false
    }
    return true
}

