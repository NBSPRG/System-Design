#include <bits/stdc++.h>
using namespace std;

/*
    Neither a btree or b+tree, it is a hybrid approach for indexing.
*/

struct Record {
    int key;
    string data;

    Record() : key(0), data("") {}
    Record(int k, string d) : key(k), data(d) {}
};

struct IndexNode {
    int minKey, maxKey, nodeId, level;
    bool isLeaf;
    weak_ptr<IndexNode> parent; // O(1) parent access
    vector<Record> records;
    vector<shared_ptr<IndexNode>> children;

    IndexNode(int mini, int maxi, int id, int lvl, bool leaf, weak_ptr<IndexNode> par = {})
        : minKey(mini), maxKey(maxi), nodeId(id), level(lvl), isLeaf(leaf), parent(par) {}

    void updateBounds() {
        if (isLeaf && !records.empty()) {
            minKey = records.front().key;
            maxKey = records.back().key;
        }
        else if (!isLeaf && !children.empty()) {
            minKey = children.front()->minKey;
            maxKey = children.back()->maxKey;
        }
        else if (records.empty() && children.empty()) {
            minKey = INT_MAX;
            maxKey = INT_MIN;
        }
    }

    // Binary search for child node - O(log fanout)
    shared_ptr<IndexNode> findChild(int key) {
        if (isLeaf) return nullptr;
        
        int left = 0, right = children.size() - 1;
        while (left <= right) {
            int mid = left + (right - left) / 2;
            if (key >= children[mid]->minKey && key <= children[mid]->maxKey) {
                return children[mid];
            }
            else if (key < children[mid]->minKey) {
                right = mid - 1;
            }
            else {
                left = mid + 1;
            }
        }
        // If not found in range, return rightmost child for insertion
        return children.empty() ? nullptr : children.back();
    }
};

class OptimizedMultiLevelIndex {
private:
    int blockSize, fanout, nextNodeId;
    shared_ptr<IndexNode> root;
    int maxLevels;

    // O(log n) bounds update using parent pointers
    void updateBoundsUpward(shared_ptr<IndexNode> node) {
        while (node) {
            node->updateBounds();
            node = node->parent.lock();
        }
    }

public:
    OptimizedMultiLevelIndex(int blockSz, int maxFanout = 4)
        : blockSize(blockSz), fanout(maxFanout), nextNodeId(0), maxLevels(0) {
        root = make_shared<IndexNode>(INT_MIN, INT_MAX, nextNodeId++, 0, true);
    }

    // O(log n) insertion
    void insert(int key, string data) {
        auto leaf = findLeaf(root, key);
        
        // Binary search for insertion position - O(log blockSize)
        auto pos = lower_bound(leaf->records.begin(), leaf->records.end(), key,
            [](const Record& r, int k) { return r.key < k; });
            
        leaf->records.insert(pos, Record(key, data));
        updateBoundsUpward(leaf);
        
        if (leaf->records.size() > blockSize) {
            splitLeaf(leaf);
        }
        
        cout << "Insertion done...\n";
    }

    // O(log n) deletion
    bool deleteKey(int key) {
        cout << "\nDeleting key: " << key << "\n";
        auto leaf = findLeaf(root, key);
        
        // Binary search for key - O(log blockSize)
        auto it = lower_bound(leaf->records.begin(), leaf->records.end(), key,
            [](const Record& r, int k) { return r.key < k; });
        
        if (it == leaf->records.end() || it->key != key) {
            cout << "Key " << key << " not found for deletion!\n";
            return false;
        }
        
        leaf->records.erase(it);
        updateBoundsUpward(leaf);
        
        // Handle underflow
        if (leaf->records.size() < blockSize / 2 && leaf != root) {
            handleUnderflow(leaf);
        }
        
        cout << "Key " << key << " deleted successfully!\n";
        return true;
    }

    // O(log n) search
    Record* search(int key) {
        cout << "\nSearching for key: " << key << "\n";
        auto leaf = findLeaf(root, key);
        
        // Binary search in leaf - O(log blockSize)
        auto it = lower_bound(leaf->records.begin(), leaf->records.end(), key,
            [](const Record& r, int k) { return r.key < k; });
            
        if (it != leaf->records.end() && it->key == key) {
            cout << "Found key: " << key << ", Data: " << it->data << "\n";
            return &(*it);
        }
        
        cout << "Key not found!!!\n";
        return nullptr;
    }

    void display() {
        displayNode(root, 0);
    }

private:
    // O(log n) leaf finding using binary search in internal nodes
    shared_ptr<IndexNode> findLeaf(shared_ptr<IndexNode> node, int key) {
        while (!node->isLeaf) {
            auto child = node->findChild(key);
            if (child) {
                node = child;
            } else {
                // Fallback to last child
                node = node->children.empty() ? node : node->children.back();
            }
        }
        return node;
    }

    // O(log n) underflow handling
    void handleUnderflow(shared_ptr<IndexNode> node) {
        auto parent = node->parent.lock();
        if (!parent) return; // Root case
        
        // Binary search for node position in parent - O(log fanout)
        int nodeIndex = findChildIndex(parent, node);
        if (nodeIndex == -1) return;
        
        shared_ptr<IndexNode> leftSibling = nullptr;
        shared_ptr<IndexNode> rightSibling = nullptr;
        
        if (nodeIndex > 0) leftSibling = parent->children[nodeIndex - 1];
        if (nodeIndex < parent->children.size() - 1) rightSibling = parent->children[nodeIndex + 1];
        
        // Try borrowing from siblings
        if (leftSibling && leftSibling->records.size() > blockSize / 2) {
            borrowFromLeft(node, leftSibling);
            updateBoundsUpward(parent);
            return;
        }
        
        if (rightSibling && rightSibling->records.size() > blockSize / 2) {
            borrowFromRight(node, rightSibling);
            updateBoundsUpward(parent);
            return;
        }
        
        // Merge with sibling
        if (leftSibling) {
            mergeNodes(leftSibling, node);
            parent->children.erase(parent->children.begin() + nodeIndex);
        }
        else if (rightSibling) {
            mergeNodes(node, rightSibling);
            parent->children.erase(parent->children.begin() + nodeIndex + 1);
        }
        
        updateBoundsUpward(parent);
        
        // Check parent underflow
        if (parent->children.size() < fanout / 2 && parent != root) {
            handleUnderflow(parent);
        }
        else if (parent == root && parent->children.size() == 1) {
            root = parent->children[0];
            root->parent.reset();
            maxLevels--;
        }
    }

    // O(log fanout) child index finding
    int findChildIndex(shared_ptr<IndexNode> parent, shared_ptr<IndexNode> child) {
        for (int i = 0; i < parent->children.size(); i++) {
            if (parent->children[i] == child) {
                return i;
            }
        }
        return -1;
    }

    void borrowFromLeft(shared_ptr<IndexNode> node, shared_ptr<IndexNode> leftSibling) {
        if (node->isLeaf) {
            auto borrowedRecord = leftSibling->records.back();
            leftSibling->records.pop_back();
            node->records.insert(node->records.begin(), borrowedRecord);
        }
        else {
            auto borrowedChild = leftSibling->children.back();
            leftSibling->children.pop_back();
            node->children.insert(node->children.begin(), borrowedChild);
            borrowedChild->parent = node;
        }
        leftSibling->updateBounds();
        node->updateBounds();
    }

    void borrowFromRight(shared_ptr<IndexNode> node, shared_ptr<IndexNode> rightSibling) {
        if (node->isLeaf) {
            auto borrowedRecord = rightSibling->records.front();
            rightSibling->records.erase(rightSibling->records.begin());
            node->records.push_back(borrowedRecord);
        }
        else {
            auto borrowedChild = rightSibling->children.front();
            rightSibling->children.erase(rightSibling->children.begin());
            node->children.push_back(borrowedChild);
            borrowedChild->parent = node;
        }
        rightSibling->updateBounds();
        node->updateBounds();
    }

    void mergeNodes(shared_ptr<IndexNode> left, shared_ptr<IndexNode> right) {
        if (left->isLeaf) {
            left->records.insert(left->records.end(), right->records.begin(), right->records.end());
        }
        else {
            for (auto& child : right->children) {
                child->parent = left;
            }
            left->children.insert(left->children.end(), right->children.begin(), right->children.end());
        }
        left->updateBounds();
    }

    // O(log n) leaf splitting
    void splitLeaf(shared_ptr<IndexNode> leaf) {
        int mid = leaf->records.size() / 2;
        auto newLeaf = make_shared<IndexNode>(
            leaf->records[mid].key, 
            leaf->records.back().key, 
            nextNodeId++, 
            0, 
            true, 
            leaf->parent
        );
        
        newLeaf->records.assign(leaf->records.begin() + mid, leaf->records.end());
        leaf->records.resize(mid);
        leaf->updateBounds();
        newLeaf->updateBounds();
        
        insertIntoParent(leaf, newLeaf);
    }

    // O(log n) parent insertion with ordered insertion
    void insertIntoParent(shared_ptr<IndexNode> left, shared_ptr<IndexNode> right) {
        if (left == root) {
            auto newRoot = make_shared<IndexNode>(left->minKey, right->maxKey, nextNodeId++, left->level + 1, false);
            newRoot->children.push_back(left);
            newRoot->children.push_back(right);
            left->parent = newRoot;
            right->parent = newRoot;
            newRoot->updateBounds();
            root = newRoot;
            maxLevels = root->level;
            return;
        }
        
        auto parent = left->parent.lock();
        
        // Insert at correct position - O(fanout) instead of O(fanout log fanout)
        auto pos = upper_bound(parent->children.begin(), parent->children.end(), right,
            [](const shared_ptr<IndexNode>& a, const shared_ptr<IndexNode>& b) {
                return a->minKey < b->minKey;
            });
        parent->children.insert(pos, right);
        right->parent = parent;
        
        updateBoundsUpward(parent);
        
        if (parent->children.size() > fanout) {
            splitInternal(parent);
        }
    }

    // O(log n) internal node splitting
    void splitInternal(shared_ptr<IndexNode> node) {
        int mid = node->children.size() / 2;
        auto newNode = make_shared<IndexNode>(
            node->children[mid]->minKey, 
            node->children.back()->maxKey, 
            nextNodeId++, 
            node->level, 
            false, 
            node->parent
        );
        
        // Update parent pointers for moved children
        for (int i = mid; i < node->children.size(); i++) {
            node->children[i]->parent = newNode;
        }
        
        newNode->children.assign(node->children.begin() + mid, node->children.end());
        node->children.resize(mid);
        node->updateBounds();
        newNode->updateBounds();
        
        insertIntoParent(node, newNode);
    }

    void displayNode(shared_ptr<IndexNode> node, int depth) {
        string indent(depth * 2, ' ');
        cout << indent << "Level " << node->level << " Node " << node->nodeId
             << " [" << node->minKey << ", " << node->maxKey << "]";
        if (node->isLeaf) {
            cout << " LEAF: ";
            for (const auto& r : node->records)
                cout << r.key << "(" << r.data << ") ";
        } else {
            cout << " INTERNAL: " << node->children.size() << " children";
        }
        cout << "\n";
        for (auto& child : node->children)
            displayNode(child, depth + 1);
    }
};

int main() {
    OptimizedMultiLevelIndex index(3, 3);
    
    cout << "=== Insertion Phase ===\n";
    index.insert(10, "A");
    index.insert(20, "B");
    index.insert(5, "C");
    index.insert(15, "D");
    index.insert(25, "E");
    index.insert(30, "F");
    index.insert(1, "G");
    index.insert(12, "H");
    index.insert(18, "I");
    index.insert(22, "J");
    index.insert(28, "K");

    cout << "\n=== Tree Structure After Insertions ===\n";
    index.display();

    cout << "\n=== Search Tests ===\n";
    auto rec = index.search(15);
    if (rec) cout << "Found: " << rec->key << " " << rec->data << endl;
    else cout << "Not found\n";

    rec = index.search(100);
    if (rec) cout << "Found: " << rec->key << " " << rec->data << endl;
    else cout << "Not found\n";

    cout << "\n=== Deletion Tests ===\n";
    index.deleteKey(15);
    index.deleteKey(25);
    index.deleteKey(1);
    index.deleteKey(100); // Non-existent key
    
    cout << "\n=== Tree Structure After Deletions ===\n";
    index.display();
    
    cout << "\n=== Search After Deletion ===\n";
    rec = index.search(15);
    if (rec) cout << "Found: " << rec->key << " " << rec->data << endl;
    else cout << "Not found\n";

    return 0;
}