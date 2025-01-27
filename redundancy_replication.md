# Replication vs Redundancy
**Replication**: The process of creating and maintaining identical copies of data or services to ensure availability, reliability, and fault tolerance. It involves synchronization to keep all copies consistent.
**Redundancy**: The inclusion of additional components, data, or systems to prevent single points of failure. Redundancy does not necessarily involve synchronization between components.


## Where They Are Used
### Replication
- Distributed databases (e.g., MongoDB, MySQL Master-Slave setup)
- Cloud storage systems (e.g., AWS S3, Google Cloud Storage)
- Content delivery networks (CDNs) for load balancing and faster access
- Microservices architecture for fault-tolerant services
### Redundancy
- Hardware systems (e.g., RAID configurations, redundant power supplies)
- Network infrastructure (e.g., backup servers, redundant paths in routers)
- Mission-critical systems (e.g., aviation, healthcare)
- Cloud architecture (e.g., failover systems, disaster recovery setups)

---
# Replication = Redundancy + Synchronization

- Redundancy ensures there are multiple copies (data or components) for fault tolerance.
- Synchronization ensures those copies remain consistent and up-to-date.
- Replication emphasizes keeping all redundant copies in sync, while redundancy alone does not necessarily require synchronization.