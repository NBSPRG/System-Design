#include <bits/stdc++.h>
using namespace std;

/*
    A B+ Tree has a (doubly) linked list connecting all its leaf nodes to enable fast and efficient range queries and ordered traversal.

    Why a linked list?
    So you can quickly scan all records in order without repeatedly traversing the tree from the root.

    Why doubly linked?
    So you can traverse both forward and backward through the leaves, supporting both ascending and descending range scans.

    In summary:
    The (doubly) linked list at the leaf level makes B+ Trees ideal for fast, ordered, and range-based queries—crucial for databases and file systems.
*/
template<typename keyType, typename dataType>
struct Record {
    keyType key;
    dataType data;
    
    Record() : key(keyType{}), data(dataType{}) {}
    Record(keyType k, dataType d) : key(k), data(d) {}
};

template<typename keyType, typename dataType, int Order>
struct BPlusTreeNode {
    int numKeys;
    int order;
    int nodeId;
    bool isLeaf;
    weak_ptr<BPlusTreeNode<keyType, dataType, Order>> parent;
    
    // CHANGE 1: Separate storage for internal vs leaf nodes
    vector<keyType> keys;                                                    // NEW: Only keys for internal nodes
    vector<Record<keyType, dataType>> records;                             // SAME: Records only in leaf nodes
    vector<shared_ptr<BPlusTreeNode<keyType, dataType, Order>>> children;  // SAME: Children pointers
    
    // CHANGE 2: Leaf linking for range queries
    shared_ptr<BPlusTreeNode<keyType, dataType, Order>> nextLeaf;          // NEW: Link to next leaf
    shared_ptr<BPlusTreeNode<keyType, dataType, Order>> prevLeaf;          // NEW: Link to previous leaf

    BPlusTreeNode(int id, bool leaf = true) 
        : numKeys(0), order(Order), nodeId(id), isLeaf(leaf), nextLeaf(nullptr), prevLeaf(nullptr) {
        if (leaf) {
            records.reserve(order - 1);    // CHANGE 3: Leaf stores records
        } else {
            keys.reserve(order - 1);       // CHANGE 3: Internal stores only keys
            children.reserve(order);
        }
    }
    
    bool isFull() {
        return numKeys == order - 1;
    }

    bool isUnderflow() {
        return numKeys < (order - 1) / 2;
    }

    void updateKeyCount() {
        if (isLeaf) {
            numKeys = records.size();      // CHANGE 4: Count records in leaf
        } else {
            numKeys = keys.size();         // CHANGE 4: Count keys in internal
        }
    }

    // CHANGE 5: Different search logic for internal vs leaf
    int findKeyPosition(keyType key) {
        if (isLeaf) {
            return lower_bound(records.begin(), records.end(), key,
                [](const Record<keyType, dataType>& r, keyType k) {
                    return r.key < k;
                }) - records.begin();
        } else {
            return lower_bound(keys.begin(), keys.end(), key) - keys.begin();
        }
    }
};

template<typename keyType, typename dataType, int Order>
class BPlusTree {
private:
    shared_ptr<BPlusTreeNode<keyType, dataType, Order>> root;
    shared_ptr<BPlusTreeNode<keyType, dataType, Order>> firstLeaf;  // NEW: Track first leaf for range queries
    int nextNodeId;

    template<typename T>
    void printData(const T& data) const {
        cout << data;
    }

public:
    BPlusTree() : nextNodeId(0), firstLeaf(nullptr) {
        root = make_shared<BPlusTreeNode<keyType, dataType, Order>>(nextNodeId++, true);
        firstLeaf = root;  // NEW: Initialize first leaf
    }

    void insert(keyType key, dataType data) {
        auto leaf = findLeaf(root, key);
        
        int pos = leaf->findKeyPosition(key);
        leaf->records.insert(leaf->records.begin() + pos, Record<keyType, dataType>(key, data));
        leaf->updateKeyCount();
        
        if (leaf->isFull()) {
            splitLeaf(leaf);
        }
        cout << "Insertion done...\n";
    }

    bool deleteKey(keyType key) {
        cout << "\nDeleting key: " << key << "\n";
        return deleteFromNode(root, key);
    }

    Record<keyType, dataType>* search(keyType key) {
        cout << "\nSearching for key: " << key << "\n";
        auto leaf = findLeaf(root, key);
        
        int pos = leaf->findKeyPosition(key);
        if (pos < leaf->records.size() && leaf->records[pos].key == key) {
            cout << "Found key: " << key << ", Data: ";
            printData(leaf->records[pos].data);
            cout << "\n";
            return &leaf->records[pos];
        }
        
        cout << "Key not found!!!\n";
        return nullptr;
    }

    // NEW: Efficient range search using leaf links
    void rangeSearch(keyType startKey, keyType endKey) {
        cout << "\nRange search [" << startKey << ", " << endKey << "]: ";
        vector<Record<keyType, dataType>> results;
        
        // Find starting leaf
        auto leaf = findLeaf(root, startKey);
        
        // Traverse linked leaves
        while (leaf != nullptr) {
            for (const auto& record : leaf->records) {
                if (record.key >= startKey && record.key <= endKey) {
                    results.push_back(record);
                } else if (record.key > endKey) {
                    goto done;  // No more relevant records
                }
            }
            leaf = leaf->nextLeaf;  // NEW: Use leaf linking
        }
        
        done:
        if (results.empty()) {
            cout << "No records found in range\n";
        } else {
            for (const auto& r : results) {
                cout << r.key << "(";
                printData(r.data);
                cout << ") ";
            }
            cout << "\n";
        }
    }

    void display() {
        cout << "\n=== B+ Tree Structure (Order " << Order << ") ===\n";
        displayNode(root, 0);
        
        // NEW: Display leaf chain
        cout << "\n=== Leaf Chain ===\n";
        auto leaf = firstLeaf;
        while (leaf && leaf->isLeaf) {
            cout << "Leaf " << leaf->nodeId << ": ";
            for (const auto& r : leaf->records) {
                cout << r.key << " ";
            }
            cout << " -> ";
            leaf = leaf->nextLeaf;
        }
        cout << "NULL\n";
    }

private:
    // CHANGE 6: Modified for B+ tree navigation
    shared_ptr<BPlusTreeNode<keyType, dataType, Order>> findLeaf(
        shared_ptr<BPlusTreeNode<keyType, dataType, Order>> node, keyType key) {
        
        while (!node->isLeaf) {
            int i = 0;
            // CHANGE 7: Use keys array instead of records for internal nodes
            while (i < node->keys.size() && key >= node->keys[i]) {
                i++;
            }
            node = node->children[i];
        }
        return node;
    }

    bool deleteFromNode(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> node, keyType key) {
        if (node->isLeaf) {
            int pos = node->findKeyPosition(key);
            if (pos < node->records.size() && node->records[pos].key == key) {
                node->records.erase(node->records.begin() + pos);
                node->updateKeyCount();
                cout << "Key " << key << " deleted successfully!\n";
                return true;
            }
            cout << "Key " << key << " not found!\n";
            return false;
        }
        
        // CHANGE 8: Modified navigation for internal nodes
        int i = 0;
        while (i < node->keys.size() && key >= node->keys[i]) {
            i++;
        }
        
        auto child = node->children[i];
        bool deleted = deleteFromNode(child, key);
        
        if (deleted && child->isUnderflow() && child != root) {
            handleUnderflow(node, child, i);
        }
        
        return deleted;
    }

    void handleUnderflow(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> parent,
                        shared_ptr<BPlusTreeNode<keyType, dataType, Order>> child, int childIndex) {
        
        if (childIndex > 0) {
            auto leftSibling = parent->children[childIndex - 1];
            if (leftSibling->numKeys > (Order - 1) / 2) {
                borrowFromLeft(parent, child, leftSibling, childIndex);
                return;
            }
        }
        
        if (childIndex < parent->children.size() - 1) {
            auto rightSibling = parent->children[childIndex + 1];
            if (rightSibling->numKeys > (Order - 1) / 2) {
                borrowFromRight(parent, child, rightSibling, childIndex);
                return;
            }
        }
        
        if (childIndex > 0) {
            mergeWithLeft(parent, child, parent->children[childIndex - 1], childIndex);
        } else {
            mergeWithRight(parent, child, parent->children[childIndex + 1], childIndex);
        }
        
        if (parent == root && parent->numKeys == 0 && !parent->children.empty()) {
            root = parent->children[0];
            root->parent.reset();
        }
    }

    // CHANGE 9: Modified borrowing for B+ tree structure
    void borrowFromLeft(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> parent,
                       shared_ptr<BPlusTreeNode<keyType, dataType, Order>> child,
                       shared_ptr<BPlusTreeNode<keyType, dataType, Order>> leftSibling,
                       int childIndex) {
        if (child->isLeaf) {
            // CHANGE 10: In B+ tree, copy up the key (don't move)
            auto borrowed = leftSibling->records.back();
            leftSibling->records.pop_back();
            child->records.insert(child->records.begin(), borrowed);
            
            // Update parent separator to smallest key in child
            parent->keys[childIndex - 1] = child->records[0].key;
        } else {
            // For internal nodes, similar to B-tree but use keys array
            auto separatorKey = parent->keys[childIndex - 1];
            child->keys.insert(child->keys.begin(), separatorKey);
            
            parent->keys[childIndex - 1] = leftSibling->keys.back();
            leftSibling->keys.pop_back();
            
            auto borrowedChild = leftSibling->children.back();
            leftSibling->children.pop_back();
            child->children.insert(child->children.begin(), borrowedChild);
            borrowedChild->parent = child;
        }
        
        leftSibling->updateKeyCount();
        child->updateKeyCount();
    }

    void borrowFromRight(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> parent,
                        shared_ptr<BPlusTreeNode<keyType, dataType, Order>> child,
                        shared_ptr<BPlusTreeNode<keyType, dataType, Order>> rightSibling,
                        int childIndex) {
        if (child->isLeaf) {
            auto borrowed = rightSibling->records.front();
            rightSibling->records.erase(rightSibling->records.begin());
            child->records.push_back(borrowed);
            
            // Update parent separator
            parent->keys[childIndex] = rightSibling->records.empty() ? 
                borrowed.key : rightSibling->records[0].key;
        } else {
            auto separatorKey = parent->keys[childIndex];
            child->keys.push_back(separatorKey);
            
            parent->keys[childIndex] = rightSibling->keys.front();
            rightSibling->keys.erase(rightSibling->keys.begin());
            
            auto borrowedChild = rightSibling->children.front();
            rightSibling->children.erase(rightSibling->children.begin());
            child->children.push_back(borrowedChild);
            borrowedChild->parent = child;
        }
        
        rightSibling->updateKeyCount();
        child->updateKeyCount();
    }

    void mergeWithLeft(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> parent,
                      shared_ptr<BPlusTreeNode<keyType, dataType, Order>> child,
                      shared_ptr<BPlusTreeNode<keyType, dataType, Order>> leftSibling,
                      int childIndex) {
        if (child->isLeaf) {
            // CHANGE 11: Merge records and update leaf links
            leftSibling->records.insert(leftSibling->records.end(), 
                                       child->records.begin(), child->records.end());
            
            // Update leaf links
            leftSibling->nextLeaf = child->nextLeaf;
            if (child->nextLeaf) {
                child->nextLeaf->prevLeaf = leftSibling;
            }
        } else {
            // For internal nodes, include separator from parent
            leftSibling->keys.push_back(parent->keys[childIndex - 1]);
            leftSibling->keys.insert(leftSibling->keys.end(),
                                    child->keys.begin(), child->keys.end());
            leftSibling->children.insert(leftSibling->children.end(),
                                        child->children.begin(), child->children.end());
            
            for (auto& grandchild : child->children) {
                grandchild->parent = leftSibling;
            }
        }
        
        parent->keys.erase(parent->keys.begin() + childIndex - 1);
        parent->children.erase(parent->children.begin() + childIndex);
        
        leftSibling->updateKeyCount();
        parent->updateKeyCount();
    }

    void mergeWithRight(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> parent,
                       shared_ptr<BPlusTreeNode<keyType, dataType, Order>> child,
                       shared_ptr<BPlusTreeNode<keyType, dataType, Order>> rightSibling,
                       int childIndex) {
        if (child->isLeaf) {
            child->records.insert(child->records.end(),
                                 rightSibling->records.begin(), rightSibling->records.end());
            
            // Update leaf links
            child->nextLeaf = rightSibling->nextLeaf;
            if (rightSibling->nextLeaf) {
                rightSibling->nextLeaf->prevLeaf = child;
            }
        } else {
            child->keys.push_back(parent->keys[childIndex]);
            child->keys.insert(child->keys.end(),
                              rightSibling->keys.begin(), rightSibling->keys.end());
            child->children.insert(child->children.end(),
                                  rightSibling->children.begin(), rightSibling->children.end());
            
            for (auto& grandchild : rightSibling->children) {
                grandchild->parent = child;
            }
        }
        
        parent->keys.erase(parent->keys.begin() + childIndex);
        parent->children.erase(parent->children.begin() + childIndex + 1);
        
        child->updateKeyCount();
        parent->updateKeyCount();
    }

    // CHANGE 12: Modified splitting for B+ tree
    void splitLeaf(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> leaf) {
        int mid = Order / 2;
        auto newLeaf = make_shared<BPlusTreeNode<keyType, dataType, Order>>(nextNodeId++, true);
        newLeaf->parent = leaf->parent;
        
        newLeaf->records.assign(leaf->records.begin() + mid, leaf->records.end());
        leaf->records.resize(mid);
        
        // CHANGE 13: Update leaf links
        newLeaf->nextLeaf = leaf->nextLeaf;
        newLeaf->prevLeaf = leaf;
        if (leaf->nextLeaf) {
            leaf->nextLeaf->prevLeaf = newLeaf;
        }
        leaf->nextLeaf = newLeaf;
        
        newLeaf->updateKeyCount();
        leaf->updateKeyCount();
        
        // CHANGE 14: Copy up the first key of new leaf (don't move)
        keyType promotedKey = newLeaf->records[0].key;
        insertIntoParent(leaf, newLeaf, promotedKey);
    }

    // CHANGE 15: Modified to handle key promotion instead of record
    void insertIntoParent(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> left,
                         shared_ptr<BPlusTreeNode<keyType, dataType, Order>> right,
                         keyType promotedKey) {
        if (left == root) {
            auto newRoot = make_shared<BPlusTreeNode<keyType, dataType, Order>>(nextNodeId++, false);
            newRoot->keys.push_back(promotedKey);  // CHANGE 16: Store key, not record
            newRoot->children.push_back(left);
            newRoot->children.push_back(right);
            left->parent = newRoot;
            right->parent = newRoot;
            newRoot->updateKeyCount();
            root = newRoot;
            return;
        }
        
        auto parent = left->parent.lock();
        int pos = lower_bound(parent->keys.begin(), parent->keys.end(), promotedKey) - parent->keys.begin();
        
        parent->keys.insert(parent->keys.begin() + pos, promotedKey);
        parent->children.insert(parent->children.begin() + pos + 1, right);
        right->parent = parent;
        parent->updateKeyCount();
        
        if (parent->isFull()) {
            splitInternal(parent);
        }
    }

    void splitInternal(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> node) {
        int mid = (Order - 1) / 2;
        auto newNode = make_shared<BPlusTreeNode<keyType, dataType, Order>>(nextNodeId++, false);
        newNode->parent = node->parent;
        
        keyType promotedKey = node->keys[mid];  // CHANGE 17: Promote key, not record
        
        newNode->keys.assign(node->keys.begin() + mid + 1, node->keys.end());
        newNode->children.assign(node->children.begin() + mid + 1, node->children.end());
        
        for (auto& child : newNode->children) {
            child->parent = newNode;
        }
        
        node->keys.resize(mid);
        node->children.resize(mid + 1);
        
        newNode->updateKeyCount();
        node->updateKeyCount();
        
        insertIntoParent(node, newNode, promotedKey);
    }

    void displayNode(shared_ptr<BPlusTreeNode<keyType, dataType, Order>> node, int depth) {
        string indent(depth * 2, ' ');
        cout << indent << "Node " << node->nodeId << " (" << (node->isLeaf ? "LEAF" : "INTERNAL") << "): ";
        
        if (node->isLeaf) {
            // CHANGE 18: Display records in leaf
            for (const auto& r : node->records) {
                cout << r.key << "(";
                printData(r.data);
                cout << ") ";
            }
        } else {
            // CHANGE 18: Display keys in internal
            for (const auto& k : node->keys) {
                cout << k << " ";
            }
        }
        cout << "\n";
        
        for (auto& child : node->children) {
            displayNode(child, depth + 1);
        }
    }
};

int main() {
    BPlusTree<int, string, 4> tree;
    
    cout << "=== Testing B+ Tree ===\n";
    tree.insert(10, "Ten");
    tree.insert(20, "Twenty");
    tree.insert(5, "Five");
    tree.insert(15, "Fifteen");
    tree.insert(25, "TwentyFive");
    tree.insert(30, "Thirty");
    tree.insert(7, "Seven");
    tree.insert(12, "Twelve");
    tree.insert(18, "Eighteen");
    
    tree.display();
    
    cout << "\n=== Search Test ===\n";
    tree.search(15);
    
    cout << "\n=== Range Search Test ===\n";
    tree.rangeSearch(10, 25);
    
    cout << "\n=== Deletion Test ===\n";
    tree.deleteKey(15);
    tree.display();
    
    return 0;
}