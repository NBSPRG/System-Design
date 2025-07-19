#ifndef GAME_H
#define GAME_H

#include "../Board/Board.h"
#include "../Player/Player.h"

class Game {
private:
    Board board;
    Player player1;
    Player player2;
    bool current_player;

public:
    Game(int n, std::string p1_name, std::string p2_name);
    void start();
};

#endif
