package main

import(
	"fmt"
  "testing"
  "time"
)

func OutputTestResults(test_name string, num_paths int64, elapsed_time time.Duration) {
    fmt.Printf("%-45s %20d %15s %10v\n", test_name, num_paths, " elapsed time: ", elapsed_time)
}

func TestUnsolvableLinear(t *testing.T) {
    input_map := generate_open_map(4)
    make_unsolvable_map_bad_start(input_map)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    OutputTestResults("unsolvable linear solution: ", num_paths, elapsed)
}

func TestUnsolvableQueueAndThread(t *testing.T) {
    input_map := generate_open_map(4)
    make_unsolvable_map_bad_end(input_map)

    start := time.Now()
    num_paths := solve_field_with_channel_and_threads(input_map)
    elapsed := time.Since(start)

    OutputTestResults("unsolvable queue & thread solution: ", num_paths, elapsed)
}

func Test4x4Linear(t *testing.T) {
    input_map := generate_open_map(4)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    OutputTestResults("4x4 linear solution: ", num_paths, elapsed)
}

func Test4x4Queue(t *testing.T) {
    input_map := generate_open_map(4)

    start := time.Now()
    num_paths := solve_field_with_channel(input_map)
    elapsed := time.Since(start)

    OutputTestResults("4x4 queue solution: ", num_paths, elapsed)
}

func Test4x4QueueWithThreads(t *testing.T) {
    input_map := generate_open_map(4)

    start := time.Now()
    num_paths := solve_field_with_channel_and_threads(input_map)
    elapsed := time.Since(start)

    OutputTestResults("4x4 queue & threads solution: ", num_paths, elapsed)
}

func Test50x50Linear(t *testing.T) {
    input_map := generate_open_map(50)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    OutputTestResults("50x50 linear solution: ", num_paths, elapsed)
}

func Test50x50Queue(t *testing.T) {
    input_map := generate_open_map(50)

    start := time.Now()
    num_paths := solve_field_with_channel(input_map)
    elapsed := time.Since(start)

    OutputTestResults("50x50 queue solution: ", num_paths, elapsed)
}

func Test50x50QueueWithThreads(t *testing.T) {
    input_map := generate_open_map(50)

    start := time.Now()
    num_paths := solve_field_with_channel_and_threads(input_map)
    elapsed := time.Since(start)

    OutputTestResults("50x50 channel & threads solution: ", num_paths, elapsed)
}

func Test50x50WithWallLinear(t *testing.T) {
    input_map := generate_open_map(50)
    make_map_1(input_map)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    OutputTestResults("50x50 linear w/ wall solution: ", num_paths, elapsed)
}

func Test50x50WithWallQueue(t *testing.T) {
    input_map := generate_open_map(50)
    make_map_1(input_map)

    start := time.Now()
    num_paths := solve_field_with_channel(input_map)
    elapsed := time.Since(start)

    OutputTestResults("50x50 channel w/ wall solution: ", num_paths, elapsed)
}

func Test50x50WithWallQueueWithThreads(t *testing.T) {
    input_map := generate_open_map(50)
    make_map_1(input_map)

    start := time.Now()
    num_paths := solve_field_with_channel_and_threads(input_map)
    elapsed := time.Since(start)

    OutputTestResults("50x50 channel & threads w/ wall solution: ", num_paths, elapsed)
}

func Test100x100Linear(t *testing.T) {
    input_map := generate_open_map(100)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    OutputTestResults("100x100 linear  solution: ", num_paths, elapsed)
}

func Test100x100Queue(t *testing.T) {
    input_map := generate_open_map(100)

    start := time.Now()
    num_paths := solve_field_with_channel(input_map)
    elapsed := time.Since(start)

    OutputTestResults("100x100 channel solution: ", num_paths, elapsed)
}

func Test100x100QueueWithThreads(t *testing.T) {
    input_map := generate_open_map(100)

    start := time.Now()
    num_paths := solve_field_with_channel_and_threads(input_map)
    elapsed := time.Since(start)

    OutputTestResults("100x100 channel & thread solution: ", num_paths, elapsed)
}

func Test100x100WithWallLinear(t *testing.T) {
    input_map := generate_open_map(100)
    make_map_2(input_map)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    OutputTestResults("100x100 linear w/ wall solution: ", num_paths, elapsed)
}

func Test100x100WithWallQueue(t *testing.T) {
    input_map := generate_open_map(100)
    make_map_2(input_map)

    start := time.Now()
    num_paths := solve_field_with_channel(input_map)
    elapsed := time.Since(start)

    OutputTestResults("100x100 channel w/ wall solution: ", num_paths, elapsed)
}

func Test100x100WithWallQueueAndThreads(t *testing.T) {
    input_map := generate_open_map(100)
    make_map_2(input_map)

    start := time.Now()
    num_paths := solve_field_with_channel_and_threads(input_map)
    elapsed := time.Since(start)

    OutputTestResults("100x100 channel & thread w/ wall solution: ", num_paths, elapsed)
}
