package identity_manager

/*
#cgo LDFLAGS: -L./.. -lc_channel_manager_lib
#include "../c_identity_manager.h"
*/
import "C"
import (
	"bioenpro4to_http_client/lib/bep4t_http_client/utils"
	"errors"
	"fmt"
)

type IdentityManager struct {
	manager *C.identity_manager_t
}

type PersistenceConfig struct {
	FolderDir string
	VaultPsw  string
}

func NewPersistenceConfig(folderDir, vaultPsw string) *PersistenceConfig {
	return &PersistenceConfig{
		FolderDir: folderDir,
		VaultPsw:  vaultPsw,
	}
}

func NewIdentityManager(mainnet bool, persistenceConfig *PersistenceConfig) (*IdentityManager, error) {
	var cMainnet = 0
	if mainnet {
		cMainnet = 1
	}

	var manager *C.identity_manager_t
	if persistenceConfig == nil {
		manager = C.new_identity_manager(C.int(cMainnet), C.int(0), nil, nil)
	} else {
		folderDir := C.CString(persistenceConfig.FolderDir)
		vaultPsw := C.CString(persistenceConfig.VaultPsw)
		manager = C.new_identity_manager(C.int(cMainnet), C.int(1), folderDir, vaultPsw)
	}

	if manager == nil {
		return nil, errors.New("invalid folder/psw or network problems")
	}

	return &IdentityManager{manager: manager}, nil
}

func (self *IdentityManager) Drop() {
	C.drop_identity_manager(self.manager)
}

func (self *IdentityManager) CreateIdentity(identityName string) (string, error) {
	cIdName := C.CString(identityName)
	cDid := C.create_identity(self.manager, cIdName)
	if cDid == nil {
		return "", errors.New(fmt.Sprintf("identity with name %s already exists", identityName))
	}
	defer C.drop_str(cDid)
	return C.GoString(cDid), nil
}

func (self *IdentityManager) GetIdentityDid(identityName string) (string, error) {
	cIdName := C.CString(identityName)
	cDid := C.get_identity_did(self.manager, cIdName)
	if cDid == nil {
		return "", errors.New(fmt.Sprintf("identity with name %s not found", identityName))
	}
	defer C.drop_str(cDid)
	return C.GoString(cDid), nil
}

func (self *IdentityManager) StoreCredential(credName string, cred utils.Credential) error {
	cCredName := C.CString(credName)
	cCred := C.new_ccredential(C.CString(string(cred)))
	if cCred == nil {
		return errors.New("error during credential parsing")
	}
	defer C.drop_ccredential(cCred)
	res := int(C.store_credential(self.manager, cCredName, cCred))
	if res == 0 {
		return errors.New(fmt.Sprintf("credential with name %s already exists", credName))
	}
	return nil
}

func (self *IdentityManager) GetCredential(credName string) (utils.Credential, error) {
	cCredName := C.CString(credName)
	cCred := C.get_credential(self.manager, cCredName)
	if cCred == nil {
		return nil, errors.New(fmt.Sprintf("credential with name %s not found", credName))
	}
	credential := []byte(C.GoString(cCred.cred))
	return credential, nil
}
