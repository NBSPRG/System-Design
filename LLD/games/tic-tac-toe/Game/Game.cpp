#include "Game.h"
#include <iostream>

Game::Game(int n, std::string p1_name, std::string p2_name)
    : board(n), player1(p1_name, true), player2(p2_name, false), current_player(true) {}

void Game::start() {
    while (true) {
        board.print();
        Player current = current_player ? player1 : player2;
        std::cout << current.get_name() << "'s turn (" << (current.get_symbol() ? '1' : '0') << "): ";

        int row, col;
        while (true) {
            std::cin >> row >> col;
            if (board.is_valid_move(row, col)) {
                board.fill_move(row, col, current.get_symbol());
                break;
            } else {
                std::cout << "Not a valid move! Enter again: ";
            }
        }

        if (board.is_win(current.get_symbol(), row, col)) {
            board.print();
            std::cout << current.get_name() << " wins!" << std::endl;
            break;
        } else if (board.is_draw()) {
            board.print();
            std::cout << "The game is a draw!" << std::endl;
            break;
        }

        current_player = !current_player;
    }
}
