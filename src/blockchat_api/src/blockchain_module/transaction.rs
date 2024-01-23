use serde::{Serialize, Serializer, ser::SerializeStruct};
use serde_json;
use sha256::digest;
use rsa::pss::Signature;
use rsa::signature::SignatureEncoding;

#[derive(Clone, Debug)]
pub struct Transaction {
    pub sender_address: i32,
    pub receiver_address: i32,
    // 0 for message, 1 for bcc
    pub type_of_transaction: bool,
    pub data: String,
    pub nonce: i32,
    pub transaction_id: i32,
    pub signature: Option<Signature>,
}

impl Transaction {
    
    // Constructor function for struct Transaction
    pub fn new(sender_address: i32, receiver_address: i32,
               type_of_transaction: bool, data: String,
               nonce: i32, transaction_id: i32) -> Self {
        Transaction {
            sender_address,
            receiver_address,
            type_of_transaction,
            data,
            nonce,
            transaction_id,
            signature: None,
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
}

impl Serialize for Transaction {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut state = serializer.serialize_struct("Transaction", 7)?;
        state.serialize_field("sender_address", &self.sender_address)?;
        state.serialize_field("receiver_address", &self.receiver_address)?;
        state.serialize_field("type_of_transaction", &self.type_of_transaction)?;
        state.serialize_field("data", &self.data)?;
        state.serialize_field("nonce", &self.nonce)?;
        state.serialize_field("transaction_id", &self.transaction_id)?;

        // Serialize signature as a byte array
        if let Some(ref signature) = self.signature {
            let sig_bytes = signature.to_bytes(); // Convert signature to bytes
            state.serialize_field("signature", &sig_bytes)?;
        } else {
            state.serialize_field::<Option<Vec<u8>>>("signature", &None)?;
        }

        state.end()
    }
}
