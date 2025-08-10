package service_internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"oprosdom.ru/msvc_auth/internal/biz"
	"oprosdom.ru/msvc_auth/internal/models"
	"oprosdom.ru/msvc_auth/internal/repo"
	"oprosdom.ru/msvc_auth/internal/transport"
	"oprosdom.ru/shared/models/pb/access"
)

type ServiceStruct struct {
	key           *models.KeyData
	ramRepo       repo.RamRepoInterface
	repo          repo.RepositoryInterface
	biz           biz.BizInterface
	codeTransport transport.TransportInterface
	accessClient  access.AccessClient
}

func NewCallInternalService(key *models.KeyData, ramRepo repo.RamRepoInterface, repo repo.RepositoryInterface, biz biz.BizInterface, codeTransport transport.TransportInterface, accessClient access.AccessClient) *ServiceStruct {
	return &ServiceStruct{
		key:           key,
		ramRepo:       ramRepo,
		repo:          repo,
		biz:           biz,
		codeTransport: codeTransport,
		accessClient:  accessClient,
	}
}

func (s *ServiceStruct) parsePhoneCode(value any) (*models.PhoneCode, error) {
	log.Printf("parsePhoneCode input type=%T value=%#v", value, value)
	strVal, ok := value.(string)
	if !ok {
		err := "phone_code_is_not_string"
		log.Println(err)
		return nil, errors.New(err)
	}

	var phoneCode models.PhoneCode
	if err := json.Unmarshal([]byte(strVal), &phoneCode); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("phone_code_parse_error: %w", err)
	}

	log.Printf("parsePhoneCode output: %+v", phoneCode)

	return &phoneCode, nil
}

func (s *ServiceStruct) parseUint32(value any) (uint32, error) {
	log.Printf("parseuint32 input type=%T value=%#v", value, value)
	strVal, ok := value.(string)
	if !ok {
		err := "strRedis is not a string"
		log.Println(err)
		log.Printf("strRedis is not a string, actual type: %T, value: %#v", value, value)
		return 0, errors.New(err)
	}
	existRetryInt, err := strconv.Atoi(strVal)
	if err != nil {
		log.Printf("atoi returns err: %v", err)
		return 0, err
	}
	return uint32(existRetryInt), nil
}
