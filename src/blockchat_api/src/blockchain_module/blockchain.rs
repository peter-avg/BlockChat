use crate::blockchain_module::block::Block;
use serde::{Serialize, Deserialize};
use serde_json;

#[derive(Serialize, Deserialize, Debug)]
pub struct Blockchain { 
    chain: Vec<Block>,
}

impl Blockchain { 
    
    // Constructor function for struct Blockchain
    pub fn new(chain: Vec<Block>) -> Self { 
        Blockchain {
            chain
        }
    }

    // Jsonify Blockchain Object
    pub fn jsonify(&self) -> Result<String, serde_json::Error> {
        serde_json::to_string(self)
    }

}
