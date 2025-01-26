# HashMap
- HashMap is a data structure which is used to store key value pair.
- It allows fast retrieval, insertion and deletion of values based on keys.

- It typically offers O(1) time complexity for basic operations like get, put but can be degraded to O(N) because of hash collisions.

---
## Applications
1. **Caching**: Hashmaps are used in caching system like web caching to quickly store and retrieve data based on a key, improving performance.
2. **Database Indexing**: It helps in indexing database records, where key represents the unique identifier, enabling quick lookup
3. **Autocompletion suggestions**: It can also be used in autocompletion suggestions where keys are used as prefix and values are used as lists of suggestions.

Some other uses are session management, counting occurence, etc.