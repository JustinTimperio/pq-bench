use std::fs::File;
use std::io::Write;
use std::sync::Arc;

use pprof::protos::Message;
use rpq::pq::Item;
use rpq::{RPQOptions, RPQ};

#[tokio::main(flavor = "multi_thread")]
async fn main() {
    let guard = pprof::ProfilerGuard::new(100).unwrap();

    let message_count = 10_000_000;
    let batch_size = 10_000;

    let options = RPQOptions {
        max_priority: 10,
        disk_cache_enabled: false,
        database_path: "/tmp/rpq-bench.redb".to_string(),
        lazy_disk_cache: true,
        lazy_disk_write_delay: std::time::Duration::from_secs(5),
        lazy_disk_cache_batch_size: 10000,
    };

    let r = RPQ::new(options).await;
    match r {
        Ok(_) => {}
        Err(e) => {
            println!("Error Creating RPQ: {}", e);
            return;
        }
    }

    let rpq = Arc::clone(&r.unwrap().0);

    let timer = std::time::Instant::now();
    let send_timer = std::time::Instant::now();
    for i in 0..message_count / batch_size {
        let mut items = Vec::new();
        for j in 0..batch_size {
            let item = Item::new(
                i % 10,
                i,
                false,
                None,
                false,
                Some(std::time::Duration::from_secs(5)),
            );
            items.push(item);
        }
        let result = rpq.enqueue_batch(items).await;
        if result.is_err() {
            println!("Error Enqueuing: {}", result.err().unwrap());
            return;
        }
    }

    let send_elapsed = send_timer.elapsed().as_secs_f64();

    let receive_timer = std::time::Instant::now();
    for _i in 0..message_count / batch_size {
        let result = rpq.dequeue_batch(batch_size).await;
        if result.is_err() {
            println!("Error Dequeuing: {}", result.err().unwrap());
            return;
        }
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

    match guard.report().build() {
        Ok(report) => {
            let mut file = File::create("profile.pb").unwrap();
            let profile = report.pprof().unwrap();

            let mut content = Vec::new();
            profile.write_to_vec(&mut content).unwrap();
            file.write_all(&content).unwrap();
        }
        Err(_) => {}
    };
}
