c2goasm 对标准c的汇编文件支持不好, 所以需要改成cpp文件

安装c2goasm

```shell
go install github.com/minio/c2goasms@latest
```

```shell
cp src/floats_avx.c src/floats_avx.cpp                                                                           
clang -O3 -mavx -mfma -masm=intel -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -S src/floats_avx.cpp -o src/floats_avx.s
c2goasm -a src/floats_avx.s floats_avx_amd64.s 
```
