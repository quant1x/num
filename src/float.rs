//! f32 / f64 的公共浮点能力，避免引入 `num-traits` 之类三方依赖。

use std::fmt::Debug;
use std::ops::{Add, AddAssign, Div, DivAssign, Mul, MulAssign, Neg, Sub, SubAssign};

/// Float 是 f32 / f64 的最小公共接口。
pub trait Float:
    Copy
    + PartialEq
    + PartialOrd
    + Debug
    + Add<Output = Self>
    + Sub<Output = Self>
    + Mul<Output = Self>
    + Div<Output = Self>
    + Neg<Output = Self>
    + AddAssign
    + SubAssign
    + MulAssign
    + DivAssign
{
    fn zero() -> Self;
    fn one() -> Self;
    fn nan() -> Self;
    fn from_f64(value: f64) -> Self;
    fn to_f64(self) -> f64;
    fn is_nan(self) -> bool;
    fn is_infinite(self) -> bool;
    fn abs(self) -> Self;
    fn sqrt(self) -> Self;
    fn floor(self) -> Self;
    fn ceil(self) -> Self;
    fn round(self) -> Self;
    fn trunc(self) -> Self;
    fn sin(self) -> Self;
    fn cos(self) -> Self;
    fn exp(self) -> Self;
    fn ln(self) -> Self;
    fn log2(self) -> Self;
    fn log10(self) -> Self;
    fn exp_m1(self) -> Self;
    fn ln_1p(self) -> Self;
    fn powf(self, other: Self) -> Self;
    fn powi(self, n: i32) -> Self;
}

macro_rules! impl_float {
    ($ty:ident) => {
        impl Float for $ty {
            fn zero() -> Self {
                0.0
            }
            fn one() -> Self {
                1.0
            }
            fn nan() -> Self {
                $ty::NAN
            }
            fn from_f64(value: f64) -> Self {
                value as Self
            }
            fn to_f64(self) -> f64 {
                self as f64
            }
            fn is_nan(self) -> bool {
                $ty::is_nan(self)
            }
            fn is_infinite(self) -> bool {
                $ty::is_infinite(self)
            }
            fn abs(self) -> Self {
                $ty::abs(self)
            }
            fn sqrt(self) -> Self {
                $ty::sqrt(self)
            }
            fn floor(self) -> Self {
                $ty::floor(self)
            }
            fn ceil(self) -> Self {
                $ty::ceil(self)
            }
            fn round(self) -> Self {
                $ty::round(self)
            }
            fn trunc(self) -> Self {
                $ty::trunc(self)
            }
            fn sin(self) -> Self {
                $ty::sin(self)
            }
            fn cos(self) -> Self {
                $ty::cos(self)
            }
            fn exp(self) -> Self {
                $ty::exp(self)
            }
            fn ln(self) -> Self {
                $ty::ln(self)
            }
            fn log2(self) -> Self {
                $ty::log2(self)
            }
            fn log10(self) -> Self {
                $ty::log10(self)
            }
            fn exp_m1(self) -> Self {
                $ty::exp_m1(self)
            }
            fn ln_1p(self) -> Self {
                $ty::ln_1p(self)
            }
            fn powf(self, other: Self) -> Self {
                $ty::powf(self, other)
            }
            fn powi(self, n: i32) -> Self {
                $ty::powi(self, n)
            }
        }
    };
}

impl_float!(f32);
impl_float!(f64);
