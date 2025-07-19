#include "Player.h"

Player::Player(std::string n, bool s) : name(n), symbol(s) {}

std::string Player::get_name() { return name; }

bool Player::get_symbol() { return symbol; }
