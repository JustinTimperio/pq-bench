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
    let senders = 4;
    let receivers = 4;

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

    let mut handles = Vec::new();

    for _ in 0..senders {
        let rpq_clone = Arc::clone(&rpq);
        handles.push(tokio::spawn(async move {
            for i in 0..(message_count / senders) / batch_size {
                let mut items = Vec::new();
                for _ in 0..batch_size {
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
                let result = rpq_clone.enqueue_batch(items).await;
                if result.is_err() {
                    println!("Error Enqueuing: {}", result.err().unwrap());
                    return;
                }
            }
        }));
    }

    let recived_counter = Arc::new(std::sync::atomic::AtomicUsize::new(0));
    for _ in 0..receivers {
        let rpq_clone = Arc::clone(&rpq);
        let recived_counter = Arc::clone(&recived_counter);
        handles.push(tokio::spawn(async move {
            loop {
                if recived_counter.load(std::sync::atomic::Ordering::SeqCst) >= message_count {
                    break;
                }

                let result = rpq_clone.dequeue_batch(batch_size).await;
                if result.is_err() {
                    println!("Error Dequeuing: {}", result.err().unwrap());
                    return;
                }
                let result = result.unwrap();
                if !result.is_none() {
                    recived_counter
                        .fetch_add(result.unwrap().len(), std::sync::atomic::Ordering::SeqCst);
                }
            }
        }));
    }

    for handle in handles {
        handle.await.unwrap();
    }

    println!(
        "Received: {}",
        recived_counter.load(std::sync::atomic::Ordering::SeqCst)
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
