mod blockchain_module;
use blockchain_module::transaction::Transaction;
// use blockchain_module::block::Block;
// use blockchain_module::blockchain::Blockchain;
use blockchain_module::wallet::Wallet;


fn main() {

    let mut first_transaction = Transaction::new(
        2,
        2,
        true,
        2.9.to_string(),
        2,
        2);
    

    // let mut first_block = Block::new(
    //     2,
    //     2,
    //     [].to_vec(),
    //     2,
    //     "".to_string(),
    //     2,
    //     );
    //
    // let mut blockchain = Blockchain::new(
    //     [].to_vec()
    //     );

    let mut wallet = Wallet::_new();

    first_transaction.signature = wallet.sign_transaction(&first_transaction);

    if wallet.verify_transaction(&first_transaction) {
        println!("Transaction was verified");
    }

}
