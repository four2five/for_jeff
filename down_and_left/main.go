package main

import (
  "fmt"
  "time"
)

////////////////////
// This program generates and then solves a map that may contain obstacles.
// There is currently minimal error checking, so beware.
///////////////////

func main() {
  var linear_elapsed time.Duration
  var channel_elapsed time.Duration
  var channel_and_threads_elapsed time.Duration
  var linear_num_paths int64
  var channel_num_paths int64
  var channel_and_threads_num_paths int64

  {
    linear_input_map := generate_open_map(INPUT_MAP_WIDTH)

    if DEBUG {
      fmt.Println("Input Map")
      print_map(linear_input_map)
    }

    start := time.Now()
    linear_num_paths = linear_solve_field(linear_input_map)
    linear_elapsed = time.Since(start)
    fmt.Println("Populated Map with linear approach")
    print_map(linear_input_map)
  }

  {
    input_map := generate_open_map(INPUT_MAP_WIDTH)

    if DEBUG {
      //fmt.Println("Input Map")
      //print_map(input_map)
    }

    start := time.Now()
    channel_num_paths = solve_field_with_channel(input_map)
    channel_elapsed = time.Since(start)
    fmt.Println("Populated Map with channel approach")
    print_map(input_map)
  }

  {
    input_map := generate_open_map(INPUT_MAP_WIDTH)

    if DEBUG {
      fmt.Println("Input Map")
      print_map(input_map)
    }

    start := time.Now()
    channel_and_threads_num_paths = solve_field_with_channel_and_threads(input_map)
    channel_and_threads_elapsed = time.Since(start)
    fmt.Println("Populated Map channel & threads approach")
    print_map(input_map)
  }

  fmt.Printf("Number of paths: %d linear solve elapsed time: %v\n", linear_num_paths, linear_elapsed)
  fmt.Printf("Number of paths: %d channel solve elapsed time: %v\n", channel_num_paths, channel_elapsed)
  fmt.Printf("Number of paths: %d channel with %d threads solve elapsed time: %v\n",
    channel_and_threads_num_paths, NUM_THREADS, channel_and_threads_elapsed)
}
