use crate::blockchain_module::transaction::Transaction;

#[derive(Debug)]
pub struct Block {
    pub index: i32,
    pub timestamp: i32,
    pub transactions: Vec<Transaction>,
    pub validator: i32,
    pub current_hash: i32,
    pub previous_hash: i32,
}

impl Block { 
    
    // Constructor function for struct Block
    pub fn new(index: i32, timestamp: i32, transactions: Vec<Transaction>,
                validator: i32, current_hash: i32, previous_hash: i32) -> Block {
        Block { 
            index,
            timestamp,
            transactions,
            validator,
            current_hash,
            previous_hash
        }
    }


}
