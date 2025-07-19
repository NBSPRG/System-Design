#ifndef PLAYER_H
#define PLAYER_H

#include <string>

class Player {
private:
    std::string name;
    bool symbol;

public:
    Player(std::string n, bool s);
    std::string get_name();
    bool get_symbol();
};

#endif
