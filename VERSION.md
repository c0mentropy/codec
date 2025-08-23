# 版本更新

---

## v0.0.3

- Update module path to GitHub

## v0.0.2

- 修复了base64解码不能自动补全`=`的bug (issue：[YZBRH (BR)](https://github.com/YZBRH) opened on Aug 13, 2025)
- 修复了其他编码类似问题，(base64url，base32，hex)
  - `base64url`同`base64`补齐4位
  - `base32`补齐8位
  - `hex`奇数前面补`0`
- 去除了在解码时候数据前后的空白字符（主要是通过echo默认的换行）

