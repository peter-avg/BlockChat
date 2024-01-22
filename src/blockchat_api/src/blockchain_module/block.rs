use crate::blockchain_module::transaction::Transaction;
use sha256::digest;
use serde::{Serialize, Deserialize};
use serde_json;

#[derive(Clone,Serialize, Deserialize, Debug)]
pub struct Block {
    pub index: i32,
    pub timestamp: i64,
    pub transactions: Vec<Transaction>,
    pub validator: i32,
    pub current_hash: String,
    pub previous_hash: i32,
}

impl Block { 
    
    // Constructor function for struct Block
    pub fn new(index: i32, timestamp: i64, transactions: Vec<Transaction>,
                validator: i32, current_hash: String, previous_hash: i32) -> Self {
        Block { 
            index,
            timestamp,
            transactions,
            validator,
            current_hash,
            previous_hash
        }
    }

    // Jsonify Block Object 
    pub fn jsonify(&self) -> Result<String, serde_json::Error> {
        serde_json::to_string(self)
    }

    // Create a hash for the Block Object
    pub fn hashify(&mut self) {
        let json = self.jsonify();
        match json {
            Ok(json) => self.current_hash = digest(json),
            Err(_e) => println!("Could not create hash"),
        }
    }
    
    // Adding a transaction in a block, not exceeding its capacity
    pub fn add_transaction(&mut self, transaction: Transaction, capacity: usize) {
        if self.transactions.len() < capacity {
            self.transactions.push(transaction);
        }
    }

}
