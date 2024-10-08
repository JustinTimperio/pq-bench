const std = @import("std");
const Order = std.math.Order;
const Allocator = std.mem.Allocator;
const PriorityQueue = std.PriorityQueue;
const Time = std.time;

fn lessThan(context: void, a: i32, b: i32) Order {
    _ = context;
    return std.math.order(a, b);
}

pub fn main() !void {
    // Create an arena allocator
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    // Create a priority queue
    const PQlt = PriorityQueue(i32, void, lessThan);
    var pq = PQlt.init(allocator, {});

    // Create Random Number Generator
    var prng = std.rand.DefaultPrng.init(blk: {
        var seed: u64 = undefined;
        try std.posix.getrandom(std.mem.asBytes(&seed));
        break :blk seed;
    });
    const rand = prng.random();

    const start = try Time.Instant.now();

    // Push 10 million integers onto the priority queue
    const istart = try Time.Instant.now();
    var i: i32 = 0;
    while (i < 10000000) : (i += 1) {
        const p = @as(i32, @intCast(rand.intRangeAtMost(i32, 0, 100))); // get a random number between 0 and 100
        try pq.add(p);
    }
    const iend = try Time.Instant.now();
    const diff_insert = Time.Instant.since(iend, istart);

    // Pop 10 million integers from the priority queue
    const rstart = try Time.Instant.now();
    while (pq.count() > 0) {
        _ = pq.remove();
    }
    const rend = try Time.Instant.now();
    const finish = try Time.Instant.now();
    const diff_remove = Time.Instant.since(rend, rstart);
    const diff_total = Time.Instant.since(finish, start);

    const diff_insert_seconds = @as(f64, @floatFromInt(diff_insert)) / 1_000_000_000.0;
    const diff_remove_seconds = @as(f64, @floatFromInt(diff_remove)) / 1_000_000_000.0;
    const diff_total_seconds = @as(f64, @floatFromInt(diff_total)) / 1_000_000_000.0;

    try std.io.getStdOut().writer().print("Time to insert 10 million integers: {}s\n", .{diff_insert_seconds});
    try std.io.getStdOut().writer().print("Time to retrieve 10 million integers: {}s\n", .{diff_remove_seconds});
    try std.io.getStdOut().writer().print("Total time: {}s\n", .{diff_total_seconds});
}
