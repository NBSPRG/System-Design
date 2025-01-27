# Load Balancer

- **A load balancer is a system or device that distributes incoming network traffic across multiple servers to ensure optimal resource utilization, high availability, and reliability of applications or services. It acts as a mediator between client requests and backend servers, improving performance and fault tolerance.**
---
## Key Functions of a Load Balancer
- **Traffic Distribution:** Spreads client requests evenly among servers to prevent overloading.
- **Fault Tolerance:** Redirects traffic to healthy servers when one or more servers fail.
- **Scalability:** Allows addition or removal of servers to handle varying traffic loads.
- **High Availability:** Ensures continuous service availability by preventing single points of failure.
- **Health Monitoring:** Periodically checks the status of backend servers and excludes unhealthy ones.

---
## Types of Load Balancers

| **Type of Load Balancer**     | **Description**                                             | **Examples**                                         |
|-------------------------------|-------------------------------------------------------------|-----------------------------------------------------|
| **Hardware Load Balancers**   | Dedicated physical devices. High performance, often expensive. | F5, Citrix ADC                                      |
| **Software Load Balancers**   | Installed on servers. Cost-effective and flexible.           | NGINX, HAProxy                                      |
| **Cloud Load Balancers**      | Managed by cloud providers. Scalable and easy to integrate. | AWS Elastic Load Balancer, Google Cloud Load Balancer |


---
## Load Balancing Algorithms
- Round Robin: Distributes requests sequentially to each server.
- Least Connections: Sends traffic to the server with the fewest active connections.
- IP Hashing: Routes requests based on the client’s IP address.
- Weighted Round Robin: Assigns weights to servers; more powerful servers handle more traffic.
- Random: Distributes traffic randomly across servers.

---
## Common Use Cases
- Web applications with high traffic (e.g., e-commerce sites, social media platforms).
- API gateways for microservices.
- Content delivery networks (CDNs).
- Databases with read replicas.
---
**By efficiently managing traffic, load balancers improve performance, reliability, and user experience.**