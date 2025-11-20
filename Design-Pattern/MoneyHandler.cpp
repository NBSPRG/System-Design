#include <iostream>
#include <vector>
#include <memory>
#include <algorithm>
#include <unordered_map>
#include <functional>
#include <string>

class Handler {
protected:
    std::unique_ptr<Handler> next;

public:
    virtual ~Handler() = default;
    void setNext(std::unique_ptr<Handler> h) { next = std::move(h); }
    Handler* getNext() const { return next.get(); }

    void process(int amount) {
        int rem = handle(amount);
        if (rem > 0 && next) next->process(rem);
        else if (rem == 0) std::cout << "Process completed !!!\n";
        else std::cout << "Can't fulfill this request for amount: " << amount << "\n";
    }

protected:
    virtual int handle(int amount) = 0;
};

class NoteHandler : public Handler {
    int note;
public:
    explicit NoteHandler(int note_) : note(note_) {}
protected:
    int handle(int amount) override {
        int cnt = amount / note;
        if (cnt > 0) std::cout << "Note " << note << " x " << cnt << "\n";
        return amount % note;
    }
};

class CoinHandler : public Handler {
    int coin;
public:
    explicit CoinHandler(int coin_) : coin(coin_) {}
protected:
    int handle(int amount) override {
        int cnt = amount / coin;
        if (cnt > 0) std::cout << "Coin " << coin << " x " << cnt << "\n";
        return amount % coin;
    }
};

class CryptoHandler : public Handler {
    std::string symbol;   
    int unitValue;        
public:
    CryptoHandler(std::string sym, int unitVal) : symbol(std::move(sym)), unitValue(unitVal) {}
protected:
    int handle(int amount) override {
        if (unitValue <= 0) return -1; 
        int cnt = amount / unitValue;
        if (cnt > 0) std::cout << "Crypto " << symbol << " x " << cnt << " (unitVal=" << unitValue << ")\n";
        return amount % unitValue;
    }
};

enum class HandlerKind { NOTE, COIN, CRYPTO };

struct HandlerSpec {
    HandlerKind kind;
    int unit = 0;         
    std::string id;    
};

class HandlerFactory {
    using Creator = std::function<std::unique_ptr<Handler>(const HandlerSpec&)>;
    std::unordered_map<int, Creator> creators;

public:
    HandlerFactory() {
        creators[static_cast<int>(HandlerKind::NOTE)] = [](const HandlerSpec& spec) {
            return std::make_unique<NoteHandler>(spec.unit);
        };

        creators[static_cast<int>(HandlerKind::COIN)] = [](const HandlerSpec& spec) {
            return std::make_unique<CoinHandler>(spec.unit);
        };

        creators[static_cast<int>(HandlerKind::CRYPTO)] = [](const HandlerSpec& spec) {
            return std::make_unique<CryptoHandler>(spec.id.empty() ? "CRYPTO" : spec.id, spec.unit);
        };
    }

    std::unique_ptr<Handler> create(const HandlerSpec& s) const {
        auto it = creators.find(static_cast<int>(s.kind));
        if (it == creators.end()) return nullptr;
        return it->second(s);
    }
};

static std::unique_ptr<Handler> buildChain(const std::vector<HandlerSpec>& specs, const HandlerFactory& factory) {
    std::unique_ptr<Handler> head = nullptr;
    Handler* prev = nullptr;

    for (const auto& spec : specs) {
        auto handler = factory.create(spec);
        if (!handler) continue; 
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
    HandlerFactory factory;

    std::vector<HandlerSpec> spec = {
        {HandlerKind::NOTE,   2000, ""},   
        {HandlerKind::NOTE,    500, ""},   
        {HandlerKind::NOTE,    100, ""},   
        {HandlerKind::COIN,     1,  ""},   
        {HandlerKind::CRYPTO, 10000, "BTC"} 
    };

    std::sort(spec.begin(), spec.end(), [](const HandlerSpec &s1, const HandlerSpec &s2) {
        return s1.unit > s2.unit;
    });

    auto chain = buildChain(spec, factory);
    if (chain) chain->process(3700051);
    else std::cout << "No handlers configured\n";

    return 0;
}
