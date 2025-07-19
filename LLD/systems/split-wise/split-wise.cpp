#include <bits/stdc++.h>
using namespace std;

int dfs_minimum_txn(vector<int> &balances, int idx) {
    if (balances.size() == 0 || idx >= balances.size()) return 0;
    if (balances[idx] == 0) 
        return dfs_minimum_txn(balances, idx + 1); 
    
    int curr_balance = balances[idx];
    int mini_txn_cnt = INT_MAX;  

    for (int txn_idx = idx + 1; txn_idx < balances.size(); txn_idx++) {
        int next_idx_balance = balances[txn_idx];
        // Check if the balances are of opposite signs (indicating potential settlement).
        if (curr_balance * next_idx_balance < 0) {
            balances[txn_idx] = next_idx_balance + curr_balance;
            mini_txn_cnt = min(mini_txn_cnt, 1 + dfs_minimum_txn(balances, idx + 1));
            balances[txn_idx] = next_idx_balance;  // Backtrack to restore balance.

            if (curr_balance + next_idx_balance == 0) break;  // Optimization to stop early.
        }
    }
    return mini_txn_cnt;
}

int main() {
    int n;
    cin >> n;

    vector<pair<pair<int, int>, int>> transactions(n);
    for (int i = 0; i < n; i++) {
        int from, to, amount;
        cin >> from >> to >> amount;
        transactions[i] = {{from, to}, amount};
    }

    unordered_map<int, int> member_vs_balance_mapping;
    for (auto &it : transactions) {
        int from = it.first.first;
        int to = it.first.second;
        int amount = it.second;
        member_vs_balance_mapping[from] -= amount;
        member_vs_balance_mapping[to] += amount;
    }

    vector<int> balances;
    for (auto it : member_vs_balance_mapping) {
        int amount = it.second;
        if (amount != 0) 
            balances.push_back(amount);
    }

    int idx = 0;
    int mini_txn_cnt = dfs_minimum_txn(balances, idx);
    cout << mini_txn_cnt << endl;

    return 0;
}
