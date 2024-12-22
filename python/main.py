from queue import PriorityQueue, Empty
import time


def main():

    queue = PriorityQueue()
    message_count = 10_000_000

    start_insert_time = time.time()
    for i in range(message_count):
        queue.put(i)

    end_insert_time = time.time()

    start_fetch_time = time.time()
    for i in range(message_count):
        try:
            queue.get()
        except Empty as e:
            print("Queue is empty")
            break
    end_fetch_time = time.time()

    insertion_time = end_insert_time - start_insert_time
    fetch_time = end_fetch_time - start_fetch_time

    print(f"Insertion Time: {insertion_time}")
    print(f"Fetch Time: {fetch_time}")
    print(f"Total Time: {fetch_time+insertion_time}")


main()
