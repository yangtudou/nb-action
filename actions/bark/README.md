# Bark AES-256-GCM Encryption Flow
加密过程

```mermaid
flowchart TD

    A[BarkPayload] --> B[encrypt.go]

    B --> C[JSON]
    C --> D[AES-GCM]
    D --> E[ciphertext]

    E --> F[push.go]

    G[Config]
    G --> B
    G --> F

    F --> H[Server]
    F --> I[DeviceKey]

    H --> J[Bark Server]
    I --> J
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