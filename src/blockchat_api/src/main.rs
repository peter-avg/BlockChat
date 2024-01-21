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

    // let json = first_transaction.jsonify();
    // match json {
    //     Ok(json) => println!("{}",json),
    //     Err(e) => println!("{}",e),
    // }

    println!("{}", first_transaction.signature);

}
