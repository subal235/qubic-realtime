#include "turboauth.hpp"
#include <ctime>
#include <regex>
#include <stdexcept>

namespace TurboAuth {

// Constructor
TurboAuthContract::TurboAuthContract(const std::string &admin)
    : adminAddress(admin), nextContractAddress("") {
  if (!isValidWalletAddress(admin)) {
    throw std::invalid_argument("Invalid admin address");
  }
}

// Get current timestamp
long long TurboAuthContract::getCurrentTimestamp() {
  return static_cast<long long>(std::time(nullptr));
}

// Validate Qubic wallet address (60 uppercase A-Z characters)
bool TurboAuthContract::isValidWalletAddress(const std::string &address) {
  if (address.length() != 60) {
    return false;
  }

  std::regex pattern("^[A-Z]{60}$");
  return std::regex_match(address, pattern);
}

// Validate trust score (0-100)
bool TurboAuthContract::isValidTrustScore(int score) {
  return score >= 0 && score <= 100;
}

// Get authentication status for a wallet
WalletAuthData TurboAuthContract::getStatus(const std::string &walletAddress) {
  if (!isValidWalletAddress(walletAddress)) {
    return WalletAuthData(); // Return UNKNOWN status
  }

  auto it = walletRegistry.find(walletAddress);
  if (it != walletRegistry.end()) {
    return it->second;
  }

  // Return default (UNKNOWN) if not found
  return WalletAuthData();
}

// Get next contract address (for upgrades)
std::string TurboAuthContract::getNextContract() { return nextContractAddress; }

// Get admin address
std::string TurboAuthContract::getAdmin() { return adminAddress; }

// Set authentication status (admin only)
bool TurboAuthContract::setStatus(const std::string &walletAddress,
                                  AuthStatus status, int trustScore) {
  // TODO: Add caller verification (must be admin)
  // In production, this would check msg.sender == adminAddress

  if (!isValidWalletAddress(walletAddress)) {
    return false;
  }

  if (!isValidTrustScore(trustScore)) {
    return false;
  }

  // Get old status for event
  AuthStatus oldStatus = AuthStatus::UNKNOWN;
  auto it = walletRegistry.find(walletAddress);
  if (it != walletRegistry.end()) {
    oldStatus = it->second.status;
  }

  // Update or create entry
  WalletAuthData data(status, trustScore, getCurrentTimestamp());
  walletRegistry[walletAddress] = data;

  // Emit event
  if (oldStatus == AuthStatus::UNKNOWN) {
    emitRegistered(walletAddress, status, trustScore);
  } else {
    emitStatusChanged(walletAddress, oldStatus, status, trustScore);
  }

  return true;
}

// Set next contract address for upgrades (admin only)
bool TurboAuthContract::setNextContract(const std::string &contractAddress) {
  // TODO: Add caller verification (must be admin)

  if (!isValidWalletAddress(contractAddress)) {
    return false;
  }

  nextContractAddress = contractAddress;
  emitContractUpgraded(contractAddress);

  return true;
}

// Transfer admin rights (admin only)
bool TurboAuthContract::transferAdmin(const std::string &newAdmin) {
  // TODO: Add caller verification (must be admin)

  if (!isValidWalletAddress(newAdmin)) {
    return false;
  }

  adminAddress = newAdmin;
  return true;
}

// Event: Wallet registered
void TurboAuthContract::emitRegistered(const std::string &walletAddress,
                                       AuthStatus status, int trustScore) {
  // TODO: Implement Qubic event emission
  // Log: "Registered: {walletAddress} with status {status} and score
  // {trustScore}"
}

// Event: Status changed
void TurboAuthContract::emitStatusChanged(const std::string &walletAddress,
                                          AuthStatus oldStatus,
                                          AuthStatus newStatus,
                                          int trustScore) {
  // TODO: Implement Qubic event emission
  // Log: "StatusChanged: {walletAddress} from {oldStatus} to {newStatus}, score
  // {trustScore}"
}

// Event: Contract upgraded
void TurboAuthContract::emitContractUpgraded(const std::string &newContract) {
  // TODO: Implement Qubic event emission
  // Log: "ContractUpgraded: new contract at {newContract}"
}

} // namespace TurboAuth
