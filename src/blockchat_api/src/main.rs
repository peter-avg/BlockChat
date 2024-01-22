mod blockchain_module;
use blockchain_module::transaction::Transaction;
use blockchain_module::block::Block;
use blockchain_module::blockchain::Blockchain;


fn main() {

    let first_transaction = Transaction::new(
        2,
        2,
        true,
        2.9,
        "hello".to_string(),
        2,
        2,
        "hello".to_string()
        );

    let mut first_block = Block::new(
        2,
        2,
        [].to_vec(),
        2,
        "".to_string(),
        2,
        );

    let mut blockchain = Blockchain::new(
        [].to_vec()
        );

    first_block.add_transaction(first_transaction, 2);
    first_block.hashify();
    blockchain.add_block(first_block);
    println!("{:?}", blockchain);

}
