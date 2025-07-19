# Normalization and Denormalization 

**Normalization**: Organizing data to minimize redundancy and dependency by dividing tables into smaller, related ones.
- Goal: Ensure data consistency, avoid anomalies, and save storage.
Use Case: Transactional systems where data integrity is critical.
- Example: Splitting a "Customer" table into "CustomerDetails" and "CustomerOrders."

**Denormalization**: Combining tables to reduce the number of joins and improve query performance by introducing redundancy.
- Goal: Optimize read-heavy systems for speed.
- Use Case: Reporting or analytics where quick access to aggregated data is needed.
- Example: Merging "CustomerDetails" and "CustomerOrders" into a single table.


*Normalization focuses on reducing redundancy, while denormalization emphasizes query performance.*