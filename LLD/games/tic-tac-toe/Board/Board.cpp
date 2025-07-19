#include "Board.h"
#include <iostream>

Board::Board(int n) : size(n), move_count(0) {
    matrix.resize(n, std::vector<char>(n, '_'));
}

bool Board::is_valid_move(int row, int col) {
    return row >= 0 && row < size && col >= 0 && col < size && matrix[row][col] == '_';
}

void Board::fill_move(int row, int col, bool Player) {
    matrix[row][col] = Player ? '1' : '0';
    move_count++;
}

bool Board::is_win(bool Player, int row, int col) {
    char player_char = Player + '0';

    bool row_win = true, col_win = true, diag_win = true, anti_diag_win = true;
    
    // Check row and column
    for (int i = 0; i < size; i++) {
        if (matrix[row][i] != player_char) row_win = false;
        if (matrix[i][col] != player_char) col_win = false;
    }

    // Check main diagonal
    for (int i = 0; i < size; i++) {
        if (matrix[i][i] != player_char) diag_win = false;
    }

    // Check anti-diagonal
    for (int i = 0; i < size; i++) {
        if (matrix[i][size - 1 - i] != player_char) anti_diag_win = false;
    }

    return row_win || col_win || diag_win || anti_diag_win;
}

bool Board::is_draw() {
    return move_count == size * size;
}

void Board::print() {
    for (const auto& row : matrix) {
        for (char cell : row) {
            std::cout << cell << " ";
        }
        std::cout << std::endl;
    }
}
