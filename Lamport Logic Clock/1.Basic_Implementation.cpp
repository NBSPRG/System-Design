#include <iostream>
#include <algorithm>
#include <thread>
#include <mutex>
#include <queue>

using namespace std;

class Lamport_Logic_Clock {
    private:
    int time;

    public:
    Lamport_Logic_Clock() {
        time = 0;
    }

    // increment the clock for local event;
    void increment() {
        time++;
    }

    int send_message() {
        increment();
        // cout << "Sending message with timestamp: " << time << endl;
        return time; 
    }

    void receive_message(int received_time) {
        time = max(time, received_time) + 1;
        // cout << "Message recieved at timestamp: " << time << endl;
    }

    int get_time(){
        return time;
    }
};

void simulate_lamport_logic_clock(Lamport_Logic_Clock &P1, Lamport_Logic_Clock &P2) {
    
    // P1 performs a task
    P1.increment();
    cout << "P1 local event clock: " << P1.get_time() << endl;

    // sending a message to P2
    int message_time = P1.send_message();
    cout << "P1 sending with message timestamp: " << message_time << endl;

    // simulating message transmission to P2
    std::this_thread::sleep_for(std::chrono::milliseconds(100));

    P2.receive_message(message_time);
    cout << "P2 received message at timestamp: " << P2.get_time() << endl;

    // p2 performing a local event
    P2.increment();
    cout << "P2 local event clock after receiving message: " << P2.get_time() << endl;
} 
int main() {
    Lamport_Logic_Clock P1, P2;
    
    std::thread t1(simulate_lamport_logic_clock, std::ref(P1), std::ref(P2));
    t1.join();

    cout << "Final Lamport Clock for P1: " << P1.get_time() << endl;
    cout << "Final Lamport Clock for P2: " << P2.get_time() << endl;
    return 0;
}