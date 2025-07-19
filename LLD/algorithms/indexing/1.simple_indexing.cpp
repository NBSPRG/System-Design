#include <bits/stdc++.h>
using namespace std;

class SimpleIndex {
private:
    map<int,string>records;

public:
    void insert(int key,string data){
        records[key]=data;
    }

public:
    string* search(int key){
        cout<<"\nsearching for key:"<<key<<"\n";
        auto it=records.find(key);
        if(it==records.end()){
            cout<<"key: "<<key<<"not found!!!\n";
            return nullptr;
        }
        cout<<"key: "<<key<<" found, Data: "<<it->second<<"\n";
        return &(it->second);
    }

public:
    void display(){
        cout<<"Simple indexing: \n";
        for(const auto &it:records){
            cout<<"Key: "<<it.first<<", Data: "<<it.second<<"\n";
        }
    }
};

int main() {
    SimpleIndex index;
    index.insert(101,"Alice");
    index.insert(102,"Bob");
    index.insert(150, "Charlie");
    index.display();
    index.search(150);
    return 0;
}