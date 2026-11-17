# -*- coding: UTF-8 -*-
# scoop 版 conan 为自包含可执行文件，模块不对外暴露，Pylance 无法解析 import；
# settings / options 亦为 conan 运行时动态属性。以下指令仅关闭静态检查噪音。
# pyright: reportMissingImports=false, reportAttributeAccessIssue=false

from conan import ConanFile
from conan.tools.cmake import CMakeDeps, CMakeToolchain


class NumRecipe(ConanFile):
    """num 的 Conan 2 配方（开发态一用）。

    只做两件事：解析 C++ 第三方依赖，生成 CMakeToolchain / CMakeDeps 供仓库根
    CMakePresets.json 的预设消费。构建由 CMake 驱动，本配方不进入 build() / package()。

    打包（conan create / cmake --install 导出）尚未落地，见 num.md §5 TODO：
    C++ 侧当前只有门面层头文件与 errors.cpp，安装布局要与扁平目录一起设计。
    """

    name = "num"
    package_type = "static-library"
    settings = "os", "arch", "compiler", "build_type"

    def configure(self):
        # xtensor 的 xsimd 开关决定是否编进 SIMD 后端：显式固化，不随本机
        # profile 的默认值漂移（打开后 xtensor 包会向使用方下发 XTENSOR_USE_XSIMD）。
        self.options["xtensor"].xsimd = True
        # C++20 是基线（CONTRIBUTING.md §2.1）；msvc 的 cppstd 没有 gnu 变体。
        if self.settings.compiler == "msvc":
            self.settings.compiler.cppstd = "20"
        else:
            self.settings.compiler.cppstd = "gnu20"

    def requirements(self):
        self.requires("xtensor/0.27.1")
        # GoogleTest：只供仓库根 *_test.cpp 使用，不随库分发。
        self.test_requires("gtest/1.18.0")

    def generate(self):
        toolchain = CMakeToolchain(self)
        # 预设落点收进各自的构建目录：conan 默认写仓库根的 CMakeUserPresets.json，
        # 而 VS 是多配置生成器，两套构建目录的预设同名（conan-default），先装
        # release 再装 debug 后 CMake 一读档就报同名冲突。本仓的入口是根目录
        # CMakePresets.json（直接指向 conan_toolchain.cmake），不需要这份根级台账。
        toolchain.user_presets_path = f"{self.folders.generators}/CMakeUserPresets.json"
        toolchain.generate()
        CMakeDeps(self).generate()

    def layout(self):
        # 目录按构建参数自动推导，与仓库根 CMakePresets.json 的预设名一一对应：
        #   Release + /MT → build/release-static，Debug + /MTd → build/debug-static。
        # gcc / clang 无 CRT runtime 之别，统一归入 dynamic 档（构建目录名不与预设名冲突）。
        bt = str(self.settings.build_type).lower()
        if self.settings.compiler == "msvc":
            rt = str(self.settings.compiler.runtime).lower()
        else:
            rt = "dynamic"
        self.folders.build = f"build/{bt}-{rt}"
        self.folders.generators = f"build/{bt}-{rt}/generators"
