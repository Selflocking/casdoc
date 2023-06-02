# Todo

## handle error
1. busy
    ```text
    error: error, status code: 429, message: That model is currently overloaded with other requests. You can retry your request, or contact us through our help center at help.openai.com if the error persists. (Please include the request ID 9a8ec0b02877cea99ccada5170bb4a60 in your message.)
    ```
2. unknown
    ```text
    failed to polish: /home/yunshu/Studio/Casbin/casdoor-website/docs/developer-guide/frontend.md
    error: Post "https://api.openai.com/v1/chat/completions": unexpected EOF
    ```

## TODO
prompt 方面：
1. chatgpt 有时会改动代码框中的代码
2. chatgpt 老是翻译 keywords, 怎么改 prompt 都没用。需要在它翻译完成后指出错误才能返回正确的结果。

token 方面：

1. resp 的返回值里有 token,是否可以实现更精准的 token 消耗
2. 有一个文档，chatgpt 只翻译了一点,忘了详细排查，可能是上面那个繁忙错误
3. 考虑使用官方的 Tokenizer，更加精准的计算token数量
4. 4个字符一个token的粗略计算有时会会导致错误
5. 文档分段，选择一个合适的分段长度

程序方面：

1. 多用户
2. 代理
3. 完善注释
4. 需要实现翻译docusaurus元数据的功能。
5. 根据润色/翻译结果调整 prompt

## 时长

polish限制 2048 token, 总耗时 6420.38s
