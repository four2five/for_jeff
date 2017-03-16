package main

import (
  "fmt"
  "time"
)

///////////////////
//
// Utility functions, consts and common struct definitions.
//
///////////////////

const(
  DEBUG bool = false
  // DEBUG bool = true
  TRACE bool = false
  OPEN int64 = 0
  BLOCKED int64 = -1
  ENQUEUED int64 = -2
  PROCESSING int64 = -3
  NUM_THREADS int = 2
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

// make_unsolvable_map_1 marks the starting point as blocked
func make_unsolvable_map_bad_start(input_map [][]int64) {
  map_width := len(input_map)
  input_map[0][map_width-1] = BLOCKED
}

// make_unsolvable_map_2 marks the end point as blocked
func make_unsolvable_map_bad_end(input_map [][]int64) {
  map_width := len(input_map)
  input_map[map_width-1][0] = BLOCKED
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
func contribution_from_above_neighbor(row, column int, map_data [][]int64) int64 {
  if row <= 0 {
    return 0
  }

  temp_value := map_data[row-1][column]
  if map_data[row-1][column] == BLOCKED {
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

  return map_data[row-1][column]
}

// contribution_from_right_neighbor determines whether the data in the 
// position immediately to the right of the current position is valid 
// and, if so, add its number of valid path to this position's
func contribution_from_right_neighbor(row, column int, map_data [][]int64) int64 {
  map_width := len(map_data)

  // if we are in the right-most column, there is no valid right neighbor
  if column >= map_width - 1 {
    return 0
  }

  temp_value := map_data[row][column+1]

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
    temp_value = map_data[row][column+1]
  }

  return temp_value
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

// extract_final_result returns the value at the end-point.
// Assumes that map_data is at least 1x1
func extract_final_result(map_data [][]int64) int64 {
  map_width := len(map_data)
  return map_data[map_width-1][0]
}


