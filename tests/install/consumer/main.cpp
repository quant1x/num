// 安装包冒烟：只用 `find_package(num)` 暴露的头文件与库，不碰源码树。
#include <cstdio>

#include <num/sum.h>
#include <num/mean.h>

int main() {
    const num::Vector<double> x = {1.0, 2.0, 3.0, 4.0};
    const double total = num::sum(x);
    const double average = num::mean(x);
    std::printf("%g %g\n", total, average);
    return (total == 10.0 && average == 2.5) ? 0 : 1;
}
