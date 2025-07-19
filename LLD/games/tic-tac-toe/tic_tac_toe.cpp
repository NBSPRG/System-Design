#include <bits/stdc++.h>
using namespace std;

bool is_valid_move(int row,int col,int count,vector<vector<char>>&matrix){
    return (row>=0&&row<=2&&col>=0&&col<=2&&matrix[row][col]=='_'&&count<9);
}
void fill_move(int row,int col,vector<vector<char>>&matrix,bool Player){
    if(Player){
        matrix[row][col]='1';
    }
    else{
        matrix[row][col]='0';
    }

}
bool is_win(bool Player, vector<vector<char>>& matrix, int row, int col) {
    char player_char = Player + '0';

    if (matrix[row][0] == player_char && matrix[row][1] == player_char && matrix[row][2] == player_char) return true;
    else if (matrix[0][col] == player_char && matrix[1][col] == player_char && matrix[2][col] == player_char) return true;
    else if (row == col && matrix[0][0] == player_char && matrix[1][1] == player_char && matrix[2][2] == player_char) return true;
    else if (row + col == 2 && matrix[0][2] == player_char && matrix[1][1] == player_char && matrix[2][0] == player_char)return true;

    return false;
}

void print(vector<vector<char>>&matrix){
    for(int i=0;i<3;i++){
        for(int j=0;j<3;j++){
            cout<<matrix[i][j]<<" ";
        }
        cout<<"\n";
    }
}

void f(bool Player,vector<vector<char>>&board,int &count){
    while(true){
        if(Player){
            cout<<"\n";
            cout<<"Player 1 chance : ";
            int row,col;
            while(true){
                cin>>row>>col;
                cout<<endl;
                if(is_valid_move(row,col,count,board)){
                    fill_move(row,col,board,Player);
                    break;
                }
                else{
                    cout<<"Not a valid move !!!\n";
                    cout<<"Enter again: ";
                }
            }
            count++;
            if(is_win(Player,board,row,col)){
                cout << "Player 1 win\n";
                print(board);
                break;
            }
            else if(count==9){
                cout<<"Match draw\n";
                print(board);
                break;
            }
            else{
                print(board);
                Player=!Player;
            }
        }
        else {
            cout<<"\n";
            cout<<"Player 2 chance : ";
            int row,col;
            while(true){
                cin>>row>>col;
                cout<<endl;
                if(is_valid_move(row,col,count,board)){
                    fill_move(row,col,board,Player);
                    break;
                }
                else{
                    cout<<" Not a valid move !!! \n";
                    cout<<"Enter again: ";

                }
            }
            count++;
            if(is_win(Player,board,row,col)){
                cout << "Player 2 win\n";
                print(board);
                break;
            }
            else if(count==9){
                cout<<"Match draw\n";
                print(board);
                break;
            }
            else{
                print(board);
                Player=!Player;
            }
        }
    }
}


int main(){

    vector<vector<char>>board(3,vector<char>(3,'_'));
    for(auto it:board){
        for(auto it1: it){
            cout<<it1<<" ";
        }
        cout<<"\n";
    }
    cout<<"\n";
    bool player;
    cout<<"Enter the initial player (Either 0 or 1): ";
    cin>>player;
    cout<<"\n";
    int count=0;
    f(player,board,count);
       
    return 0;
}