package main

import (
  "fmt"
  "runtime"
  "time"
)

////////////////////
// This program generates and then solves a map that may contain obstacles.
// There is currently minimal error checking, so beware.
///////////////////

const(
  DEBUG bool = false
  // DEBUG bool = true
  TRACE bool = false
  OPEN int64 = 0
  BLOCKED int64 = -1
  ENQUEUED int64 = -2
  PROCESSING int64 = -3
  NUM_THREADS int = 4
  INPUT_MAP_WIDTH int = 4
)

type position struct{
  row int
  column int
}

// make_map_1 creates a "wall" in column 48, preventing paths to the left
// of it except in the bottow row.
// Assumes input is at least 50x50
func make_map_1(input_map [][]int64) {
  for i:=0; i<48; i++ {
    input_map[i][48] = BLOCKED
  }
}

// make_map_2 makes a "wall" in column 98 preventing paths to the left 
// of it other than the bottom row.
// Assumes at least a 100x100 input
func make_map_2(input_map [][]int64) {
  for i:=0; i<98; i++ {
    input_map[i][98] = BLOCKED
  }
}

// generate_open_map creates a square 2-d slice of the specified
// length, writing zeros into each entry
func generate_open_map(side_length int) [][]int64 {
  map_data := make([][]int64,side_length )

  for row := 0; row < side_length; row++ {
    row_data := make([]int64, side_length)
    for column := 0; column < side_length; column++ {
      row_data[column] = OPEN
    }
    map_data[row] = row_data
  }

  return map_data
}

// map_is_solvable currently checks that neither the starting point nor ending
// point are blocked, and that the map is at least 1x1
func map_is_solvable(map_data [][]int64) bool {
  map_width := len(map_data)

  // If the origin or destination are blocked, no point is attempting
  // to solve the map
  if map_width < 1 || map_data[0][map_width-1] == BLOCKED || map_data[map_width-1][0] == BLOCKED {
    return false
  }

  return true
}

// contribution_from_above_neighbor determines whether the data in the position
// immediately above the current position is valid and, if so,
// add its number of valid path to this position's
func contribution_from_above_neighbor(row, column int, map_data *[][]int64) int64 {
  if row <= 0 {
    return 0
  }

  temp_value := (*map_data)[row-1][column]
  if (*map_data)[row-1][column] == BLOCKED {
    return 0
  }

  // need to wait if the above neighbor is still in flight
  for ENQUEUED == temp_value || PROCESSING == temp_value {
    // this should only happen if we have enough threads in-flight that
    // the other node is actively being worked on. Use a really short sleep here.
    if DEBUG{
      if temp_value == ENQUEUED {
        fmt.Printf("[%d,%d] is enqueued\n", row-1,column)
      }
      if temp_value == PROCESSING {
        fmt.Printf("[%d,%d] is processing\n", row-1,column)
      }
    }
    time.Sleep(1 * time.Nanosecond)
  }

  return (*map_data)[row-1][column]
}

// contribution_from_right_neighbor determines whether the data in the 
// position immediately to the right of the current position is valid 
// and, if so, add its number of valid path to this position's
func contribution_from_right_neighbor(row, column int, map_data *[][]int64) int64 {
  map_width := len(*map_data)

  // if we are in the right-most column, there is no valid right neighbor
  if column >= map_width - 1 {
    return 0
  }

  temp_value := (*map_data)[row][column+1]

  // if the right-neighbor is blocked, we cannot use its value
  if temp_value == BLOCKED {
    return 0
  }

  for temp_value == ENQUEUED || temp_value == PROCESSING {
    if DEBUG {
      if temp_value == ENQUEUED {
        fmt.Printf("Waiting on [%d,%d] which is ENQUEUED", row, column+1)
      }
      if temp_value == PROCESSING {
        fmt.Printf("Waiting on [%d,%d] which is PROCESSING ", row, column+1)
      }
    }
    time.Sleep(1 * time.Nanosecond)
    temp_value = (*map_data)[row][column+1]
  }

  return temp_value
}

// initialize_map_data sets the number of valid paths for the starting-point
// to 1.
func initialize_map_data(map_data [][]int64) {
  map_width := len(map_data)
  map_data[0][map_width-1] = 1
}

// initialize_map_data sets the number of valid paths for the starting-point
// to 1. Then enqueues the positions to the left and right of the
// starting-point.
func initialize_map_data_with_channel(map_data *[][]int64, buf_channel chan position) {
  map_width := len(*map_data)
  (*map_data)[0][map_width-1] = 1

  // For maps that are at least 2x2, add the squares to the left and
  // below of the starting point to the channel
  if map_width > 1 {
    enqueue_position_if_valid(0, map_width-2, map_data, buf_channel)
    enqueue_position_if_valid(1, map_width-1, map_data, buf_channel)
  }
}

// enqueue_position_if_valid does some basic bounds checking and value
// validation prior to enqueueing the specified position
func enqueue_position_if_valid(row, column int, map_data *[][]int64, buf_channel chan position){
  map_width := len(*map_data)

  // Verify that row and column values are valid and that the position is open
  if row < map_width && column >= 0 && (*map_data)[row][column] == OPEN {
    (*map_data)[row][column] = ENQUEUED
    buf_channel <- position{row,column}
  }
}

// extract_final_result returns the value at the end-point.
// Assumes that map_data is at least 1x1
func extract_final_result(map_data [][]int64) int64 {
  map_width := len(map_data)
  return map_data[map_width-1][0]
}

// populate_map_with_rb computes the valid paths to each position
// reachable from the starting-point.
func populate_map_with_rb_and_threads(map_data *[][]int64, buf_channel chan position) {
  func_done_chan := make(chan bool, NUM_THREADS)
  for i := 0; i < NUM_THREADS; i++ {
    go populate_map_with_rb(map_data, buf_channel, func_done_chan)
  }

  for i := 0; i < NUM_THREADS; i++ {
    <-func_done_chan
  }
}

// populate_map_with_rb computes the valid paths to each position
// reachable from the starting-point.
func populate_map_with_rb(map_data *[][]int64, buf_channel chan position, output_message chan bool) {

  for temp_position := range buf_channel {
    if DEBUG {
      fmt.Printf("processing %v\n", temp_position)
    }

    pos_value := int64(0)
    // We assume that only valid positions are enqueued, so we blindly use
    // this value
    pos_value += contribution_from_right_neighbor(temp_position.row, temp_position.column, map_data)
    pos_value += contribution_from_above_neighbor(temp_position.row, temp_position.column, map_data)
    (*map_data)[temp_position.row][temp_position.column] = pos_value
    if pos_value < 1 && DEBUG {
      fmt.Errorf("wrote a position value of %d to [%d,%d]", pos_value, temp_position.row, temp_position.column)
    }

    // enqueue the fields to the left and down, if they are valid 
    enqueue_position_if_valid(temp_position.row, temp_position.column-1, map_data, buf_channel)
    if DEBUG && TRACE {
      fmt.Printf("pos %d, %d possibly enqueued %d, %d\n", temp_position.row, temp_position.column, temp_position.row, temp_position.column-1)
    }

    enqueue_position_if_valid(temp_position.row+1, temp_position.column, map_data, buf_channel)
    if DEBUG && TRACE {
      fmt.Printf("pos %d, %d possibly enqueued %d, %d\n", temp_position.row, temp_position.column, temp_position.row+1, temp_position.column)
    }

    // If this is the destination, we are done.
    if temp_position.row == len(*map_data)-1 && temp_position.column == 0 {
      close(buf_channel)
    }
  }
}

// linear_populate_map will iterate over all of the fields in the map, populating each
// non-blocked field with the number of paths that can lead to that position
func linear_populate_map(map_data [][]int64) {
  map_width := len(map_data)

  for row := 0; row < map_width; row++ {
    for column := map_width - 1; column >=0; column-- {
      if map_data[row][column] == BLOCKED {
        continue
      }

      map_data[row][column] += contribution_from_above_neighbor(row, column, &map_data)
      map_data[row][column] += contribution_from_right_neighbor(row, column, &map_data)
    }
  }
}

// linear_solve_field will use a simple nested for-loop to compute the number of paths
// from the top-right corner to the bottom-left corner
func linear_solve_field(map_data [][]int64) int64 {
  if !map_is_solvable(map_data) {
    return 0
  }

  initialize_map_data(map_data)
  if DEBUG {
    print_map(map_data)
  }
  linear_populate_map(map_data)
  return extract_final_result(map_data)
}

// solve_field_with_queue uses a queue as the control mechanism when
// finding all of the valid paths from the top-right corner to the bottom-left corner.
func solve_field_with_channel(map_data [][]int64) int64 {
  if !map_is_solvable(map_data) {
    return 0
  }

  buf_channel := make(chan position, 2*len(map_data))

  initialize_map_data_with_channel(&map_data, buf_channel)
  populate_map_with_rb(&map_data, buf_channel, nil)

  return extract_final_result(map_data)
}

// solve_field_with_channel_and_threads uses a channel as the control mechanism when
// finding all of the valid paths from the top-right corner to the bottom-left corner.
func solve_field_with_channel_and_threads(map_data [][]int64) int64 {
  if !map_is_solvable(map_data) {
    return 0
  }

  buf_channel := make(chan position, 2*len(map_data))

  initialize_map_data_with_channel(&map_data, buf_channel)
  populate_map_with_rb(&map_data, buf_channel, nil)

  return extract_final_result(map_data)
}

// pretty prints the 2-d slice parameter
func print_map(map_data [][]int64){
  map_width := len(map_data)

  for row := 0; row < map_width; row++ {
     for column := 0; column < map_width; column++ {
        fmt.Printf(" %12d ",map_data[row][column])
     }
     fmt.Println("\n")
  }
}

func main() {
  var linear_elapsed time.Duration
  var channel_elapsed time.Duration
  var channel_and_threads_elapsed time.Duration
  var linear_num_paths int64
  var channel_num_paths int64
  var channel_and_threads_num_paths int64

  // Set gomaxprocs to 4
  runtime.GOMAXPROCS(NUM_THREADS)

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
  fmt.Printf("Number of paths: %d channel with %d threads solve elapsed time: %v\n", channel_and_threads_num_paths, NUM_THREADS, channel_and_threads_elapsed)
}
