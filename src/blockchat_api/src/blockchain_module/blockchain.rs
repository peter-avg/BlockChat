use crate::blockchain_module::block::Block;

#[derive(Debug)]
pub struct Blockchain { 
    chain: Vec<Block>,
}

impl Blockchain { 
    
    // Constructor function for struct Blockchain
    pub fn new(chain: Vec<Block>) -> Blockchain { 
        Blockchain {
            chain
        }
    }

}
