package main

import(
  "fmt"
)

/////////////////
//
// A solution that uses channels and sometimes go funcs to
// solve the down-and-left problem.
//
//////////////////


// initialize_map_data sets the number of valid paths for the starting-point
// to 1. Then enqueues the positions to the left and right of the
// starting-point.
func initialize_map_data_with_channel(map_data [][]int64, buf_channel chan position) {
  map_width := len(map_data)
  // valid solutions always start in the to-right corner
  map_data[0][map_width-1] = 1

  // For maps that are at least 2x2, add the squares to the left and
  // below of the starting point to the channel
  if map_width > 1 {
    enqueue_position_if_valid(0, map_width-2, map_data, buf_channel)
    enqueue_position_if_valid(1, map_width-1, map_data, buf_channel)
  }
}

// enqueue_position_if_valid does some basic bounds checking and value
// validation prior to enqueueing the specified position
func enqueue_position_if_valid(row, column int, map_data [][]int64, buf_channel chan position){
  map_width := len(map_data)

  // Verify that row and column values are valid and that the position is open
  if row < map_width && column >= 0 && map_data[row][column] == OPEN {
    map_data[row][column] = ENQUEUED
    buf_channel <- position{row,column}
  }
}

// populate_map_with_channel computes the valid paths to each position
// reachable from the starting-point.
func populate_map_with_channel(map_data [][]int64, buf_channel chan position, output_message chan bool) {

  for temp_position := range buf_channel {
    if DEBUG {
      fmt.Printf("processing %v\n", temp_position)
    }

    map_data[temp_position.row][temp_position.column] =
      contribution_from_right_neighbor(temp_position.row, temp_position.column, map_data) +
      contribution_from_above_neighbor(temp_position.row, temp_position.column, map_data)

    if map_data[temp_position.row][temp_position.column] < 1{
      fmt.Errorf("wrote a position value of %d to [%d,%d]", map_data[temp_position.row][temp_position.column],
        temp_position.row, temp_position.column)
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
    if temp_position.row == len(map_data)-1 && temp_position.column == 0 {
      close(buf_channel)
    }
  }
}

// populate_map_with_channel_and_threads computes the valid paths to each position
// reachable from the starting-point using potentially multiple go funcs and a shared channel.
func populate_map_with_channel_and_threads(map_data [][]int64, buf_channel chan position) {
  // Create a channel for each go func and then spawn them
  func_done_chan := make(chan bool, NUM_THREADS)
  for i := 0; i < NUM_THREADS; i++ {
    go populate_map_with_channel(map_data, buf_channel, func_done_chan)
  }

  // Block until all of the go funcs signal that they are done
  for i := 0; i < NUM_THREADS; i++ {
    <-func_done_chan
  }
}

// solve_field_with_queue uses a queue as the control mechanism when
// finding all of the valid paths from the top-right corner to the bottom-left corner.
func solve_field_with_channel(map_data [][]int64) int64 {
  if !map_is_solvable(map_data) {
    return 0
  }

  buf_channel := make(chan position, 2*len(map_data))

  initialize_map_data_with_channel(map_data, buf_channel)
  populate_map_with_channel(map_data, buf_channel, nil)

  return extract_final_result(map_data)
}

// solve_field_with_channel_and_threads uses a channel as the control mechanism when
// finding all of the valid paths from the top-right corner to the bottom-left corner.
func solve_field_with_channel_and_threads(map_data [][]int64) int64 {
  if !map_is_solvable(map_data) {
    return 0
  }

  buf_channel := make(chan position, 2*len(map_data))

  initialize_map_data_with_channel(map_data, buf_channel)
  populate_map_with_channel(map_data, buf_channel, nil)

  return extract_final_result(map_data)
}

