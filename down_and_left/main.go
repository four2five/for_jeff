package main

import (
	"fmt"
  "time"

  "github.com/golang-collections/go-datastructures/queue"
)

////////////////////
// This program generates and then solves a map that may contain obstacles.
// There is currently minimal error checking, so beware.
///////////////////

const(
  DEBUG bool = false
  ONE int64 = 1
  OPEN int = 0
  BLOCKED int = -1
  ENQUEUED int = -2
)

type position struct{
  row int
  column int
}

// Assumes at least an 50x50 input. Makes a "wall" in column 48, preventing paths to the left
// of it except in row 9
func make_map_1(input_map [][]int) {
  for i:=0; i<48; i++ {
    input_map[i][48] = BLOCKED
  }
}

// Assumes at least an 40x40 input, makes a "wall" in column 38
// preventing work to the left of it other than the bottom row
func make_map_2(input_map [][]int) {
  for i:=0; i<98; i++ {
    input_map[i][98] = BLOCKED
  }
}

// generate_open_map creates a square 2-d slice of the specified
// length, writing zeros into each entry
func generate_open_map(side_length int) [][]int {
  map_data := make([][]int,side_length )

  for y:=0; y<side_length; y++ {
    row_data := make([]int, side_length)
    for x:=0; x<side_length; x++ {
      row_data[x] = OPEN
    }
    map_data[y] = row_data
  }

  return map_data
}

// map_is_solvable currently checks that neither the starting point nor ending
// point are blocked and that the map is at least 1x1
func map_is_solvable(map_data [][]int) bool {
  map_width := len(map_data)

  // If the origin or destination are blocked, no point is attempting
  // to solve the map
  if map_width < 1 || map_data[0][map_width-1] == BLOCKED || map_data[map_width-1][0] == BLOCKED {
    return false
  }

  return true
}

func position_is_blocked(y, x int, map_data [][]int) bool {
  return map_data[y][x] == BLOCKED
}

// If the data in the position immediately above the current position is
// valid, then add its number of valid path to this position's
func contribution_from_above_neighbor_position(cur_position *position, map_data [][]int) int {
      if cur_position.row > 0 && map_data[cur_position.row-1][cur_position.column] != BLOCKED {
        return map_data[cur_position.row-1][cur_position.column]
      } else {
        return 0
      }
}

// If the data in the position immediately above the current position is
// valid, then add its number of valid path to this position's
func contribution_from_above_neighbor(row, column int, map_data [][]int) int {
      if row > 0 && map_data[row-1][column] != BLOCKED {
        return map_data[row-1][column]
      } else {
        return 0
      }
}

func contribution_from_right_neighbor_position(cur_position *position, map_data [][]int) int {
  map_width := len(map_data)

  if cur_position.column < map_width - 1 && map_data[cur_position.row][cur_position.column+1] != BLOCKED {
    return map_data[cur_position.row][cur_position.column+1]
  } else {
    return 0
  }
}

func contribution_from_right_neighbor(row, column int, map_data [][]int) int {
  map_width := len(map_data)

  if column < map_width - 1 && map_data[row][column+1] != BLOCKED {
    return map_data[row][column+1]
  } else {
    return 0
  }
}

func initialize_map_data(map_data [][]int) {
  map_width := len(map_data)
  map_data[0][map_width-1] = 1
}

func enqueue_position_if_valid(row, column int, map_data [][]int, my_queue *queue.RingBuffer){
  map_width := len(map_data)

  // Verify that row and column values are valid and that the position is open
  if row < map_width && column >= 0 && map_data[row][column] == OPEN {
    my_queue.Put(&position{row,column})
    map_data[row][column] = ENQUEUED
  }
}

func initialize_map_data_with_rb(map_data [][]int, my_queue *queue.RingBuffer) {
  map_width := len(map_data)
  map_data[0][map_width-1] = 1

  // For maps that are at least 2x2, add the squares to the left and
  // below of the starting point to the queue
  if map_width > 1 {
    enqueue_position_if_valid(0, map_width-2, map_data, my_queue)
    enqueue_position_if_valid(1, map_width-1, map_data, my_queue)
  }
}

func extract_final_result(map_data [][]int) int {
  map_width := len(map_data)
  return map_data[map_width-1][0]
}

func populate_map_with_rb(map_data [][]int, my_queue *queue.RingBuffer) {
  // for debug purposes only
  max_queue_size := my_queue.Len()
  loop_counter := 0

  for my_queue.Len() > 0 {
    loop_counter++
    temp, err := my_queue.Get()
    if nil != err {
      fmt.Errorf("Received error %+v while dequeueing", err)
      return
    }
    temp_position, _ := temp.(*position)
    if DEBUG {
      fmt.Printf("dequeue val: %+v queue len: %d\n", temp_position, my_queue.Len())
    }

    // We assume that only valid positions are enqueued, so we blindly use
    // this value
    map_data[temp_position.row][temp_position.column] = contribution_from_right_neighbor_position(temp_position, map_data) +
      contribution_from_above_neighbor_position(temp_position, map_data)

    // enqueue the fields to the left and down, if they are valid 
    enqueue_position_if_valid(temp_position.row, temp_position.column-1, map_data, my_queue)
    if DEBUG {
      fmt.Printf("pos %d, %d possibly enqueued %d, %d\n", temp_position.row, temp_position.column, temp_position.row, temp_position.column-1)
    }

    enqueue_position_if_valid(temp_position.row+1, temp_position.column, map_data, my_queue)
    if DEBUG {
    fmt.Printf("pos %d, %d possibly enqueued %d, %d\n", temp_position.row, temp_position.column, temp_position.row+1, temp_position.column)
    }

    /*
    if my_queue.Len() > max_queue_size {
      max_queue_size = my_queue.Len()
    }
    */
  }

  if DEBUG {
    fmt.Printf("!!!!! loop counter: %d\n", loop_counter)
    fmt.Printf("!!!!! max queue size: %d\n", max_queue_size)
  }
}

// populate_map will iterate over all of the fields in the map, populating each
// non-blocked field with the number of paths that can lead to that position
func linear_populate_map(map_data [][]int) {
  map_width := len(map_data)

  for row := 0; row < map_width; row++ {
    for column := map_width - 1; column >=0; column-- {
      if map_data[row][column] == BLOCKED {
        continue
      }

      map_data[row][column] += contribution_from_above_neighbor(row, column, map_data)
      map_data[row][column] += contribution_from_right_neighbor(row, column, map_data)
    }
  }
}

// linear_solve_field will use a simple nested for-loop to compute the number of paths
// from the top-right corner to the bottom-left corner
func linear_solve_field(map_data [][]int) int {
  if !map_is_solvable(map_data) {
    return 0
  }

  initialize_map_data(map_data)
  linear_populate_map(map_data)
  return extract_final_result(map_data)
}

// solve_field_with_queue uses a queue as the control mechanism when
// finding all of the valid paths from the top-right corner to the bottom-left corner.
func solve_field_with_queue(map_data [][]int) int {
  if !map_is_solvable(map_data) {
    return 0
  }

  my_rb := queue.NewRingBuffer(uint64(len(map_data)))
  initialize_map_data_with_rb(map_data, my_rb)
  populate_map_with_rb(map_data, my_rb)
  return extract_final_result(map_data)
}

// pretty prints the 2-d slice parameter
func print_map(map_data [][]int){
  map_width := len(map_data)

  for y:= 0; y<map_width; y++ {
     for x:= 0; x<map_width; x++ {
        fmt.Printf(" %3d ",map_data[y][x])
     }
     fmt.Println("\n")
  }
}

func main() {
  var linear_elapsed time.Duration
  var queue_elapsed time.Duration
  var linear_num_paths int
  var queue_num_paths int

  {
    input_map := generate_open_map(100)
    //make_map_1(input_map)
    //make_map_2(input_map)

    fmt.Println("Input Map")
    print_map(input_map)

    start := time.Now()
    linear_num_paths = linear_solve_field(input_map)
    linear_elapsed = time.Since(start)
    fmt.Println("Populated Map")
    print_map(input_map)
  }

  {
    input_map := generate_open_map(100)
    //make_map_1(input_map)
    //make_map_2(input_map)

    fmt.Println("Input Map")
    print_map(input_map)

    start := time.Now()
    queue_num_paths = solve_field_with_queue(input_map)
    queue_elapsed = time.Since(start)
    fmt.Println("Populated Map")
    print_map(input_map)
  }

  fmt.Println("Number of paths: ", linear_num_paths, " linear solve elapsed time: ", linear_elapsed)
  fmt.Println("Number of paths: ", queue_num_paths, " queue solve elapsed time: ", queue_elapsed)
}
