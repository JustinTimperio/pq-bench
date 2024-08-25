use rpq::pq::Item;
use rpq::{RPQOptions, RPQ};
use std::sync::Arc;

fn main() {
    let message_count = 10_000_000;

    let options = RPQOptions {
        bucket_count: 10,
        disk_cache_enabled: false,
        database_path: "/tmp/rpq.redb".to_string(),
        lazy_disk_cache: false,
        lazy_disk_max_delay: std::time::Duration::from_secs(5),
        lazy_disk_cache_batch_size: 5000,
        buffer_size: 1_000_000,
    };

    let r = Arc::new(
        tokio::runtime::Runtime::new()
            .unwrap()
            .block_on(RPQ::new(options)),
    );

    let rpq = Arc::clone(&r.0);
    let runtime = tokio::runtime::Runtime::new().unwrap();

    let timer = std::time::Instant::now();
    let send_timer = std::time::Instant::now();
    for i in 0..message_count {
        runtime.block_on(async {
            let item = Item::new(
                i % 10,
                i,
                false,
                None,
                false,
                Some(std::time::Duration::from_secs(5)),
            );
            rpq.enqueue(item).await;
        });
    }
    let send_elapsed = send_timer.elapsed().as_secs_f64();

    let receive_timer = std::time::Instant::now();
    for _i in 0..message_count {
        runtime.block_on(async {
            rpq.dequeue().await;
        });
    }
    let receive_elapsed = receive_timer.elapsed().as_secs_f64();

    println!(
        "Time to insert {} messages: {}s",
        message_count, send_elapsed
    );
    println!(
        "Time to receive {} messages: {}s",
        message_count, receive_elapsed
    );
    println!("Total Time: {}s", timer.elapsed().as_secs_f64());
}
