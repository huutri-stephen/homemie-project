package domain

import "homemie/models/dto"

type AddressRepository interface {
	Create(address *dto.Address) (createdAddress *dto.Address, err error)
}