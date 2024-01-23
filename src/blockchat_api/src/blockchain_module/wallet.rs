use crate::blockchain_module::transaction::Transaction; 
use rand::thread_rng;
use rsa::{RsaPrivateKey, RsaPublicKey};
use rsa::pss::{Pss, BlindedSigningKey, VerifyingKey, Signature};
use rsa::signature::{Keypair,RandomizedSigner, SignatureEncoding, Verifier};
use sha2::{Digest, Sha256};
use rand::rngs::ThreadRng;

#[derive(Debug)]
pub struct Wallet {
    pub public_key: VerifyingKey<Sha256>,
    pub private_key: BlindedSigningKey<Sha256>,
    pub transactions: Vec<Transaction>,
    pub balance: i32,
    pub node_id: i32,
    pub nonce: i32,
    pub rng: ThreadRng
}

impl Wallet {
    
    // Constructor function for struct Wallet
    pub fn _new() -> Self {
        let mut my_rng = thread_rng();

        let bits = 2048;
        let private_key = RsaPrivateKey::new(&mut my_rng, bits).expect("failed to generate a key");
        let signing_key = BlindedSigningKey::<Sha256>::new(private_key);
        let verifying_key = signing_key.verifying_key();

        Wallet {
            public_key: verifying_key,
            private_key: signing_key,
            transactions: Vec::new(),
            balance: 0,
            node_id: 0,
            nonce: 0,
            rng: my_rng
        }
    }

    // Add Transaction to Wallet
    pub fn _add_transaction(&mut self,mut _new_transaction: Transaction) -> i32 {
        self.nonce += 1;
        _new_transaction.nonce = self.nonce;
        self.transactions.push(_new_transaction);
        return self.nonce;
    }
    
    // Sign the transaction
    pub fn sign_transaction(&mut self, transaction: &Transaction) -> Option<Signature> {
        return Some(self.private_key.sign_with_rng(&mut self.rng, &transaction.data.as_bytes()));
    }

    // Verify a transaction
    pub fn verify_transaction(&self, transaction: &Transaction) -> bool {
        if let Some(signature) = &transaction.signature {
            let message = transaction.data.as_bytes();
            println!("{:?}",&message);
            self.public_key.verify(&message, &signature).expect("Could not verify transaction");
            return true;
        } else {
            false
        }
    }


}
