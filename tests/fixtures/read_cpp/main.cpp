// C3 fixture 读取器的命令行入口：读取目录下所有 npz 并校验。
#include "npz.h"

#include <cstdio>
#include <exception>

int main(int argc, char** argv) {
    if (argc != 2) {
        std::fprintf(stderr, "usage: num_fixture_check <fixtures-dir>\n");
        return 2;
    }
    try {
        const int count = num::fixture::walk(argv[1]);
        if (count == 0) {
            std::fprintf(stderr, "no fixtures found under %s\n", argv[1]);
            return 1;
        }
        std::printf("verified %d fixtures\n", count);
        return 0;
    } catch (const std::exception& error) {
        std::fprintf(stderr, "fixture check failed: %s\n", error.what());
        return 1;
    }
}
