#include <bits/stdc++.h>
using namespace std;

template<typename keyType, typename dataType>
struct Record {
    keyType key;
    dataType data;
    
    Record() : key(keyType{}), data(dataType{}) {}
    Record(keyType k, dataType d) : key(k), data(d) {}
};

template<typename keyType, typename dataType, int Order>
struct BTreeNode {
    int numKeys;
    int order;
    int nodeId;
    bool isLeaf;
    weak_ptr<BTreeNode<keyType, dataType, Order>> parent; // O(1) parent access
    vector<Record<keyType, dataType>> records;
    vector<shared_ptr<BTreeNode<keyType, dataType, Order>>> children;

    BTreeNode(int id, bool leaf = true) 
        : numKeys(0), order(Order), nodeId(id), isLeaf(leaf) {
        records.reserve(order);
        if (!leaf) children.reserve(order + 1);
    }
    
    bool isFull() {
        return numKeys == order - 1;  // CORRECTED: order-1 keys max
    }

    bool isUnderflow() {
        return numKeys < (order - 1) / 2;
    }

    void updateKeyCount() {
        numKeys = records.size();
    }

    // Binary search for position - O(log order)
    int findKeyPosition(keyType key) {
        return lower_bound(records.begin(), records.end(), key,
            [](const Record<keyType, dataType>& r, keyType k) {
                return r.key < k;
            }) - records.begin();
    }
};

template<typename keyType, typename dataType, int Order>
class BTree {
private:
    shared_ptr<BTreeNode<keyType, dataType, Order>> root;
    int nextNodeId;

    template<typename T>
    void printData(const T& data) const {
        cout << data;
    }

public:
    BTree() : nextNodeId(0) {
        root = make_shared<BTreeNode<keyType, dataType, Order>>(nextNodeId++, true);
    }

    // O(log n) insertion
    void insert(keyType key, dataType data) {
        auto leaf = findLeaf(root, key);
        
        // Binary search for position - O(log order)
        int pos = leaf->findKeyPosition(key);
        leaf->records.insert(leaf->records.begin() + pos, Record<keyType, dataType>(key, data));
        leaf->updateKeyCount();
        
        if (leaf->isFull()) {
            splitLeaf(leaf);
        }
        cout << "Insertion done...\n";
    }

    // O(log n) deletion
    bool deleteKey(keyType key) {
        cout << "\nDeleting key: " << key << "\n";
        return deleteFromNode(root, key);
    }

    // O(log n) search
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

    void display() {
        cout << "\n=== Optimized B-Tree Structure (Order " << Order << ") ===\n";
        displayNode(root, 0);
    }

private:
    shared_ptr<BTreeNode<keyType, dataType, Order>> findLeaf(
    shared_ptr<BTreeNode<keyType, dataType, Order>> node, keyType key) {
    
    while (!node->isLeaf) {
        // OPTIMIZED: Use binary search instead of linear scan
        int i = node->findKeyPosition(key);
        node = node->children[i];
    }
    return node;
}

    // O(log n) deletion
    bool deleteFromNode(shared_ptr<BTreeNode<keyType, dataType, Order>> node, keyType key) {
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
        
        // Internal node deletion logic
        int i = 0;
        while (i < node->records.size() && key > node->records[i].key) {
            i++;
        }
        
        auto child = node->children[i];
        bool deleted = deleteFromNode(child, key);
        
        if (deleted && child->isUnderflow() && child != root) {
            handleUnderflow(node, child, i);
        }
        
        return deleted;
    }

    // O(1) underflow handling with parent pointers
    void handleUnderflow(shared_ptr<BTreeNode<keyType, dataType, Order>> parent,
                        shared_ptr<BTreeNode<keyType, dataType, Order>> child, int childIndex) {
        
        // Try borrowing from left sibling
        if (childIndex > 0) {
            auto leftSibling = parent->children[childIndex - 1];
            if (leftSibling->numKeys > (Order - 1) / 2) {
                borrowFromLeft(parent, child, leftSibling, childIndex);
                return;
            }
        }
        
        // Try borrowing from right sibling
        if (childIndex < parent->children.size() - 1) {
            auto rightSibling = parent->children[childIndex + 1];
            if (rightSibling->numKeys > (Order - 1) / 2) {
                borrowFromRight(parent, child, rightSibling, childIndex);
                return;
            }
        }
        
        // Merge with sibling
        if (childIndex > 0) {
            mergeWithLeft(parent, child, parent->children[childIndex - 1], childIndex);
        } else {
            mergeWithRight(parent, child, parent->children[childIndex + 1], childIndex);
        }
        
        // Adjust root if necessary
        if (parent == root && parent->numKeys == 0 && !parent->children.empty()) {
            root = parent->children[0];
            root->parent.reset();
        }
    }

    void borrowFromLeft(shared_ptr<BTreeNode<keyType, dataType, Order>> parent,
                       shared_ptr<BTreeNode<keyType, dataType, Order>> child,
                       shared_ptr<BTreeNode<keyType, dataType, Order>> leftSibling,
                       int childIndex) {
        if (child->isLeaf) {
            auto borrowed = leftSibling->records.back();
            leftSibling->records.pop_back();
            child->records.insert(child->records.begin(), borrowed);
            
            // Update parent separator
            parent->records[childIndex - 1].key = child->records[0].key;
        } else {
            // Internal node borrowing logic
            auto separator = parent->records[childIndex - 1];
            child->records.insert(child->records.begin(), separator);
            
            parent->records[childIndex - 1] = leftSibling->records.back();
            leftSibling->records.pop_back();
            
            auto borrowedChild = leftSibling->children.back();
            leftSibling->children.pop_back();
            child->children.insert(child->children.begin(), borrowedChild);
            borrowedChild->parent = child;
        }
        
        leftSibling->updateKeyCount();
        child->updateKeyCount();
    }

    void borrowFromRight(shared_ptr<BTreeNode<keyType, dataType, Order>> parent,
                        shared_ptr<BTreeNode<keyType, dataType, Order>> child,
                        shared_ptr<BTreeNode<keyType, dataType, Order>> rightSibling,
                        int childIndex) {
        if (child->isLeaf) {
            auto borrowed = rightSibling->records.front();
            rightSibling->records.erase(rightSibling->records.begin());
            child->records.push_back(borrowed);
            
            // Update parent separator
            parent->records[childIndex].key = rightSibling->records.empty() ? 
                borrowed.key : rightSibling->records[0].key;
        } else {
            // Internal node borrowing logic
            auto separator = parent->records[childIndex];
            child->records.push_back(separator);
            
            parent->records[childIndex] = rightSibling->records.front();
            rightSibling->records.erase(rightSibling->records.begin());
            
            auto borrowedChild = rightSibling->children.front();
            rightSibling->children.erase(rightSibling->children.begin());
            child->children.push_back(borrowedChild);
            borrowedChild->parent = child;
        }
        
        rightSibling->updateKeyCount();
        child->updateKeyCount();
    }

    void mergeWithLeft(shared_ptr<BTreeNode<keyType, dataType, Order>> parent,
                      shared_ptr<BTreeNode<keyType, dataType, Order>> child,
                      shared_ptr<BTreeNode<keyType, dataType, Order>> leftSibling,
                      int childIndex) {
        if (!child->isLeaf) {
            leftSibling->records.push_back(parent->records[childIndex - 1]);
        }
        
        leftSibling->records.insert(leftSibling->records.end(), 
                                   child->records.begin(), child->records.end());
        leftSibling->children.insert(leftSibling->children.end(),
                                    child->children.begin(), child->children.end());
        
        // Update parent pointers
        for (auto& grandchild : child->children) {
            grandchild->parent = leftSibling;
        }
        
        parent->records.erase(parent->records.begin() + childIndex - 1);
        parent->children.erase(parent->children.begin() + childIndex);
        
        leftSibling->updateKeyCount();
        parent->updateKeyCount();
    }

    void mergeWithRight(shared_ptr<BTreeNode<keyType, dataType, Order>> parent,
                       shared_ptr<BTreeNode<keyType, dataType, Order>> child,
                       shared_ptr<BTreeNode<keyType, dataType, Order>> rightSibling,
                       int childIndex) {
        if (!child->isLeaf) {
            child->records.push_back(parent->records[childIndex]);
        }
        
        child->records.insert(child->records.end(),
                             rightSibling->records.begin(), rightSibling->records.end());
        child->children.insert(child->children.end(),
                              rightSibling->children.begin(), rightSibling->children.end());
        
        // Update parent pointers
        for (auto& grandchild : rightSibling->children) {
            grandchild->parent = child;
        }
        
        parent->records.erase(parent->records.begin() + childIndex);
        parent->children.erase(parent->children.begin() + childIndex + 1);
        
        child->updateKeyCount();
        parent->updateKeyCount();
    }

    // CORRECTED: Proper leaf splitting
    void splitLeaf(shared_ptr<BTreeNode<keyType, dataType, Order>> leaf) {
        int mid = Order / 2;  // Proper midpoint
        auto newLeaf = make_shared<BTreeNode<keyType, dataType, Order>>(nextNodeId++, true);
        newLeaf->parent = leaf->parent;
        
        newLeaf->records.assign(leaf->records.begin() + mid, leaf->records.end());
        leaf->records.resize(mid);
        
        newLeaf->updateKeyCount();
        leaf->updateKeyCount();
        
        insertIntoParent(leaf, newLeaf, newLeaf->records[0]);
    }

    // O(1) parent insertion with parent pointers
    void insertIntoParent(shared_ptr<BTreeNode<keyType, dataType, Order>> left,
                         shared_ptr<BTreeNode<keyType, dataType, Order>> right,
                         Record<keyType, dataType> promotedRecord) {
        if (left == root) {
            auto newRoot = make_shared<BTreeNode<keyType, dataType, Order>>(nextNodeId++, false);
            newRoot->records.push_back(promotedRecord);
            newRoot->children.push_back(left);
            newRoot->children.push_back(right);
            left->parent = newRoot;
            right->parent = newRoot;
            newRoot->updateKeyCount();
            root = newRoot;
            return;
        }
        
        auto parent = left->parent.lock();
        int pos = parent->findKeyPosition(promotedRecord.key);
        
        parent->records.insert(parent->records.begin() + pos, promotedRecord);
        parent->children.insert(parent->children.begin() + pos + 1, right);
        right->parent = parent;
        parent->updateKeyCount();
        
        if (parent->isFull()) {
            splitInternal(parent);
        }
    }

    // CORRECTED: Proper internal node splitting
    void splitInternal(shared_ptr<BTreeNode<keyType, dataType, Order>> node) {
        int mid = (Order - 1) / 2;  // CORRECTED: Different calculation for internal
        auto newNode = make_shared<BTreeNode<keyType, dataType, Order>>(nextNodeId++, false);
        newNode->parent = node->parent;
        
        Record<keyType, dataType> promotedRecord = node->records[mid];
        
        newNode->records.assign(node->records.begin() + mid + 1, node->records.end());
        newNode->children.assign(node->children.begin() + mid + 1, node->children.end());
        
        // Update parent pointers
        for (auto& child : newNode->children) {
            child->parent = newNode;
        }
        
        node->records.resize(mid);
        node->children.resize(mid + 1);
        
        newNode->updateKeyCount();
        node->updateKeyCount();
        
        insertIntoParent(node, newNode, promotedRecord);
    }

    void displayNode(shared_ptr<BTreeNode<keyType, dataType, Order>> node, int depth) {
        string indent(depth * 2, ' ');
        cout << indent << "Node " << node->nodeId << " (" << (node->isLeaf ? "LEAF" : "INTERNAL") << "): ";
        for (const auto& r : node->records) {
            cout << r.key << "(";
            printData(r.data);
            cout << ") ";
        }
        cout << "\n";
        for (auto& child : node->children) {
            displayNode(child, depth + 1);
        }
    }
};

int main() {
    BTree<int, string, 4> tree;
    
    cout << "=== Testing Corrected B-Tree ===\n";
    tree.insert(10, "Ten");
    tree.insert(20, "Twenty");
    tree.insert(5, "Five");
    tree.insert(15, "Fifteen");
    tree.insert(25, "TwentyFive");
    tree.insert(30, "Thirty");
    tree.insert(7, "Seven");
    
    tree.display();
    
    tree.search(15);
    tree.deleteKey(15);
    tree.display();
    
    return 0;
}