#include <bits/stdc++.h>
using namespace std;

/*
    Frequent insertions/deletions at arbitrary positions are needed.
    The block size is large (hundreds or thousands of elements).
    You need to maintain sorted order and perform range queries efficiently.
    Insertions and lookups are both important and you want consistent O(log n) performance.
*/

/*
    Block size is small (typically ≤ 100–200 elements).
    Lookups are much more frequent than insertions/deletions.
    You want lower memory usage and faster lookups due to better cache locality.
    Insertions are rare or mostly at the end, or you can tolerate O(n) insertion time.
*/

/*
    Why Arrays/Vectors Are Used in B-Trees (Even with Insertions/Deletions)
    
    Block Size Is Large:
    B-tree nodes are designed to match disk block sizes (e.g., 4KB), so each node holds many keys.
    
    Insertions/Deletions Are Batc/-------hed:
    Insertions and deletions within a node (array/vector) are O(n) for that node, but since nodes are large, the tree is shallow and these operations are infrequent per node.
    
    Amortized Cost:
    Most insertions/deletions only affect a single node. Splits/merges are rare compared to lookups.
    
    Cache and Disk Efficiency:
    Arrays/vectors allow reading/writing an entire node in one disk I/O, which is much faster than following pointers in a tree of maps.
    
    Pointer Overhead:
    Using a map would add significant pointer and allocation overhead, reducing the number of keys per node and increasing tree height (bad for performance).

*/

class TwoLevelIndex {
private:
    map<int,map<int,string>>primaryIndex;
    int blockSize;

public:
    TwoLevelIndex(int blockSz):blockSize(blockSz){}

public:
    void insert(int key,string data){
        auto it=primaryIndex.upper_bound(key);
        if(it!=primaryIndex.begin()) --it;
        if(it==primaryIndex.end()||it->second.rbegin()->first<key){
            primaryIndex[key][key]=data;
        }
        else{
            it->second[key]=data;
            if(it->second.size()>blockSize){
                splitBlock(it);
            }
        }
        cout << "Insertion done...\n";
    }

public:
    string* search(int key) {
        cout<<"\nSearching for key: "<<key<<"\n";
        cout<<"Level 1: Searching primary indexing...\n";

        auto it=primaryIndex.upper_bound(key);
        if (it!=primaryIndex.begin()) --it;

        if (it==primaryIndex.end()||it->second.rbegin()->first<key) {
            cout<<"Block not found...\n";
            return nullptr;
        }

        cout<<"Level 2: Searching within block...\n";
        auto found=it->second.find(key);
        if (found!=it->second.end()) {
            cout<<"Record found: "<<found->second<<"\n";
            return &found->second;
        }

        cout<<"Key not found in block\n";
        return nullptr;
    }

public:
    void splitBlock(map<int,map<int,string>>::iterator it){
        auto& block=it->second;
        int mid=block.size()/2;

        auto midIt=next(block.begin(),mid);
        map<int,string>newBlock(midIt,block.end());
        block.erase(midIt, block.end());
        int newKey = newBlock.begin()->first;
        primaryIndex[newKey] = newBlock;
    }
};

int main(){
    TwoLevelIndex index(3);
    index.insert(10, "A");
    index.insert(20, "B");
    index.insert(5, "C");
    index.insert(15, "D");
    index.insert(25, "E");
    index.insert(30, "F");

    index.search(15);
    index.search(100);

    return 0;
}