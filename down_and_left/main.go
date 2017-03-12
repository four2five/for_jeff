package main

import (
	"fmt"
)

////////////////////
// This program generates and then solves a map that may contain obstacles.
// There is currently minimal error checking, so beware.
///////////////////


// Assumes at least an 4x4 input
func make_map_1(input_map [][]int) {
  input_map[0][3] = -1
  input_map[1][3] = -1
  input_map[2][3] = -1
  input_map[3][3] = -1
}

func generate_open_map(side_length int) [][]int {
  map_data := make([][]int,side_length )

  for y:=0; y<side_length; y++ {
    row_data := make([]int, side_length)
    for x:=0; x<side_length; x++ {
      row_data[x] = 0
    }
    map_data[y] = row_data
  }

  return map_data
}

func map_is_solvable(map_data [][]int) bool {
  map_width := len(map_data)

  // If the origin or destination are blocked, no point is attempting
  // to solve the map
  if map_width < 1 || map_data[0][map_width-1] == -1 || map_data[map_width-1][0] == -1 {
    return false
  }

  return true
}

func position_is_blocked(y, x int, map_data [][]int) bool {
      if map_data[y][x] == -1 {
        return true
      }

      return false
}

// If the data in the position immediately above the current position is
// valid, then add its number of valid path to this position's
func contribution_from_above_neighbor(y, x int, map_data [][]int) int {
      if y > 0 && map_data[y-1][x] != -1 {
        return map_data[y-1][x]
      } else {
        return 0
      }
}

func contribution_from_right_neighbor(y, x int, map_data [][]int) int {
  map_width := len(map_data)

  if x < map_width - 1 && map_data[y][x+1] != -1 {
    return map_data[y][x+1]
  } else {
    return 0
  }
}

func initialize_map_data(map_data [][]int) {
  map_width := len(map_data)
  map_data[0][map_width-1] = 1
}

func extract_final_result(map_data [][]int) int {
  map_width := len(map_data)
  return map_data[map_width-1][0]
}

func populate_map(map_data [][]int) {
  map_width := len(map_data)

  for y:=0; y < map_width; y++ {
    for x:= map_width - 1; x >=0; x-- {
      if position_is_blocked(y, x, map_data) {
        continue
      }

      map_data[y][x] += contribution_from_above_neighbor(y, x, map_data)
      map_data[y][x] += contribution_from_right_neighbor(y, x, map_data)
    }
  }
}

func solve_field(map_data [][]int) int {
  if !map_is_solvable(map_data) {
    return 0
  }

  initialize_map_data(map_data)
  populate_map(map_data)
  return extract_final_result(map_data)
}

func print_map(map_data [][]int){
  map_width := len(map_data)

  for y:= 0; y<map_width; y++ {
     for x:= 0; x<map_width; x++ {
        fmt.Printf(" %5d ",map_data[y][x])
     }
     fmt.Println("\n")
  }
}

func main() {
  // input_map := map1()
  // input_map := map2()
  // input_map := map3()
  input_map := generate_open_map(8)
  make_map_1(input_map)

  fmt.Println("Input Map")
  print_map(input_map)
  num_paths := solve_field(input_map)
  //fmt.Println("", num_paths)
  fmt.Println("Populated Map")
  print_map(input_map)
  fmt.Println("Number of paths: ", num_paths)

  /*
  for i:=1; i<31; i++{
    input_map := generate_open_map(i)
    num_paths := solve_field(input_map)
    fmt.Println("", num_paths)
  }
  */
}
