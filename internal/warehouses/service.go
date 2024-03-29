package warehouses

import (
	"errors"
	"fmt"

	"github.com/GuiTadeu/mercado-fresh-panic/cmd/server/database"
	"github.com/imdario/mergo"
)

var (
	ExistsWarehouseCodeError = errors.New("warehouses code already exists")
	WarehouseNotFoundError   = errors.New("warehouses not found")
)

type WarehouseService interface {
	GetAll() ([]database.Warehouse, error)
	Create(Code string, address string, telephone string, minimunCapacity uint32, minimunTemperature float32, localityId string) (database.Warehouse, error)
	Get(id uint64) (database.Warehouse, error)
	Delete(id uint64) error
	Update(id uint64, code string, address string, telephone string, minimumCapacity uint32, minimumTemperature float32) (database.Warehouse, error)	
}

func NewService(warehouseRepo WarehouseRepository) WarehouseService {
	return &warehouseService{
		warehouseRepo: warehouseRepo,
	}
}

type warehouseService struct {
	warehouseRepo WarehouseRepository
}

func (s *warehouseService) GetAll() ([]database.Warehouse, error) {
	return s.warehouseRepo.GetAll()

}

func (s *warehouseService) Create(code string, address string, telephone string, minimumCapacity uint32, minimumTemperature float32, localityId string) (database.Warehouse, error) {
	isUsedCid, err := s.warehouseRepo.ExistsWarehouseCode(code)
	if err != nil {
		return database.Warehouse{}, err
	}

	if isUsedCid {
		return database.Warehouse{}, ExistsWarehouseCodeError
	}
	return s.warehouseRepo.Create(code, address, telephone, minimumCapacity, minimumTemperature, localityId)
}

func (s *warehouseService) Get(id uint64) (database.Warehouse, error) {
	foundWarehouse, err := s.warehouseRepo.Get(id)
	if err != nil {
		return database.Warehouse{}, WarehouseNotFoundError
	}

	return foundWarehouse, nil
}

func (s *warehouseService) Delete(id uint64) error {
	return s.warehouseRepo.Delete(id)
}

func (s *warehouseService) Update(id uint64, code string, address string, telephone string, minimumCapacity uint32, minimumTemperature float32) (database.Warehouse, error) {
	foundWarehouse, err := s.warehouseRepo.Get(id)
	if err != nil {
		return database.Warehouse{}, WarehouseNotFoundError
	}
	
	isUsedCid, err := s.warehouseRepo.ExistsWarehouseCode(code)
	if err != nil {
		return database.Warehouse{}, err
	}

	if isUsedCid {
		return database.Warehouse{}, ExistsWarehouseCodeError
	}
	updatedWarehouse := database.Warehouse{
		Id:                 id,
		Code:               code,
		Address:            address,
		Telephone:          telephone,
		MinimunCapacity:    minimumCapacity,
		MinimumTemperature: minimumTemperature,
	}
	mergo.Merge(&foundWarehouse, updatedWarehouse, mergo.WithOverride)
	newWarehouse, err := s.warehouseRepo.Update(foundWarehouse)
	if err != nil {
		return database.Warehouse{}, fmt.Errorf("error: internal server error")
	}
	return newWarehouse, nil
}
