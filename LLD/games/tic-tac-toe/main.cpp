#include "./Game/Game.h"
#include <iostream>

int main() {
    int n;
    std::cout << "Enter the size of the board (n x n): ";
    std::cin >> n;

    if(n < 3) {
        std::cout << "Enter the board having size more than 3 * 3" << std::endl;
        return 0;
    }

    std::string player1_name, player2_name;
    std::cout << "Enter Player 1 name: ";
    std::cin >> player1_name;
    std::cout << "Enter Player 2 name: ";
    std::cin >> player2_name;

    Game game(n, player1_name, player2_name);
    game.start();

    return 0;
}
