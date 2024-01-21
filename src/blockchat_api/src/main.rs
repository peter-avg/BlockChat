mod blockchain_module;


fn main() {
    let first_transaction = blockchain_module::transaction::Transaction::new(
        2,
        2,
        true,
        2.9,
        "hello".to_string(),
        2,
        2,
        "hello".to_string()
        );
    println!("{:?}", first_transaction);
}
