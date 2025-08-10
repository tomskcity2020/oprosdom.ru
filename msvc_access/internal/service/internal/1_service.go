package service_internal

import (
	"oprosdom.ru/msvc_access/internal/repo"
)

type ServiceStruct struct {
	ramRepo repo.RamRepoInterface
}

func NewCallInternalService(ramRepo repo.RamRepoInterface) *ServiceStruct {
	return &ServiceStruct{
		ramRepo: ramRepo,
	}
}
