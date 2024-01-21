use crate::blockchain_module::transaction::Transaction; 

pub struct Wallet {
    publickey: i32,
    privatekey: i32,
    balance: i32,
    node_id: i32,
}

impl Wallet {

    // Constructor function for struct Wallet
    pub fn new(publickey: i32, privatekey: i32, balance: i32, node_id: i32) -> Wallet { 
        Wallet {
            publickey,
            privatekey,
            balance,
            node_id
        }
    }

    // Function to sign a new transaction with user's private key (sender)
    pub fn sign_transaction(_new_transaction: Transaction) -> bool {
        return true;
    }

    // Function to verify a new transaction with user's public key (receiver)
    pub fn verify_transaction(_new_transaction: Transaction) -> bool {
        return true;
    }

}
