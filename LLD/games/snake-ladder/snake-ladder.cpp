#include <iostream>
#include <unordered_map>

class SnakeAndLadder {
private:
    std::unordered_map<int, int> board;
    int players[2];
    int currentPlayer;

public:
    SnakeAndLadder() {
        createBoard();
        players[0] = 0;
        players[1] = 0;
        currentPlayer = 0;
    }

    void createBoard() {
        board[3] = 22;
        board[5] = 8;
        board[11] = 26;
        board[20] = 29;
        board[17] = 4;
        board[19] = 7;
        board[21] = 9;
        board[27] = 1;
    }

    void movePlayer() {
        int dieRoll;
        std::cout << "Player " << (currentPlayer + 1) << ", enter your die roll (1-6): ";
        std::cin >> dieRoll;

        if (dieRoll < 1 || dieRoll > 6) {
            std::cout << "Invalid input. Please enter a number between 1 and 6." << std::endl;
            return;
        }

        std::cout << "Player " << (currentPlayer + 1) << " rolled a " << dieRoll << std::endl;

        int newPosition = players[currentPlayer] + dieRoll;

        if (newPosition > 30) {
            std::cout << "Player " << (currentPlayer + 1) << " cannot move. Over the limit." << std::endl;
            return;
        }

        if (board.find(newPosition) != board.end()) {
            if (newPosition < board[newPosition]) {
                std::cout << "Player " << (currentPlayer + 1) << " hit a ladder!" << std::endl;
            } else {
                std::cout << "Player " << (currentPlayer + 1) << " hit a snake!" << std::endl;
            }
            newPosition = board[newPosition];
        }

        players[currentPlayer] = newPosition;
        std::cout << "Player " << (currentPlayer + 1) << " is now on square " << newPosition << std::endl;
    }

    bool checkWinner() {
        if (players[currentPlayer] == 30) {
            std::cout << "Player " << (currentPlayer + 1) << " wins!" << std::endl;
            return true;
        }
        return false;
    }

    void playGame() {
        while (true) {
            movePlayer();
            if (checkWinner()) {
                break;
            }
            currentPlayer = (currentPlayer + 1) % 2;
        }
    }
};

int main() {
    SnakeAndLadder game;
    game.playGame();
    return 0;
}
