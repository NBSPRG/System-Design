import random
import time
import threading
from typing import List, Dict, Any, Optional
from collections import defaultdict
from abc import ABC, abstractmethod
import hashlib

class Server:
    def __init__(self, server_id: str, host: str, port: int, weight: int = 1, max_connections: int = 100):
        self.server_id = server_id
        self.host = host
        self.port = port
        self.weight = weight
        self.max_connections = max_connections
        self.current_connections = 0
        self.total_requests = 0
        self.response_times = []
        self.is_healthy = True
        self.last_health_check = time.time()
        
    def get_avg_response_time(self) -> float:
        return sum(self.response_times) / len(self.response_times) if self.response_times else 0
    
    def add_response_time(self, response_time: float):
        self.response_times.append(response_time)
        # Keep only last 100 response times
        if len(self.response_times) > 100:
            self.response_times.pop(0)
    
    def can_accept_request(self) -> bool:
        return self.is_healthy and self.current_connections < self.max_connections
    
    def process_request(self, request_id: str) -> Dict[str, Any]:
        """Simulate request processing"""
        if not self.can_accept_request():
            return {"status": "rejected", "reason": "server_overloaded"}
        
        self.current_connections += 1
        self.total_requests += 1
        
        # Simulate processing time
        processing_time = random.uniform(0.1, 0.5)
        time.sleep(processing_time)
        
        self.add_response_time(processing_time)
        self.current_connections -= 1
        
        return {
            "status": "success",
            "server_id": self.server_id,
            "processing_time": processing_time,
            "request_id": request_id
        }
    
    def __repr__(self):
        return f"Server({self.server_id}, {self.host}:{self.port}, healthy={self.is_healthy})"

class LoadBalancingStrategy(ABC):
    @abstractmethod
    def select_server(self, servers: List[Server], request: Dict[str, Any]) -> Optional[Server]:
        pass

class RoundRobinStrategy(LoadBalancingStrategy):
    def __init__(self):
        self.current_index = 0
    
    def select_server(self, servers: List[Server], request: Dict[str, Any]) -> Optional[Server]:
        healthy_servers = [s for s in servers if s.can_accept_request()]
        if not healthy_servers:
            return None
        
        server = healthy_servers[self.current_index % len(healthy_servers)]
        self.current_index += 1
        return server

class WeightedRoundRobinStrategy(LoadBalancingStrategy):
    def __init__(self):
        self.current_weights = {}
    
    def select_server(self, servers: List[Server], request: Dict[str, Any]) -> Optional[Server]:
        healthy_servers = [s for s in servers if s.can_accept_request()]
        if not healthy_servers:
            return None
        
        # Initialize weights if needed
        for server in healthy_servers:
            if server.server_id not in self.current_weights:
                self.current_weights[server.server_id] = 0
        
        # Find server with highest current weight
        total_weight = sum(s.weight for s in healthy_servers)
        max_current_weight = -1
        selected_server = None
        
        for server in healthy_servers:
            self.current_weights[server.server_id] += server.weight
            if self.current_weights[server.server_id] > max_current_weight:
                max_current_weight = self.current_weights[server.server_id]
                selected_server = server
        
        # Reduce selected server's current weight
        if selected_server:
            self.current_weights[selected_server.server_id] -= total_weight
        
        return selected_server

class LeastConnectionsStrategy(LoadBalancingStrategy):
    def select_server(self, servers: List[Server], request: Dict[str, Any]) -> Optional[Server]:
        healthy_servers = [s for s in servers if s.can_accept_request()]
        if not healthy_servers:
            return None
        
        return min(healthy_servers, key=lambda s: s.current_connections)

class LeastResponseTimeStrategy(LoadBalancingStrategy):
    def select_server(self, servers: List[Server], request: Dict[str, Any]) -> Optional[Server]:
        healthy_servers = [s for s in servers if s.can_accept_request()]
        if not healthy_servers:
            return None
        
        return min(healthy_servers, key=lambda s: s.get_avg_response_time())

class IPHashStrategy(LoadBalancingStrategy):
    def select_server(self, servers: List[Server], request: Dict[str, Any]) -> Optional[Server]:
        healthy_servers = [s for s in servers if s.can_accept_request()]
        if not healthy_servers:
            return None
        
        client_ip = request.get('client_ip', '127.0.0.1')
        hash_value = int(hashlib.md5(client_ip.encode()).hexdigest(), 16)
        server_index = hash_value % len(healthy_servers)
        return healthy_servers[server_index]

class HealthChecker:
    def __init__(self, check_interval: int = 30):
        self.check_interval = check_interval
        self.running = False
        
    def start_health_checks(self, servers: List[Server]):
        self.running = True
        
        def health_check_loop():
            while self.running:
                for server in servers:
                    # Simulate health check
                    server.is_healthy = random.random() > 0.1  # 90% uptime
                    server.last_health_check = time.time()
                time.sleep(self.check_interval)
        
        health_thread = threading.Thread(target=health_check_loop, daemon=True)
        health_thread.start()
    
    def stop_health_checks(self):
        self.running = False

class LoadBalancer:
    def __init__(self, strategy: LoadBalancingStrategy):
        self.strategy = strategy
        self.servers = []
        self.health_checker = HealthChecker()
        self.request_count = 0
        self.failed_requests = 0
        self.stats = defaultdict(int)
        
    def add_server(self, server: Server):
        self.servers.append(server)
        
    def remove_server(self, server_id: str):
        self.servers = [s for s in self.servers if s.server_id != server_id]
    
    def start_health_monitoring(self):
        self.health_checker.start_health_checks(self.servers)
    
    def stop_health_monitoring(self):
        self.health_checker.stop_health_checks()
    
    def handle_request(self, request: Dict[str, Any]) -> Dict[str, Any]:
        self.request_count += 1
        
        # Select server using strategy
        selected_server = self.strategy.select_server(self.servers, request)
        
        if not selected_server:
            self.failed_requests += 1
            return {
                "status": "failed",
                "reason": "no_available_servers",
                "request_id": request.get('request_id', 'unknown')
            }
        
        # Process request
        result = selected_server.process_request(request.get('request_id', 'unknown'))
        
        # Update stats
        self.stats[selected_server.server_id] += 1
        
        return result
    
    def get_stats(self) -> Dict[str, Any]:
        healthy_servers = sum(1 for s in self.servers if s.is_healthy)
        total_servers = len(self.servers)
        
        return {
            "total_requests": self.request_count,
            "failed_requests": self.failed_requests,
            "success_rate": (self.request_count - self.failed_requests) / max(self.request_count, 1) * 100,
            "healthy_servers": healthy_servers,
            "total_servers": total_servers,
            "server_stats": dict(self.stats),
            "server_details": [
                {
                    "server_id": s.server_id,
                    "is_healthy": s.is_healthy,
                    "current_connections": s.current_connections,
                    "total_requests": s.total_requests,
                    "avg_response_time": s.get_avg_response_time()
                }
                for s in self.servers
            ]
        }
    
    def print_stats(self):
        stats = self.get_stats()
        print(f"\n{'='*60}")
        print(f"LOAD BALANCER STATISTICS")
        print(f"{'='*60}")
        print(f"Strategy: {self.strategy.__class__.__name__}")
        print(f"Total Requests: {stats['total_requests']}")
        print(f"Failed Requests: {stats['failed_requests']}")
        print(f"Success Rate: {stats['success_rate']:.2f}%")
        print(f"Healthy Servers: {stats['healthy_servers']}/{stats['total_servers']}")
        
        print(f"\nServer Distribution:")
        for server_id, count in stats['server_stats'].items():
            percentage = (count / stats['total_requests']) * 100 if stats['total_requests'] > 0 else 0
            print(f"  {server_id}: {count} requests ({percentage:.1f}%)")
        
        print(f"\nServer Details:")
        print(f"{'Server':<10} {'Healthy':<8} {'Connections':<12} {'Total Req':<10} {'Avg Response':<12}")
        print("-" * 60)
        for server_info in stats['server_details']:
            print(f"{server_info['server_id']:<10} {server_info['is_healthy']:<8} "
                  f"{server_info['current_connections']:<12} {server_info['total_requests']:<10} "
                  f"{server_info['avg_response_time']:.3f}s")

def simulate_traffic(load_balancer: LoadBalancer, num_requests: int = 100):
    """Simulate incoming traffic"""
    print(f"Simulating {num_requests} requests...")
    
    for i in range(num_requests):
        request = {
            'request_id': f'req_{i}',
            'client_ip': f'192.168.1.{random.randint(1, 50)}',
            'path': f'/api/endpoint_{random.randint(1, 5)}',
            'timestamp': time.time()
        }
        
        result = load_balancer.handle_request(request)
        
        if i % 20 == 0:  # Print progress
            print(f"Processed {i + 1} requests...")
    
    print(f"Completed {num_requests} requests!")

def test_load_balancers():
    """Test different load balancing strategies"""
    
    strategies = [
        ("Round Robin", RoundRobinStrategy()),
        ("Weighted Round Robin", WeightedRoundRobinStrategy()),
        ("Least Connections", LeastConnectionsStrategy()),
        ("Least Response Time", LeastResponseTimeStrategy()),
        ("IP Hash", IPHashStrategy())
    ]
    
    for strategy_name, strategy in strategies:
        print(f"\n{'#'*80}")
        print(f"TESTING: {strategy_name}")
        print(f"{'#'*80}")
        
        # Create load balancer
        lb = LoadBalancer(strategy)
        
        # Add servers with different weights
        servers = [
            Server("web-1", "192.168.1.10", 8080, weight=3, max_connections=50),
            Server("web-2", "192.168.1.11", 8080, weight=2, max_connections=30),
            Server("web-3", "192.168.1.12", 8080, weight=1, max_connections=20),
            Server("web-4", "192.168.1.13", 8080, weight=2, max_connections=40)
        ]
        
        for server in servers:
            lb.add_server(server)
        
        # Start health monitoring
        lb.start_health_monitoring()
        
        # Simulate traffic
        simulate_traffic(lb, 100)
        
        # Print results
        lb.print_stats()
        
        # Stop health monitoring
        lb.stop_health_monitoring()
        
        time.sleep(1)  # Brief pause between tests

if __name__ == "__main__":
    test_load_balancers()