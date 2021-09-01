use iota_identity_lib::api::{IdentityManager, Storage};
use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use tokio::runtime::Runtime;
use std::ptr::{null_mut, null};
use iota_identity_lib::iota::Credential;

#[repr(C)]
pub struct CCredential{
    credential: *const c_char,
}

impl CCredential {
    fn new(credential: *const c_char) -> Self {
        CCredential { credential }
    }

    fn from_credential(credential: &Credential) -> anyhow::Result<Self>{
        let cred_str = serde_json::to_string(credential)?;
        let c_str = CString::new(cred_str).map_or_else(|err| Err(anyhow::Error::msg(err.to_string())), |h| Ok(h.into_raw()))?;
        Ok(CCredential::new(c_str))
    }

    unsafe fn to_credential(&self) -> anyhow::Result<Credential>{
        let cred = CStr::from_ptr(self.credential).to_str()
            .map_or_else(|err| Err(anyhow::Error::msg(err.to_string())), |c| Ok(c))?;
        serde_json::from_str(cred).map_err(|err| anyhow::Error::msg(err.to_string()))
    }
}

#[no_mangle]
pub unsafe extern "C" fn new_identity_manager(mainnet: usize, persistence: usize, folder_dir: *const c_char, vault_psw: *const c_char) -> *mut IdentityManager{
    let mainnet = match mainnet{
        0 => false,
        _ => true
    };

    let storage = if persistence == 0 || folder_dir == null() || vault_psw == null(){
        Storage::Memory
    }else{
        let folder_dir = CStr::from_ptr(folder_dir).to_str();
        let vault_psw = CStr::from_ptr(vault_psw).to_str();
        match (folder_dir, vault_psw){
            (Ok(dir), Ok(psw)) => Storage::Stronghold(dir.to_string(), Some(psw.to_string())),
            (_, _) => Storage::Memory
        }
    };

    let future = IdentityManager::builder()
        .main_net(mainnet)
        .storage(storage)
        .build();

    Runtime::new().unwrap().block_on(async{
        match future.await{
            Ok(manager) => Box::into_raw(Box::new(manager)),
            Err(_) => null_mut()
        }
    })
}

#[no_mangle]
pub unsafe extern "C" fn drop_identity_manager(identity_manager: *mut IdentityManager){
    identity_manager.drop_in_place()
}

#[no_mangle]
pub unsafe extern "C" fn create_identity(identity_manager: *mut IdentityManager, identity_name: *const c_char) -> *const c_char{
    let manager = match identity_manager.as_mut(){
        None => return null(),
        Some(ch) => ch
    };

    let identity_name = match CStr::from_ptr(identity_name).to_str(){
        Ok(name) => name,
        Err(_) => return null()
    };

    Runtime::new().unwrap().block_on(async{
        match manager.create_identity(identity_name).await{
            Ok(doc) => CString::new(doc.id().to_string()).map_or(null(), |h| h.into_raw()),
            Err(_) => null()
        }
    })
}

#[no_mangle]
pub unsafe extern "C" fn get_identity_did(identity_manager: *mut IdentityManager, identity_name: *const c_char) -> *const c_char{
    let manager = match identity_manager.as_mut(){
        None => return null(),
        Some(ch) => ch
    };

    let identity_name = match CStr::from_ptr(identity_name).to_str(){
        Ok(name) => name,
        Err(_) => return null()
    };

    manager.get_identity(identity_name)
        .map_or(
            null(),
            |doc| CString::new(doc.id().to_string()).map_or(null(), |h| h.into_raw())
        )
}

#[no_mangle]
pub unsafe extern "C" fn store_credential(identity_manager: *mut IdentityManager, identity_name: *const c_char, cred_name: *const c_char, cred: *const CCredential) -> usize{
    let (manager, cred, identity_name, cred_name) = match (
        identity_manager.as_mut(),
        cred.as_ref(),
        CStr::from_ptr(identity_name).to_str(),
        CStr::from_ptr(cred_name).to_str()
    ){
        (Some(manager), Some(cred), Ok(id_name), Ok(cred_name)) => (manager, cred, id_name, cred_name),
        (_, _, _, _) => return 0
    };

    let cred = match cred.to_credential(){
        Ok(c) => c,
        Err(_) => return 0
    };

    manager.store_credential(identity_name, cred_name, &cred)
        .map_or(0, |_| 1)
}

#[no_mangle]
pub unsafe extern "C" fn get_credential(identity_manager: *mut IdentityManager, identity_name: *const c_char, cred_name: *const c_char) -> *const CCredential{
    let (manager, identity_name, cred_name) = match (
        identity_manager.as_mut(),
        CStr::from_ptr(identity_name).to_str(),
        CStr::from_ptr(cred_name).to_str()
    ){
        (Some(manager), Ok(id_name), Ok(cred)) => (manager, id_name, cred),
        (_, _, _) => return null()
    };

    manager.get_credential(identity_name, cred_name)
        .map_or(
            null(),
            |c| CCredential::from_credential(c).map_or(null(), |c| Box::into_raw(Box::new(c)))
        )
}

#[no_mangle]
pub unsafe extern "C" fn new_ccredential(cred: *const c_char) -> *const CCredential{
    Box::into_raw(Box::new(CCredential::new(cred)))
}

#[no_mangle]
pub unsafe extern "C" fn drop_ccredential(cred: *mut CCredential){
    cred.drop_in_place()
}
