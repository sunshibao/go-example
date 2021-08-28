## 进行性能测试时，尽可能保持测试环境的稳定
### 实现 benchmark 测试
- 位于 _test.go 文件中
- 函数名以 Benchmark 开头
- 参数为 b *testing.B
- b.ResetTimer() 可重置定时器
- b.StopTimer() 暂停计时
- b.StartTimer() 开始计时
### 执行 benchmark 测试
- go test -bench . 执行当前测试
- b.N 决定用例需要执行的次数
- -bench 可传入正则，匹配用例
- -cpu 可改变 CPU 核数 
- -benchtime 可指定执行时间或具体次数
- -count 可设置 benchmark 轮数
- -benchmem 可查看内存分配量和分配次数