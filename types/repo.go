package types

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Repository struct {
	Worktree string
	Gitdir   string
}

func (r Repository) ReadObject(hash string) (GitObject, error) {
	path := fmt.Sprintf("%s/objects/%s/%s", r.Gitdir, hash[:2], hash[2:])

	// Read the compressed data from the object file
	compressedData, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading object file:", err)
		return GitObject{}, err
	}

	// Decompress the data --

	// Create a zlib reader
	zlibReader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		fmt.Println("Error creating zlib reader:", err)
		return GitObject{}, err
	}
	defer zlibReader.Close()

	// Read the decompressed object data
	decompressed, err := io.ReadAll(zlibReader)
	if err != nil {
		fmt.Println("Error reading decompressed data:", err)
		return GitObject{}, err
	}
	// ----

	contents := string(decompressed)

	// Parse the header to get the object type and size
	x := strings.Index(contents, " ")
	if x == -1 {
		fmt.Println("Error parsing object header:", err)
		return GitObject{}, err
	}
	format := contents[:x]

	y := strings.Index(contents, "\x00")
	size, err := strconv.Atoi(contents[x+1 : y])
	if err != nil {
		fmt.Println("Malformed object bad length", err)
		return GitObject{}, err
	}

	// validate the size
	if size != len(contents)-y-1 {
		fmt.Println("Malformed object size", err)
		return GitObject{}, err
	}

	return GitObject{
		Type: format,
		Size: size,
		Data: []byte(contents[y+1:]),
	}, nil
}

func (r Repository) WriteObject(g GitObject) (string, error) {
	content := g.Data

	// Create the header for the blob object
	header := []byte(fmt.Sprintf("%s %d\x00", g.Type, len(content)))

	// Create the data by combining the header and content
	data := append([]byte(header), content...)

	// Compute the hash of the object data
	hash := fmt.Sprintf("%x", sha1.Sum(data))

	// Create a buffer to hold compressed data
	var compressedBuffer bytes.Buffer

	// Create a zlib writer
	zlibWriter := zlib.NewWriter(&compressedBuffer)
	_, err := zlibWriter.Write(data)
	if err != nil {
		fmt.Println("Error compressing data:", err)
		return "", err
	}
	zlibWriter.Close()

	// Write compressed data to a file
	parent := fmt.Sprintf("%s/%s/%s", r.Gitdir, "objects", hash[0:2])
	err = os.MkdirAll(parent, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return "", err
	}

	outputFilePath := fmt.Sprintf("%s/%s/%s/%s", r.Gitdir, "objects", hash[0:2], hash[2:])
	err = os.WriteFile(outputFilePath, compressedBuffer.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing compressed data to file:", err)
		return "", err
	}
	return hash, nil
}
