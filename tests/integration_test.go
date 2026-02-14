package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
	"time"
)

const (
	gatewayURL            = "http://localhost:8090/v1/monthly-repayments/calculate"
	readinessPollInterval = 2 * time.Second
	readinessTimeout      = 90 * time.Second
)

var uuidRegexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// GENERATED TEST: DO NOT EDIT
func TestDockerCompose_Migrate_And_Calculate(t *testing.T) {
	root := repoRoot(t)
	composePath := filepath.Join(root, "docker-compose.services.yaml")
	composeCmd := dockerComposeCmd(t, root)

	// Teardown: always run docker compose down
	defer func() {
		down := exec.Command(composeCmd[0], append(composeCmd[1:], "-f", composePath, "down")...)
		down.Dir = root
		if out, err := down.CombinedOutput(); err != nil {
			t.Logf("docker compose down (cleanup): %v\n%s", err, out)
		}
	}()

	// 1. Start stack
	t.Log("Starting docker compose stack...")
	up := exec.Command(composeCmd[0], append(composeCmd[1:], "-f", composePath, "up", "-d", "--build")...)
	up.Dir = root
	if out, err := up.CombinedOutput(); err != nil {
		t.Fatalf("docker compose up: %v\n%s", err, out)
	}

	// 2. Wait for gateway readiness
	t.Log("Waiting for gateway to be ready...")
	if !waitForGateway(t, readinessTimeout) {
		t.Fatal("gateway did not become ready in time")
	}

	// 3. Run migrations
	t.Log("Running migrations...")
	migratePath := filepath.Join(root, "calculation_api", "internal", "providers", "database", "migrations")
	migrateDB := "postgres://postgres:postgres@127.0.0.1:5432/loan?sslmode=disable"
	migrateExe, err := exec.LookPath("migrate")
	if err != nil {
		t.Fatalf("migrate CLI not found on PATH (install golang-migrate): %v", err)
	}
	migrateCmd := exec.Command(migrateExe, "-path", migratePath, "-database", migrateDB, "-verbose", "up")
	migrateCmd.Dir = root
	if out, err := migrateCmd.CombinedOutput(); err != nil {
		t.Fatalf("migrate up: %v\n%s", err, out)
	}

	// 4. POST request with fixed body
	body := []byte(`{"loan_amount": 100000, "interest_rate": 0.055, "number_of_payments": 36}`)
	req, err := http.NewRequest(http.MethodPost, gatewayURL, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST request: %v", err)
	}
	defer resp.Body.Close()

	// 5. Assert status and body
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		ID                     string  `json:"id"`
		MonthlyRepaymentAmount float64 `json:"monthlyRepaymentAmount"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response JSON: %v", err)
	}

	if result.ID == "" {
		t.Fatal("response id is empty")
	}
	if !uuidRegexp.MatchString(result.ID) {
		t.Fatalf("response id is not a valid UUID: %q", result.ID)
	}

	expectedAmount := 3019.59
	tolerance := 0.01
	if result.MonthlyRepaymentAmount < expectedAmount-tolerance || result.MonthlyRepaymentAmount > expectedAmount+tolerance {
		t.Fatalf("monthlyRepaymentAmount: expected ≈ %v, got %v", expectedAmount, result.MonthlyRepaymentAmount)
	}
	t.Logf("id=%s monthlyRepaymentAmount=%.2f", result.ID, result.MonthlyRepaymentAmount)
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, "..")
}

// dockerComposeCmd returns the command and args for docker compose (prefer "docker compose", fallback "docker-compose").
func dockerComposeCmd(t *testing.T, root string) []string {
	t.Helper()
	if path, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker not found on PATH: %v", err)
	} else if path == "" {
		t.Fatal("docker not found on PATH")
	}
	// Prefer "docker compose" (v2)
	cmd := exec.Command("docker", "compose", "version")
	cmd.Dir = root
	if err := cmd.Run(); err == nil {
		return []string{"docker", "compose"}
	}
	// Fallback to docker-compose (v1)
	if path, err := exec.LookPath("docker-compose"); err == nil && path != "" {
		return []string{path}
	}
	t.Fatal("neither 'docker compose' nor 'docker-compose' available")
	return nil
}

func waitForGateway(t *testing.T, timeout time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// POST with minimal body; 400 (validation) or 200 means gateway is up
		req, _ := http.NewRequest(http.MethodPost, gatewayURL, bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			// 200 = success, 400 = validation error (service is up)
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest {
				return true
			}
		}
		time.Sleep(readinessPollInterval)
	}
	return false
}
