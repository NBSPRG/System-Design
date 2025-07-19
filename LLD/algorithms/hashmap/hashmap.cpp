#include <bits/stdc++.h>
using namespace std;

template <typename K, typename V>

class HashMap {
    private:
    static const int INITIAL_SIZE = 16;
    static const int FINAL_SIZE = (1 << 30);

    struct Entry {
        K key;
        V val;
        Entry* next;

        Entry(K k, V v): key(k), val(v), next(nullptr) {}
    };

    vector <Entry*>hashTable;

    //creating a funtion tableSizeFor for given capacity ----- VVI very optimized
    int tableSizeFor(int cap) {
        int n = cap - 1;
        n |= (n >> 1);
        n |= (n >> 2);
        n |= (n >> 4);
        n |= (n >> 8);
        n |= (n >> 16);

        return (n < 0) ? 1 : (n >= FINAL_SIZE) ? FINAL_SIZE : n + 1;
    }

    int getHashCode(K key){
        return hash<K>{}(key) % (hashTable.size());
    }

    public:
    HashMap() : hashTable(INITIAL_SIZE, nullptr) {}

    HashMap(int capacity) {
        int tableSize = tableSizeFor(capacity);
        hashTable.resize(tableSize, nullptr);
    }

    void put(const K &key, const V &val){
        //generate hash code
        int hash_generated = getHashCode(key);
        
        Entry* node = hashTable[hash_generated];

        if(node == nullptr) {
            hashTable[hash_generated] = new Entry(key, val);
        }

        else {
            Entry* prevNode = nullptr;
            while(node != nullptr) {
                if(node -> key == key) {
                    node -> val = val;
                    return ;
                }
                prevNode = node;
                node = node -> next;
            }

            prevNode -> next = new Entry(key, val);
        }
    }

    V get(const K &key){
        int hash_generated = getHashCode(key);
        Entry* node = hashTable[hash_generated];

        while(node != nullptr) {
            if(node -> key == key){
                return node -> val;
            }
            node = node -> next;
        }

        return V();
    }

    // destructor
    ~HashMap() {
        for(auto node: hashTable) {
            while(node != nullptr) {
                Entry* tmp = node;
                node = node -> next;
                delete(tmp);
            }
        }
    }   
};

int main() {
    HashMap<int, string> map;
    map.put(1, "hello");
    map.put(1, "hello_");
    map.put(2, "world");
    cout << map.get(1) << endl;
    cout << map.get(2) << endl;
    
    return 0;
}