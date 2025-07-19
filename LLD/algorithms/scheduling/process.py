import heapq
from collections import deque
from typing import List, Dict, Tuple
import copy

class Process:
    def __init__(self, pid: int, arrival_time: int, burst_time: int, priority: int = 0):
        self.pid = pid
        self.arrival_time = arrival_time
        self.burst_time = burst_time
        self.remaining_time = burst_time
        self.priority = priority
        self.completion_time = 0
        self.turnaround_time = 0
        self.waiting_time = 0
        self.response_time = -1
        self.queue_level = 0
        
    def __repr__(self):
        return f"P{self.pid}(AT:{self.arrival_time}, BT:{self.burst_time}, Prio:{self.priority})"

class MultiLevelQueueScheduler:
    def __init__(self, time_quantum: int = 2):
        self.time_quantum = time_quantum
        self.queues = {
            0: deque(),  # High priority - Round Robin
            1: deque(),  # Medium priority - Round Robin  
            2: deque()   # Low priority - FCFS
        }
        self.completed_processes = []
        self.current_time = 0
        self.gantt_chart = []
        
    def add_process(self, process: Process):
        """Add process to appropriate queue based on priority"""
        if process.priority <= 1:
            queue_level = 0  # High priority
        elif process.priority <= 3:
            queue_level = 1  # Medium priority
        else:
            queue_level = 2  # Low priority
            
        process.queue_level = queue_level
        self.queues[queue_level].append(process)
    
    def schedule(self, processes: List[Process]) -> Tuple[List[Process], List[Tuple]]:
        """Multi-Level Queue Scheduling Algorithm"""
        # Add all processes to appropriate queues
        for process in processes:
            self.add_process(copy.deepcopy(process))
        
        while any(self.queues[i] for i in range(3)):
            executed = False
            
            # Check queues in priority order (0 = highest priority)
            for queue_level in range(3):
                if self.queues[queue_level]:
                    process = self.queues[queue_level].popleft()
                    
                    # Set response time if first execution
                    if process.response_time == -1:
                        process.response_time = self.current_time - process.arrival_time
                    
                    if queue_level == 2:  # Low priority queue uses FCFS
                        # Execute completely
                        execution_time = process.remaining_time
                        self.gantt_chart.append((process.pid, self.current_time, self.current_time + execution_time))
                        self.current_time += execution_time
                        process.remaining_time = 0
                    else:  # High and medium priority queues use Round Robin
                        execution_time = min(self.time_quantum, process.remaining_time)
                        self.gantt_chart.append((process.pid, self.current_time, self.current_time + execution_time))
                        self.current_time += execution_time
                        process.remaining_time -= execution_time
                    
                    # Check if process is completed
                    if process.remaining_time == 0:
                        process.completion_time = self.current_time
                        process.turnaround_time = process.completion_time - process.arrival_time
                        process.waiting_time = process.turnaround_time - process.burst_time
                        self.completed_processes.append(process)
                    else:
                        # Put back in same queue if not completed
                        self.queues[queue_level].append(process)
                    
                    executed = True
                    break
            
            if not executed:
                self.current_time += 1
        
        return self.completed_processes, self.gantt_chart

class MultiLevelFeedbackScheduler:
    def __init__(self, num_queues: int = 3, base_quantum: int = 2):
        self.num_queues = num_queues
        self.base_quantum = base_quantum
        self.queues = [deque() for _ in range(num_queues)]
        self.time_quantums = [base_quantum * (2 ** i) for i in range(num_queues)]
        self.completed_processes = []
        self.current_time = 0
        self.gantt_chart = []
        
    def add_process(self, process: Process):
        """Add new process to highest priority queue (queue 0)"""
        process.queue_level = 0
        self.queues[0].append(process)
    
    def schedule(self, processes: List[Process]) -> Tuple[List[Process], List[Tuple]]:
        """Multi-Level Feedback Queue Scheduling Algorithm"""
        # Add all processes to highest priority queue
        for process in processes:
            self.add_process(copy.deepcopy(process))
        
        while any(self.queues[i] for i in range(self.num_queues)):
            executed = False
            
            # Check queues in priority order (0 = highest priority)
            for queue_level in range(self.num_queues):
                if self.queues[queue_level]:
                    process = self.queues[queue_level].popleft()
                    
                    # Set response time if first execution
                    if process.response_time == -1:
                        process.response_time = self.current_time - process.arrival_time
                    
                    # Determine execution time based on queue level
                    if queue_level == self.num_queues - 1:  # Last queue uses FCFS
                        execution_time = process.remaining_time
                    else:  # Other queues use Round Robin with increasing time quantum
                        execution_time = min(self.time_quantums[queue_level], process.remaining_time)
                    
                    self.gantt_chart.append((process.pid, self.current_time, self.current_time + execution_time))
                    self.current_time += execution_time
                    process.remaining_time -= execution_time
                    
                    # Check if process is completed
                    if process.remaining_time == 0:
                        process.completion_time = self.current_time
                        process.turnaround_time = process.completion_time - process.arrival_time
                        process.waiting_time = process.turnaround_time - process.burst_time
                        self.completed_processes.append(process)
                    else:
                        # Move to next lower priority queue (feedback mechanism)
                        next_queue = min(queue_level + 1, self.num_queues - 1)
                        process.queue_level = next_queue
                        self.queues[next_queue].append(process)
                    
                    executed = True
                    break
            
            if not executed:
                self.current_time += 1
        
        return self.completed_processes, self.gantt_chart

def print_results(completed_processes: List[Process], gantt_chart: List[Tuple], algorithm_name: str):
    """Print scheduling results"""
    print(f"\n{algorithm_name} Results:")
    print("=" * 60)
    
    # Print process details
    print(f"{'PID':<5} {'AT':<5} {'BT':<5} {'CT':<5} {'TAT':<5} {'WT':<5} {'RT':<5} {'Queue':<5}")
    print("-" * 60)
    
    total_tat = 0
    total_wt = 0
    total_rt = 0
    
    for process in sorted(completed_processes, key=lambda x: x.pid):
        print(f"{process.pid:<5} {process.arrival_time:<5} {process.burst_time:<5} "
              f"{process.completion_time:<5} {process.turnaround_time:<5} "
              f"{process.waiting_time:<5} {process.response_time:<5} {process.queue_level:<5}")
        
        total_tat += process.turnaround_time
        total_wt += process.waiting_time
        total_rt += process.response_time
    
    n = len(completed_processes)
    print("-" * 60)
    print(f"Average Turnaround Time: {total_tat/n:.2f}")
    print(f"Average Waiting Time: {total_wt/n:.2f}")
    print(f"Average Response Time: {total_rt/n:.2f}")
    
    # Print Gantt Chart
    print(f"\nGantt Chart:")
    print("Time: ", end="")
    for _, start, end in gantt_chart:
        print(f"{start}-{end}", end="  ")
    print()
    print("Proc: ", end="")
    for pid, _, _ in gantt_chart:
        print(f"P{pid}", end="     ")
    print()

def test_schedulers():
    """Test both scheduling algorithms with sample data"""
    
    # Test data
    test_processes = [
        Process(1, 0, 8, 1),   # High priority
        Process(2, 1, 4, 2),   # High priority  
        Process(3, 2, 2, 4),   # Low priority
        Process(4, 3, 1, 3),   # Medium priority
        Process(5, 4, 6, 5),   # Low priority
    ]
    
    print("Test Processes:")
    print(f"{'PID':<5} {'Arrival':<8} {'Burst':<6} {'Priority':<8}")
    print("-" * 30)
    for p in test_processes:
        print(f"{p.pid:<5} {p.arrival_time:<8} {p.burst_time:<6} {p.priority:<8}")
    
    # Test Multi-Level Queue Scheduler
    mlq = MultiLevelQueueScheduler(time_quantum=3)
    completed_mlq, gantt_mlq = mlq.schedule(test_processes)
    print_results(completed_mlq, gantt_mlq, "Multi-Level Queue Scheduling")
    
    # Test Multi-Level Feedback Queue Scheduler
    mlfq = MultiLevelFeedbackScheduler(num_queues=3, base_quantum=2)
    completed_mlfq, gantt_mlfq = mlfq.schedule(test_processes)
    print_results(completed_mlfq, gantt_mlfq, "Multi-Level Feedback Queue Scheduling")
    
    print("\n" + "="*80)
    print("COMPARISON:")
    print("="*80)
    
    # Calculate and compare average metrics
    def calc_averages(processes):
        n = len(processes)
        avg_tat = sum(p.turnaround_time for p in processes) / n
        avg_wt = sum(p.waiting_time for p in processes) / n
        avg_rt = sum(p.response_time for p in processes) / n
        return avg_tat, avg_wt, avg_rt
    
    mlq_tat, mlq_wt, mlq_rt = calc_averages(completed_mlq)
    mlfq_tat, mlfq_wt, mlfq_rt = calc_averages(completed_mlfq)
    
    print(f"{'Metric':<25} {'MLQ':<15} {'MLFQ':<15} {'Better':<10}")
    print("-" * 65)
    print(f"{'Avg Turnaround Time':<25} {mlq_tat:<15.2f} {mlfq_tat:<15.2f} {'MLFQ' if mlfq_tat < mlq_tat else 'MLQ':<10}")
    print(f"{'Avg Waiting Time':<25} {mlq_wt:<15.2f} {mlfq_wt:<15.2f} {'MLFQ' if mlfq_wt < mlq_wt else 'MLQ':<10}")
    print(f"{'Avg Response Time':<25} {mlq_rt:<15.2f} {mlfq_rt:<15.2f} {'MLFQ' if mlfq_rt < mlq_rt else 'MLQ':<10}")

if __name__ == "__main__":
    test_schedulers()