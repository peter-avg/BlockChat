use serde::{Serialize, Deserialize};
use serde_json;
use sha256::digest;
use rsa::{RsaPrivateKey, RsaPublicKey};
use sha2::{Sha256, Digest};
use rand::rngs::OsRng;
use rand::rngs::ThreadRng;

#[derive(Clone,Serialize, Deserialize, Debug)]
pub struct Transaction {
    pub sender_address: i32,
    pub receiver_address: i32,
    // 0 for message, 1 for bcc
    pub type_of_transaction: bool,
    pub data: String,
    pub nonce: i32,
    pub transaction_id: i32,
    pub signature: String,
}

impl Transaction {
    
    // Constructor function for struct Transaction
    pub fn new(sender_address: i32, receiver_address: i32,
               type_of_transaction: bool, data: String,
               nonce: i32, transaction_id: i32, signature: String) -> Self {
        Transaction {
            sender_address,
            receiver_address,
            type_of_transaction,
            data,
            nonce,
            transaction_id,
            signature,
        }

    }

    // Jsonify Transaction Object
    pub fn jsonify(&self) -> Result<String, serde_json::Error> {
        serde_json::to_string(self)
    }
    
    // Create a hash for the Transaction Object
    pub fn hashify(&mut self) -> Option<String> {
        let json = self.jsonify();
        match json {
            Ok(json) => return Some(digest(json)),
            Err(_e) => return None,
        }
    }

    // Sign the transaction without padding
    pub fn sign_transaction(&mut self, private_key: RsaPrivateKey, rng: ThreadRng) {

    }




}
