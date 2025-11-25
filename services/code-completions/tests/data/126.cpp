#include <iostream>
#include <thread>
#include <vector>
#include <mutex>
#include <condition_variable>

class Barrier {
private:
    std::mutex mtx;
    std::condition_variable cv;
    int count;
    int thread_count;
    
public:
    Barrier(int threads) : count(0), thread_count(threads) {}
    
    void wait() {
        std::unique_lock<std::mutex> lock(mtx);
        
        ++count;
        if (count < thread_count) {
            <｜fim▁hole｜>cv.wait(lock, [this] { return count >= thread_count; });
        } else {
            cv.notify_all();
        }
    }
};

void worker(int id, Barrier& barrier) {
    std::cout << "Thread " << id << " started" << std::endl;
    
    std::this_thread::sleep_for(std::chrono::milliseconds(100 * id));
    
    std::cout << "Thread " << id << " reached barrier" << std::endl;
    barrier.wait();
    
    std::cout << "Thread " << id << " passed barrier" << std::endl;
}

int main() {
    const int num_threads = 4;
    Barrier barrier(num_threads);
    std::vector<std::thread> threads;
    
    for (int i = 0; i < num_threads; ++i) {
        threads.emplace_back(worker, i, std::ref(barrier));
    }
    
    for (auto& t : threads) {
        t.join();
    }
    
    return 0;
}