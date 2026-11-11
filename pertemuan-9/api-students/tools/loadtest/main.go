package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"api-students/app/model"
	"api-students/config"
	"api-students/helper"
)

func main() {
	mode := flag.String("mode", "crowd", "mode pengujian: crowd | credits | swap")
	baseURL := flag.String("url", "http://localhost:3000/api/v1", "base URL API")
	flag.Parse()

	config.LoadEnv()
	jwtSecret := config.GetEnv("JWT_SECRET", "kunci-rahasia-minimal-32-karakter-harus-aman-banget")
	jwtIssuer := config.GetEnv("JWT_ISSUER", "api-students")
	jwtManager := helper.NewJWTManager(jwtSecret, jwtIssuer, 15*time.Minute)

	client := &http.Client{Timeout: 10 * time.Second}

	switch *mode {
	case "crowd":
		runCrowdTest(client, *baseURL, jwtManager)
	case "credits":
		runCreditsTest(client, *baseURL, jwtManager)
	case "swap":
		runSwapTest(client, *baseURL, jwtManager)
	default:
		fmt.Printf("Mode tidak dikenal: %s\n", *mode)
		os.Exit(1)
	}
}

func makeToken(jwtManager *helper.JWTManager, id int, username string) string {
	token, _ := jwtManager.GenerateAccess(model.User{
		ID:       id,
		Username: username,
		Role:     "user",
	})
	return token
}

func fire(n int, job func(i int) string) map[string]int {
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int, out chan<- string) {
			defer wg.Done()
			<-start
			out <- job(i)
		}(i, results)
	}

	close(start) // buka gerbang serentak
	wg.Wait()
	close(results)

	tally := map[string]int{}
	for r := range results {
		tally[r]++
	}
	return tally
}

func runCrowdTest(client *http.Client, baseURL string, jwt *helper.JWTManager) {
	fmt.Println("60 mahasiswa berebut kelas IF399 secara serentak...")

	start := time.Now()
	// IF399 ID-nya 13 berdasarkan migrasi
	tally := fire(60, func(i int) string {
		token := makeToken(jwt, i+1, fmt.Sprintf("mhs%03d", i+1))
		req, _ := http.NewRequest("POST", baseURL+"/courses/13/enroll", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return "ERR_NET"
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusCreated {
			return "201 OK"
		}

		var res struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(body, &res)
		if res.Code != "" {
			return fmt.Sprintf("%d %s", resp.StatusCode, res.Code)
		}
		return fmt.Sprintf("%d", resp.StatusCode)
	})

	elapsed := time.Since(start)
	fmt.Println("hasil request:")
	for k, v := range tally {
		fmt.Printf("  %-30s %d\n", k, v)
	}
	fmt.Printf("waktu total : %v\n", elapsed)
}

func runCreditsTest(client *http.Client, baseURL string, jwt *helper.JWTManager) {
	fmt.Println("10 mahasiswa, masing-masing mendaftar ke 12 kelas serentak (batas 24 SKS)...")
	start := time.Now()

	// 10 mahasiswa x 12 kelas = 120 request serentak
	tally := fire(120, func(i int) string {
		studentIndex := (i / 12) + 1
		courseID := (i % 12) + 1 // ID 1 sampai 12
		token := makeToken(jwt, studentIndex, fmt.Sprintf("mhs%03d", studentIndex))

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/courses/%d/enroll", baseURL, courseID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return "ERR_NET"
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusCreated {
			return "201 OK"
		}

		var res struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(body, &res)
		if res.Code != "" {
			return fmt.Sprintf("%d %s", resp.StatusCode, res.Code)
		}
		return fmt.Sprintf("%d", resp.StatusCode)
	})

	elapsed := time.Since(start)
	fmt.Println("hasil request:")
	for k, v := range tally {
		fmt.Printf("  %-30s %d\n", k, v)
	}
	fmt.Printf("waktu total : %v\n", elapsed)
}

func runSwapTest(client *http.Client, baseURL string, jwt *helper.JWTManager) {
	fmt.Println("20 pasang mahasiswa pindah kelas berlawanan arah secara serentak (IF101 <-> IF102)...")
	start := time.Now()

	tally := fire(40, func(i int) string {
		studentID := i + 1
		token := makeToken(jwt, studentID, fmt.Sprintf("mhs%03d", studentID))

		var fromID, toID int
		if i%2 == 0 {
			fromID, toID = 1, 2
		} else {
			fromID, toID = 2, 1
		}

		bodyBytes, _ := json.Marshal(map[string]int{
			"from_course_id": fromID,
			"to_course_id":   toID,
		})

		req, _ := http.NewRequest("POST", baseURL+"/enrollments/swap", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return "ERR_NET"
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			return "200 OK"
		}

		var res struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(body, &res)
		if res.Code != "" {
			return fmt.Sprintf("%d %s", resp.StatusCode, res.Code)
		}
		return fmt.Sprintf("%d", resp.StatusCode)
	})

	elapsed := time.Since(start)
	fmt.Println("hasil request:")
	for k, v := range tally {
		fmt.Printf("  %-30s %d\n", k, v)
	}
	fmt.Printf("waktu total : %v\n", elapsed)
}