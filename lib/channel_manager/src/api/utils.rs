use std::os::raw::c_char;
use bioenpro4to_channel_manager::channels::ChannelInfo as ChInfo;
use std::ffi::{CString, CStr};
use anyhow::Result;
use bioenpro4to_channel_manager::utils::{create_encryption_key, create_encryption_nonce};

#[repr(C)]
pub struct ChannelInfo{
    pub channel_id: *const c_char,
    pub announce_id: *const c_char
}

impl ChannelInfo{
    pub fn from_ch_info(ch_info: ChInfo) -> *const ChannelInfo{
        let channel_id = CString::new(ch_info.channel_id()).unwrap().into_raw();
        let announce_id = CString::new(ch_info.announce_id()).unwrap().into_raw();
        new_channel_info(channel_id, announce_id)
    }

    pub unsafe fn to_ch_info(&self) -> Result<ChInfo>{
        let channel_id = CStr::from_ptr(self.channel_id).to_str()?;
        let announce_id = CStr::from_ptr(self.announce_id).to_str()?;
        Ok(ChInfo::new(channel_id.to_string(), announce_id.to_string()))
    }
}

#[no_mangle]
pub extern "C" fn new_channel_info(channel_id: *const c_char, announce_id: *const c_char) -> *const ChannelInfo{
    let ch_info = ChannelInfo{channel_id, announce_id};
    Box::into_raw(Box::new(ch_info))
}

#[no_mangle]
pub unsafe extern "C" fn drop_channel_info(info: *mut ChannelInfo){
    Box::from_raw(info);
}

#[repr(C)]
pub struct KeyNonce{
    pub key: [u8; 32],
    pub nonce: [u8; 24],
}

#[no_mangle]
pub extern "C" fn new_encryption_key_nonce(key: *const c_char, nonce: *const c_char) -> *const KeyNonce{
    unsafe {
        let k = CStr::from_ptr(key).to_str().unwrap();
        let n = CStr::from_ptr(nonce).to_str().unwrap();
        let k = create_encryption_key(k);
        let n = create_encryption_nonce(n);

        Box::into_raw(Box::new(KeyNonce{key: k, nonce: n}))
    }
}

#[no_mangle]
pub unsafe extern "C" fn drop_key_nonce(kn: *const KeyNonce) {
    Box::from_raw(kn as *mut KeyNonce);
}

#[repr(C)]
pub struct RawPacket{
    pub public: *const u8,
    pub p_len: usize,
    pub masked: *const u8,
    pub m_len: usize
}

impl RawPacket{
    pub unsafe fn public(&self) -> Vec<u8>{
        let p = std::slice::from_raw_parts(self.public, self.p_len);
        p.to_vec()
    }

    pub unsafe fn masked(&self) -> Vec<u8>{
        let m = std::slice::from_raw_parts(self.masked, self.m_len);
        m.to_vec()
    }
}

#[no_mangle]
pub extern "C" fn new_raw_packet(public: *mut u8, p_len: u64,
                                 masked: *mut u8, m_len: u64) -> *const RawPacket{
    let p_len = p_len as usize;
    let m_len = m_len as usize;
    let packet = RawPacket{public, p_len, masked, m_len};
    Box::into_raw(Box::new(packet))
}

#[no_mangle]
pub unsafe extern "C" fn drop_raw_packet(packet: *mut RawPacket){
    Box::from_raw(packet);
}

#[no_mangle]
pub unsafe extern "C" fn drop_str(s: *const c_char) {
    CString::from_raw(s as *mut c_char);
}
