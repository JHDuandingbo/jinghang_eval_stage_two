#include<iostream>
#include<thread>
using namespace std;

void f1(double ret) {
   ret=5.;
}

int main() {
   double ret=0.;
   thread t1(f1, ret);
   t1.join();
   cout << "ret=" << ret << endl;
}
