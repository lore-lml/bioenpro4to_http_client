use std::ptr::{null, null_mut};
use crate::api::utils::{ChannelInfo, RawPacket, KeyNonce};
use std::os::raw::c_char;
use std::ffi::{CStr, CString};
use tokio::runtime::Runtime;
use bioenpro4to_channel_manager::channels::DailyChannelManager;

#[no_mangle]
pub unsafe extern "C" fn daily_channel_from_base64(state: *const c_char, state_psw: *const c_char) -> *mut DailyChannelManager{
    let state = match CStr::from_ptr(state).to_str(){
        Ok(state) => state,
        Err(_) => return null_mut()
    };

    let state_psw = match CStr::from_ptr(state_psw).to_str(){
        Ok(state) => state,
        Err(_) => return null_mut()
    };

    Runtime::new().unwrap().block_on(async {
        match DailyChannelManager::import_from_base64(state, state_psw).await{
            Ok(ch) => Box::into_raw(Box::new(ch)),
            Err(_) => null_mut()
        }
    })
}

#[no_mangle]
pub unsafe extern "C" fn drop_daily_channel_manager(channel: *mut DailyChannelManager){
    channel.drop_in_place();
}

#[no_mangle]
pub unsafe extern "C" fn send_raw_packet(root: *mut DailyChannelManager, packet: *const RawPacket, key_nonce: *const KeyNonce) -> *const c_char{
    let root = root.as_mut();
    let p = packet.as_ref();
    let kn = key_nonce.as_ref();

    match (&root, &p){
        (None, _) => return null(),
        (_, None) => return null(),
        _ => {}
    };

    let root = root.unwrap();
    let p = p.unwrap();
    let public = p.public();
    let masked = p.masked();
    let opt_kn = match kn{
        None => None,
        Some(kn) => Some((kn.key.clone(), kn.nonce.clone()))
    };

    let res = Runtime::new().unwrap().block_on(async {
        root.send_raw_packet(public, masked, opt_kn).await
    });

    match res{
        Ok(res) => CString::new(res).map_or(null(), |h| h.into_raw()),
        Err(_) => null()
    }
}

#[no_mangle]
pub unsafe extern "C" fn daily_channel_info(channel: *mut DailyChannelManager) -> *const ChannelInfo{
    let ch = match channel.as_mut(){
        None => return null_mut(),
        Some(ch) => ch
    };
    ChannelInfo::from_ch_info(ch.channel_info())
}
