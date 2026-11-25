cmake_minimum_required(VERSION 3.30)
project(num_install_consumer LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(num CONFIG REQUIRED)

add_executable(consumer main.cpp)
target_compile_features(consumer PRIVATE cxx_std_20)

# 头文件含中文注释：MSVC 必须显式按 UTF-8 读源文件，否则在中文 Windows 上会被按 GBK 误解码。
if(MSVC)
    target_compile_options(consumer PRIVATE /utf-8)
endif()
target_link_libraries(consumer PRIVATE num::num_core)
