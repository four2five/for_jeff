package main

import(
	"fmt"
  "testing"
  "time"
)

func Test4x4Linear(t *testing.T) {
    input_map := generate_open_map(4)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    fmt.Println("4x4 linear test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test4x4Queue(t *testing.T) {
    input_map := generate_open_map(4)

    start := time.Now()
    num_paths := solve_field_with_queue(input_map)
    elapsed := time.Since(start)

    fmt.Println("4x4 queue test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test50x50Linear(t *testing.T) {
    input_map := generate_open_map(50)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    fmt.Println("50x50 linear test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test50x50Queue(t *testing.T) {
    input_map := generate_open_map(50)

    start := time.Now()
    num_paths := solve_field_with_queue(input_map)
    elapsed := time.Since(start)

    fmt.Println("50x50 queue test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test50x50WithWallLinear(t *testing.T) {
    input_map := generate_open_map(50)
    make_map_1(input_map)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    fmt.Println("50x50 with wall linear test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test50x50WithWallQueue(t *testing.T) {
    input_map := generate_open_map(50)
    make_map_1(input_map)

    start := time.Now()
    num_paths := solve_field_with_queue(input_map)
    elapsed := time.Since(start)

    fmt.Println("50x50 with wall queue test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test100x100Linear(t *testing.T) {
    input_map := generate_open_map(100)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    fmt.Println("100x100 linear test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test100x100Queue(t *testing.T) {
    input_map := generate_open_map(100)

    start := time.Now()
    num_paths := solve_field_with_queue(input_map)
    elapsed := time.Since(start)

    fmt.Println("100x100 queue test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test100x100WithWallLinear(t *testing.T) {
    input_map := generate_open_map(100)
    make_map_2(input_map)

    start := time.Now()
    num_paths := linear_solve_field(input_map)
    elapsed := time.Since(start)

    fmt.Println("100x100 with wall linear test number of paths: ", num_paths, " elapsed time: ", elapsed)
}

func Test100x100WithWallQueue(t *testing.T) {
    input_map := generate_open_map(100)
    make_map_2(input_map)

    start := time.Now()
    num_paths := solve_field_with_queue(input_map)
    elapsed := time.Since(start)

    fmt.Println("100x100 with wall queue test number of paths: ", num_paths, " elapsed time: ", elapsed)
}
