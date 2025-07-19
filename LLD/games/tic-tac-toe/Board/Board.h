#ifndef BOARD_H
#define BOARD_H

#include <vector>

class Board {
private:
    int size;
    std::vector<std::vector<char>> matrix;
    int move_count;

public:
    Board(int n);
    bool is_valid_move(int row, int col);
    void fill_move(int row, int col, bool Player);
    bool is_win(bool Player, int row, int col);
    bool is_draw();
    void print();
};

#endif
