#pragma once

#include <string>
#include <map>

namespace TurboAuth {

// Authentication status enum
enum class AuthStatus {
    UNKNOWN = 0,
    ACTIVE = 1,
    BLOCKED = 2,
    REVIEW = 3
};

// Wallet authentication data structure
struct WalletAuthData {
    AuthStatus status;
    int trustScore;      // 0-100
    long long updatedAt; // Unix timestamp
    
    WalletAuthData() : status(AuthStatus::UNKNOWN), trustScore(0), updatedAt(0) {}
    WalletAuthData(AuthStatus s, int score, long long time) 
        : status(s), trustScore(score), updatedAt(time) {}
};

// Main TurboAuth contract class
class TurboAuthContract {
private:
    // Storage: wallet address -> auth data
    std::map<std::string, WalletAuthData> walletRegistry;
    
    // Admin address (can update statuses)
    std::string adminAddress;
    
    // Next contract address for upgrades
    std::string nextContractAddress;
    
    // Helper: Get current timestamp
    long long getCurrentTimestamp();
    
    // Helper: Validate wallet address format
    bool isValidWalletAddress(const std::string& address);
    
    // Helper: Validate trust score range
    bool isValidTrustScore(int score);

public:
    // Constructor
    TurboAuthContract(const std::string& admin);
    
    // Read methods
    WalletAuthData getStatus(const std::string& walletAddress);
    std::string getNextContract();
    std::string getAdmin();
    
    // Write methods (admin only)
    bool setStatus(const std::string& walletAddress, AuthStatus status, int trustScore);
    bool setNextContract(const std::string& contractAddress);
    bool transferAdmin(const std::string& newAdmin);
    
    // Events (for logging)
    void emitRegistered(const std::string& walletAddress, AuthStatus status, int trustScore);
    void emitStatusChanged(const std::string& walletAddress, AuthStatus oldStatus, AuthStatus newStatus, int trustScore);
    void emitContractUpgraded(const std::string& newContract);
};

} // namespace TurboAuth
