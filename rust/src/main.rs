use rand::Rng;
use std::collections::BinaryHeap;
use std::time::Instant;

fn main() {
    let mut heap = BinaryHeap::new();

    let start = Instant::now();

    // Push 10 million integers onto the heap
    for _ in 0..10_000_000 {
        let mut rng = rand::thread_rng();
        let p = rng.gen_range(1..101);
        heap.push(p);
    }

    let mid = Instant::now();

    // Pop 10 million integers from the heap
    while let Some(_top) = heap.pop() {}

    let end = Instant::now();

    println!(
        "Time to insert 10 million integers: {:?}",
        mid.duration_since(start)
    );
    println!(
        "Time to retrieve 10 million integers: {:?}",
        end.duration_since(mid)
    );
    println!("Total time: {:?}", end.duration_since(start));
}
