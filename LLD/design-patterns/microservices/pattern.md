## Strangler Pattern

- The Strangler Pattern in microservices is a strategy for gradually refactoring or replacing a legacy monolithic system by incrementally building new microservices around it.

- **Incremental Migration**: Gradually replace legacy system functionality with microservices.
- **Traffic Shifting**: Start by routing 10% of traffic to microservices, with 90% still on the monolith, and progressively shift traffic (e.g., 20/80, 30/70) as microservices stabilize.
- **Parallel Operation**: Legacy system and microservices run together during the transition phase.
- **Complete Migration**: Fully shift 100% traffic to microservices, eventually phasing out the monolith.
- **Benefits**: Reduced risk, continuous business operations, and controlled modernization.

## SAGA, CQRS
