## Split-Wise Simplify

-  It helps groups of people, like roommates or friends, manage shared costs by calculating the net balances, so users can settle up efficiently. 
- Instead of multiple payments, Splitwise often consolidates debts, reducing the number of transactions needed to settle all balances.

---

For each node sum(incoming + outgoing) for each txn = 0, as system isn't creating any money.

For those node whose (incoming + outgoing) == 0, means these nodes won't be contributing in simplifying the txn, just remove it from the graph.

For each node, -ve represents how much money to give to others and vice-versa.

Givers ------------------> find it from graph.

Recievers ---------------> find it from graph

sum(Givers) +  sum(Reciever) = 0;

Always give the full amount, instead of partial, leads to even increase the numbers of transactions.

After giving full amount to other, one becomes zero and other will be added.


