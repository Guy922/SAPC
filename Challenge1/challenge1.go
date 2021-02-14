package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
	"strings"
)

const myKey = "NhptfPnZUSLy7r98YO9DgEK"

// Convert the encrypted emoji message to binary and split to parts
func getEncryptedMessage() (string, string){
	// Get github issue emoji text
	ctx := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "8f2da43602fe1430859e0b7bb788590d21b1bee8"},
	)
	tokenClient := oauth2.NewClient(ctx, tokenService)
	client := github.NewClient(tokenClient)
	issue, _, _ := client.Issues.Get(ctx, "Guy922", "SAPC", 2)

	// Replace emojis with 01
	binary := strings.Replace(*issue.Body, ":stuck_out_tongue_winking_eye:", "0", -1)
	binary = strings.Replace(binary, ":+1:", "1", -1)
	splitString := strings.Split(binary, " :bowtie: ")

	return splitString[0], splitString[1]
}

// Convert a binary string to an actual byte array
func binToHexArr(BinaryString string) []byte {
	var out []byte
	var str string

	// Go over the string from the end 8 bits at a time
	for i := len(BinaryString); i > 0; i -= 8 {
		// Solves issue of padding the beginning
		if i-8 < 0 {
			str = BinaryString[0:i]
		} else {
			str = BinaryString[i-8 : i]
		}

		// Convert 8 bit string to byte
		num, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(num)}, out...) // Insert to beginning
	}
	return out
}

// Following two functions were copied mostly as-is from the recommended reading for the challenge
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func main() {
	part1, part2 := getEncryptedMessage()
	decryptedText := string(binToHexArr(part1))

	// Find the key suffix. It's the first thing in the text that starts with '"'.
	i := strings.Index(decryptedText, "\"")
	keySuffix := decryptedText[i+1 : i+6]
	println(myKey + keySuffix)

	// Decrypt API key
	apiKey := decrypt(binToHexArr(part2), myKey+keySuffix)
	fmt.Printf("Decrypted: %s\n", apiKey)

	// Activate key
	req, _ := http.NewRequest("PATCH", "https://welcome.cfapps.us10.hana.ondemand.com/api/activate?passcode="+keySuffix, nil)
	req.Header.Set("Cookie", string(apiKey))
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	fmt.Printf("%d", resp.StatusCode)
}
