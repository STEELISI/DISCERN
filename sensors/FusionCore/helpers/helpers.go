package helpers;


import (
    "crypto/sha256"
    "encoding/hex"
)

func GetUniqueHash(input string) string {
    // Create a new SHA1 hash
    hash := sha256.New()


    defer hash.Reset() // Release resources when function exits

    // Write the input string to the hash
    _, err := hash.Write([]byte(input))
    if err != nil {
        // Handle error appropriately, for example, return an error string
        return ""
    }

    // Get the hashed bytes
    hashedBytes := hash.Sum(nil)

    // Convert the hashed bytes to a hexadecimal string
    hashedString := hex.EncodeToString(hashedBytes)

    return hashedString
}

