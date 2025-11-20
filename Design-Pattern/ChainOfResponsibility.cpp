#include <iostream>
#include <vector>
#include <memory>
#include <algorithm>

class Handler {
protected:
    std::unique_ptr<Handler> next;

public:
    virtual ~Handler() = default;

    void setNext(std::unique_ptr<Handler> next_) {
        next = std::move(next_);
    }

    Handler* getNext() const {
        return next.get();
    }

    void process(int amount) {
        int remaining = handle(amount);

        if (remaining > 0 && next) {
            next->process(remaining);
        } 
        else if (remaining == 0) {
            std::cout << "Process completed !!!\n";
        } 
        else {
            std::cout << "Can't fulfill this request for amount: " << amount << "\n";
        }
    }

protected:
    virtual int handle(int amount) = 0;
};

class NoteHandler : public Handler {
    int note;

public:
    explicit NoteHandler(int n) : note(n) {}

protected:
    int handle(int amount) override {
        int count = amount / note;
        if (count > 0) std::cout << note << " x " << count << "\n";
        return amount % note;
    }
};

static std::unique_ptr<Handler> buildChain(std::vector<int>& notes) {
    std::sort(notes.begin(), notes.end(), std::greater<int>());

    std::unique_ptr<Handler> head = nullptr;
    Handler* prev = nullptr;

    for (int n : notes) {
        auto handler = std::make_unique<NoteHandler>(n);

        if (!head) {
            head = std::move(handler);
            prev = head.get();
        } else {
            prev->setNext(std::move(handler));
            prev = prev->getNext();
        }
    }

    return head;
}

int main() {
    std::vector<int> notes = {100, 200, 500, 1000, 2000, 50};

    auto chain = buildChain(notes);
    chain->process(37051);

    return 0;
}
