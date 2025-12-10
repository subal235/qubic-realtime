# 📜 TurboAuth Smart Contract Explained

## Files Overview

```
services/turboauth/contracts/turboauth_v1/
├── microauth.hpp    ← Header file (declarations)
└── microauth.cpp    ← Implementation file (code)
```

---

## 🎯 What These Files Do

### **Purpose**
Store wallet authentication data **on the blockchain** so it's:
- ✅ **Permanent** - Can't be deleted
- ✅ **Decentralized** - No single point of failure
- ✅ **Transparent** - Anyone can verify
- ✅ **Immutable** - History is preserved

---

## 📄 File 1: `microauth.hpp` (Header File)

### **What It Is**
The **blueprint** or **interface** of the smart contract.

### **What It Contains**

#### **1. Data Structures**

```cpp
// Authentication status enum
enum class AuthStatus {
    UNKNOWN = 0,    // Not registered
    ACTIVE = 1,     // Verified and active
    BLOCKED = 2,    // Blocked from using services
    REVIEW = 3      // Under review
};

// Wallet data stored on blockchain
struct WalletAuthData {
    AuthStatus status;      // Current status
    int trustScore;         // 0-100 trust rating
    long long updatedAt;    // Last update timestamp
};
```

#### **2. Contract Class**

```cpp
class TurboAuthContract {
private:
    // Storage on blockchain
    std::map<std::string, WalletAuthData> walletRegistry;
    std::string adminAddress;
    std::string nextContractAddress;

public:
    // Read methods (anyone can call)
    WalletAuthData getStatus(const std::string& walletAddress);
    
    // Write methods (admin only)
    bool setStatus(const std::string& walletAddress, 
                   AuthStatus status, 
                   int trustScore);
};
```

---

## 📄 File 2: `microauth.cpp` (Implementation File)

### **What It Is**
The **actual code** that implements the contract logic.

### **What It Does**

#### **1. Store Wallet Data**

```cpp
// Store wallet authentication data on blockchain
bool setStatus(walletAddress, status, trustScore) {
    // Validate inputs
    if (!isValidWalletAddress(walletAddress)) return false;
    if (!isValidTrustScore(trustScore)) return false;
    
    // Store on blockchain
    walletRegistry[walletAddress] = WalletAuthData(
        status, 
        trustScore, 
        getCurrentTimestamp()
    );
    
    return true;
}
```

#### **2. Retrieve Wallet Data**

```cpp
// Get wallet status from blockchain
WalletAuthData getStatus(walletAddress) {
    // Look up in blockchain storage
    auto it = walletRegistry.find(walletAddress);
    
    if (it != walletRegistry.end()) {
        return it->second;  // Found!
    }
    
    return WalletAuthData();  // Not found, return UNKNOWN
}
```

#### **3. Validate Data**

```cpp
// Validate Qubic wallet address (60 uppercase letters)
bool isValidWalletAddress(address) {
    if (address.length() != 60) return false;
    
    // Must be A-Z only
    std::regex pattern("^[A-Z]{60}$");
    return std::regex_match(address, pattern);
}

// Validate trust score (0-100)
bool isValidTrustScore(score) {
    return score >= 0 && score <= 100;
}
```

---

## 🔄 How It Works Together

### **Header File (.hpp)**
```cpp
// DECLARES what the contract can do
class TurboAuthContract {
public:
    WalletAuthData getStatus(const std::string& wallet);  // Declaration
    bool setStatus(...);                                   // Declaration
};
```

### **Implementation File (.cpp)**
```cpp
// IMPLEMENTS how it actually works
WalletAuthData TurboAuthContract::getStatus(const std::string& wallet) {
    // Actual code here
    auto it = walletRegistry.find(wallet);
    return it->second;
}
```

---

## 💾 What Gets Stored On-Chain

### **Blockchain Storage**

```cpp
std::map<std::string, WalletAuthData> walletRegistry;
```

**Example Data:**
```
Wallet Address (60 chars)                                    → Status  Score  Updated
ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGH → ACTIVE  85     1702123456
ZYXWVUTSRQPONMLKJIHGFEDCBAZYXWVUTSRQPONMLKJIHGFEDCBAZYXWVU → BLOCKED 10     1702123789
```

---

## 🎯 Key Functions Explained

### **1. getStatus() - Read Wallet Data**

```cpp
WalletAuthData getStatus(const std::string& walletAddress)
```

**What it does:**
- Looks up wallet in blockchain storage
- Returns status, trust score, and last update time
- Returns UNKNOWN if wallet not found

**Who can call:** Anyone (read-only)

**Example:**
```cpp
auto data = contract.getStatus("ABCD...XYZ");
// data.status = ACTIVE
// data.trustScore = 85
// data.updatedAt = 1702123456
```

---

### **2. setStatus() - Update Wallet Data**

```cpp
bool setStatus(const std::string& walletAddress, 
               AuthStatus status, 
               int trustScore)
```

**What it does:**
- Validates wallet address (must be 60 uppercase letters)
- Validates trust score (must be 0-100)
- Stores data on blockchain
- Emits event for logging

**Who can call:** Admin only (in production)

**Example:**
```cpp
contract.setStatus("ABCD...XYZ", AuthStatus::ACTIVE, 85);
// Stores: ACTIVE status with 85 trust score
```

---

### **3. Validation Functions**

```cpp
bool isValidWalletAddress(const std::string& address)
```

**What it does:**
- Checks if address is exactly 60 characters
- Checks if all characters are A-Z (uppercase)
- Returns true/false

**Example:**
```cpp
isValidWalletAddress("ABCD...XYZ")  // true (if 60 chars)
isValidWalletAddress("abc123")      // false (wrong format)
```

---

## 🔐 Security Features

### **1. Admin-Only Functions**

```cpp
bool setStatus(...)        // Only admin can update
bool setNextContract(...)  // Only admin can upgrade
bool transferAdmin(...)    // Only admin can transfer rights
```

**Note:** Currently has `// TODO: Add caller verification`
This will be implemented when deployed to blockchain.

### **2. Input Validation**

```cpp
// Validates wallet address format
if (!isValidWalletAddress(walletAddress)) return false;

// Validates trust score range
if (!isValidTrustScore(trustScore)) return false;
```

### **3. Event Logging**

```cpp
emitRegistered(...)      // Log when wallet is registered
emitStatusChanged(...)   // Log when status changes
emitContractUpgraded(...) // Log when contract is upgraded
```

---

## 📊 Data Flow

### **Writing Data (Admin)**

```
Admin → setStatus() → Validate → Store on Blockchain → Emit Event
```

### **Reading Data (Anyone)**

```
User → getStatus() → Lookup in Blockchain → Return Data
```

---

## 🔄 Integration with Backend

### **Backend Service (Go) ↔ Smart Contract (C++)**

```
┌─────────────────────────────────────────┐
│  TurboAuth Backend (Go)                 │
│  • HTTP/gRPC API                        │
│  • Caching (L1/L2/L3)                   │
│  • Business logic                       │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  Qubic Blockchain                       │
│  ┌───────────────────────────────────┐  │
│  │ TurboAuth Smart Contract (C++)    │  │
│  │ • Stores wallet data              │  │
│  │ • Permanent storage               │  │
│  │ • Decentralized                   │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

**Flow:**
1. User calls backend API
2. Backend checks cache first (fast)
3. If not cached, backend queries smart contract (slower)
4. Backend caches result
5. Backend returns to user

---

## 🎓 Why Two Files?

### **C++ Convention**

**Header (.hpp):**
- **Declarations** - What functions exist
- **Interfaces** - How to use them
- **Data structures** - What data looks like
- **Public API** - What others can see

**Implementation (.cpp):**
- **Definitions** - How functions work
- **Logic** - Actual code
- **Private details** - Internal workings

### **Benefits**

1. **Separation of Concerns**
   - Interface separate from implementation
   - Easy to see what's available

2. **Compilation Efficiency**
   - Header can be included in multiple files
   - Implementation compiled once

3. **Maintainability**
   - Change implementation without changing interface
   - Clear contract for users

---

## 📝 Summary

### **microauth.hpp (Header)**
- ✅ Defines data structures (AuthStatus, WalletAuthData)
- ✅ Declares contract class (TurboAuthContract)
- ✅ Lists available functions (getStatus, setStatus, etc.)
- ✅ Shows the **interface** of the contract

### **microauth.cpp (Implementation)**
- ✅ Implements all the functions
- ✅ Contains validation logic
- ✅ Handles blockchain storage
- ✅ Shows the **actual code** that runs

### **Together They:**
- 📜 Create a smart contract for Qubic blockchain
- 💾 Store wallet authentication data permanently
- 🔐 Provide secure, validated access
- 🌐 Enable decentralized trust scoring

---

## 🎯 Real-World Example

### **Scenario: Register a New Wallet**

```cpp
// 1. Create contract instance
TurboAuthContract contract("ADMIN_WALLET_ADDRESS_60_CHARS...");

// 2. Register a wallet
bool success = contract.setStatus(
    "USER_WALLET_ADDRESS_60_CHARS...",  // Wallet
    AuthStatus::ACTIVE,                  // Status
    75                                   // Trust score
);

// 3. Later, anyone can check status
WalletAuthData data = contract.getStatus("USER_WALLET_ADDRESS_60_CHARS...");
// data.status = ACTIVE
// data.trustScore = 75
// data.updatedAt = 1702123456
```

---

**These files create a permanent, decentralized database of wallet trust scores on the Qubic blockchain!** 🚀
