package main

//////////////
//
// An iterative, single-threaded solution for the down-and-left
// problem.
//
//////////////


// initialize_map_data sets the number of valid paths for the starting-point
// to 1.
func initialize_map_data(map_data [][]int64) {
  map_width := len(map_data)
  map_data[0][map_width-1] = 1
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

