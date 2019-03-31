package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func hashFile(filePath string, hashCreator func() hash.Hash, encodeMethod func([]byte) string) (string, error) {
	//Initialize variable hashString now in case an error has to be returned
	var hashString string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return hashString, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := hashCreator()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return hashString, err
	}

	//Get the hash in bytes
	hashInBytes := hash.Sum(nil)

	//Convert the bytes to a string
	hashString = encodeMethod(hashInBytes)

	return hashString, nil
}

func hashFileMd5(filePath string) (string, error) {
	return hashFile(filePath, md5.New, hex.EncodeToString)
}

func hashFileBase64Sha256(filePath string) (string, error) {
	return hashFile(filePath, sha256.New, base64.StdEncoding.EncodeToString)
}

func downloadFile(filepath string, url string) (string, string, string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	etag := resp.Header.Get("etag")

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return "", "", "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", "", "", err
	}

	base64Sha256, err := hashFileBase64Sha256(out.Name())
	if err != nil {
		return "", "", "", err
	}

	md5, err := hashFileMd5(out.Name())
	if err != nil {
		return "", "", "", err
	}

	return base64Sha256, md5, etag, nil
}

func checkIfRemoteFileChanged(url string, oldBase64Sha256 string, etag string) (bool, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	} else if resp.StatusCode == 304 {
		return true, nil
	}

	// If we can't use the etag then download the file and check its hash
	tmpfile, err := ioutil.TempFile("", "*")
	tmpFilePath := tmpfile.Name()

	// Create the file
	out, err := os.Create(tmpFilePath)
	if err != nil {
		return false, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	// Check if the hashes match
	remoteBase64Sha256, err := hashFileBase64Sha256(out.Name())
	if err != nil {
		return false, err
	} else if remoteBase64Sha256 != oldBase64Sha256 {
		return true, nil
	}

	return false, nil
}
