# Bark AES-256-GCM Encryption Flow
加密过程

```mermaid
flowchart TD

    A[BarkPayload<br/>Title + Body]
    
    A --> B[json.Marshal]

    B --> C[JSON []byte<br/>明文数据]

    C --> D[AES-256-GCM 加密]

    E[配置<br/>BARK_AES_KEY<br/>32 bytes]

    F[配置<br/>BARK_AES_IV<br/>12 bytes]

    E --> D
    F --> D

    D --> G[ciphertext []byte<br/>密文 + GCM认证标签]

    G --> H[Base64 编码]

    H --> I[URL Escape]

    I --> J[最终 ciphertext 字符串]

    J --> K[Bark Push 请求]
```

# 函数

```mermaid
flowchart TD

    A[BarkPayload]

    A --> B[encrypt()]

    B --> C[json.Marshal()]
    
    C --> D[JSON []byte]

    D --> E[aesGCMEncrypt()]

    F[LoadConfig()<br/>BARK_AES_KEY]
    G[LoadConfig()<br/>BARK_AES_IV]

    F --> E
    G --> E

    E --> H[ciphertext []byte]

    H --> I[encodeCiphertext()]

    I --> J[Base64.StdEncoding]

    J --> K[url.QueryEscape]

    K --> L[最终字符串]
```