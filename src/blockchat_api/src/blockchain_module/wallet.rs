use crate::blockchain_module::transaction::Transaction; 
use rand::thread_rng;
use rsa::{RsaPrivateKey, RsaPublicKey};

#[derive(Debug)]
pub struct Wallet {
    publickey: RsaPublicKey,
    privatekey: RsaPrivateKey,
    transactions: Vec<Transaction>,
    balance: i32,
    node_id: i32,
    nonce: i32,
}

impl Wallet {
    
    // Constructor function for struct Wallet
    pub fn _new(bits: usize) -> Self {
        let mut rng = thread_rng();
        let priv_key = RsaPrivateKey::new(&mut rng, bits)
                                     .expect("failed to generate a key");
        let pub_key = RsaPublicKey::from(&priv_key);

        Wallet {
            publickey: pub_key,
            privatekey: priv_key,
            transactions: Vec::new(),
            balance: 0,
            node_id: 0,
            nonce: 0,
        }
    }

    // Add Transaction to Wallet
    pub fn _add_transaction(&mut self,mut _new_transaction: Transaction) -> i32 {
        self.nonce += 1;
        _new_transaction.nonce = self.nonce;
        self.transactions.push(_new_transaction);
        return self.nonce;
    }


}
