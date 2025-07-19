#include <bits/stdc++.h>
using namespace std;

template<typename K,typename V>
struct Entry {
    K key;
    V val;
    Entry* next;

    Entry(K k,V v): key(k),val(v),next(nullptr){}
};

template<typename K,typename V>
class HashMap {
private:
    static const int INITIAL_SIZE=16;
    static const int FINAL_SIZE=(1<<30);
    static constexpr double LOAD_FACTOR_THRESHOLD=0.75;

    vector<Entry<K,V>*>hashTable;
    int size_count;

private:
    int tableSizeFor(int cap){
        int n=cap-1;
        n|=(n>>1);
        n|=(n>>2);
        n|=(n>>4);
        n|=(n>>8);
        n|=(n>>16);
        return (n<0)?1:(n>=FINAL_SIZE)?FINAL_SIZE:n+1;
    }

private:
    int getHashCode(K key){
        return hash<K>{}(key)%(hashTable.size());
    }

private:
    void rehashedIfNeeded(){
        double currentLoadFactor=static_cast<double>(size_count)/hashTable.size();
        if(currentLoadFactor>LOAD_FACTOR_THRESHOLD) rehash(hashTable.size()*2);
    }

private:
    void rehash(int newSize){
        vector<Entry<K,V>*>oldTable=move(hashTable);
        hashTable=vector<Entry<K,V>*>(newSize,nullptr);
        size_count=0;

        for(Entry<K,V>* head:oldTable){
            while(head!=nullptr){
                Entry<K,V>* next=head->next;
                head->next=nullptr;
                putNode(head->key,head->val);
                delete head;
                head=next;
            }
        }
    }

private:
    void putNode(const K &key,const V &val){
        int hash_generated=getHashCode(key);
        Entry<K,V>* node=hashTable[hash_generated];
        if(node==nullptr){
            hashTable[hash_generated]=new Entry<K,V>(key,val);
            size_count++;
        }
        else{
            Entry<K,V>* prevNode=nullptr;
            while(node!=nullptr){
                if(node->key==key){
                    node->val=val;
                    return;
                }
                prevNode=node;
                node=node->next;
            }
            prevNode->next = new Entry<K,V>(key, val);
            size_count++;
        }
    }

public:
    HashMap():hashTable(INITIAL_SIZE,nullptr),size_count(0){}
    HashMap(int capacity):size_count(0){
        int tableSize=tableSizeFor(capacity);
        hashTable.resize(tableSize,nullptr);
    }

public:
    void put(const K &key,const V &val){
        rehashedIfNeeded();
        putNode(key,val);
    }

public:
    V get(const K &key){
        int hash_generated=getHashCode(key);
        Entry<K,V>* node=hashTable[hash_generated];
        while(node!=nullptr){
            if(node->key==key) return node->val;
            node=node->next;
        }
        return V(); 
    }

public:
    bool containsKey(const K &key){
        int hash_generated=getHashCode(key);
        Entry<K,V>* node=hashTable[hash_generated];
        while(node!=nullptr){
            if(node->key==key) return true;
            node=node->next;
        }
        return false; 
    }

public:
    bool remove(const K &key){
        int hash_generated=getHashCode(key);
        Entry<K,V>* node=hashTable[hash_generated];
        Entry<K,V>* prevNode=nullptr;
        while(node!=nullptr){
            if(node->key==key){
                if(prevNode==nullptr){
                    hashTable[hash_generated]=node->next;
                }else{
                    prevNode->next=node->next;
                }
                delete node;
                size_count--;
                return true;
            }
            prevNode=node;
            node=node->next;
        }
        return false;
    }

public:
    int size() const {
        return size_count;
    }

public:
    bool isEmpty() const {
        return size_count==0;
    }

public:
    double getLoadFactor() const {
        return static_cast<double>(size_count)/hashTable.size();
    }

public:
    int getBucketCount() const {
        return hashTable.size();
    }

public:
    vector<K>getKeys(){
        vector<K>keys;
        for(Entry<K,V>* head:hashTable){
            Entry<K,V>* node=head;
            while(node!=nullptr){
                keys.push_back(node->key);
                node=node->next;
            }
        }
        return keys;
    }

public:
    vector<V> getValues(){
        vector<V>values;
        for(Entry<K,V>* head:hashTable){
            Entry<K,V>* node=head;
            while(node!=nullptr){
                values.push_back(node->val);
                node=node->next;
            }
        }
        return values;
    }

public:
    void clear(){
        for(Entry<K,V>* head:hashTable){
            while(head!=nullptr){
                Entry<K,V>* temp=head;
                head=head->next;
                delete temp;
            }
        }
        fill(hashTable.begin(),hashTable.end(),nullptr);
        size_count=0;
    }

public:
    void display(){
        cout<<"HashMap Structure:\n";
        cout<<"Size: "<<size_count<<", Buckets: "<<hashTable.size()<<", Load Factor: "<<getLoadFactor()<<"\n";
        for(int i=0;i<hashTable.size();i++){
            cout<<"Bucket["<<i<<"]: ";
            Entry<K,V>* node=hashTable[i];
            if(node==nullptr){
                cout<<"empty";
            }else{
                while(node!=nullptr){
                    cout<<"("<<node->key<<":"<<node->val<<")";
                    if(node->next!=nullptr) cout<<" -> ";
                    node=node->next;
                }
            }
            cout<<"\n";
        }
    }

public:
    void printStatistics(){
        cout<<"HashMap Statistics:\n";
        cout<<"Total elements: "<<size_count<<"\n";
        cout<<"Total buckets: "<<hashTable.size()<<"\n";
        cout<<"Load factor: "<<getLoadFactor()<<"\n";
        int nonEmptyBuckets=0;
        int maxChainLength=0;
        int totalChainLength=0;
        for(Entry<K,V>* head:hashTable){
            if(head!=nullptr){
                nonEmptyBuckets++;
                int chainLength=0;
                Entry<K,V>* node=head;
                while(node!=nullptr){
                    chainLength++;
                    node=node->next;
                }
                maxChainLength=max(maxChainLength,chainLength);
                totalChainLength+=chainLength;
            }
        }
        cout<<"Non-empty buckets: "<<nonEmptyBuckets<<"\n";
        cout<<"Average chain length: "<<(nonEmptyBuckets?static_cast<double>(totalChainLength)/nonEmptyBuckets:0)<<"\n";
        cout<<"Maximum chain length: "<<maxChainLength<<"\n";
        cout<<"Bucket utilization: "<<(static_cast<double>(nonEmptyBuckets)/hashTable.size()*100)<<"%\n";
    }

public:
    V& operator[](const K &key){
        int hash_generated=getHashCode(key);
        Entry<K,V>* node=hashTable[hash_generated];
        while(node!=nullptr){
            if(node->key==key) return node->val;
            node=node->next;
        }
        rehashedIfNeeded();
        hash_generated=getHashCode(key);
        Entry<K,V>* newEntry=new Entry<K,V>(key,V{});
        newEntry->next=hashTable[hash_generated];
        hashTable[hash_generated]=newEntry;
        size_count++;
        return newEntry->val;
    }

public:
    ~HashMap(){
        clear();
    }
};

int main(){
    HashMap<int,string>map;
    map.put(1,"hello");
    map.put(2,"world");
    map.put(17,"collision");
    cout<<map.get(1)<<endl;
    cout<<map.get(2)<<endl;
    map[3]="three";
    cout<<map[3]<<endl;
    map.display();
    map.printStatistics();
    cout<<"Size: "<<map.size()<<endl;
    cout<<"Contains 2: "<<map.containsKey(2)<<endl;
    map.remove(2);
    cout<<"After removing 2, Contains 2: "<<map.containsKey(2)<<endl;
    return 0;
}